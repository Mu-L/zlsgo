package zhttp

import (
	"bytes"
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/sohaha/zlsgo/zerror"
	"github.com/sohaha/zlsgo/zlog"
	"github.com/sohaha/zlsgo/zstring"
)

type (
	SSEEngine struct {
		ctx        context.Context
		eventCh    chan *SSEEvent
		errCh      chan error
		ctxCancel  context.CancelFunc
		method     string
		readyState int
	}

	SSEEvent struct {
		ID        string
		Event     string
		Undefined []byte
		Data      []byte
	}
)

var (
	delim   = []byte{':', ' '}
	ping    = []byte("ping")
	dataEnd = byte('\n')
)

func (sse *SSEEngine) Event() <-chan *SSEEvent {
	return sse.eventCh
}

func (sse *SSEEngine) Close() {
	sse.ctxCancel()
}

func (sse *SSEEngine) Done() <-chan struct{} {
	return sse.ctx.Done()
}

func (sse *SSEEngine) Error() <-chan error {
	return sse.errCh
}

func (sse *SSEEngine) ResetMethod(method string) {
	sse.method = method
}

func (sse *SSEEngine) OnMessage(fn func(*SSEEvent, error)) {
	for {
		select {
		case <-sse.Done():
			zlog.Debug("sse done")
			return
		case err := <-sse.Error():
			fn(nil, err)
		case v := <-sse.Event():
			fn(v, nil)
		}
	}
}

func SSE(url string, v ...interface{}) *SSEEngine {
	return std.SSE(url, v...)
}

func (e *Engine) sseReq(method, url string, v ...interface{}) (*Res, error) {
	r, err := e.Do(method, url, v...)
	if err != nil {
		return nil, err
	}
	statusCode := r.resp.StatusCode
	if statusCode == http.StatusNoContent {
		return nil, nil
	}

	if statusCode != http.StatusOK {
		return nil, zerror.With(zerror.New(zerror.ErrCode(statusCode), r.String()), "status code is "+strconv.Itoa(statusCode))
	}
	return r, nil
}

func (e *Engine) SSE(url string, v ...interface{}) (sse *SSEEngine) {
	var (
		retry     = 3000
		currEvent = &SSEEvent{}
	)

	ctx, cancel := context.WithCancel(context.TODO())
	sse = &SSEEngine{
		readyState: 0,
		ctx:        ctx,
		ctxCancel:  cancel,
		eventCh:    make(chan *SSEEvent),
		errCh:      make(chan error),
	}

	lastID := ""

	go func() {
		for {
			if sse.ctx.Err() != nil {
				break
			}

			data := v
			if lastID != "" {
				data = append(data, Header{"Last-Event-ID": lastID})
			}
			data = append(data, Header{"Accept": "text/event-stream"})
			data = append(data, Header{"Connection": "keep-alive"})
			data = append(data, sse.ctx)

			r, err := e.sseReq(sse.method, url, data...)

			if err == nil {
				if r == nil {
					sse.readyState = 2
					cancel()
					return
				}

				sse.readyState = 1

				isPing := false
				_ = r.Stream(func(line []byte) error {
					i := len(line)
					if i == 1 && line[0] == dataEnd {
						if !isPing {
							sse.eventCh <- currEvent
							currEvent = &SSEEvent{}
							isPing = false
						} else {
							currEvent = &SSEEvent{}
						}

						return nil
					}

					if i < 2 {
						return nil
					}

					spl := bytes.SplitN(line, delim, 2)
					if len(spl) < 2 {
						return nil
					}

					val := bytes.TrimSuffix(spl[1], []byte{'\n'})

					switch zstring.Bytes2String(spl[0]) {
					case "":
						isPing = bytes.Equal(ping, val)
						if !isPing {
							currEvent.Undefined = val
						}
					case "id":
						lastID = zstring.Bytes2String(val)
						currEvent.ID = lastID
					case "event":
						currEvent.Event = zstring.Bytes2String(val)
					case "data":
						if len(currEvent.Data) > 0 {
							currEvent.Data = append(currEvent.Data, '\n')
						}
						currEvent.Data = append(currEvent.Data, val...)
					case "retry":
						if t, err := strconv.Atoi(zstring.Bytes2String(val)); err == nil {
							retry = t
						}
					}
					return nil
				})
			} else {
				sse.errCh <- err
			}

			sse.readyState = 0
			time.Sleep(time.Millisecond * time.Duration(retry))
		}
	}()

	return
}

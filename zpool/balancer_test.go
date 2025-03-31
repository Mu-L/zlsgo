//go:build go1.18
// +build go1.18

package zpool

import (
	"errors"
	"testing"
	"time"

	"github.com/sohaha/zlsgo"
	"github.com/sohaha/zlsgo/zsync"
	"github.com/sohaha/zlsgo/zutil"
)

func TestNewBalancer(t *testing.T) {
	tt := zlsgo.NewTest(t)
	b := NewBalancer[string]()

	err := b.Run(func(node string) (bool, error) {
		return true, nil
	}, 0)
	tt.NotNil(err)

	err = b.Run(nil, 0)
	tt.NotNil(err)

	b.Add("n1", "n1", func(opts *BalancerNodeOptions) {
		opts.MaxConns = 1
		opts.Weight = 10
	})
	b.Add("n2", "n2", func(opts *BalancerNodeOptions) {
		opts.MaxConns = 2
		opts.Cooldown = 5000
		opts.Weight = 5
	})
	b.Add("n5", "n5", func(opts *BalancerNodeOptions) {
		opts.MaxConns = 5
		opts.Cooldown = 5000
		opts.Weight = 50
	})

	tt.Run("limit", func(tt *zlsgo.TestUtil) {
		var wg zsync.WaitGroup
		success := zutil.NewInt32(0)
		for i := 0; i < 10; i++ {
			wg.Go(func() {
				if err = b.Run(func(node string) (bool, error) {
					tt.Log(node)
					time.Sleep(time.Second / 10)
					return true, nil
				}, 0); err == nil {
					success.Add(1)
				}
			})
		}
		wg.Wait()
		tt.Equal(int32(8), success.Load())
	})

	tt.Run("least", func(tt *zlsgo.TestUtil) {
		var wg zsync.WaitGroup
		success := zutil.NewInt32(0)
		for i := 0; i < 10; i++ {
			wg.Go(func() {
				err = b.Run(func(node string) (bool, error) {
					tt.Log(node)
					time.Sleep(time.Second / 10)
					return true, nil
				}, 1, StrategyLeastConn)
				tt.NoError(err)
				if err == nil {
					success.Add(1)
				}
			})
		}
		wg.Wait()
		tt.Equal(int32(10), success.Load())
	})

	tt.Run("round", func(tt *zlsgo.TestUtil) {
		var wg zsync.WaitGroup
		success := zutil.NewInt32(0)
		nodes := []string{}
		for i := 0; i < 10; i++ {
			err = b.Run(func(node string) (bool, error) {
				tt.Log(node)
				nodes = append(nodes, node)
				return true, nil
			}, 0)
			if err == nil {
				success.Add(1)
			}
		}
		wg.Wait()
		tt.Equal(int32(10), success.Load())
		tt.Equal([]string{"n1", "n2", "n5", "n1", "n2", "n5", "n1", "n2", "n5", "n1"}, nodes)
	})

	tt.Run("error", func(tt *zlsgo.TestUtil) {
		var wg zsync.WaitGroup
		success := zutil.NewInt32(0)
		for i := 0; i < 10; i++ {
			err = b.Run(func(node string) (bool, error) {
				tt.Log(i, node)
				if node == "n5" || node == "n2" {
					return false, errors.New("n error")
				}

				return true, nil
			}, 1, StrategyRandom)
			if err == nil {
				success.Add(1)
			}
		}
		wg.Wait()
		tt.Equal(int32(10), success.Load())
	})
}

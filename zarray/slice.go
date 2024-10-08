//go:build go1.18
// +build go1.18

package zarray

import (
	"math/rand"
	"strings"

	"github.com/sohaha/zlsgo/zstring"
	"github.com/sohaha/zlsgo/zsync"
	"github.com/sohaha/zlsgo/ztype"
	"github.com/sohaha/zlsgo/zutil"
)

// CopySlice copy a slice.
func CopySlice[T any](l []T) []T {
	nl := make([]T, len(l))
	copy(nl, l)
	return nl
}

// Rand A random eents.
func Rand[T any](collection []T) T {
	l := len(collection)
	if l == 0 {
		var zero T
		return zero
	}

	i := zstring.RandInt(0, l-1)
	return collection[i]
}

// Map manipulates a slice and transforms it to a slice of another type.
func Map[T any, R any](collection []T, iteratee func(int, T) R, parallel ...uint) []R {
	colLen := len(collection)
	res := make([]R, colLen)

	if len(parallel) == 0 {
		for i := range collection {
			res[i] = iteratee(i, collection[i])
		}
		return res
	}

	var (
		idx = zutil.NewInt64(0)
		wg  zsync.WaitGroup
	)

	task := func() {
		i := int(idx.Add(1) - 1)
		for ; i < colLen; i = int(idx.Add(1) - 1) {
			res[i] = iteratee(i, collection[i])
		}
	}

	for i := 0; i < int(parallel[0]); i++ {
		wg.Go(task)
	}

	wg.Wait()

	return res
}

// ParallelMap Parallel manipulates a slice and transforms it to a slice of another type.
// If the calculation does not involve time-consuming operations, we recommend using a Map.
// Deprecated: please use Map
func ParallelMap[T any, R any](collection []T, iteratee func(int, T) R, workers uint) []R {
	return Map(collection, iteratee, workers)
}

// Shuffle creates a slice of shuffled values.
func Shuffle[T any](collection []T) []T {
	n := CopySlice(collection)
	rand.Shuffle(len(n), func(i, j int) {
		n[i], n[j] = n[j], n[i]
	})

	return n
}

// Reverse creates a slice of reversed values.
func Reverse[T any](collection []T) []T {
	n := CopySlice(collection)
	l := len(n)
	for i := 0; i < l/2; i++ {
		n[i], n[l-i-1] = n[l-i-1], n[i]
	}

	return n
}

// Filter iterates over eents of collection.
func Filter[T any](slice []T, predicate func(index int, item T) bool) []T {
	slice = CopySlice(slice)

	j := 0
	for i := range slice {
		if !predicate(i, slice[i]) {
			continue
		}
		slice[j] = slice[i]
		j++
	}

	return slice[:j:j]
}

// Contains returns true if an eent is present in a collection.
func Contains[T comparable](collection []T, v T) bool {
	for _, item := range collection {
		if item == v {
			return true
		}
	}

	return false
}

// Find search an eent in a slice based on a predicate. It returns eent and true if eent was found.
func Find[T any](collection []T, predicate func(index int, item T) bool) (res T, ok bool) {
	for i := range collection {
		item := collection[i]
		if predicate(i, item) {
			return item, true
		}
	}

	return
}

// Unique returns a duplicate-free version of an array.
func Unique[T comparable](collection []T) []T {
	repeat := make(map[T]struct{}, len(collection))

	return Filter(collection, func(_ int, item T) bool {
		if _, ok := repeat[item]; ok {
			return false
		}
		repeat[item] = struct{}{}
		return true
	})
}

// Diff returns the difference between two slices.
func Diff[T comparable](list1 []T, list2 []T) ([]T, []T) {
	l, r := []T{}, []T{}

	rl, rr := map[T]struct{}{}, map[T]struct{}{}

	for _, e := range list1 {
		rl[e] = struct{}{}
	}

	for _, e := range list2 {
		rr[e] = struct{}{}
	}

	for _, e := range list1 {
		if _, ok := rr[e]; !ok {
			l = append(l, e)
		}
	}

	for _, e := range list2 {
		if _, ok := rl[e]; !ok {
			r = append(r, e)
		}
	}

	return l, r
}

// Pop returns an eent and removes it from the slice.
func Pop[T comparable](list *[]T) (v T) {
	l := len(*list)
	if l == 0 {
		return
	}

	v = (*list)[l-1]
	*list = (*list)[:l-1]
	return
}

// Shift returns an eent and removes it from the slice.
func Shift[T comparable](list *[]T) (v T) {
	l := len(*list)
	if l == 0 {
		return
	}

	v = (*list)[0]
	*list = (*list)[1:]
	return
}

// Slice converts a string to a slice.
// If n is not empty, the string will be split into n parts.
func Slice[T comparable](s string, n ...int) []T {
	if s == "" {
		return []T{}
	}

	var ss []string
	if len(n) > 0 {
		ss = strings.SplitN(s, ",", n[0])
	} else {
		ss = strings.Split(s, ",")
	}
	res := make([]T, len(ss))
	for i := range ss {
		ztype.To(zstring.TrimSpace(ss[i]), &res[i])
	}

	return res
}

// Join slice to string.
// If n is not empty, the string will be split into n parts.
func Join[T comparable](s []T, sep string) string {
	if len(s) == 0 {
		return ""
	}

	b := zstring.Buffer(len(s))
	for i := 0; i < len(s); i++ {
		v := ztype.ToString(s[i])
		if v == "" {
			continue
		}
		b.WriteString(v)
		if i < len(s)-1 {
			b.WriteString(sep)
		}
	}

	return b.String()
}

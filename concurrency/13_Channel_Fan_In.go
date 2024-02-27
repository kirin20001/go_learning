package concurrency

import (
	"fmt"
	"math/rand"
	"reflect"
	"time"
)

func sigIn(in int) <-chan interface{} {
	c := make(chan interface{})
	go func(in int) {
		defer close(c)
		time.Sleep(time.Duration(rand.Intn(in)) * time.Second)
		c <- in
	}(in)
	return c
}

func fanIn() {

	fanInCha := fanInReflect(
		sigIn(1),
		sigIn(3),
		sigIn(5),
		sigIn(7),
		sigIn(9),
		sigIn(1),
	)

	for v := range fanInCha {
		fmt.Println(v)
	}
}

func fanInReflect(chans ...<-chan interface{}) <-chan interface{} {
	out := make(chan interface{})
	go func() {
		defer close(out)
		// 构造SelectCase slice
		var cases []reflect.SelectCase
		for _, c := range chans {
			cases = append(cases, reflect.SelectCase{
				Dir:  reflect.SelectRecv,
				Chan: reflect.ValueOf(c),
			})
		}

		// 循环，从cases中选择一个可用的
		for len(cases) > 0 {
			i, v, ok := reflect.Select(cases)
			if !ok { // 此channel已经close
				cases = append(cases[:i], cases[i+1:]...)
				continue
			}
			out <- v.Interface()
		}
	}()
	return out
}

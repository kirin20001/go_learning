package concurrency

import (
	"context"
	"fmt"
	"time"
)

func channel1() {
	ch := make(chan int, 4)
	for {
		for i := 1; i < 5; i++ {
			ch <- i
			go func() {
				st := <-ch
				fmt.Println(st)
			}()
			time.Sleep(1 * time.Second)
		}
	}
}

func channel2() {
	var chans [4]chan struct{} // 如果没close channel有什么影响，最终会被gc
	stopCh := make(chan struct{})
	for i, _ := range chans {
		chans[i] = make(chan struct{})
	}

	// receiver
	for i := 0; i < 4; i++ {
		go func(i int) {
			for {
				select {
				case <- chans[i]:
					fmt.Println(i, time.Now())
					time.Sleep(1 * time.Second)
					chans[(i+1)%4] <- struct{}{}
				case <- stopCh: // 如果在第一个case无法接收stopChannel
					fmt.Println(i, "terminated")
					return
				}
			}
		}(i)
	}

	chans[0] <- struct{}{}
	select {
	case <-time.After(3 * time.Second):
		close(stopCh)
	}
	time.Sleep(1 * time.Second)
	//for i, _ := range chans {
	//	close(chans[i])
	//}
}


func channel2_1() {
	var chans [4]chan struct{}
	stopCh := make(chan struct{})
	for i, _ := range chans {
		chans[i] = make(chan struct{})
	}

	for i := 0; i < 4; i++ {
		go func(i int) {
			for {
				select {
				case <- chans[i]:
					fmt.Println(i, time.Now())
					select {
					case <- time.After(1 * time.Second):
						chans[(i+1)%4] <- struct{}{}
					case <- stopCh:
						fmt.Println(i, "terminated")
						return
					}
				case <- stopCh: // 如果在第一个case无法接收stopChannel
					fmt.Println(i, "terminated")
					return
				}
			}
		}(i)
	}

	chans[0] <- struct{}{}
	select {
	case <-time.After(5 * time.Second):
		close(stopCh)
	}
	time.Sleep(1 * time.Second)
	for i, _ := range chans {
		close(chans[i])
	}
}

func channel3() {
	var chans [4]chan struct{}
	for i, _ := range chans {
		chans[i] = make(chan struct{})
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	for i := 0; i < 4; i++ {
		go func(ctx context.Context, i int,) {
			for {
				select {
				case <- chans[i]:
					fmt.Println(i, time.Now())
					select {
					case <- time.After(1 * time.Second):
						chans[(i+1)%4] <- struct{}{}
					case <- ctx.Done():
						fmt.Println(i, "terminated")
						return
					}
				case <- ctx.Done(): // 如果在第一个case无法接收stopChannel
					fmt.Println(i, "terminated")
					return
				}
			}
		}(ctx, i)
	}

	chans[0] <- struct{}{}
	time.Sleep(5 * time.Second)
	cancel()
	for i, _ := range chans {
		close(chans[i])
	}
}
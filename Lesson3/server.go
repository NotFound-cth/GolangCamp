package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"golang.org/x/sync/errgroup"
)

/*

1. 基于 errgroup 实现一个 http server 的启动和关闭 ，以及 linux signal 信号的注册和处理，要保证能够一个退出，全部注销退出。

*/
type svr struct {
	start  func() error
	stop   func() error
	cancel func()
}

func serverApp() error {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(resp http.ResponseWriter, req *http.Request) {
		fmt.Fprintln(resp, "hello world")
	})
	http.ListenAndServe("0.0.0.0:8080", mux)
	return nil
}

func serverDebug() error {
	http.ListenAndServe("0.0.0.0:8081", http.DefaultServeMux)
	return nil
}

func serverStop() error {
	fmt.Println("server exit")
	return nil
}

func debugStop() error {
	fmt.Println("debug exit")
	return nil
}

func run() error {
	var fg [2]svr
	fg[0].start = serverApp
	fg[0].stop = serverStop
	fg[1].start = serverDebug
	fg[1].stop = debugStop
	ctx0, cancel := context.WithCancel(context.Background())
	g, ctx := errgroup.WithContext(ctx0)
	for i := 0; i < 2; i++ {
		s := fg[i]
		g.Go(func() error {
			return s.start()
		})
		g.Go(func() error {
			<-ctx.Done()
			return s.stop()
		})
	}
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT)
	g.Go(func() error {
		for {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-c:
				cancel()
				log.Fatal("exit")
				return nil
			}
		}
	})
	if err := g.Wait(); err != nil && !errors.Is(err, context.Canceled) {
		return err
	}
	return nil
}

func main() {
	run()
}

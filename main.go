package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/RicardoLinck/simple-pool/server"

	"golang.org/x/sync/errgroup"
)

func main() {
	apiConfig := server.NewAPIConfig(server.GenerateSampleItems())
	s := apiConfig.Init()
	go http.ListenAndServe("localhost:3000", s)

	numRequests := 30
	results := make(chan *http.Response)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	go writeResponses(results)
	ctx, work, group := initPool(ctx, numRequests, results)

	for i := 0; i < numRequests; i++ {
		work <- fmt.Sprintf("http://localhost:3000/?count=%d", i)
	}
	close(work)

	err := group.Wait()
	if err != nil {
		fmt.Printf("Request error: %v\n", err)
	}
	close(results)
	fmt.Println("Finished")
}

func writeResponses(results <-chan *http.Response) {
	for resp := range results {
		b, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Printf("Response error: %v", err)
			continue
		}
		fmt.Println(string(b))
		resp.Body.Close()
	}
}

func initPool(ctx context.Context, size int, results chan<- *http.Response) (context.Context, chan<- string, *errgroup.Group) {
	c := make(chan string, size)
	g, ctx := errgroup.WithContext(ctx)

	for i := 0; i < 3; i++ {
		i := i
		g.Go(func() error {
			return worker(ctx, i, c, results)
		})
	}

	return ctx, c, g
}

func worker(ctx context.Context, num int, work <-chan string, results chan<- *http.Response) error {
	for url := range work {
		req, err := http.NewRequestWithContext(ctx, "GET", url, nil)

		if err != nil {
			return err
		}

		fmt.Printf("Requesting %s via worker %d\n", url, num)
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return err
		}
		results <- resp
	}
	return nil
}

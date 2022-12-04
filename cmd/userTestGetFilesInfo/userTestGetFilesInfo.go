package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"sync"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pbapi "grpc-api/pkg/api"
)

func main() {

	address := ":8080"
	countGoroutines := 103

	flag.StringVar(&address, "a", address, "grpc api server address")
	flag.IntVar(&countGoroutines, "c", countGoroutines, "count goroutines")

	flag.Parse()

	conn, err := grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	client := pbapi.NewFileMngrClient(conn)

	wg := sync.WaitGroup{}

	for i := 0; i < countGoroutines; i++ {

		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			res, err := client.GetFilesInfo(context.Background(), &pbapi.GetFilesInfoRequest{})
			if err != nil {
				fmt.Printf("goroutine %d:\n   err: %s\n\n", i, err)
			} else {
				fmt.Printf("goroutine %d:\n   result: %s\n", i, res)
			}
		}(i)

	}

	wg.Wait()

}

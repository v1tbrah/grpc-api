package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"sync"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pbapi "grpc-api/pkg/api"
)

func main() {

	address := ":8080"
	dirWithFiles := "filesForSendingToGRPCServer"

	flag.StringVar(&address, "a", address, "grpc api server address")
	flag.StringVar(&dirWithFiles, "d", dirWithFiles, "dir with files for sending to server")

	flag.Parse()

	conn, err := grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	client := pbapi.NewFileMngrClient(conn)

	wg := sync.WaitGroup{}

	dirData, err := os.ReadDir(dirWithFiles)
	if err != nil {
		log.Fatal(err)
	}

	for _, val := range dirData {

		wg.Add(1)
		go func(val os.DirEntry) {
			defer wg.Done()

			data, err := readFile(dirWithFiles + "/" + val.Name())
			if err != nil {
				log.Fatal(err)
			}

			_, err = client.SaveFile(context.Background(), &pbapi.SaveFileRequest{Name: val.Name(), Data: data})
			if err != nil {
				fmt.Printf("saving file %s:\n   err: %s\n\n", val.Name(), err)
			} else {
				fmt.Printf("file %s saved\n", val.Name())
			}

		}(val)
	}

	wg.Wait()

}

func readFile(name string) ([]byte, error) {
	f, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	fr := bufio.NewReader(f)
	data, err := io.ReadAll(fr)
	if err != nil {
		return nil, err
	}

	return data, err
}

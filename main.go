package main

import (
	"fmt"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/tsemach/go-fuse/server"
	// "gitlab.mobileye.com/iot/iot-upload-go/config"
	// cfg "gitlab.mobileye.com/iot/iot-upload-go/config"
)

var port = 8080

// func main_init() {
// 	var wg sync.WaitGroup

// 	config.BuildConfig()
// 	wg.Add(1)
// 	go iotaws.AWS.Init(&wg)
// 	wg.Add(1)
// 	go ops.OS.Connect(&wg)
// 	wg.Add(1)
// 	go redis.Redis.Connect(&wg)
// 	wg.Add(1)
// 	go upload.Upload.Init(&wg)
// 	wg.Wait()
// }

func createServer() *http.Server {
	r := server.CreateGin()

	return &http.Server{
		Addr:    fmt.Sprintf(":%v", port),
		Handler: r,

		ReadTimeout:  60 * time.Second,
		WriteTimeout: 60 * time.Second,
	}
}

func main() {
	// main_init()
	// config.BuildConfig()

	mydir, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(mydir)

	wg := new(sync.WaitGroup)
	wg.Add(1)

	// goroutine to launch a server on port 8080
	go func() {
		server := createServer()
		fmt.Println(server.ListenAndServe())
		wg.Done()
	}()

	// wait until WaitGroup is done
	wg.Wait()
}

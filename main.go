package main

import (
	"fmt"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/tsemach/go-fuse/common"
	"github.com/tsemach/go-fuse/config"
	"github.com/tsemach/go-fuse/fusefs"
	"github.com/tsemach/go-fuse/server"
)

var port = 8080

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
	config.BuildConfig(common.First(config.GetConfigPath("fusefs.yaml")))
	for i := 0; i < len(config.GetConfig().Filesystems); i++ {
		mountpoint := config.GetConfig().Filesystems[i].Mountpoint
		targetpath := config.GetConfig().Filesystems[i].Targetpath

		go fusefs.CreateFuseFSWatchDog(mountpoint, targetpath)
	}

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

package fusefs

import (
	"fmt"
	"log"
	"os"
	"time"
)

func CreateFuseFS(mountpoint string, targetpath string) {
	log.Println("[CreateFuseFS] pathpoint:", mountpoint, "targetpath:", targetpath)

		if !isDirectoryExist(targetpath) {
		log.Fatalln("Error unable to find target directory:", targetpath)
		return
	}

	if !isDirectoryExist(mountpoint) {
		err := os.MkdirAll(mountpoint, os.ModePerm)
		if err != nil {
			log.Fatalln("Error creating mountpoint directories:", err)
			return
		}
	}

	_, err := os.Stat(mountpoint)
	if os.IsNotExist(err) {
		log.Fatalf("failed to create %s: %v", mountpoint, err)
		os.Exit(2)
	}
	log.Println("going to mount on:", mountpoint)

	fs, err := NewFuseFS(mountpoint, targetpath)
	if err != nil {
		log.Fatalf("failed to new filesystem: %s", err)
		os.Exit(3)
	}

	if err = fs.Mount(); err != nil {
		log.Fatalf("failed to mount: %s", err)
	}

	fmt.Printf("mount fs: %v\n", fs)

	if err = fs.Serve(); err != nil {
		log.Fatalf("failed to serve: %s", err)
	}

	// if err = fs.Unmount(); err != nil {
	// 	log.Fatalf("failed to unmount: %s", err)
	// }
}

func CreateFuseFSWatchDog(mountpoint string, targetpath string) {
	for {
		func() {
			defer func() {
				if r := recover(); r != nil {
					log.Fatalln("CreateFuseFS panicked, restarting:", r)
				}
			}()

			CreateFuseFS(mountpoint, targetpath)

		}()

		time.Sleep(1 * time.Second)
	}
}
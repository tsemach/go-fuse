package fusefs

import (
	"fmt"
	"log"
	"os"
)

func CreateFuseFS() {
	path := "/tmp/fusefs"

	fmt.Println("path:", path)

	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		log.Fatalf("failed to create %s: %v", path, err)
		os.Exit(2)
	}
	log.Println("going to mount on:", path)

	fs, err := NewFuseFS(path)
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

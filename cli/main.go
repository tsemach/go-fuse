package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/tsemach/go-fuse/common"
	"github.com/tsemach/go-fuse/config"
)

var (
	configFileName string
)

func main() {
	flag.StringVar(&configFileName, "config", "fusefs.yaml", "fusefs config file")
	flag.Parse()

	fmt.Println("config is:", configFileName)
	fmt.Println(common.First(os.Getwd()))
	config.BuildConfig(common.GetRootDir()+"/"+configFileName)
	fmt.Println(config.GetConfig().Filesystems)
	fmt.Println(config.GetConfig().Stam)


	// file, err := os.Open("/tmp/fusefs-1/dir-1/file-11")

	// if err != nil {
	// 	fmt.Println("ERROR opening file: /tmp/fusefs-1/dir-1/file-11")
	// }

	// defer file.Close();

	// fmt.Println("file /tmp/fusefs-1/dir-1/file-11 open ok")
}

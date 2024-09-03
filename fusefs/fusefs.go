package fusefs

import (
	"fmt"
	"os"

	"bazil.org/fuse"
	"bazil.org/fuse/fs"
)

type FuseFS interface {
	Mount() error
	Serve() error
	Unmount() error
	Mountpoint() string

	fs.FS
	fs.FSInodeGenerator
}

type fuseFS struct {
	mountpoint string
	conn       *fuse.Conn
	rootNode   FuseFSNode
	lastInode  uint64
}

func NewFuseFS(mountpoint string) FuseFS {
	rfs := &fuseFS{mountpoint: mountpoint, lastInode: 1}
	rfs.rootNode = &fuseFSNode{
		FS:    rfs,
		Inode: 1,
		Mode:  os.ModeDir | 0o555,
	}
	return rfs
}

func (rfs *fuseFS) Mount() error {
	c, err := fuse.Mount(
		rfs.mountpoint,
		fuse.FSName("fusefs"),
		fuse.Subtype("fusefs"),
	)
	
	if err != nil {
		return err
	}
	rfs.conn = c

	return nil
}

func (rfs *fuseFS) Serve() error {
	server := fs.New(rfs.conn, &fs.Config{Debug: func(msg interface{}) { fmt.Println(msg) }})
	return server.Serve(rfs)
}

func (rfs *fuseFS) Unmount() error {
	err := fuse.Unmount(rfs.mountpoint)
	if err != nil {
		return err
	}

	return rfs.conn.Close()
}

func (rfs *fuseFS) Mountpoint() string {
	return rfs.mountpoint
}

/* fs.FS */
func (rfs fuseFS) Root() (fs.Node, error) {
	return rfs.rootNode, nil
}

/* fs.FSInodeGenerator */
func (rfs fuseFS) GenerateInode(parentInode uint64, name string) uint64 {
	rfs.lastInode++
	return rfs.lastInode
}

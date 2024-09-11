package fusefs

import (
	"fmt"
	"log"
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
	rootNode   *FuseFSNode // FuseFSNode
	lastInode  uint64
}

func NewFuseFS(mountpoint string, targetpath string) (FuseFS, error) {
	rfs := &fuseFS{mountpoint: mountpoint, lastInode: 1}
	rfs.rootNode = &FuseFSNode{
		FS:    rfs,
		Inode: 1,
		Path: targetpath,
		Mode:  os.ModeDir | 0o555,
	}

	// nodes, err := buildFSNodesTree("/home/tsemach/tmp/fusefs", rfs, "", rfs.rootNode.Inode)
	nodes, err := buildFSNodesTree(targetpath, rfs, "", rfs.rootNode.Inode)
	if err != nil {
		log.Fatalln("error building fs nodes tree", err)

		return nil, err
	}

	rfs.rootNode.Nodes = nodes

	// rfs.rootNode.Nodes[0] =  &fuseFSNode{
	// 	FS:    rfs,
	// 	Name:  "stam-file",
	// 	Path:  "/",
	// 	Inode: 1,
	// 	Mode:  420,
	// }

	// var nodes = make([]*fuseFSNode, 2)
	// nodes[0] = &fuseFSNode{
	// 	FS:    rfs,
	// 	Name:  "auto-file",
	// 	Path:  "/",
	// 	Inode: 1,
	// 	Mode:  420,
	// }

	// nodes[1] = &fuseFSNode{
	// 	FS:    rfs,
	// 	Name:  "auto-dir",
	// 	Path:  "/",
	// 	Inode: 1,
	// 	Mode:  0x800001ed,
	// }

	// rfs.rootNode.Nodes = append(rfs.rootNode.Nodes, &fuseFSNode{
	// 	FS:    rfs,
	// 	Name:  "auto-file",
	// 	Path:  "/",
	// 	Inode: 1,
	// 	Mode:  420,
	// })

	// rfs.rootNode.Nodes = append(rfs.rootNode.Nodes, &fuseFSNode{
	// 	FS:    rfs,
	// 	Name:  "auto-dir",
	// 	Path:  "/",
	// 	Inode: 1,
	// 	Mode:  0x800001ed,
	// })

	// 2147484141 = 0x800001ed

	// var nodes = append(rfs.rootNode.Nodes, &fuseFSNode{
	// 	FS:    rfs,
	// 	Name:  "stam-file-2",
	// 	Path:  "/",
	// 	Inode: 1,
	// 	Mode:  420,
	// })

	// rfs.rootNode.Nodes = nodes
	// rfs.rootNode.Nodes = make([]*fuseFSNode, 0)

	return rfs, nil
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
func (rfs *fuseFS) Root() (fs.Node, error) {
	return rfs.rootNode, nil
}

/* fs.FSInodeGenerator */
func (rfs *fuseFS) GenerateInode(parentInode uint64, name string) uint64 {
	rfs.lastInode++
	return rfs.lastInode
}

package fusefs

import (
	"context"
	"fmt"
	"os"
	"syscall"

	"bazil.org/fuse"
	"bazil.org/fuse/fs"
)

type FuseFSNode interface {
	fs.Node
	// fs.NodeGetattrer
	fs.NodeSetattrer
	// fs.NodeSymlinker
	// fs.NodeReadlinker
	// fs.NodeLinker
	fs.NodeRemover
	// fs.NodeAccesser
	fs.NodeStringLookuper
	fs.NodeMkdirer
	// fs.NodeOpener <-
	fs.NodeCreater
	// fs.NodeForgetter
	// fs.NodeRenamer
	// fs.NodeMknoder
	// fs.NodeFsyncer
	fs.NodeGetxattrer
	// fs.NodeListxattrer
	// fs.NodeSetxattrer
	// fs.NodeRemovexattrer
	// fs.NodePoller // fs.HandlePoller <-

	// fs.HandleFlusher <-
	// fs.HandleReadAller
	fs.HandleReadDirAller
	// fs.HandleReader
	fs.HandleWriter
	// fs.HandleReleaser <-
}

func NewFuseFSNode() FuseFSNode {
	return &fuseFSNode{}
}

type fuseFSNode struct {
	FS     FuseFS
	Name   string
	Path 	 string	
	Inode  uint64
	Mode   os.FileMode
	Nodes  []*fuseFSNode
	Data   []byte
	Xattrs map[string]string
}

// fs.Node */
func (n fuseFSNode) Attr(ctx context.Context, attr *fuse.Attr) error {
	attr.Inode = n.Inode
	attr.Mode = n.Mode
	attr.Size = uint64(len(n.Data))
	
	return nil
}

func (n *fuseFSNode) Setattr(ctx context.Context, req *fuse.SetattrRequest, resp *fuse.SetattrResponse) error {
	// NOTE: res.Atrr is filled by Attr method
 
	if req.Mode&os.ModeIrregular != 0 {
		fmt.Println("call to Setattr with mode irregular")
		return nil
	}

	n.Mode = req.Mode
	return nil
}

func (n *fuseFSNode) Open(ctx context.Context, req *fuse.OpenRequest, resp *fuse.OpenResponse) (fs.Handle, error) {
	return n, nil
}

// fs.NodeRemover
func (n *fuseFSNode) Remove(ctx context.Context, req *fuse.RemoveRequest) error {
	for i, node := range n.Nodes {
		if node.Name == req.Name {
			// TODO: Test if rmdir fills req.Dir
			if req.Dir {
				if !node.Mode.IsDir() {
					return syscall.ENOTDIR
				}
				if len(req.Name) != 0 && req.Name[len(req.Name)-1] == '.' {
					return syscall.EINVAL
				}
			} else {
				if node.Mode.IsDir() {
					return syscall.EISDIR
				}
			}

			n.Nodes = append(n.Nodes[:i], n.Nodes[i+1:]...)
			return nil
		}
	}
	return syscall.ENOENT
}

// fs.NodeStringLookuper
func (n fuseFSNode) Lookup(ctx context.Context, name string) (fs.Node, error) {
	for _, n := range n.Nodes {
		if n.Name == name {
			return n, nil
		} else if n.Mode.IsDir() {
			// TODO: Check if this is needed
			if lookupNode, err := n.Lookup(ctx, name); err == nil {
				return lookupNode, nil
			}
		}
	}
	return nil, syscall.ENOENT
}

func (n *fuseFSNode) Create(ctx context.Context, req *fuse.CreateRequest, res *fuse.CreateResponse) (fs.Node, fs.Handle, error) {
	if !n.Mode.IsDir() {
		return nil, nil, syscall.ENOTDIR
	}

	// var path string
	
	// if n.Path != "" {
	// 	path = n.Path+"/"+req.Name
	// } else {
	// 	path = req.Name
	// }

	newNode := &fuseFSNode{
		FS:    n.FS,
		Name:  req.Name,
		Path:  n.Path,
		Inode: n.FS.GenerateInode(n.Inode, req.Name),
		Mode:  req.Mode,
	}
	n.Nodes = append(n.Nodes, newNode)
	return newNode, newNode, nil
}

// fs.NodeMkdirer
func (n *fuseFSNode) Mkdir(ctx context.Context, req *fuse.MkdirRequest) (fs.Node, error) {
	if !n.Mode.IsDir() {
		return nil, syscall.ENOTDIR
	}

	var path string

	if n.Path != "" {
		path = n.Path+"/"+req.Name
	} else {
		path = req.Name
	}

	newNode := &fuseFSNode{
		FS:    n.FS,
		Name:  req.Name,
		Path:  path,
		Inode: n.FS.GenerateInode(n.Inode, req.Name),
		Mode:  req.Mode,
	}
	n.Nodes = append(n.Nodes, newNode)
	return newNode, nil
}

// fs.NodeGetxattrer
func (n fuseFSNode) Getxattr(ctx context.Context, req *fuse.GetxattrRequest, res *fuse.GetxattrResponse) error {
	// NOTE: req.Size is the size of res.Xattr. Size check is performed by fuse library

	if n.Xattrs == nil {
		return syscall.ENODATA
	}

	value, found := n.Xattrs[req.Name]
	if !found {
		return syscall.ENODATA
	}

	res.Xattr = []byte(value)
	return nil
}

// fs.HandleReadDirAller
func (n *fuseFSNode) ReadDirAll(ctx context.Context) ([]fuse.Dirent, error) {
	ents := make([]fuse.Dirent, len(n.Nodes))
	for i, node := range n.Nodes {
		typ := fuse.DT_File
		if node.Mode.IsDir() {
			typ = fuse.DT_Dir
		}
		ents[i] = fuse.Dirent{Inode: node.Inode, Type: typ, Name: node.Name}
	}
	return ents, nil
}

// fs.HandleWriter
/*
Unused fields
	type WriteRequest struct {
		Handle    HandleID
		Offset    int64
		Flags     WriteFlags
		LockOwner LockOwner
	}
*/

func (n *fuseFSNode) Read(ctx context.Context, req *fuse.ReadRequest, resp *fuse.ReadResponse) error {
	if req.Offset >= int64(len(n.Data)) {
			return nil
	}

	data := n.Data[req.Offset:]
	// dirname, errno := getHomeDir()
	// if errno != 0 {
	// 	return errno
	// }

	// filesize, err := getFileSize(dirname+"/"+n.Name)
	// if err != nil {
	// 	fmt.Println("[Read] unable to get size of file:", dirname+"/"+n.Name, "err:", err)
	// 	return nil
	// }

	// file, err := os.Open(dirname+"/"+n.Name)
	// if err != nil {
	// 	fmt.Println("[Read] error opening file:", dirname+"/"+n.Name, "err:", err)
	// 	return nil
	// }
	// defer file.Close()

	// _, err = file.Seek(req.Offset, 0)
	// if err != nil {
	// 	fmt.Println("[Read] error seeking to offset:", err)
	// 	return nil
	// }

	// data := make([]byte, filesize)
	// nbytes, err := file.Read(data)
	// if err != nil {
	// 	fmt.Println("Error reading file:", err)
	// 	return syscall.EBADR
	// }

	// if (int64(nbytes) < filesize) {
	// 	fmt.Println("[Read] not enough bytes read, nbytes:", nbytes, "req.Size:", req.Size)
	// 	return syscall.EBADR
	// }

	if len(data) > req.Size {
		data = data[:req.Size]
	}

	resp.Data = data
	return nil
}

func (n *fuseFSNode) Write(ctx context.Context, req *fuse.WriteRequest, res *fuse.WriteResponse) error {
	if req.FileFlags.IsReadOnly() {
		return syscall.EBADF
	}

	// TODO: Get request GID+UID and file UID+GID and check if user or group is allowed to write to the file. If not return EPERM
	// dirname, errno := getHomeDir()
	// if errno != 0 {
	// 	return errno
	// }

	// os.WriteFile(dirname+"/"+n.Name, req.Data, n.Mode)
	n.Data = req.Data
	res.Size = len(req.Data)

	return nil
}

func (n *fuseFSNode) getNodeDir() string {
	return n.FS.Mountpoint() + "/" + n.Path
}

func getHomeDir() (string, syscall.Errno) {
	dirname, err := os.UserHomeDir()
	if err != nil {
		return "", syscall.EBADR
	}

	return dirname + "/tmp/fusefs", syscall.F_OK
}

func getFileSize(filepath string) (int64, error) {
	fileInfo, err := os.Stat(filepath)

	if err != nil {
		fmt.Println("Error:", err)
		return -1, err
	}

	return fileInfo.Size(), nil
}

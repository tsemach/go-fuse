package fusefs

import (
	"log"
	"os"
)

/*
 * read the directory contents (non-recursive)
 */
func buildFSNodesTree(rootDir string, rfs *fuseFS, path string, parentInode uint64) ([]*fuseFSNode, error) {	
	entries, err := os.ReadDir(getPath(rootDir, path))
	if err != nil {
		log.Println(" [INFO] [buildFSNodesTree] error reading directory:", err)
		return nil, err
	}

	var nodes = make([]*fuseFSNode, 0)
	for _, entry := range entries {
		if entry.IsDir() {
			log.Println("[INFO] [buildFSNodesTree]", entry.Name())

			info, _ := entry.Info()
			node := &fuseFSNode{
				FS:    rfs.rootNode.FS,
				Name:  entry.Name(),
				Path:  path,
				Inode: rfs.GenerateInode(parentInode, entry.Name()),
				Mode:  info.Mode(),
			}

			entryNodes, err := buildFSNodesTree(rootDir, rfs, getPath(path, entry.Name()), node.Inode)
			if err != nil {
				log.Fatalln("[ERROR] [buildFSNodesTree] unable to build node tree of", path+"/"+entry.Name())
			}
			node.Nodes = entryNodes

			nodes = append(nodes, node)
		}
	}

	return nodes, nil
}

func getPath(path string, name string) string {
	if path == "" {
		return name
	}
	return path + "/" + name
}

// func getDirectoriesCount(rootDir string /*, nodes *fuseFSNode*/) (int, error) {
// 	// Walk the directory tree and process directories only
// 	var count int
// 	err := filepath.WalkDir(rootDir, func(path string, d os.DirEntry, err error) error {
// 		if err != nil {
// 			return err
// 		}
// 		// Check if it's a directory
// 		if d.IsDir() {
// 			count++
// 		}
// 		return nil
// 	})

// 	return count, err
// }

// append(rfs.rootNode.Nodes, &fuseFSNode{
// 	FS:    rfs,
// 	Name:  "stam-file",
// 	Path:  "/",
// 	Inode: 1,
// 	Mode:  420,
// })
// return rfs

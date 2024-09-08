# go-fuse

A go program that implement a fuse file system control by a web api.

## Cache Managment
Cache is manage by score calculate by the time inode is access or not access and the size of its data buffer
Cache is a link list of items. 
Item is:
````golang
type struct {
  score: int
  inode: *fuseFSNode
}
````
1. cache link list is sort by score in ascending order.
2. score access time is calculate by:

   Current time - last access time
   ------------------------------- ==> this will result number berween 0 - 100 on a period of 10 secs.
            100 miliseconds

3. score data size is calculate by defining thresholds:
   0    - 10K  => 
   10K  - 128K => 
   128K - 512K => 
   512K - 1M   =>    

   *Note: file that are larger then 1M do not cached.
4. file that is not use for 10 sec (configurable) is remove from cache.
5. file that are larger then 1M do not cache at all.
6. remove file from cache will result remove it inode from fuseFS

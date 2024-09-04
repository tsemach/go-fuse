# go-fuse

A go program that implement a fuse file system control by a web api.

In a FUSE filesystem, each method (like `Lookup`, `ReadAll`, etc.) is called in response to system calls made by the operating system when users or applications interact with the mounted filesystem. FUSE acts as an intermediary between the kernel and your userspace program, handling requests for file operations. Let's break down how it works:

### 1. **Mounting the Filesystem**
When you run your Go FUSE program and mount the filesystem (e.g., `./hellofs /tmp/mountpoint`), the FUSE library registers your filesystem with the operating system, telling it that any file-related operations under `/tmp/mountpoint` should be handled by your FUSE filesystem.

### 2. **System Calls and Kernel Interaction**
When a user or program interacts with the mounted filesystem (e.g., `ls /tmp/mountpoint` or `cat /tmp/mountpoint/hello.txt`), the operating system makes **system calls** like `open`, `read`, `write`, `stat`, etc., to access files and directories.

These system calls are intercepted by the FUSE kernel module, which then forwards the request to your Go FUSE daemon (your Go program). The daemon translates these requests into Go method calls, like `Lookup`, `ReadDirAll`, `ReadAll`, and others, depending on the operation.

### 3. **Method Dispatching in Go**
The FUSE library (`bazil.org/fuse`) dispatches incoming filesystem requests to the appropriate methods in your Go code. For example:

- **Directory Listing (`ls`)**: When a user runs `ls /tmp/mountpoint`, the kernel asks your FUSE filesystem for the contents of the directory. The `ReadDirAll` method is called to provide a list of files and subdirectories.
  
- **File Lookup (`open` or `cat`)**: When the system tries to access a specific file (e.g., `cat /tmp/mountpoint/hello.txt`), the kernel first asks your filesystem to look up that file in the directory using the `Lookup` method. If `Lookup` returns a valid node (like a `File`), the kernel continues with the operation.

### Example of the Lifecycle of a `Lookup` Call:

Suppose you run the following command in the terminal:

```bash
cat /tmp/mountpoint/hello.txt
```

This triggers a sequence of events:

1. **The `cat` command calls `open` on `hello.txt`**:
   - The operating system looks for the file `hello.txt` in `/tmp/mountpoint`.
   - The FUSE kernel module intercepts this call and asks your FUSE filesystem to "look up" the file using the `Lookup` method.

   In your Go code, the `Lookup` method of the `Dir` struct gets called:
   ```go
   func (Dir) Lookup(ctx context.Context, name string) (fs.Node, error) {
       if name == "hello.txt" {
           return File{}, nil
       }
       return nil, fuse.ENOENT
   }
   ```
   - If the file is found (`name == "hello.txt"`), it returns a `File` object representing the file. Otherwise, it returns `ENOENT` (file not found).

2. **The kernel asks for the file's attributes**:
   - Once the `File` node is returned by `Lookup`, the kernel asks for the file's attributes (like size, permissions) by calling the `Attr` method on the `File` struct.
   
   This happens in your Go code:
   ```go
   func (File) Attr(ctx context.Context, a *fuse.Attr) error {
       a.Inode = 2
       a.Mode = 0444 // Read-only file
       a.Size = uint64(len("Hello, world!\n"))
       return nil
   }
   ```

3. **Reading the file (`read` system call)**:
   - After opening the file, the kernel now requests to read its contents by calling the `ReadAll` method of the `File` struct. This method returns the file content.

   In your Go code:
   ```go
   func (File) ReadAll(ctx context.Context) ([]byte, error) {
       return []byte("Hello, world!\n"), nil
   }
   ```

   The kernel receives the data, and the `cat` command outputs "Hello, world!" in the terminal.

### Breakdown of Method Calls:

- **`Lookup`**: Called when the system needs to locate a file or directory by name. This is typically invoked during operations like `open`, `stat`, or `ls`.
  
- **`Attr`**: Called to retrieve the attributes (metadata) of a file or directory. This method is used for both files and directories to report information like permissions, size, inode number, etc.

- **`ReadDirAll`**: Called when the system requests the contents of a directory, such as when you run `ls`. It returns a list of directory entries (files, directories, etc.).

- **`ReadAll`**: Called when the system reads the entire content of a file. This method is used to return the file's data when a file is read (e.g., `cat` or `less`).

### How FUSE Handles These Calls Internally:

1. **File Operations**: The FUSE library listens for incoming file operation requests from the kernel.
2. **Dispatching to Methods**: The FUSE library inspects the request type (e.g., `open`, `read`, `write`, `stat`) and dispatches the corresponding method (`Lookup`, `ReadDirAll`, `Attr`, etc.) to handle the request.
3. **Return Values**: After handling the request in your Go code, you return data (e.g., file content or attributes) to the kernel, which then continues the file operation (e.g., displaying content or listing files).

### Summary of Events for a `cat hello.txt` Command:

1. Kernel calls `Lookup("hello.txt")` → Go code returns `File`.
2. Kernel calls `Attr(File)` → Go code returns file attributes (size, permissions).
3. Kernel calls `ReadAll(File)` → Go code returns file data (`Hello, world!\n`).

### Serving Multiple Requests
In a real-world application, many of these requests (like `Lookup`, `Attr`, `ReadAll`, etc.) are served concurrently by FUSE. The `fs.Serve` function handles the dispatching of multiple filesystem requests to your Go methods as needed.

This is the basic flow of how the file system methods like `Lookup` are called in a FUSE-based filesystem written in Go.
package main

import (
	"fmt"
	"os"
)

// A FileInfo describes a file and is returned by Stat and Lstat.
// type FileInfo interface {
// 	Name() string       // base name of the file
// 	Size() int64        // length in bytes for regular files; system-dependent for others
// 	Mode() FileMode     // file mode bits
// 	ModTime() time.Time // modification time
// 	IsDir() bool        // abbreviation for Mode().IsDir()
// 	Sys() interface{}   // underlying data source (can return nil)
// }

func main() {
	//The file has to be opened first
	f, err := os.Open("test.txt")
	// The file descriptor (File*) has to be used to get metadata
	fi, err := f.Stat()
	// The file can be closed
	f.Close()
	if err != nil {
		fmt.Println(err)
		return
	}
	// fi is a fileInfo interface returned by Stat
	fmt.Println(fi.Name(), fi.Size())
	fmt.Println(fi.ModTime(), fi.IsDir())
	fmt.Println(fi.Sys())

}

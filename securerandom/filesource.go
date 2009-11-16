package main

import "fmt"
import "os"
import "syscall"

func main() {
	fmt.Printf("Hello, bitches!\n")
}

type fileBasedSource struct {
	filename string;
	fd int;
	seed int;
}

func (src *fileBasedSource) Seed(seed int64) {
	// it's a noop.
}

func (rng *fileBasedSource) Int63() int64 {
	return 1;
}

func NewFileBasedSource(filename string) (src *fileBasedSource, err os.Error) {
	fileSource := new (fileBasedSource);
	fd, e := syscall.Open(filename, 0, 0);
	if e != 0 {
		err = os.Errno(e);
	}
	fileSource.filename = filename;
	fileSource.fd = fd;
	fileSource.seed = 0;
	return fileSource, err
}

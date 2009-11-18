package filesource
// Author rschonberger@gomail.com (Robert Schonberger) 
// Released under the Creative Commons - Attribution, Non commercial Usage OK. 
// http://creativecommons.org/licenses/by-nc-nd/3.0 2009


import "encoding/binary"
import "os"
import "rand"
import "syscall"

type fileBasedSource struct {
	filename string;
	fd int;
}

type FileBasedSource interface {
	rand.Source;
	Close();
}

func (src *fileBasedSource) Seed(seed int64) {
	// it's a noop.
}

// I wish there was a nicer way of doing this.
func convertToInt(randomSlice []byte) int64 {
     return int64(binary.BigEndian.Uint64(randomSlice) &^ ( 1<< 63))
}

func (rng *fileBasedSource) Int63() int64 {
	var randomData [8]byte;
	randomSlice := randomData[0 : 8];
	ret, e := syscall.Read(rng.fd, randomSlice);
	if e != 0 || ret < 8 {
		return 0
	} else {
		return convertToInt(randomSlice)
	}
	return 1
}

func (rng *fileBasedSource) Close() {
	if rng.fd > 0 {
		syscall.Close(rng.fd)
	}
}

func NewFileBasedSource(filename string) (src FileBasedSource, err os.Error) {
	fileSource := new (fileBasedSource);
	fd, e := syscall.Open(filename, 0, 0);
	if e != 0 {
		err = os.Errno(e)
	}
	fileSource.filename = filename;
	fileSource.fd = fd;
	return fileSource, err
}

func NewFileSeededSource(filename string) (src rand.Source, err os.Error) {
	fileSource, err := NewFileBasedSource(filename);
	defer fileSource.Close();
	if err != nil {
		return src, err
	}

	var seed int64 = fileSource.Int63();
	src = rand.NewSource(seed);
	return src, err
}

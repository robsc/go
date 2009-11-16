package filesource

import "fmt"
import "os"
import "rand"
import "syscall"

func main() {
	src, _ := NewFileSeededSource("/dev/random");
	rgn := rand.New(src);
	var sum int64 = 0;
	for i:= 0; i < 5000000; i++ {
		sum += int64(rgn.Int31());
	}
	fmt.Printf("Hello, bitches!%f \n", sum / 5000000);
}

type fileBasedSource struct {
	filename string;
	fd int;
	seed int64;
}

type FileBasedSource interface {
	rand.Source;
	Close();
}

func (src *fileBasedSource) Seed(seed int64) {
	// it's a noop.
	src.seed = seed;
}

// todo(robsc): make this pretty
func convertToInt(randomSlice [8]byte) int64 {
	randomSlice[0] = randomSlice[0] &^ 0x80;
	var retVal int64;
	retVal = (int64(randomSlice[7]) |
		int64(randomSlice[6]) << 8 |
		int64(randomSlice[5]) << 16 |
		int64(randomSlice[4]) << 24 |
		int64(randomSlice[3]) << 32 |
		int64(randomSlice[2]) << 40 |
		int64(randomSlice[1]) << 48 |
		int64(randomSlice[0]) << 56);
	return retVal;
	
}

func (rng *fileBasedSource) Int63() int64 {
	var randomData [8]byte;
	randomSlice := randomData[0 : 8];
	ret, e := syscall.Read(rng.fd, randomSlice);
	if e != 0 || ret < 8 {
		return 0
	} else {
		return convertToInt(randomData)
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
		err = os.Errno(e);
	}
	fileSource.filename = filename;
	fileSource.fd = fd;
	fileSource.seed = 0;
	return fileSource, err
}

func NewFileSeededSource(filename string) (src rand.Source, err os.Error) {
	fileSource, err := NewFileBasedSource(filename);
	defer fileSource.Close();
	if err != nil {
		return src, err
	}
	// NOTE(robsc): note the limited seed space that only goes through half.
	var seed int64 = fileSource.Int63();
	src = rand.NewSource(seed);
	return src, err
}

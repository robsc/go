package main
// Author rschonberger@gomail.com (Robert Schonberger) 
// Released under the Creative Commons - Attribution, Non commercial Usage OK. 
// http://creativecommons.org/licenses/by-nc-nd/3.0 2009

import "./filesource"
import "fmt"
import "rand"

func main() {
       src, _ := filesource.NewFileSeededSource("/dev/random");
       rgn := rand.New(src);
       var sum int64 = 0;
       for i:= 0; i < 5000000; i++ {
               sum += int64(rgn.Int31());
       }
       fmt.Printf("Hello, world. Random average sum is %f \n", sum / 5000000);
}


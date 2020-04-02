package main

// #cgo LDFLAGS: -L . -lccode
//#include "ccode.h"
import "C"

import "fmt"

func main() {

fmt.Println("Start in the world of go");

C.my_func(10)

fmt.Println("Finish in the world of go");
}

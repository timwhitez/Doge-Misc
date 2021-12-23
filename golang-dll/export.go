package main

import "C"
import (
	"fmt"
	"os/exec"
)

//export DCall
func DCall() {
	exec.Command("calc.exe").Start()
	fmt.Println("test")
}



func main() {
}
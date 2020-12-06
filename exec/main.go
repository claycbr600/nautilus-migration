package main

import (
	"fmt"
	"log"
	"os/exec"
	"strings"
)

func main() {
	out, err := exec.Command("ls").Output()
	if err != nil {
		log.Println(err)
	}

	s := strings.Trim(string(out), "\n")
	arr := strings.Split(s, "\n")
	fmt.Printf("%v len:%d cap:%d\n", arr, len(arr), cap(arr))
	fmt.Printf("arr[0] = %s\n", arr[0])
	fmt.Printf("The LS output is \n%s\n", out)
}

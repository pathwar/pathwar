package main

import (
	"log"
	"os"
	"os/exec"
	"syscall"
)

func main() {
	log.Println("Hello world from entrypoint, os.Args:", os.Args)
	binary, err := exec.LookPath(os.Args[1])
	if err != nil {
		panic(err)
	}
	args := os.Args[1:]
	env := os.Environ()
	if err := syscall.Exec(binary, args, env); err != nil {
		panic(err)
	}
}

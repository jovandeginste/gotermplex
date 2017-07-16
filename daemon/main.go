package main

import (
	"fmt"
	"github.com/kr/pty"
	"io"
	"log"
	"net"
	"os/exec"
)

var clients []net.Conn

func dataHandler(c net.Conn, input chan<- []byte) {
	fmt.Println("\n## Start of transmission.")
	clients = append(clients, c)
	var err error

	for err != io.EOF {
		buf := make([]byte, 1)
		_, r_err := c.Read(buf)

		if r_err != nil {
			fmt.Println("\n## Error during transmission.")
			return
		}

		input <- buf
	}
	fmt.Println("\n## End of transmission.")
}

func bash(input <-chan []byte, output chan<- []byte) {
	c := exec.Command("/bin/bash", "-l")
	f, err := pty.Start(c)
	if err != nil {
		panic(err)
	}

	go func() {
		for {
			v := <-input
			f.Write(v)
		}
	}()

	v := make([]byte, 1)
	for {
		n, err := f.Read(v)
		if n > 0 && err == nil {
			fmt.Print(string(v))
			output <- v
		}
	}
}

func outputMuxer(output <-chan []byte) {
	v := make([]byte, 1)

	for {
		v = <-output

		fmt.Print(string(v))
		for _, c := range clients {
			c.Write(v)
		}
	}
}

func main() {
	input := make(chan []byte, 0)
	output := make(chan []byte, 0)
	l, err := net.Listen("unix", "/tmp/example.sock")
	if err != nil {
		log.Fatal(err)
		return
	}

	go bash(input, output)
	go outputMuxer(output)

	fmt.Println("## Starting listening for connections.")
	for {
		fd, err := l.Accept()
		if err != nil {
			log.Fatal(err)
			return
		}
		go dataHandler(fd, input)
	}
}

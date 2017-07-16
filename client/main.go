package main

import (
	"fmt"
	"github.com/nsf/termbox-go"
	"io"
	"net"
)

func Output(c net.Conn) {
	var err error
	for err != io.EOF {
		buf := make([]byte, 1)
		_, r_err := c.Read(buf)
		fmt.Printf(string(buf))

		if r_err != nil {
			fmt.Println("\n## Error during transmission.")
			return
		}
	}
}

func main() {
	c, err := net.Dial("unix", "/tmp/example.sock")

	if err != nil {
		panic(err)
	}

	err = termbox.Init()
	if err != nil {
		panic(err)
	}
	defer termbox.Close()

	//termbox.SetInputMode(termbox.InputAlt)
	go Output(c)

	Input(c)
}

func Input(c net.Conn) {
	var current string
	d := make([]byte, 1)
	for {
		switch ev := termbox.PollRawEvent(d); ev.Type {
		case termbox.EventRaw:
			current = fmt.Sprintf("%q", d)
			c.Write(d)
			if current == `"q"` {
				return
			}
		case termbox.EventError:
			panic(ev.Err)
		}
	}
}

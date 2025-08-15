package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"strings"
)

const port = ":42069"

func getLinesChannel(f io.ReadCloser) <-chan string {
	out := make(chan string, 1)

	go func() {
		defer f.Close()
		defer close(out)

		buf := make([]byte, 8)
		currentLineContents := ""
		for {
			n, err := f.Read(buf)
			if err != nil {
				if currentLineContents != "" {
					out <- currentLineContents
				}
				if err == io.EOF {
					break
				}
				log.Fatalf("Error: %s", err)
				return
			}

			str := string(buf[:n])
			parts := strings.Split(str, "\n")
			for i := 0; i < len(parts)-1; i++ {
				out <- currentLineContents + parts[i]
				currentLineContents = ""
			}
			currentLineContents += parts[len(parts)-1]
		}
	}()
	return out
}

func main() {
	listener, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("error %s: %s\n", "error", err)
	}

	defer listener.Close()
	fmt.Println("Listening for TCP traffic on", port)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatalf("error %s: %s\n", "error", err)
		}

		fmt.Println("Accepted connection from", conn.RemoteAddr())

		lines := getLinesChannel(conn)
		for line := range lines {
			fmt.Println(line)
		}
		fmt.Println("Connection to ", conn.RemoteAddr(), "closed")
	}
}

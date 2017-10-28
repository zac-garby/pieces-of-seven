package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/Zac-Garby/pieces-of-seven/server/lib"
)

const DefaultPort = "12358"

func main() {
	fmt.Printf("server's port [%s]? :", DefaultPort)

	reader := bufio.NewReader(os.Stdin)
	port, err := reader.ReadString('\n')
	if err != nil {
		fmt.Printf("io error: %s\n", err.Error())
	}

	port = strings.TrimSpace(port)

	if len(port) == 0 {
		port = DefaultPort
	}

	port = ":" + port

	fmt.Println("listening on", port)

	server := lib.New(port)
	fmt.Println(server.Listen())
}

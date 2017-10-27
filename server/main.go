package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/Zac-Garby/pieces-of-seven/server/lib"
)

const DefaultAddress = "localhost:12358"

func main() {
	fmt.Printf("server's ip [%s]? ", DefaultAddress)

	reader := bufio.NewReader(os.Stdin)
	addr, err := reader.ReadString('\n')
	if err != nil {
		fmt.Printf("io error: %s\n", err.Error())
	}

	addr = strings.TrimSpace(addr)

	if len(addr) == 0 {
		addr = DefaultAddress
	}

	fmt.Println("listening on", addr)

	server := lib.New(addr)
	fmt.Println(server.Listen())
}

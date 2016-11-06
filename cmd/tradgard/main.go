package main

import (
	"fmt"

	"github.com/onlyafly/tradgard/pkg/server"
)

func main() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("PANIC RECOVERED IN MAIN", r)
		}
	}()

	server.Start()
}

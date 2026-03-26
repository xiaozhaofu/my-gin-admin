package main

import (
	"log"

	"go_sleep_admin/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		log.Fatalf("execute cobra command: %v", err)
	}
}

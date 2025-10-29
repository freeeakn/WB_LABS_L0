package main

import (
	"log"

	"wb_labs_l0_backend/internal/app"
)

func main() {
	if err := app.Run(); err != nil {
		log.Fatalf("app terminated: %v", err)
	}
}

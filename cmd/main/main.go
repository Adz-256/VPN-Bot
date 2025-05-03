package main

import (
	"context"
	"github.com/Adz-256/cheapVPN/internal/app"
	"log"
)

func main() {
	ctx := context.Background()

	a, err := app.New(ctx)
	if err != nil {
		log.Fatal(err)
	}

	log.Fatal(a.Run())
}

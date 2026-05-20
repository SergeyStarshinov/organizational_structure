package main

import (
	"fmt"
	"orgstructure/internal/config"
)

func main() {
	cfg := config.New()
	fmt.Println(cfg)

	// TODO: init logger

	// TODO: connect DB: postgreSQL

	// TODO: migrations: goose

	// TODO: init router: net/http

	// TODO: run server
}

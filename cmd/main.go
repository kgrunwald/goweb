package main

import (
	"log"
	"net/http"

	"github.com/joho/godotenv"
	"github.com/kgrunwald/goweb/framework"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	c := framework.BuildContainer()
	c.Invoke(func(handler http.Handler) {
		http.ListenAndServe(":80", handler)
	})
}

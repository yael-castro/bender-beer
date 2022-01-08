package main

import (
	"github.com/yael-castro/bender-beer/internal/dependency"
	"github.com/yael-castro/bender-beer/internal/handler"
	"log"
	"os"
)

const defaultPort = "8080"

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	log.SetFlags(log.Flags() | log.Lshortfile)

	h := handler.New()

	err := dependency.NewInjector(dependency.Default).Inject(h)
	if err != nil {
		log.Fatal(err)
	}

	h.Init()
	log.Fatal(h.Start(":" + port))
}

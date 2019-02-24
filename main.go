package main

import (
	"flag"
	"log"
	"path/filepath"

	"github.com/traed/screenshooter/api"
)

func main() {
	fp := flag.String("d", "./uploads", "Path to the directory where screenshots will be saved")
	addr := flag.String("a", ":8080", "Address for the server")

	path, err := filepath.Abs(*fp)
	if err != nil {
		log.Fatal("Invalid save path.")
		return
	}

	server := new(api.Server)
	server.SavePath = path
	server.ThrottleLimit = 15
	server.Addr = *addr
	server.Start()
}

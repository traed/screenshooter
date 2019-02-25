package main

import (
	"flag"
	"log"
	"path/filepath"
	"strconv"

	"github.com/traed/screenshooter/api"
	"github.com/traed/screenshooter/pkg/worker"
)

var (
	// MaxWorker - os.Getenv("MAX_WORKERS")
	MaxWorker = 3
	// MaxQueue - os.Getenv("MAX_QUEUE")
	MaxQueue = 20
	// MaxLength -
	MaxLength int64 = 2048
)

func main() {
	fp := flag.String("d", "./uploads", "Path to the directory where screenshots will be saved.")
	addr := flag.String("a", ":8080", "Address for the server.")
	thr := flag.String("t", "15", "Max number of simultaneous connections.")
	flag.Parse()

	path, err := filepath.Abs(*fp)
	if err != nil {
		log.Fatal("Invalid save path.")
		return
	}

	server := new(api.Server)
	server.SavePath = path
	server.Addr = *addr

	jobQueue := make(chan worker.Job, MaxQueue)
	dispatcher := worker.NewDispatcher(MaxWorker, jobQueue)

	server.Dispatcher = dispatcher

	if i, err := strconv.Atoi(*thr); err != nil {
		server.ThrottleLimit = i
	}

	server.Start()
}

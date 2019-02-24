package main

import "../../api"

func main() {
	server := new(api.Server)
	server.Start()
}

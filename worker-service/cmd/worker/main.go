package main

import "log"

func main() {
	log.Println("Worker service started")

	select {}
}

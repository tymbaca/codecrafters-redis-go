package main

import (
	"log"

	"github.com/codecrafters-io/redis-starter-go/cmd/redis"
)

func main() {
	if err := redis.Run(); err != nil {
		log.Fatal(err)
	}
}

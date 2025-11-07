package listconsumer

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

const listName = "demo-list"

func ListConsumerMain() {
	fmt.Println("list consumer application started")

	client := redis.NewClient(&redis.Options{Addr: "localhost:6379"})
	_, err := client.Ping(context.Background()).Result()
	if err != nil {
		log.Fatal("failed to connect", err)
	}

	for {
		data, err := client.BRPop(context.Background(), 2*time.Second, listName).Result()
		if err != nil {
			if errors.Is(err, redis.Nil) {
				continue
			}
			log.Println("brpop operation failed", err)
		}
		fmt.Println("received data from the list - ", data[1])
	}
}

package main

import (
	"context"
	"fmt"

	"github.com/DroidZed/go_lance/config"
	"github.com/DroidZed/go_lance/db"
)

func main() {

	port := config.EnvDbPORT()

	client := db.GetConnection()

	defer func() {
		if err := client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()

	fmt.Printf("Listening to port %s", port)
}

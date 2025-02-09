package main

import (
	"context"
	"fmt"
	"github.com/anpotashev/mpdgo/pkg/mpdapi"
)

func main() {
	ctx := context.Background()
	api, err := mpdapi.NewMpdApi(ctx, "192.168.0.110", 6600, "12345678")
	if err != nil {
		return
	}
	if err := api.Connect(); err != nil {
		fmt.Println("error")
		return
	}
	status, err := api.Status()
	if err != nil {
		return
	}
	fmt.Printf("status %s", status)
	fmt.Println("connected")
}

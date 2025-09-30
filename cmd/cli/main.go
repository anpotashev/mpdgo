package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/anpotashev/mpdgo/pkg/mpdapi"
	"time"
)

func main() {
	ctx := context.Background()
	api, err := mpdapi.NewMpdApi(ctx, "192.168.0.110", 6600, "12345678", false, 100, 3, time.Microsecond*200, time.Second*10)
	if err != nil {
		return
	}
	if err := api.Connect(); err != nil {
		fmt.Println("error")
		return
	}
	tree, err := api.Tree()
	if err != nil {
		return
	}
	fmt.Println(tree)
	status, err := api.Status()
	if err != nil {
		return
	}
	jsonStatus, _ := json.Marshal(status)
	fmt.Printf("status: %s\n", jsonStatus)
	fmt.Println("connected")
}

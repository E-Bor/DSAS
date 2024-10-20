package main

import (
	"fmt"
	"time"
)

func Report() error {
	fmt.Println("Instagram report start")
	time.Sleep(5 * time.Second)
	fmt.Println("Instagram report process")
	time.Sleep(5 * time.Second)
	fmt.Println("Instagram report end")
	return nil
}

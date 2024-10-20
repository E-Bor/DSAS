package main

import (
	"fmt"
	"time"
)

func Report() error {
	fmt.Println("Facebook report start")
	time.Sleep(5 * time.Second)
	fmt.Println("Facebook report process")
	time.Sleep(5 * time.Second)
	fmt.Println("Facebook report end")
	return nil
}

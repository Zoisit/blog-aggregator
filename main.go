package main

import (
	"fmt"

	"github.com/Zoisit/blog-aggregator/internal/config"
)

func main() {
	conf, err := config.Read()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(conf)

	config.SetUser(&conf, "anja")

	conf, err = config.Read()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(conf)
}

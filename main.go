package main

import (
	"fmt"
	"os"

	"github.com/Zoisit/blog-aggregator/internal/cli"
	"github.com/Zoisit/blog-aggregator/internal/config"
)

func main() {
	conf, err := config.Read()
	if err != nil {
		fmt.Println(err)
		return
	}
	//fmt.Println(conf)

	s := cli.State{}
	s.Config = &conf

	cmds := cli.Commands{
		HandlerFunction: make(map[string]func(*cli.State, cli.Command) error),
	}

	cmds.Register("login", cli.HandlerLogin)

	args := os.Args

	if len(args) < 2 {
		fmt.Println("Not enough arguments provided")
		os.Exit(1)
	}

	cmd := cli.Command{
		Name:      args[1],
		Arguments: args[2:],
	}

	err = cmds.Run(&s, cmd)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	//conf.SetUser("anja")

	conf, err = config.Read()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(conf)
}

package main

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/Zoisit/blog-aggregator/internal/cli"
	"github.com/Zoisit/blog-aggregator/internal/config"
	"github.com/Zoisit/blog-aggregator/internal/database"
	_ "github.com/lib/pq"
)

func main() {
	conf, err := config.Read()
	if err != nil {
		fmt.Println(err)
		return
	}

	db, err := sql.Open("postgres", conf.DB_URL)
	if err != nil {
		fmt.Println(err)
		return
	}
	dbQueries := database.New(db)

	s := cli.State{}
	s.Config = &conf
	s.DB = dbQueries

	cmds := cli.Commands{
		HandlerFunction: make(map[string]func(*cli.State, cli.Command) error),
	}

	cmds.Register("login", cli.HandlerLogin)
	cmds.Register("register", cli.HandlerRegister)
	cmds.Register("reset", cli.HandlerReset)

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

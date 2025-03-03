package cli

import "fmt"

type Commands struct {
	HandlerFunction map[string]func(*State, Command) error
}

func (c *Commands) Register(name string, f func(*State, Command) error) {
	c.HandlerFunction[name] = f
}

func (c *Commands) Run(s *State, cmd Command) error {
	f, ok := c.HandlerFunction[cmd.Name]
	if !ok {
		return fmt.Errorf("command does not exist")
	}

	return f(s, cmd)
}

package commands

import "fmt"

func commandPokedex(c *Config, _ []string) error {
	if len(c.Pokedex) == 0 {
		return fmt.Errorf("your pokedex is empty")
	}

	for _, pokemon := range c.Pokedex {
		fmt.Fprintf(c.Out, " -%s", pokemon.Name)
	}

	return nil
}

package commands

import (
	"fmt"
	"io"
	"strings"
)

func commandInspect(c *Config, args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("usage: catch <pokemon_name>")
	}

	pokemonName := strings.ToLower(args[0])
	if pokemon, ok := c.Pokedex[pokemonName]; !ok {
		return fmt.Errorf("you have not caught that pokemon")
	} else {
		printPokemonData(c.Out, pokemon)
	}
	return nil
}

func printPokemonData(out io.Writer, p Pokemon) {
	fmt.Printf("Name: %s\n", p.Name)
	fmt.Printf("Height: %d\n", p.Height)
	fmt.Printf("Weight: %d\n", p.Weight)
	fmt.Println("Stats:")
	for _, s := range p.Stats {
		fmt.Fprintf(out, "  -%s: %d\n", s.Stat.Name, s.BaseStat)
	}
	fmt.Println("Types:")
	for _, t := range p.Types {
		fmt.Fprintf(out, "  - %s\n", t.Type.Name)
	}
}

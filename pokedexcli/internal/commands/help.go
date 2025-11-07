package commands

import "fmt"

func commandHelp(*Config, []string) error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Printf("Usage:\n\n")
	for cmd, val := range Commands {
		fmt.Printf("%s - %s\n", cmd, val.description)
	}
	return nil
}

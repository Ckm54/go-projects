package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/ckm54/go-projects/pokedexcli/internal/commands"
	"github.com/ckm54/go-projects/pokedexcli/internal/pokecache"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	cache := pokecache.NewCache(10 * time.Second)

	config := &commands.Config{
		Cache:   cache,
		Pokedex: make(map[string]commands.Pokemon),
		Out:     os.Stdout,
	}

	for {
		fmt.Print("Pokedex > ")
		if scanner.Scan() {
			input := scanner.Text()
			if input == "" {
				continue
			}

			parts := CleanInput(input)
			cmdName := parts[0]
			args := parts[1:]

			handler, exists := commands.Commands[cmdName]
			if !exists {
				fmt.Println("Unknown command")
				continue
			}

			err := handler.Callback(config, args)
			if err != nil {
				if errors.Is(err, commands.ErrExit) {
					break
				}
				fmt.Printf("%v\n", err)
			}

		}

		if err := scanner.Err(); err != nil {
			fmt.Printf("Error reading user input: %s\n", err)
		}
	}
}

func CleanInput(text string) []string {
	trimmedText := strings.Join(strings.Fields(strings.ToLower(text)), " ")
	return strings.Split(trimmedText, " ")
}

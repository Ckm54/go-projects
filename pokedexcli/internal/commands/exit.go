package commands

import (
	"errors"
	"fmt"
)

var ErrExit = errors.New("exit requested")

func commandExit(*Config, []string) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	return ErrExit
}

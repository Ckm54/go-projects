package commands

import (
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"strings"

	"github.com/ckm54/go-projects/pokedexcli/constants"
)

func commandCatch(c *Config, args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("usage: catch <pokemon_name>")
	}

	pokemonName := strings.ToLower(args[0])
	if _, ok := c.Pokedex[pokemonName]; ok {
		fmt.Printf("%s already in your pokedex!\n", pokemonName)
		return nil
	}

	fmt.Printf("Throwing a Pokeball at %s...\n", pokemonName)

	var pokemon Pokemon
	url := fmt.Sprintf("%s/pokemon/%s", constants.BASEURL, pokemonName)

	if data, found := c.Cache.Get(pokemonName); found {
		fmt.Println("Cache hit✅, using cached pokemon info")
		if err := json.Unmarshal(data, &pokemon); err != nil {
			return fmt.Errorf("failed to decode cached data: %v", err)
		}
	} else {
		fmt.Println("Cache miss❌, fetching data from api")
		res, err := http.Get(url)
		if err != nil {
			return err
		}
		defer res.Body.Close()

		if res.StatusCode != http.StatusOK {
			return fmt.Errorf("failed to fetch pokemon: %s", res.Status)
		}

		body, err := io.ReadAll(res.Body)
		if err != nil {
			return fmt.Errorf("failed to read response body: %v", err)
		}

		if err := json.Unmarshal(body, &pokemon); err != nil {
			return err
		}

		c.Cache.Add(pokemonName, body)
	}

	catchChance := 100 - (pokemon.BaseExperience / 10) // higher experience level has a lower chance to catch
	if catchChance < 10 {
		catchChance = 10
	}

	r := rand.Intn(100)
	if r < catchChance {
		fmt.Printf("%s was caught!\n", pokemon.Name)
		c.Pokedex[pokemon.Name] = pokemon
	} else {
		fmt.Printf("%s escaped!\n", pokemon.Name)
	}

	return nil
}

package commands

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/ckm54/go-projects/pokedexcli/constants"
)

func commandExplore(c *Config, args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("usage: explore <area_name>")
	}

	areaName := args[0]
	fmt.Printf("Exploring area: %s...\n", areaName)

	url := fmt.Sprintf("%s/location-area/%s", constants.BASEURL, areaName)

	var locationDetailResponse LocationAreaDetails

	if data, found := c.Cache.Get(url); found {
		fmt.Println("Cache hit✅, using cached response")
		if err := json.Unmarshal(data, &locationDetailResponse); err != nil {
			return fmt.Errorf("failed to decode cached response: %v", err)
		}
	} else {
		fmt.Println("Cache miss❌, fetching data from api")
		res, err := http.Get(url)
		if err != nil {
			return err
		}
		defer res.Body.Close()

		if res.StatusCode != http.StatusOK {
			return fmt.Errorf("failed to fetch: %s", res.Status)
		}

		body, err := io.ReadAll(res.Body)
		if err != nil {
			return fmt.Errorf("failed to read response body: %v", err)
		}

		if err = json.Unmarshal(body, &locationDetailResponse); err != nil {
			return err
		}

		c.Cache.Add(url, body)
	}

	for _, pokemon := range locationDetailResponse.PokemonEncounters {
		fmt.Println(pokemon.Pokemon.Name)
	}
	return nil
}

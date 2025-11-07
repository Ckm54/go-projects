package commands

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	constants "github.com/ckm54/go-projects/pokedexcli/constants"
)

type Action string

const (
	ActionNext     Action = "next"
	ActionPrevious Action = "previous"
)

func configUrl(c *Config, action Action) (string, error) {
	var locationAreasURL string

	switch action {
	case ActionNext:
		if len(c.Next) > 0 {
			locationAreasURL = c.Next
		} else if len(c.Next) == 0 && len(c.Previous) > 0 {
			fmt.Println("------------------------------")
			return "", fmt.Errorf("you have reached the end")
		} else {
			locationAreasURL = constants.BASEURL + "/location-area"
		}
	case ActionPrevious:
		if len(c.Previous) > 0 {
			locationAreasURL = c.Previous
		} else {
			return "", fmt.Errorf("you are on the first page")
		}
	default:
		return "", fmt.Errorf("no action specified")
	}

	return locationAreasURL, nil
}

func fetchData(c *Config, url string) error {
	var locationAreasResponse LocationAreaResponse

	if data, found := c.Cache.Get(url); found {
		fmt.Println("Cache hit✅, using cached response")
		if err := json.Unmarshal(data, &locationAreasResponse); err != nil {
			return fmt.Errorf("failed to decode cached response: %v", err)
		}
	} else {
		fmt.Println("Cache miss❌, fetching data from api")
		res, err := http.Get(url)
		if err != nil {
			return err
		}
		defer res.Body.Close()

		body, err := io.ReadAll(res.Body)
		if err != nil {
			return fmt.Errorf("failed to read body: %v", err)
		}

		if err = json.Unmarshal(body, &locationAreasResponse); err != nil {
			return err
		}

		c.Cache.Add(url, body)
	}

	c.Next = locationAreasResponse.Next
	c.Previous = locationAreasResponse.Previous
	for _, area := range locationAreasResponse.Results {
		fmt.Println(area.Name)
	}

	return nil
}

func commandMap(c *Config, _ []string) error {
	url, err := configUrl(c, ActionNext)
	if err != nil {
		return err
	}

	err = fetchData(c, url)

	return err

}

func commandMapB(c *Config, _ []string) error {
	url, err := configUrl(c, ActionPrevious)
	if err != nil {
		return err
	}

	err = fetchData(c, url)

	return err
}

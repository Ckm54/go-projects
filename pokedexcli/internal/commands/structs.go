package commands

import (
	"io"

	"github.com/ckm54/go-projects/pokedexcli/internal/pokecache"
)

type CliCommand struct {
	name        string
	description string
	Callback    func(*Config, []string) error
}

type Config struct {
	Next     string
	Previous string
	Cache    *pokecache.Cache
	Pokedex  map[string]Pokemon
	Out      io.Writer
}

type LocationAreaResponse struct {
	Count    int    `json:"count"`
	Next     string `json:"next"`
	Previous string `json:"previous"`
	Results  []LocationArea
}

type LocationArea struct {
	Name string `json:"name"`
	Url  string `json:"url"`
}

type PokemonContainer struct {
	PokemonEncounters []PokemonEncounter `json:"pokemon_encounters"`
}

type PokemonEncounter struct {
	Pokemon PokemonDetails `json:"pokemon"`
}

type PokemonDetails struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

type LocationAreaDetails struct {
	Id                int                `json:"id"`
	Name              string             `json:"name"`
	PokemonEncounters []PokemonEncounter `json:"pokemon_encounters"`
}

type Pokemon struct {
	Id             int           `json:"id"`
	Name           string        `json:"name"`
	BaseExperience int           `json:"base_experience"`
	Height         int           `json:"height"`
	Weight         int           `json:"weight"`
	Types          []PokemonType `json:"types"`
	Stats          []PokemonStat `json:"stats"`
}

type PokemonType struct {
	Slot int             `json:"slot"`
	Type PokemonTypeInfo `json:"type"`
}

type PokemonStat struct {
	BaseStat int             `json:"base_stat"`
	Effort   int             `json:"effort"`
	Stat     PokemonTypeInfo `json:"stat"`
}

type PokemonTypeInfo struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

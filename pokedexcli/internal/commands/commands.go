package commands

var Commands = map[string]CliCommand{}

func init() {
	// assign Callbacks in init to avoid initialization cycle between
	// package-level variables and functions that reference them
	Commands["help"] = CliCommand{
		name:        "help",
		description: "Displays a help message",
		Callback:    commandHelp,
	}
	Commands["exit"] = CliCommand{
		name:        "exit",
		description: "Exit the pokedox",
		Callback:    commandExit,
	}
	Commands["map"] = CliCommand{
		name:        "map",
		description: "Displays names of 20 location areas in the pokemon world",
		Callback:    commandMap,
	}
	Commands["mapb"] = CliCommand{
		name:        "mapb",
		description: "Displays names of previous 20 location areas in the pokemon world if available",
		Callback:    commandMapB,
	}
	Commands["explore"] = CliCommand{
		name:        "explore <area_name>",
		description: "Explore a specific location area, e.g. explore pallet-town",
		Callback:    commandExplore,
	}
	Commands["catch"] = CliCommand{
		name:        "catch <pokemon_name>",
		description: "Try and catch a specific pokemon, e.g. catch pikachu",
		Callback:    commandCatch,
	}
	Commands["inspect"] = CliCommand{
		name:        "inspect <pokemon_name>",
		description: "View the details of an already caught pokemon e.g. inspect pikachu",
		Callback:    commandInspect,
	}
	Commands["pokedex"] = CliCommand{
		name:        "pokedex",
		description: "View the names of all caught pokemons!",
		Callback:    commandPokedex,
	}
}

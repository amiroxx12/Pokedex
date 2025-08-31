package main

import (
	"math/rand"
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

func cleanInput(text string) []string {
	lower := strings.ToLower(text)
	words := strings.Fields(lower)

	var cleaned []string
	for _, word := range words {
		cleanedWord := strings.Trim(word, ".,!?\"';:-()[]{}")
		if cleanedWord != "" {
			cleaned = append(cleaned, cleanedWord)
		}
	}
	return cleaned
}

func commandExit() string {
	fmt.Print("Closing the Pokedex... Goodbye!\n")
	os.Exit(0)
	return ""
}

func commandHelp() string {
	helpText := "Welcome to the Pokedex!\n Usage:\n"
	for _, cmd := range commands {
		helpText += fmt.Sprintf(" - %s: %s\n", cmd.name, cmd.description)
	}
	return helpText
}

func commandMap() string {
	for _, loc := range globalLocations.Results {
		fmt.Println(loc.Name)
	}

	// Pre-fetch next page for the next call
	if globalLocations.Next != "" {
		locations, err := fetchLocations(globalLocations.Next)
		if err == nil {
			globalLocations = locations
		}
	}

	return ""
}

func commandbMap() string {
	if globalLocations.Previous == "" {
		return "No Previous Location."
	}
	locations, err := fetchLocations(globalLocations.Previous)
	if err != nil {
		return fmt.Sprintf("Error fetching previous locations: %v", err)
	}
	globalLocations = locations
	for _, loc := range globalLocations.Results {
		fmt.Println(loc.Name)
	}
	return ""
}

func commandExplore(args []string) string {
	if len(args) < 1 {
		return "Please provide a location to explore."
	}
	location := strings.Join(args, "-")
	url := fmt.Sprintf("https://pokeapi.co/api/v2/location-area/%s", location)
	exploreResp, err := fetchPokemonAtLocation(url)
	if err != nil {
		return fmt.Sprintf("Error exploring location %s: %v", location, err)
	}
	if len(exploreResp.PokemonEncounters) == 0 {
		return fmt.Sprintf("No Pokémon found at location %s.", location)
	}
	result := "Found Pokemon:\n"
	for _, encounter := range exploreResp.PokemonEncounters {
		result += fmt.Sprintf("%s\n", encounter.Pokemon.Name)
	}
	return result

}
func commandInspect(args []string) string {
	if len(args) < 1 {
		return "Please provide a Pokémon name to inspect."
	}
	poke := strings.ToLower(args[0]) // args = words[1:] therefore args[0] is the pokemon name (words[1])
	pokemon, exists := pokedex[poke]
	if !exists {
		return fmt.Sprintf("You don't have %s in your Pokedex. Try catching it first!", poke)
	}
	return fmt.Sprintf("Name: %s\nID: %d\nHeight: %d\nWeight: %d\nBase Experience: %d\n",
		pokemon.Name, pokemon.ID, pokemon.Height, pokemon.Weight, pokemon.BaseExperience)
}
func attemptCatch(pokemon Pokemon) bool {
	maxChance := 100
	factor := 1
	catchChance := maxChance - (pokemon.BaseExperience / factor)
	if catchChance < 5 {
		catchChance = 5 // Minimum 5% chance to catch any Pokémon
	}
	roll := rand.Intn(maxChance) + 1
	return roll <= catchChance
}
func addToPokedex(pokemon Pokemon) {
	name := strings.ToLower(pokemon.Name)
	if _, exists := pokedex[name]; !exists {
		pokedex[name] = pokemon

	} else {
		return
	}
}
func commandCatch(args []string) string {
	if len(args) < 1 {
		return "Please provide a Pokémon name to catch."
	}
	poke := strings.ToLower(args[0]) // args = words[1:] therefore args[0] is the pokemon name (words[1])
	fmt.Printf("Throwing a Pokeball at %s...\n", poke)
	pokemon,err := fetchPokemonDetails(fmt.Sprintf("https://pokeapi.co/api/v2/pokemon/%s", poke))
	if err != nil {
		return fmt.Sprintf("Error fetching details for Pokémon %s: %v",poke, err)
	}
	if attemptCatch(pokemon) {
		addToPokedex(pokemon)
		return fmt.Sprintf("%s was caught!", pokemon.Name)
	} else {
		return fmt.Sprintf("%s escaped!", pokemon.Name)
	}
}
func commandPokedex(args[] string) string {
	if len(pokedex) == 0 {
		return "Your Pokedex is empty. Catch some Pokémon!"
	}
	result := "Your Pokedex:\n"
	for name := range pokedex {
		result += fmt.Sprintf("- %s\n", name)
	}
	return result
}



func fetchLocations(url string) (LocationList, error) {
	var locations LocationList
	res, err := http.Get(url)
	if err != nil {
		return locations, err
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return locations, err
	}
	if err := json.Unmarshal(body, &locations); err != nil {
		return locations, err
	}
	return locations, nil
}

func fetchPokemonAtLocation(url string) (ExploreResponse, error) {
	var exploreResp ExploreResponse
	res, err := http.Get(url)
	if err != nil {
		return exploreResp, err
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return exploreResp, err
	}
	if err := json.Unmarshal(body, &exploreResp); err != nil {
		return exploreResp, err
	}
	return exploreResp, nil
}

func fetchPokemonDetails(url string) (Pokemon, error) {
	var pokemon Pokemon
	res, err := http.Get(url)
	if err != nil {
		return pokemon, err
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return pokemon, err
	}
	if err := json.Unmarshal(body, &pokemon); err != nil {
		return pokemon, err
	}
	return pokemon, nil
}

type Pokemon struct {
	Name string `json:"name"`
	ID  int    `json:"id"`
	Height int `json:"height"`
	Weight int `json:"weight"`
	BaseExperience int `json:"base_experience"`
	Types []struct {
		Type struct {
			Name string `json:"name"`
		} `json:"type"`
	} `json:"types"`
}
type ExploreResponse struct {
	PokemonEncounters []struct {
		Pokemon struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"pokemon"`
	} `json:"pokemon_encounters"`
}
type CliCommand struct {
	name        string
	description string
	callback    func(args []string) string
}

type Location struct {
	Name string `json:"name"`
	URL string `json:"url"`
}

type LocationList struct {
	Results  []Location `json:"results"`
	Next     string     `json:"next"`
	Previous string     `json:"previous"`
}


var pokedex = make(map[string]Pokemon)
// globalLocations holds the fetched locations for use in commands
var globalLocations LocationList

// commands map will be initialized in main
var commands map[string]CliCommand



func main() {
	limit := 20
	offset := 0
	url := fmt.Sprintf("https://pokeapi.co/api/v2/location-area?limit=%d&offset=%d", limit, offset)
	locations, err := fetchLocations(url)
	if err != nil {
		fmt.Printf("Error fetching locations: %v\n", err)
		return
	}
	// Assign locations to a package-level variable for access in commandMap
	globalLocations = locations

	commands = map[string]CliCommand{
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			callback:    func(args []string) string {return commandExit()},
		},
		"help": {
			name:        "help",
			description: "Displays a help message",
			callback:    func(args []string) string { return commandHelp() },
		},
		"map": {
			name:        "map",
			description: "List next locations",
			callback:    func(args []string) string { return commandMap()},
		},
		"bmap": {
			name:        "bmap",
			description: "List previous locations",
			callback:    func(args []string) string { return commandbMap() },
		},
		"explore": {
			name:        "explore",
			description: "explore a location",
			callback:    commandExplore,
		},
		"catch": {
			name:        "catch",
			description: "Catch a Pokemon",
			callback:    commandCatch,
		},
		"inspect": {
			name:        "inspect",
			description: "Inspect a caught Pokemon",
			callback:    commandInspect,
		},
		"pokedex": {
			name:        "pokedex",
			description: "List all caught Pokemon",
			callback:    commandPokedex,
		},
	}

	input := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("Pokedex >")
		input.Scan()
		line := input.Text()
		words := cleanInput(line)
		if len(words) == 0 {
			continue
		}
		cmdName := words[0]
		args := words[1:]
		if cmd, exists := commands[cmdName]; exists {
			result := cmd.callback(args)
			if result != "" {
				fmt.Println(result)
			}
		} else {
			fmt.Printf("Unknown command: %s\n", cmdName)
		}
	}
}

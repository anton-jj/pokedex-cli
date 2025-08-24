package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"
	"math/rand"
	"github.com/anton-jj/pokedex-cli/internal/pokeapi"
	"github.com/anton-jj/pokedex-cli/internal/pokecache"
)
type cliCommand struct {
	Name string
	Description string
	Callback func(string) error
}

type Config struct {
	Next *string
	Previous *string
	DefaultURL string
}

	
func main () {
	cache := pokecache.NewCache(5 * time.Second)
	client := pokeapi.NewClient(cache)
	pokedex := make(map[string]pokeapi.Pokemon)
  	var config = Config{ 
		Next: nil,
		Previous: nil ,
		DefaultURL: "https://pokeapi.co/api/v2/location-area/",
	}
	var arg string
	scanner := bufio.NewScanner(os.Stdin)
	cliMap := map[string]cliCommand {
		"exit" : {
			Name: "Exit",
			Description: "Exit the Pokedex",
			Callback: func(arg string) error {
				return commandExit()
			} ,
		},
		"map" : {
			Name : "Map", 
			Description: "Displays next 20 locations",
			Callback: func(arg string) error {
				return commandMap(&config, *client)
			},
		},
		"mapb" : {
			Name: "Mapb",
			Description: "Displays previous 20 locations",
			Callback: func(arg string) error {
				return commandMapb(&config, *client)
			} ,
		},

		"explore" : {
		Name: "explore",
		Description: "Shows pokemons in a location",
			Callback: func(arg string) error {
				return commandExplore(arg, &config, *client)
			},
		},
		"catch" : {
			Name: "catch",
			Description: "Allows you to try catch a pokemon",
			Callback: func(arg string) error {
				return commandCatch(arg, *client, &pokedex)
			},
		},
		"inspect" : {
			Name: "inspect",
			Description: "inspects a pokemon in you pokedex",
			Callback: func(arg string) error {
				return commandInspect(arg, &pokedex)
			},
		},
	}
	cliMap["help"] = cliCommand{
			Name: "Help",
			Description: "Displays a help message",
			Callback: func(arg string) error {
				return commandHelp(cliMap)
			},
		}

	for {
		fmt.Printf("Pokedex > ")
		scanner.Scan()
		input := cleanInput(scanner.Text())
		if len(input) > 1 {
			arg = input[1] 
		}
		fmt.Println(len(input))
		if cmd, ok := cliMap[input[0]]; ok {
			cmd.Callback(arg)
		} else {
			fmt.Println("Unknown command")
		}

	}
}

func cleanInput(text string) []string {
	return strings.Fields(strings.ToLower(text))
}


func commandInspect(pokemon string, pokedex *map[string]pokeapi.Pokemon) error {
	if mon, exists := (*pokedex)[pokemon]; !exists {
		fmt.Println("You haven't catched that pokemon yet")
	}else {
		for _, stat := range mon.Stats {
			fmt.Println(stat.BaseStat)
		}
	}
	return nil
}
func commandCatch(pokemon string,  client pokeapi.Client,pokedex *map[string]pokeapi.Pokemon) error {
	URL := "https://pokeapi.co/api/v2/pokemon/" + pokemon
	res, err := client.GetPokemonInfo(URL)
	if err != nil {
		return err
	}
	maxExp := 400
	baseRoll := 100

	fmt.Println("Throwing a Pokeball at "+ res.Name+"...")
	dif := (res.BaseExperience * baseRoll ) / maxExp
	if rand.Intn(baseRoll) > dif {
		fmt.Println(res.Name + " cought!")
		if _, exists := (*pokedex)[res.Name]; exists {
			return nil	
		}
		(*pokedex)[res.Name] = *res
	} else {
		fmt.Println(res.Name + " escaped!")

	}

	return nil
}
func commandExplore(location string, conf *Config, client pokeapi.Client)  error{
		URL := conf.DefaultURL + location
		res, err := client.GetLocationsInfo(URL)
		if err != nil {
		return err
		}
		for _, name := range res.PokemonEncounters{
		fmt.Println(name.Pokemons.Name)
	}
		
		return nil
}
func commandExit() error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandMap(conf *Config, client pokeapi.Client) error {
	if conf.Next == nil {
		firstPage := conf.DefaultURL + "?offset=0&limit=20"
		conf.Next = &firstPage
	}
	locations,err := client.GetLocationsInfo(*conf.Next)

	if err != nil {
		return  err
	}
	
	for _, loc := range locations.Results {
		fmt.Println(loc.Name)
	}

	conf.Next = locations.Next
	conf.Previous = locations.Previous

	return nil
}

func commandMapb(conf *Config, client pokeapi.Client ) error {
	if conf.Previous == nil	 {
		fmt.Printf("You're on the first page")
		return nil
	}
	locations, err := client.GetLocationsInfo(*conf.Previous)

	if err != nil {
		return err
	}
	for _, loc := range locations.Results {
		fmt.Println(loc.Name)
	}
	conf.Next = locations.Next
	conf.Previous = locations.Previous
	return nil
		
}
func commandHelp(commands map[string]cliCommand) error {
	fmt.Printf("Welcome to the Pokedex!\nUsage: \n\n" )
	for name, cmd := range commands {
		fmt.Printf("%s: %s\n", name, cmd.Description)
	}	
	return nil
}

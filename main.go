package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
)
type cliCommand struct {
	Name string
	Description string
	Callback func() error
}

type Config struct {
	Next *string
	Previous *string
}

type LocationArea struct {
    Name string `json:"name"`
    URL  string `json:"url"`
}

type LocationAreaResponse struct {
		Count int `json:"count"`
		Next *string `json:"next"`
		Previous *string `json:"previous"`
		Results  []LocationArea `json:"results"`
}
	
func main () {

	var config Config
	scanner := bufio.NewScanner(os.Stdin)
	cliMap := map[string]cliCommand {
		"exit" : {
			Name: "Exit",
			Description: "Exit the Pokedex",
			Callback: commandExit,
		},
		"map" : {
			Name : "Map", 
			Description: "Displays next 20 locations",
			Callback: func() error {
				return commandMap(&config)
			},
		},
		"mapb" : {
			Name: "Mapb",
			Description: "Displays previous 20 locations",
			Callback: func() error {
				return commandMapb(&config)
			} ,
		},
	}
	cliMap["help"] = cliCommand{
			Name: "Help",
			Description: "Displays a help message",
			Callback: func() error {
				return commandHelp(cliMap)
			},
		}

	for {
		fmt.Printf("Pokedex > ")
		scanner.Scan()
		input := cleanInput(scanner.Text())
		if cmd, ok := cliMap[input[0]]; ok {
			cmd.Callback()
		} else {
			fmt.Println("Unknown command")
		}

	}
}

func cleanInput(text string) []string {
	return strings.Fields(strings.ToLower(text))
}

func commandExit() error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}


func getLocations(url string) ( *LocationAreaResponse, error ) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	var res LocationAreaResponse
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return nil, err
	}

	return &res, nil
}

func commandMap(conf *Config) error {
	if conf.Next == nil {
		defaultURL := "https://pokeapi.co/api/v2/location-area?offset=0&limit=20"
		conf.Next = &defaultURL
	}
	locations,err := getLocations(*conf.Next)

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

func commandMapb(conf *Config) error {
	if conf.Previous == nil	 {
		fmt.Printf("You're on the first page")
		return nil
	}
	locations, err := getLocations(*conf.Previous)

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

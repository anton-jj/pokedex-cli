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
			Calllback: commandMap,
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


func getData(url string) (*ApiResponse, errro ) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	var res LocationAreaResponse
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return nil, err
	}

	return res, nil
}
func getLocations() error {
	url := "https://pokeapi.co/api/v2/location-area/?offset=20&limit=20"
	
}


func commandHelp(commands map[string]cliCommand) error {
	fmt.Printf("Welcome to the Pokedex!\nUsage: \n\n" )
	for name, cmd := range commands {
		fmt.Printf("%s: %s\n", name, cmd.Description)
	}	
	return nil
}

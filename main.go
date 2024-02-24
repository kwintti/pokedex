package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
    "errors"
)

func main() {
    readingInput()

}


func readingInput() {
    commands := make(map[string]cliCommand)
    commands = map[string]cliCommand{
        "help": {
            name:        "help",
            description: "Displays a help message",
            callback:     func() error { return commandHelp(commands) },
        },
        "exit": {
            name:        "exit",
            description: "Exit the Pokedex",
            callback:    func() error {os.Exit(0); return nil}, 
        },
        "map": {
            name:        "map",
            description: "Show next 20 location areas",
            callback:    getPokemons,
        },
        "mapb": {
            name:        "mapb",
            description: "Show previous 20 location areas",
            callback:    getPokemonsBack,
        },
    }

    fmt.Println("pokedex>")
    scanner := bufio.NewScanner(os.Stdin)

    for scanner.Scan() {
        err := scanner.Err()
        if err != nil {
            fmt.Println("error reading input", err)
        }
        cmd := commands[scanner.Text()]
        if err := cmd.callback(); err != nil {
            log.Println(err)
        }
    }
   

}

func commandHelp(commands map[string]cliCommand) error {
    fmt.Println("This is pokedex app")
    fmt.Println("Usage:")
    fmt.Println("")
    
    for _, cmd := range commands {
        fmt.Println(cmd.name, ": ", cmd.description)
    }
    return nil 
}

var url_ cliCommand

func getPokemons()  error {
    urlp := &url_
    if len(urlp.url) == 0 {
        url_.url = "https://pokeapi.co/api/v2/location-area/"
    }

    pokemons, err := apiCall(url_.url)
    if err != nil {
        return err
    }
    for i, _ := range pokemons.Results {
        fmt.Println(pokemons.Results[i].Name)
    }
    urlp.url = *pokemons.Next
    if pokemons.Previous != nil {
        urlp.urlb = *pokemons.Previous
    }
    return err
}

func getPokemonsBack() error {
    urlp := &url_
    if len(urlp.urlb) == 0 {
        err := errors.New("Already on the first page")
        return err
    }

    pokemons, err := apiCall(url_.urlb)
    if err != nil {
        return err
    }
    for i, _ := range pokemons.Results {
        fmt.Println(pokemons.Results[i].Name)
    }
    if pokemons.Previous == nil {
        urlp.urlb = "" 
    }
    if pokemons.Previous != nil {
        urlp.urlb = *pokemons.Previous
    }
    return err
} 

func apiCall(url string) (pokemonAPI, error) {
    res, err := http.Get(url)
    if err != nil {
        log.Fatal(err)
    }
    body, err := io.ReadAll(res.Body)
    res.Body.Close()
    if res.StatusCode > 299 {
        log.Fatalf("Response failed with status code: %d and\nbody: %s\n", res.StatusCode, body)
    }
    if err != nil {
        log.Fatal(err)
    }
    var pokemons pokemonAPI
    json.Unmarshal([]byte(body), &pokemons)

    return pokemons, nil
}


type pokemonAPI struct {
	Count    int    `json:"count"`
	Next     *string `json:"next"`
	Previous *string    `json:"previous"`
    Results []results `json:"results"`
}

type results struct {
    Name string `json:"name"`
    URL string  `json:"url"`
}

type cliCommand struct {
	name        string
	description string
	callback    func() error
    url         string
    urlb        string
}

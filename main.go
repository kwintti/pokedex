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
    "github.com/kwintti/pokecache"
    "strings"
)

func main() {
    readingInput()

}


func readingInput() {
    params := "canalave-city-area"
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
        "explore": {
            name:        "exolore",
            description: "Get pokemons in the area. Usage: explore <area-name>",
            callback:    func() error {return explorePokemons(params)},
        },
    }

    fmt.Println("pokedex>")
    scanner := bufio.NewScanner(os.Stdin)

    for scanner.Scan() {
        input := scanner.Text()
        err := scanner.Err()
        if err != nil {
            fmt.Println("error reading input", err)
        }
        splited := strings.Split(input, " ")
        if len(splited) > 1 {
            params = splited[1]
        }
        cmd, ok := commands[splited[0]] 
        if !ok {
            log.Print("Command does not exist")
        }else
         if err := cmd.callback(); err != nil {
            log.Print(err)
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
    pokecache.NewCache(5)
    return nil 
}

var url_ cliCommand
var pokemons pokemonAPI
var cache = pokecache.NewCache(60)


func getPokemons()  error {
    urlp := &url_
    if len(urlp.url) == 0 {
        url_.url = "https://pokeapi.co/api/v2/location-area/?offset=0&limit=20"
    }
    var body []byte
    var err error
    val, found := cache.Get(urlp.url)
    body = val
    if !found {
        fmt.Println("No cache found, getting pokemons from API")
        body, err = apiCall(urlp.url)
        if err != nil {
            return err
        }
        cache.Add(urlp.url, body)
    }
    json.Unmarshal([]byte(body), &pokemons)
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
    var body []byte
    var err error
    val, found := cache.Get(urlp.urlb)
    body = val
    if !found {
        fmt.Println("No cache found, getting pokemons from API")
        body, err = apiCall(urlp.urlb)
        if err != nil {
            return err
        }
        cache.Add(urlp.urlb, body)
    }

    json.Unmarshal([]byte(body), &pokemons)
    for i, _ := range pokemons.Results {
        fmt.Println(pokemons.Results[i].Name)
    }
    if pokemons.Previous == nil {
        urlp.urlb = "" 
    }
    if pokemons.Previous != nil {
        urlp.urlb = *pokemons.Previous
    }
    if pokemons.Next != nil {
        urlp.url = *pokemons.Next
    }
    return err
} 

func apiCall(url string) ([]byte, error) {
    res, err := http.Get(url)
    if err != nil {
        log.Fatal(err)
    }
    body, err := io.ReadAll(res.Body)
    res.Body.Close()
    if res.StatusCode > 299 {
        err := errors.New("Not found")
        return body, err
    }
    if err != nil {
        log.Fatal(err)
    }

    return body, nil
}

var exploringPokemons ExplorePokemons 

func explorePokemons(area string) error {
    var body []byte
    var err error
    val, found := cache.Get(area)
    body = val
    if !found {
        fmt.Println("No cache found, getting pokemons from API")
        body, err = apiCall("https://pokeapi.co/api/v2/location-area/" + area)
        if err != nil {
            return err
        }
        cache.Add(area, body)
    }
    json.Unmarshal([]byte(body), &exploringPokemons)
    for _, val := range exploringPokemons.PokemonEncounters {
        fmt.Println(val.Pokemon.Name)
    }

    
    return nil
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

type ExplorePokemons struct {
	PokemonEncounters    []PokemonEncounters    `json:"pokemon_encounters"`
}
type PokemonEncounters struct {
	Pokemon        Pokemon          `json:"pokemon"`
}
type Pokemon struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

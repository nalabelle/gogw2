package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
)

const API_BASE = "https://api.guildwars2.com/v2/"
const CHARACTER_ENDPOINT = API_BASE + "characters"
const COLORS_ENDPOINT = API_BASE + "colors"

var API_KEY string

type API struct {
	Characters map[string]APICharacter
	Colors     map[int]APIColor
}

func NewAPI(apiKey string) *API {
	api := &API{}
	API_KEY = apiKey
	api.Characters = make(map[string]APICharacter)
	api.Colors = make(map[int]APIColor)
	return api
}

type APICharacter struct {
	Name      string                  `json:"name,omitempty"`
	Equipment []APICharacterEquipment `json:"equipment"`
}

type APICharacterEquipment struct {
	ID   int    `json:"id"`
	Slot string `json:"slot"`
	Dyes []int  `json:"dyes"`
}

type APIColor struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	BaseRGB []int  `json:"base_rgb"`
}

func Get(uri string) io.ReadCloser {
	client := http.Client{}
	req, err := http.NewRequest("GET", uri, nil)
	if err != nil {
		log.Fatalln(err)
	}

	req.Header.Set("Authorization", "Bearer "+API_KEY)
	req.Header.Set("Accept", "application/json")
	res, err := client.Do(req)
	if err != nil {
		log.Fatalln(err)
	}

	debug, err := httputil.DumpResponse(res, true)
	if err != nil {
		log.Fatalln(err)
	} else {
		//log.Println(string(debug))
	}

	if res.StatusCode != http.StatusOK {
		log.Println(string(debug))
		log.Fatalln("HTTP Response:", res.StatusCode)
	}

	return res.Body
}

func (api API) ResolveColor(color int) APIColor {
	if color, ok := api.Colors[color]; ok {
		return color
	}

	return api.NewColor(color)
}

func (api API) NewColor(color int) APIColor {
	c := APIColor{}
	uri := fmt.Sprintf("%s/%d", COLORS_ENDPOINT, color)
	body := Get(uri)
	decoder := json.NewDecoder(body)
	err := decoder.Decode(&c)
	if err != nil {
		log.Fatalln(err)
	}

	api.Colors[color] = c
	return c
}

func (api API) Character(characterName string) APICharacter {
	//uri := fmt.Sprintf("%s/%s", CHARACTER_ENDPOINT, api.characterName)

	char := APICharacter{Name: characterName}
	api.Characters[characterName] = char
	return char
}

func (c APICharacter) ResolveEquipment() []APICharacterEquipment {
	if len(c.Equipment) > 0 {
		return c.Equipment
	}

	char := APICharacter{}
	uri := fmt.Sprintf("%s/%s/equipment", CHARACTER_ENDPOINT, c.Name)
	body := Get(uri)
	decoder := json.NewDecoder(body)
	err := decoder.Decode(&char)
	if err != nil {
		log.Fatalln(err)
	}

	c.Equipment = char.Equipment
	return c.Equipment
}

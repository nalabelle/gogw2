package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Command int

const (
	Dyes Command = iota
)

type Context struct {
	CharacterName string
	APIKey        string
	command       Command
}

func (ctx Context) validateArgs() error {
	if ctx.APIKey == "" {
		return errors.New("No API key provided")
	}
	if ctx.CharacterName == "" {
		return errors.New("No Character provided")
	}
	return nil
}

func main() {
	var ctx Context

	flag.StringVar(&ctx.APIKey, "key", "", "Account API Key")
	flag.StringVar(&ctx.CharacterName, "character", "", "Character Name")
	flag.Parse()

	err := ctx.validateArgs()
	if err != nil {
		fmt.Println("Error: %s", err)
		os.Exit(1)
	}

	api := NewAPI(ctx.APIKey)
	ctx.command = Dyes

	if ctx.command == Dyes {
		character := api.Character(ctx.CharacterName)
		equipment := character.ResolveEquipment()
		for _, ce := range equipment {
			if len(ce.Dyes) == 0 {
				continue
			}
			var colors []string
			for _, item := range ce.Dyes {
				var color APIColor
				if item == 0 {
					color = APIColor{ID: 0, Name: "Dye Remover"}
				} else {
					color = api.ResolveColor(item)
				}

				var rgbStrings []string
				for _, val := range color.BaseRGB {
					rgbStrings = append(rgbStrings, strconv.Itoa(val))
				}
				if len(rgbStrings) > 0 {
					colors = append(colors, fmt.Sprintf("%s (%s)", color.Name, strings.Join(rgbStrings, ",")))
				} else {
					colors = append(colors, color.Name)
				}
			}
			fmt.Printf("%s: %s\n", ce.Slot, strings.Join(colors, "; "))
		}
	}
}

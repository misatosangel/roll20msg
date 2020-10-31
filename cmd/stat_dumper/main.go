package main

import (
	"fmt"
	"encoding/base64"
	"encoding/json"
	"github.com/jessevdk/go-flags"
	"log"
	"os"
	"strings"

	"github.com/misatosangel/roll20msg/pkg/roll20msg"
	"github.com/misatosangel/roll20msg/internal/stats"
)

// Variables used for command line parameters
var settings struct {
	File   	   string `short:"f" long:"file" required:"true" description:"Chat data file to read" value-name:"<path>"`
	gaveVersion bool
	Version     func() `long:"version" required:"false" description:"Print tool version and exit."`
	Debug       bool   `short:"d" long:"debug" description:"Print all rolls as found to STDERR."`
}

var buildVersion = "dev"
var buildDate = "dev"
var buildCommit = "dev"

func init() {
}

func main() {
	os.Exit(run())
}


func run() int {
	CliParse()
	if settings.gaveVersion {
		return 0
	}
	var msgs = make( roll20msg.MsgStream, 0, 1000 )
	var err error
	if strings.HasSuffix( settings.File, ".json" ) || strings.HasSuffix( settings.File, ".js" ) {
		err = readJsonFile( settings.File, &msgs )
	} else{
		err = readB64JsonFile( settings.File, &msgs )
	}
	if err != nil {
		log.Fatalln( err )
	}

	nameCheck := make( map[string]string )
	gotByTime := make( map[string][]stats.DatedResult )

	for idx, msgBlock := range msgs {
		for k, msg := range msgBlock {
			if ! msg.HasRollResults() {
				continue
			}
			who := msg.Who
			plyId := msg.PlayerId
			if _, exists := nameCheck[plyId] ; ! exists {
				nameCheck[plyId] = who
			} else {
				who = nameCheck[plyId]
			}
			ts := msg.TimeStamp()
			passfunc := func( roll roll20msg.Roll ) bool {
				if roll.Dice == 0 {
					return true
				}
				if roll.Sides != 20 {
					return true
				}
				for _, diceResult := range roll.Results {
					gotByTime[who] = append(gotByTime[who], stats.DatedResult{Date: ts, Result: diceResult.Value } )
					if settings.Debug {
						fmt.Fprintf( os.Stderr, "%d: [%s] %s Rolled d20 and got: %d\n", idx, k, msg.BriefDesc(), diceResult.Value )
					}
				}
				return true
			}
			msg.IterateRawDiceRolls(passfunc)
		}
	}
	first := true
	for who, vals := range gotByTime {
		if first {
			first = false
		} else {
			fmt.Print("\n")
		}
		sb := stats.NewStatBlock(vals)
		fmt.Printf( "**__%s__**\n", who )
		fmt.Print( sb.FormatResultsDiscord() )
	}
	return 0

}

// stream through base64 if wanting that
func readB64JsonFile( path string, output interface{} ) error {
	input, err := os.Open(path)
	if err != nil {
		return err
	}
	defer input.Close()
	b64Decoder := base64.NewDecoder(base64.StdEncoding, input)
	jsonDecoder := json.NewDecoder(b64Decoder)
	return jsonDecoder.Decode(output)
}

// stream load JSON
func readJsonFile( path string, output interface{} ) error {
	input, err := os.Open(path)
	if err != nil {
		return err
	}
	defer input.Close()
	jsonDecoder := json.NewDecoder(input)
	return jsonDecoder.Decode(output)
}

// deal with our very limited command line options
func CliParse() {
	parser := flags.NewParser(&settings, flags.Default)
	settings.Version = func() {
		parser.SubcommandsOptional = true
		fmt.Printf("Roll20 d20 stat collector %s\nBuilt: %s\nCommit: %s\n", buildVersion, buildDate, buildCommit)
		settings.gaveVersion = true
	}
	_, err := parser.Parse()

	if err != nil {
		switch err.(type) {
		case *flags.Error:
			if err.(*flags.Error).Type == flags.ErrHelp {
				os.Exit(0)
			}
		}
		log.Fatalln(err)
	}
	if settings.gaveVersion {
		return;
	}
}

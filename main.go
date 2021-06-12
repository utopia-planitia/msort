package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/urfave/cli/v2"
	msort "github.com/utopia-planitia/msort/pkg"
)

func main() {
	err := run(os.Args)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
}

func run(args []string) error {
	app := &cli.App{
		Name:  "msort",
		Usage: "sort yaml manifests",
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:  "sort-keys",
				Usage: "sort keys within each yaml document",
				Value: false,
			},
			&cli.BoolFlag{
				Name:  "drop-tests",
				Usage: "remove yaml documents with \"test\" in its name",
				Value: false,
			},
			&cli.BoolFlag{
				Name:    "in-place",
				Aliases: []string{"i"},
				Usage:   "update files in place",
				Value:   false,
			},
		},
		Action: sortYamlFiles,
	}

	err := app.Run(args)
	if err != nil {
		return err
	}

	return nil
}

func sortYamlFiles(c *cli.Context) error {
	stdinPipe, err := detectStdinPipe()
	if err != nil {
		return fmt.Errorf("detect stdin usage: %v", err)
	}

	sortKeys := c.Bool("sort-keys")
	dropTests := c.Bool("drop-tests")
	inPlace := c.Bool("in-place")

	firstDocument := true

	if stdinPipe {
		yml, err := ioutil.ReadAll(os.Stdin)
		if err != nil {
			return fmt.Errorf("read stdin: %v", err)
		}

		output, err := msort.Sort(yml, sortKeys, dropTests)
		if err != nil {
			return fmt.Errorf("sort stdin: %v", err)
		}

		firstDocument = false

		_, err = fmt.Print(output)
		if err != nil {
			return fmt.Errorf("write to stdout: %v", err)
		}
	}

	for _, path := range c.Args().Slice() {
		yml, err := ioutil.ReadFile(path)
		if err != nil {
			return fmt.Errorf("read file %s: %v", path, err)
		}

		output, err := msort.Sort(yml, sortKeys, dropTests)
		if err != nil {
			return fmt.Errorf("sort stdin: %v", err)
		}

		if inPlace {
			err = ioutil.WriteFile(path, []byte(output), 0666)
			if err != nil {
				return fmt.Errorf("update %s in place: %v", path, err)
			}
			continue
		}

		if !firstDocument {
			_, err := fmt.Print("---\n")
			if err != nil {
				log.Fatalf("write document seperator to buffer: %v", err)
			}
		}

		firstDocument = false

		_, err = fmt.Print(output)
		if err != nil {
			return fmt.Errorf("write to stdout: %v", err)
		}
	}

	return nil
}

func detectStdinPipe() (bool, error) {
	stat, err := os.Stdin.Stat()
	if err != nil {
		return false, err
	}

	pipedToStdin := (stat.Mode() & os.ModeCharDevice) == 0

	return pipedToStdin, nil
}

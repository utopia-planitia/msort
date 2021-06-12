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
		Name:   "msort",
		Usage:  "sort yaml manifests",
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

	if stdinPipe {
		yml, err := ioutil.ReadAll(os.Stdin)
		if err != nil {
			return fmt.Errorf("read stdin: %v", err)
		}

		manifests := msort.NewManifests(string(yml))

		if os.Getenv("DISABLE_KEY_SORTING") == "" {
			manifests.SortByKeys()
		}

		if os.Getenv("KEEP_TESTS") == "" {
			manifests.DropTest()
		}

		manifests.OrderDocuments()

		_, err = fmt.Print(manifests.String())
		if err != nil {
			return fmt.Errorf("write to stdout: %v", err)
		}
	}

	for _, path := range c.Args().Slice() {
		yml, err := ioutil.ReadFile(path)
		if err != nil {
			return fmt.Errorf("read file %s: %v", path, err)
		}

		manifests := msort.NewManifests(string(yml))

		if os.Getenv("DISABLE_KEY_SORTING") == "" {
			manifests.SortByKeys()
		}

		if os.Getenv("KEEP_TESTS") == "" {
			manifests.DropTest()
		}

		manifests.OrderDocuments()

		_, err = fmt.Print(manifests.String())
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

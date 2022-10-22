package argparser

import (
	"fmt"
	"strings"
	"testing"
)

type RunOptions struct {
	CacheLoc string `arg:"cache-location|c"`
}

type BuildOptions struct {
	OutputLoc string `arg:"output-location|o"`
}

func TestDummyRun(t *testing.T) {
	err := dummyParseArg(strings.Split("dummy run --cache-location /path/to/cache target.dummy", " "))
	if err != nil {
		panic(err.Error())
	}
}

func TestDummyRunWrong(t *testing.T) {
	err := dummyParseArg(strings.Split("dummy run --output-location /path/to/output target.dummy", " "))
	if err == nil {
		panic("Should have failed but did not")
	}
}

func TestDummyBuild(t *testing.T) {
	err := dummyParseArg(strings.Split("dummy build --output-location /path/to/output target.dummy", " "))
	if err != nil {
		panic(err.Error())
	}
}

func TestDummyBuildWrong(t *testing.T) {
	err := dummyParseArg(strings.Split("dummy build --cache-location /path/to/cache target.dummy", " "))
	if err == nil {
		panic("Should have failed but did not")
	}
}

func dummyParseArg(args []string) error {
	p := New(args[1:])
	subcommand := p.NextArg()
	switch subcommand {
	case "run":
		options := RunOptions{}
		err := p.NextOptions(&options)
		if err != nil {
			return err
		}
		fmt.Println("Running sub-command 'run' with cache location", options.CacheLoc)
	case "build":
		options := BuildOptions{}
		err := p.NextOptions(&options)
		if err != nil {
			return err
		}
		fmt.Println("Running sub-command 'run' with output location", options.OutputLoc)
	default:
		return fmt.Errorf("Unrecognised sub-command")
	}
	return nil
}

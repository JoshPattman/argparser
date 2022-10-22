package argparser

import (
	"fmt"
	"strings"
	"testing"
)

// dummy command (does nothing but serves as an example)

// dummy [-s] [-f] <sub-program>
//	-s --silent bool -> silent mode
//  -f --fast bool -> fast mode

// <sub-program> = run | build
// dummy run [-c] [-r] <target>
//  -c --cache-loc string -> location to store cache at
//  -r --ram uint -> gb of ram to allow
// dummy build [-o] <target>
//  -o --output-loc string -> location to output to

// examples
// dummy run file.txt
// dummy -s build --output-loc /path/to/output file2.txt
// dummy -sf run -r 5 file3.txt
// dummy -f -s build -o /another/path file4.txt
func runDummyProgram(args []string) (string, error) {
	// Options to dummy program
	dummyOptions := struct {
		Silent bool `arg:"s|silent"`
		Fast   bool `arg:"f|fast"`
	}{}
	// Ignore the first ang as it always will be "dummy"
	p := New(args[1:])
	if err := p.NextOptions(&dummyOptions); err != nil {
		return "", err
	}
	subProgram := p.NextArg()
	switch subProgram {
	case "run":
		runOptions := struct {
			CacheLoc string `arg:"c|cache-loc"`
			RamGB    uint   `arg:"ram|r"`
		}{CacheLoc: "/default/path", RamGB: 1}
		if err := p.NextOptions(&runOptions); err != nil {
			return "", err
		}
		target := p.NextArg()
		if target == "" {
			return "", fmt.Errorf("must specify target")
		}
		return getDummyRunOutput(dummyOptions.Silent, dummyOptions.Fast, runOptions.CacheLoc, runOptions.RamGB), nil
	case "build":
		buildOptions := struct {
			OutputLoc string `arg:"o|output-loc"`
		}{OutputLoc: "/default/path"}
		if err := p.NextOptions(&buildOptions); err != nil {
			return "", err
		}
		target := p.NextArg()
		if target == "" {
			return "", fmt.Errorf("must specify target")
		}
		return getDummyBuildOutput(dummyOptions.Silent, dummyOptions.Fast, buildOptions.OutputLoc), nil
	default:
		return "", fmt.Errorf("sub program unrecognised")
	}
}

func getDummyRunOutput(silent bool, fast bool, cacheLoc string, ramGB uint) string {
	return fmt.Sprintf("%v:%v:RUN:%v:%v", silent, fast, cacheLoc, ramGB)
}
func getDummyBuildOutput(silent bool, fast bool, outputLoc string) string {
	return fmt.Sprintf("%v:%v:RUN:%v", silent, fast, outputLoc)
}

func TestDummyCorrect(t *testing.T) {
	if err := runAndCheckDummy("dummy -s run file.txt", getDummyRunOutput(true, false, "/default/path", 1)); err != nil {
		t.Error(err)
	}
	if err := runAndCheckDummy("dummy -sf run file.txt", getDummyRunOutput(true, true, "/default/path", 1)); err != nil {
		t.Error(err)
	}
	if err := runAndCheckDummy("dummy -f --silent run file.txt", getDummyRunOutput(true, true, "/default/path", 1)); err != nil {
		t.Error(err)
	}
	if err := runAndCheckDummy("dummy --fast run file.txt", getDummyRunOutput(false, true, "/default/path", 1)); err != nil {
		t.Error(err)
	}
	if err := runAndCheckDummy("dummy -s run -c new-loc -r 10 file.txt", getDummyRunOutput(true, false, "new-loc", 10)); err != nil {
		t.Error(err)
	}
	if err := runAndCheckDummy("dummy -s build -o new-loc file.txt", getDummyBuildOutput(true, false, "new-loc")); err != nil {
		t.Error(err)
	}
}

func TestDummyIncorrect(t *testing.T) {
	if err := runAndCheckDummy("dummy -s arg build -o new-loc file.txt", getDummyBuildOutput(true, false, "new-loc")); err == nil {
		t.Errorf("should not have passed")
	}
	if err := runAndCheckDummy("dummy -s prog -o new-loc file.txt", getDummyBuildOutput(true, false, "new-loc")); err == nil {
		t.Errorf("should not have passed")
	}
	if err := runAndCheckDummy("dummy -s run -o new-loc -r 10a file.txt", getDummyRunOutput(true, false, "new-loc", 10)); err == nil {
		t.Errorf("should not have passed")
	}
	if err := runAndCheckDummy("dummy -s run -o new-loc -r -10 file.txt", getDummyRunOutput(true, false, "new-loc", 10)); err == nil {
		t.Errorf("should not have passed")
	}
	if err := runAndCheckDummy("dummy run -s -o new-loc -r 10a file.txt", getDummyRunOutput(false, false, "new-loc", 10)); err == nil {
		t.Errorf("should not have passed")
	}
}

func runAndCheckDummy(commandline, expectedOutput string) error {
	s, err := runDummyProgram(strings.Split(commandline, " "))
	if err != nil {
		return err
	}
	if s != expectedOutput {
		return fmt.Errorf("Output was not as expected. Out: %s, Exp: %s", s, expectedOutput)
	}
	return nil
}

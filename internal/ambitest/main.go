package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
	"time"
)

const Name = "ambitest"

func main() {
	os.Exit(Run(os.Args[1:], os.Stdout, os.Stderr))
}

func Run(args []string, stdout io.Writer, stderr io.Writer) int {
	var all bool
	f := flag.NewFlagSet("ambitest", flag.ContinueOnError)
	f.SetOutput(stderr)
	f.BoolVar(&all, "all", false, "run all tests")

	errFlag := f.Parse(args)
	if errFlag == flag.ErrHelp {
		return 0
	}
	if errFlag != nil {
		fmt.Fprintf(stderr, "%s: %s\n", Name, errFlag.Error())
		return 2
	}

	var cfg Config
	cfg.Ch = os.Getenv("AMBITEST_CH")
	cfg.UserKey = os.Getenv("AMBITEST_USERKEY")
	cfg.ReadKey = os.Getenv("AMBITEST_READKEY")
	cfg.WriteKey = os.Getenv("AMBITEST_WRITEKEY")
	cfg.DevKey = os.Getenv("AMBITEST_DEVKEY")

	t := &T{Config: &cfg}
	fail := false
	for _, c := range TestList {
		if !all && c.Extended {
			fmt.Fprintf(stdout, "SKIP %s\n", c.Name)
			continue
		}

		t.Failed = false
		t.Output = nil
		c.Func(t)

		if t.Failed {
			fmt.Fprintf(stdout, "FAIL %s\n", c.Name)
		} else {
			fmt.Fprintf(stdout, "PASS %s\n", c.Name)
		}

		for _, s := range t.Output {
			fmt.Fprintf(stdout, "     %s\n", strings.Replace(s, "\n", "\n     ", -1))
		}

		fail = fail || t.Failed
	}

	if fail {
		fmt.Fprintf(stdout, "FAIL\n")
		return 1
	}
	fmt.Fprintf(stdout, "PASS\n")
	return 0
}

type T struct {
	Config   *Config
	Failed   bool
	Output   []string
	LastPost time.Time
}

func (t *T) Fail() {
	t.Failed = true
}

func (t *T) Log(args ...any) {
	s := fmt.Sprint(args...)
	t.Output = append(t.Output, s)
}

func (t *T) Logf(format string, args ...any) {
	s := fmt.Sprintf(format, args...)
	t.Output = append(t.Output, s)
}

func (t *T) Error(args ...any) {
	t.Fail()
	t.Log(args...)
}

func (t *T) Errorf(format string, args ...any) {
	t.Fail()
	t.Logf(format, args...)
}

func (t *T) PostWait() {
	time.Sleep(time.Until(t.LastPost.Add(5 * time.Second)))
	t.LastPost = time.Now()
}

func (t *T) PostDone() {
	t.LastPost = time.Now()
}

type Config struct {
	Ch       string
	UserKey  string
	ReadKey  string
	WriteKey string
	DevKey   string
}

type TestEntry struct {
	Name     string
	Func     TestFunc
	Extended bool
}

type TestFunc func(*T)

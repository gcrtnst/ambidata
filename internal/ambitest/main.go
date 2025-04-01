package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
)

const Name = "ambitest"

var TestList = []TestEntry{}

func main() {
	os.Exit(Run(os.Args[1:], os.Stdout, os.Stderr))
}

func Run(args []string, stdout io.Writer, stderr io.Writer) int {
	var cfgpath string
	f := flag.NewFlagSet(Name, flag.ContinueOnError)
	f.SetOutput(stderr)
	f.StringVar(&cfgpath, "c", "", "config file")

	errFlag := f.Parse(args)
	if errFlag == flag.ErrHelp {
		return 0
	}
	if errFlag != nil {
		fmt.Fprintf(stderr, "%s: %s\n", Name, errFlag.Error())
		return 2
	}

	var cfg Config
	errCfg := LoadConfig(&cfg, cfgpath)
	if errCfg != nil {
		fmt.Fprintf(stderr, "%s: %s\n", Name, errCfg.Error())
		return 2
	}

	fail := false
	for _, f := range TestList {
		t := NewT(&cfg)
		f.Func(t)

		if t.Failed {
			fmt.Fprintf(stdout, "FAIL %s\n", f.Name)
		} else {
			fmt.Fprintf(stdout, "PASS %s\n", f.Name)
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
	Config *Config
	Failed bool
	Output []string
}

func NewT(cfg *Config) *T {
	return &T{
		Config: cfg,
		Failed: false,
		Output: nil,
	}
}

func (t *T) Fail() {
	t.Failed = true
}

func (t *T) Logf(format string, args ...any) {
	s := fmt.Sprintf(format, args...)
	t.Output = append(t.Output, s)
}

func (t *T) Errorf(format string, args ...any) {
	t.Fail()
	t.Logf(format, args...)
}

type Config struct {
	Ch       string `json:"-"`
	UserKey  string `json:"-"`
	ReadKey  string `json:"-"`
	WriteKey string `json:"-"`
}

func LoadConfig(cfg *Config, path string) error {
	var err error

	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close() // ignore error

	dec := json.NewDecoder(f)
	dec.DisallowUnknownFields()
	err = dec.Decode(cfg)
	if err != nil {
		return err
	}

	cfg.Ch = os.Getenv("AMBITEST_CH")
	cfg.UserKey = os.Getenv("AMBITEST_USERKEY")
	cfg.ReadKey = os.Getenv("AMBITEST_READKEY")
	cfg.WriteKey = os.Getenv("AMBITEST_WRITEKEY")
	return nil
}

type TestEntry struct {
	Name string
	Func TestFunc
}

type TestFunc func(*T)

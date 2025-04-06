package main

import (
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
	var cfg Config
	cfg.Ch = os.Getenv("AMBITEST_CH")
	cfg.UserKey = os.Getenv("AMBITEST_USERKEY")
	cfg.ReadKey = os.Getenv("AMBITEST_READKEY")
	cfg.WriteKey = os.Getenv("AMBITEST_WRITEKEY")
	cfg.DevKey = os.Getenv("AMBITEST_DEVKEY")

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
	Config   *Config
	Failed   bool
	Output   []string
	LastPost time.Time
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
	Name string
	Func TestFunc
}

type TestFunc func(*T)

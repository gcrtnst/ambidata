package ambidata

import (
	"context"
	"net/url"
)

type Fetcher struct {
	ReadKey string
	Ch      string
	Config  *Config
}

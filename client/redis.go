package client

import (
	"bytes"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"text/template"

	"github.com/go-redis/redis/v7"
	"github.com/gobuffalo/envy"
	"github.com/gobuffalo/x/defaults"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

// Connections contains all available connections
var Connections = map[string]*redis.Options{}

// ErrConfigFileNotFound is returned when the pop config file can't be found,
// after looking for it.
var ErrConfigFileNotFound = errors.New("unable to find redis config file")

var lookupPaths = []string{"", "./config", "/config", "../", "../config", "../..", "../../config"}

// ConfigName is the name of the YAML databases config file
var ConfigName = "redis.yml"

//Connect connect redis client
func Connect(e string) (*redis.Client, error) {
	if len(Connections) == 0 {
		err := LoadConfigFile()
		if err != nil {
			return nil, err
		}
	}
	e = defaults.String(e, "development")
	o := Connections[e]
	if o == nil {
		return nil, errors.Errorf("could not find connection named %s", e)
	}
	return redis.NewClient(o), nil
}

// LoadConfigFile loads a POP config file from the configured lookup paths
func LoadConfigFile() error {
	path, err := findConfigPath()
	if err != nil {
		return err
	}
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()
	return LoadFrom(f)
}

func findConfigPath() (string, error) {
	for _, p := range LookupPaths() {
		path, _ := filepath.Abs(filepath.Join(p, ConfigName))
		if _, err := os.Stat(path); err == nil {
			return path, err
		}
	}
	return "", ErrConfigFileNotFound
}

// LookupPaths returns the current configuration lookup paths
func LookupPaths() []string {
	return lookupPaths
}

// LoadFrom reads a configuration from the reader and sets up the connections
func LoadFrom(r io.Reader) error {
	envy.Load()
	deets, err := ParseConfig(r)
	if err != nil {
		return err
	}
	for n, d := range deets {
		Connections[n] = d
	}
	return nil
}

// ParseConfig reads the pop config from the given io.Reader and returns
// the parsed ConnectionDetails map.
func ParseConfig(r io.Reader) (map[string]*redis.Options, error) {
	tmpl := template.New("test")
	tmpl.Funcs(map[string]interface{}{
		"envOr": func(s1, s2 string) string {
			return envy.Get(s1, s2)
		},
		"env": func(s1 string) string {
			return envy.Get(s1, "")
		},
	})
	b, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}
	t, err := tmpl.Parse(string(b))
	if err != nil {
		return nil, errors.Wrap(err, "couldn't parse config template")
	}

	var bb bytes.Buffer
	err = t.Execute(&bb, nil)
	if err != nil {
		return nil, errors.Wrap(err, "couldn't execute config template")
	}

	deets := map[string]*redis.Options{}
	err = yaml.Unmarshal(bb.Bytes(), &deets)
	return deets, errors.Wrap(err, "couldn't unmarshal config to yaml")
}

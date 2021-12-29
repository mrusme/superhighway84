package config

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/BurntSushi/toml"
)

type ConfigProfile struct {
  From                string
  Organization        string
}

type Config struct {
  ConfigFile          string `toml:"-"`

  ConnectionString    string
  CachePath           string
  Logfile             string

  Profile             ConfigProfile
}

func LoadConfig() (*Config, error) {
  configDir, exist := os.LookupEnv("XDG_CONFIG_HOME")
  if exist == false {
    configDir, exist = os.LookupEnv("HOME")
    if exist == false {
      return nil, errors.New("No XDG_CONFIG_HOME or HOME set!")
    }
    configDir = filepath.Join(configDir, ".config")
  }
  os.MkdirAll(configDir, 0755)

  configFile := filepath.Join(configDir, "superhighway84.toml")

  f, err := os.OpenFile(configFile, os.O_CREATE|os.O_RDWR, 0644)
  if err != nil {
    return nil, err
  }
  defer f.Close()

  configFileContent, err := ioutil.ReadAll(f)
  if err != nil {
    return nil, err
  }

  cfg := new(Config)
  _, err = toml.Decode(string(configFileContent), &cfg)

  cfg.ConfigFile = configFile
  return cfg, nil
}

func (cfg *Config) Persist() (error) {
  buf := new(bytes.Buffer)
  if err := toml.NewEncoder(buf).Encode(cfg); err != nil {
    return err
  }

  if err := ioutil.WriteFile(cfg.ConfigFile, buf.Bytes(), 0644); err != nil {
    return err
  }

  return nil
}

func (cfg *Config) WasSetup() (bool) {
  if cfg.CachePath == "" ||
     cfg.ConnectionString == "" ||
     cfg.Logfile == "" {
    return false
  }

  return true
}

func (cfg *Config) Setup() (error) {
  fmt.Printf("\nSUPERHIGHWAY84\n\nInitial Setup\n-------------\n\n")

  defaultConnectionString := "/orbitdb/bafyreifdpagppa7ve45odxuvudz5snbzcybwyfer777huckl4li4zbc5k4/superhighway84"
  fmt.Printf("Database connection string [%s]: ", defaultConnectionString)
  fmt.Scanln(&cfg.ConnectionString)
  if strings.TrimSpace(cfg.ConnectionString) == "" {
    cfg.ConnectionString = defaultConnectionString
  }

  cacheDir, exist := os.LookupEnv("XDG_CACHE_HOME")
  if exist == false {
    cacheDir = filepath.Join(os.Getenv("HOME"), ".cache")
  }

  defaultCachePath := filepath.Join(cacheDir, "superhighway84")
  fmt.Printf("Database cache path [%s]: ", defaultCachePath)
  fmt.Scanln(&cfg.CachePath)
  if strings.TrimSpace(cfg.CachePath) == "" {
    cfg.CachePath = defaultCachePath
  }
  os.MkdirAll(filepath.Dir(cfg.CachePath), 0755)

  defaultLogfile := filepath.Join(cacheDir, "superhighway84.log")
  fmt.Printf("Logfile path [%s]: ", defaultLogfile)
  fmt.Scanln(&cfg.Logfile)
  if strings.TrimSpace(cfg.Logfile) == "" {
    cfg.Logfile = defaultLogfile
  }


  fmt.Printf("\nProfile information\n-------------------\n\n")

  defaultProfileFrom := fmt.Sprintf("%s@localhost", os.Getenv("USER"))
  fmt.Printf("From [%s]: ", defaultProfileFrom)
  fmt.Scanln(&cfg.Profile.From)
  if strings.TrimSpace(cfg.Profile.From) == "" {
    cfg.Profile.From = defaultProfileFrom
  }

  fmt.Printf("Organization []: ")
  fmt.Scanln(&cfg.Profile.Organization)

  return cfg.Persist()
}


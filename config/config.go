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

  CachePath           string // Deprecated, should be removed soon
	DatabaseCachePath   string
  ProgramCachePath    string

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
  if cfg.DatabaseCachePath == "" ||
		 cfg.ProgramCachePath == "" ||
     cfg.ConnectionString == "" ||
     cfg.Logfile == "" ||
		 cfg.Profile.From == "" {
    return false
  }

  return true
}

func (cfg *Config) Setup() (error) {
  fmt.Printf("\nSUPERHIGHWAY84\n\nInitial Setup\n-------------\n\n")

  defaultConnectionString := "/orbitdb/bafyreifdpagppa7ve45odxuvudz5snbzcybwyfer777huckl4li4zbc5k4/superhighway84"
	if cfg.ConnectionString != "" {
		defaultConnectionString = cfg.ConnectionString
	}
  fmt.Printf("Database connection string [%s]: ", defaultConnectionString)
  fmt.Scanln(&cfg.ConnectionString)
  if strings.TrimSpace(cfg.ConnectionString) == "" {
    cfg.ConnectionString = defaultConnectionString
  }

  cacheDir, exist := os.LookupEnv("XDG_CACHE_HOME")
  if exist == false {
    cacheDir = filepath.Join(os.Getenv("HOME"), ".cache")
  }

  defaultDatabaseCachePath := filepath.Join(cacheDir, "superhighway84", "database")
	// Migration step from old CachePath to new DatabaseCachePath
	if cfg.CachePath != "" {
		defaultDatabaseCachePath = cfg.CachePath
	}
  fmt.Printf("Database cache path [%s]: ", defaultDatabaseCachePath)
  fmt.Scanln(&cfg.DatabaseCachePath)
  if strings.TrimSpace(cfg.DatabaseCachePath) == "" {
    cfg.DatabaseCachePath = defaultDatabaseCachePath
  }
  os.MkdirAll(filepath.Dir(cfg.DatabaseCachePath), 0755)

  defaultProgramCachePath := filepath.Join(cacheDir, "superhighway84", "program")
	// Migration step from old CachePath to new DatabaseCachePath
	if cfg.CachePath != "" {
		// If the previous CachePath was used, the folder already contains the
		// OrbitDB, hence we need to find a different place
		defaultProgramCachePath = filepath.Join(cacheDir, "superhighway84.program")
	}
  fmt.Printf("Program cache path [%s]: ", defaultProgramCachePath)
  fmt.Scanln(&cfg.ProgramCachePath)
  if strings.TrimSpace(cfg.ProgramCachePath) == "" {
    cfg.ProgramCachePath = defaultProgramCachePath
  }
  os.MkdirAll(filepath.Dir(cfg.ProgramCachePath), 0755)

  defaultLogfile := filepath.Join(cacheDir, "superhighway84.log")
	if cfg.Logfile != "" {
		defaultLogfile = cfg.Logfile
	}
  fmt.Printf("Logfile path [%s]: ", defaultLogfile)
  fmt.Scanln(&cfg.Logfile)
  if strings.TrimSpace(cfg.Logfile) == "" {
    cfg.Logfile = defaultLogfile
  }


  fmt.Printf("\nProfile information\n-------------------\n\n")

  defaultProfileFrom := fmt.Sprintf("%s@localhost", os.Getenv("USER"))
	if cfg.Profile.From != "" {
		defaultProfileFrom = cfg.Profile.From
	}
  fmt.Printf("From [%s]: ", defaultProfileFrom)
  fmt.Scanln(&cfg.Profile.From)
  if strings.TrimSpace(cfg.Profile.From) == "" {
    cfg.Profile.From = defaultProfileFrom
  }

	defaultProfileOrganization := ""
	if cfg.Profile.Organization != "" {
		defaultProfileOrganization = cfg.Profile.Organization
	}
  fmt.Printf("Organization [%s]: ", defaultProfileOrganization)
  fmt.Scanln(&cfg.Profile.Organization)

  return cfg.Persist()
}


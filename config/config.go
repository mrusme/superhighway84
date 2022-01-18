package config

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/gdamore/tcell/v2"
)

type ConfigProfile struct {
  From                string
  Organization        string
}

type ConfigShortcuts struct {
  Refresh             int64
  Quit                int64

  FocusGroups         int64
  FocusArticles       int64
  FocusPreviews       int64

  NewArticle          int64
  ReplyToArticle      int64
}

type Config struct {
  ConfigFile          string `toml:"-"`

  ConnectionString    string

  CachePath           string // Deprecated, should be removed soon
	DatabaseCachePath   string
  ProgramCachePath    string

  Logfile             string

  Profile             ConfigProfile

  Shortcuts           map[string]string
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
  cfg.Shortcuts = make(map[string]string)
  _, err = toml.Decode(string(configFileContent), &cfg)
  if err != nil {
    return nil, errors.New("The config could not be parsed, make sure it is valid TOML and you don't have double assignments.")
  }

  cfg.ConfigFile = configFile
  err = cfg.LoadDefaults()
  if err != nil {
    return nil, err
  }
  return cfg, nil
}

func (cfg *Config) LoadDefaults() (error) {
  if len(cfg.Shortcuts) == 0 {
    cfg.Shortcuts[strconv.FormatInt(int64(tcell.KeyCtrlQ), 10)] = "quit"
    cfg.Shortcuts[strconv.FormatInt(int64(tcell.KeyCtrlR), 10)] = "refresh"
    cfg.Shortcuts[strconv.FormatInt(int64(tcell.KeyCtrlH), 10)] = "focus-groups"
    cfg.Shortcuts[strconv.FormatInt(int64(tcell.KeyCtrlL), 10)] = "focus-articles"
    cfg.Shortcuts[strconv.FormatInt(int64(tcell.KeyCtrlK), 10)] = "focus-articles"
    cfg.Shortcuts[strconv.FormatInt(int64(tcell.KeyCtrlJ), 10)] = "focus-preview"

    cfg.Shortcuts[strconv.FormatInt(int64('n'), 10)]            = "article-new"
    cfg.Shortcuts[strconv.FormatInt(int64('r'), 10)]            = "article-reply"
    cfg.Shortcuts[strconv.FormatInt(int64(tcell.KeyCtrlA), 10)] = "article-mark-all-read"

    cfg.Shortcuts[strconv.FormatInt(int64('h'), 10)]            = "additional-key-left"
    cfg.Shortcuts[strconv.FormatInt(int64('j'), 10)]            = "additional-key-down"
    cfg.Shortcuts[strconv.FormatInt(int64('k'), 10)]            = "additional-key-up"
    cfg.Shortcuts[strconv.FormatInt(int64('l'), 10)]            = "additional-key-right"

    cfg.Shortcuts[strconv.FormatInt(int64('g'), 10)]            = "additional-key-home"
    cfg.Shortcuts[strconv.FormatInt(int64('G'), 10)]            = "additional-key-end"

    cfg.Shortcuts[strconv.FormatInt(int64(tcell.KeyF8), 10)]    = "play"
  }
  return cfg.Persist()
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


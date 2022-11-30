// Package config implements the configuration of the search service.
package config

import "time"

// Default configuration.
const (
	DefaultWebPort     = 9050
	DefaultWebHost     = ""
	DefaultMaxQueue    = 5
	DefaultIndexAge    = 100 * time.Millisecond
	DefaultIndexFile   = "search.bleve"
	DefaultIndexUpdate = 2 * time.Minute
	DefaultIndexBatch  = 4096
	DefaultModels      = "models.yml"
	DefaultSearch      = "search.yml"
	DefaultDB          = "openslides"
	DefaultDBUser      = "openslides"
	DefaultDBPassword  = "openslides"
	DefaultDBHost      = "localhost"
	DefaultDBPort      = 5432
)

// Web are the parameters for the web server.
type Web struct {
	Port     int
	Host     string
	MaxQueue int
}

// Index are the parameters for the indexer.
type Index struct {
	File   string
	Age    time.Duration
	Update time.Duration
	Batch  int
}

// Models are the paths to the YAML files containing the models
// and the searched collections.
type Models struct {
	Models string
	Search string
}

// Database are the credentials for the datavbase.
type Database struct {
	Database string
	User     string
	Password string
	Host     string
	Port     int
}

// Config is the configuration of the search service.
type Config struct {
	Web      Web
	Index    Index
	Models   Models
	Database Database
}

// GetConfig returns the configuration overwritten with env vars.
func GetConfig() (*Config, error) {
	cfg := &Config{
		Web: Web{
			Port: DefaultWebPort,
			Host: DefaultWebHost,
		},
		Index: Index{
			File:   DefaultIndexFile,
			Age:    DefaultIndexAge,
			Update: DefaultIndexUpdate,
			Batch:  DefaultIndexBatch,
		},
		Models: Models{
			Models: DefaultModels,
			Search: DefaultSearch,
		},
		Database: Database{
			Database: DefaultDB,
			User:     DefaultDBUser,
			Password: DefaultDBPassword,
			Host:     DefaultDBHost,
			Port:     DefaultDBPort,
		},
	}
	if err := cfg.fromEnv(); err != nil {
		return nil, err
	}
	return cfg, nil
}

// fromEnv fills the config from env vars.
func (cfg *Config) fromEnv() error {
	return storeFromEnv([]storeEnv{
		{"OPENSLIDES_SEARCH_PORT", storeInt(&cfg.Web.Port)},
		{"OPENSLIDES_SEARCH_HOST", storeString(&cfg.Web.Host)},
		{"OPENSLIDES_SEARCH_MAX_QUEUED", storeInt(&cfg.Web.MaxQueue)},
		{"OPENSLIDES_SEARCH_INDEX_AGE", storeDuration(&cfg.Index.Age)},
		{"OPENSLIDES_SEARCH_INDEX_FILE", storeString(&cfg.Index.File)},
		{"OPENSLIDES_SEARCH_INDEX_BATCH", storeInt(&cfg.Index.Batch)},
		{"OPENSLIDES_SEARCH_INDEX_UPDATE_INTERVAL", storeDuration(&cfg.Index.Update)},
		{"OPENSLIDES_MODELS_YML", storeString(&cfg.Models.Models)},
		{"OPENSLIDES_SEARCH_YML", storeString(&cfg.Models.Search)},
		{"OPENSLIDES_DB", storeString(&cfg.Database.Database)},
		{"OPENSLIDES_DB_USER", storeString(&cfg.Database.User)},
		{"OPENSLIDES_DB_PASSWORD", storeString(&cfg.Database.Password)},
		{"OPENSLIDES_DB_HOST", storeString(&cfg.Database.Host)},
		{"OPENSLIDES_DB_PORT", storeInt(&cfg.Database.Port)},
	})
}

package conf

import (
	"log"
	"os"

	"github.com/kelseyhightower/envconfig"
	"gopkg.in/yaml.v2"
)

// Pubstore configuration
// split_words true is how envconfig maps e.g. the PUBSTORE_PUBLIC_BASE_URL environment variable to PublicBaseUrl
type Config struct {
	// Port on which the server is running
	Port int `yaml:"port"`
	// Public base URL of the API
	PublicBaseUrl string `yaml:"public_base_url" split_words:"true"`
	// Data Source Name of the pubstore database
	DSN string `yaml:"dsn"`
	// OAuth seed
	OAuthSeed string `yaml:"oauth_seed" envconfig:"OAUTH_SEED"`
	// Path to static files and views
	RootDir string `yaml:"root_dir" split_words:"true"`
	// Path to resources, especially cover images
	//Resources string `yaml:"resources"`
	// Page size used in the REST API and Web interface
	PageSize int `yaml:"page_size"  split_words:"true"`
	// LCP print and copy limits set in LCP licenses generated from the associated LCP Server
	PrintLimit int `yaml:"print_limit"  split_words:"true"`
	CopyLimit  int `yaml:"copy_limit"  split_words:"true"`
	// Basic Auth credentials used by the LCP encryption tool to notify Pubstore of a new encrypted publication
	UserName string `yaml:"username"`
	Password string `yaml:"password"`
	// LCP Server
	LCPServer LCPServerAccess `yaml:"lcp_server"`
}

// LCP Server access parameters
type LCPServerAccess struct {
	Url      string `yaml:"url"`
	Version  string `yaml:"version"`
	UserName string `yaml:"username"`
	Password string `yaml:"password"`
}

func Init(configFile string) (Config, error) {

	var cfg Config

	log.Printf("Loading configuration from %s ...", configFile)

	f, err := os.Open(configFile)
	if err != nil {
		log.Println("Configuration file not found")
	}
	defer f.Close()

	if f != nil {
		decoder := yaml.NewDecoder(f)
		err = decoder.Decode(&cfg)
		if err != nil {
			return cfg, err
		}
	}

	err = envconfig.Process("pubstore", &cfg)
	if err != nil {
		return cfg, err
	}

	// Set some defaults
	if cfg.Port == 0 {
		cfg.Port = 8080
	}
	if cfg.PublicBaseUrl == "" {
		cfg.PublicBaseUrl = "http://localhost:8080"
	}
	if cfg.DSN == "" {
		cfg.DSN = "sqlite3://pubstore.sqlite"
	}
	if cfg.RootDir == "" {
		cfg.RootDir, err = os.Getwd()
		if err != nil {
			return cfg, err
		}
	}
	if cfg.LCPServer.Version == "" {
		cfg.LCPServer.Version = "v2"
	}

	return cfg, nil
}

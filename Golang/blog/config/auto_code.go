package config

type Autocode struct {
	TransferRestart bool   `mapstructure:"transfer-restart" json:"transfer-restart" yaml:"transfer-restart"`
	Root            string `mapstructure:"root" json:"root" yaml:"root"`
	Server          string `mapstructure:"service" json:"service" yaml:"service"`
	SApi            string `mapstructure:"service-api" json:"service-api" yaml:"service-api"`
	SPlug           string `mapstructure:"service-plug" json:"service-plug" yaml:"service-plug"`
	SInitialize     string `mapstructure:"service-initialize" json:"service-initialize" yaml:"service-initialize"`
	SModel          string `mapstructure:"service-model" json:"service-model" yaml:"service-model"`
	SRequest        string `mapstructure:"service-request" json:"service-request"  yaml:"service-request"`
	SRouter         string `mapstructure:"service-router" json:"service-router" yaml:"service-router"`
	SService        string `mapstructure:"service-service" json:"service-service" yaml:"service-service"`
	Web             string `mapstructure:"web" json:"web" yaml:"web"`
	WApi            string `mapstructure:"web-api" json:"web-api" yaml:"web-api"`
	WForm           string `mapstructure:"web-form" json:"web-form" yaml:"web-form"`
	WTable          string `mapstructure:"web-table" json:"web-table" yaml:"web-table"`
}

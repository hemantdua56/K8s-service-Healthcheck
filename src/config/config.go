package config


type ServiceConfig struct {
	Service				string  `yaml:"service"`
	Port				string  `yaml:"port"`
	Context_Path		string  `yaml:"context_path"`
	Path				string  `yaml:"path"`
}


type Service struct {
	Microservice []ServiceConfig
}

package configuration

type KaiConfiguration struct {
	Servers []Server `yaml:"servers"`
}

type Server struct {
	Name      string `yaml:"name"`
	URL       string `yaml:"url"`
	IsDefault bool   `yaml:"default"`
}

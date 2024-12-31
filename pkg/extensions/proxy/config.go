package proxy

type Options struct {
	Debug bool `yaml:"debug"`
}
type Config struct {
	Options Options `yaml:"proxy"`
}

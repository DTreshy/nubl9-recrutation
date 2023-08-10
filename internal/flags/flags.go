package flags

import "flag"

type Flags struct {
	ConfigPath string
}

func Parse() *Flags {
	configPath := flag.String("config", "config.yml", "Path to the config file")
	flag.Parse()

	f := Flags{
		ConfigPath: *configPath,
	}

	return &f
}

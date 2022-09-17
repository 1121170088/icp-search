package init

import (
	"flag"
	"gopkg.in/yaml.v3"
	"log"
	"os"
	"path/filepath"
)

var (
	Cfg *Config
	cfgDir string
	configFile = "config.yaml"
)

type Config struct {
	Dsn string
	Addr string
	Upstream string
}



func init()  {
	flag.StringVar(&cfgDir, "d", "./", "set config directoty, defalt ./")
	flag.Parse()
	configFile = filepath.Join(cfgDir, configFile)
	Cfg = &Config{
		Dsn:"user:password@/database",
		Addr: "127.0.0.1:9090",
		Upstream: "",
	}
	if _, err := os.Stat(configFile); err != nil {
		log.Printf("config file dosn't exist, writing it")
		bytes, err := yaml.Marshal(Cfg)
		if err != nil {
			log.Fatal(err)
		}
		err = os.WriteFile(configFile, bytes, os.ModePerm)
		if err != nil {
			log.Fatal(err)
		}
		os.Exit(1)
	} else {
		bytes, err := os.ReadFile(configFile)
		if err != nil {
			log.Fatal(err)
		}

		err = yaml.Unmarshal(bytes,Cfg)
		if err != nil {
			log.Fatal(err)
		}
	}
}

package init

import (
	"flag"
	"github.com/1121170088/find-domain/search"
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
	Dsn string `yaml:"dsn"`
	Addr string `yaml:"addr"`
	Doh string `yaml:"doh"`
	IpInfoServer string `yaml:"ip-info-server"`
	Upstream string `yaml:"upstream"`
	Code0Index  int  `yaml:"code0-index"`
	Code0ConCurrent int `yaml:"code0-concurrent"`
	IpIndex int `yaml:"ip-index"`
	IpCheckConCurrent int `yaml:"ip-check-concurrent"`
	DomainSuffixFile string `yaml:"domain-suffix-file"`
}



func init()  {
	flag.StringVar(&cfgDir, "d", "./", "set config directoty, defalt ./")
	flag.Parse()
	configFile = filepath.Join(cfgDir, configFile)
	Cfg = &Config{
		Dsn:"user:password@/database",
		Addr: "127.0.0.1:9090",
		Doh: "dns.alidns.com",
		IpInfoServer: "",
		Upstream: "",
		Code0Index: 0,
		IpIndex: 0,
		IpCheckConCurrent: 1,
		DomainSuffixFile: "",
	}
	if _, err := os.Stat(configFile); err != nil {
		log.Printf("config file dosn't exist, writing it")
		WriteConfig()
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
	search.Init(Cfg.DomainSuffixFile)

}

func WriteConfig()  {
	bytes, err := yaml.Marshal(Cfg)
	if err != nil {
		log.Fatal(err)
	}
	err = os.WriteFile(configFile, bytes, os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}
}

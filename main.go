package main

import (
	"icp-search/dao"
	init_ "icp-search/init"
	s "icp-search/server"
	"icp-search/service/geoip"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main()  {
	log.Printf("%v", init_.Cfg)
	dao.Init()
	geoip.Init(init_.Cfg.CountryDbFile)

	go s.Start()


	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	<-sigCh

	geoip.Uninit()
	dao.UnInit()
	init_.WriteConfig()
}
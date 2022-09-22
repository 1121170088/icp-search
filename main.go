package main

import (
	"icp-search/dao"
	init_ "icp-search/init"
	s "icp-search/server"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main()  {
	log.Printf("%v", init_.Cfg)
	dao.Init()

	go s.Start()


	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	<-sigCh

	dao.UnInit()
	init_.WriteConfig()
}
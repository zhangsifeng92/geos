package main

import (
	_ "github.com/eosspark/eos-go/appbase/plugin/net_plugin"
	_ "github.com/eosspark/eos-go/appbase/plugin/producer_plugin"
	"flag"
	"github.com/eosspark/eos-go/exception/try"
	"os"
	"os/signal"
	. "github.com/eosspark/eos-go/appbase/app"
	. "github.com/eosspark/eos-go/appbase/app/include"
	"fmt"
)

const (
	OTHER_FAIL              = -2
	INITIALIZEFAIL          = -1
	SUCCESS                 = 0
	BAD_ALLOC               = 1
	DATABASE_DIRTY          = 2
	FIXED_REVERSIBLE        = 3
	EXTRACTED_GENESIS       = 4
	NODE_MANAGEMENT_SUCCESS = 5
)

var (
	PluginFromConfig string
	//Name string
	//Age int
)

func init() {
	flag.StringVar(&PluginFromConfig, "plugin", "net_plugin", "Plugin(s) to enable, may be specified multiple times")
	//flag.StringVar(&name,"name","","what's your names")
	//flag.IntVar(&age,"age",22,"how old are yous")

}

//var pro producer_plugin.Producer_plugin
//var net net_plugin.Net_plugin

func main() {
	try.Try(func() {
		App.SetVersion(Version)
		App.SetDefaultDataDir()
		App.SetDefaultConfigDir()
		fmt.Println(App.My.DateDir)
		fmt.Println(App.My.ConfigDir)
		fmt.Println(App.My.Version)
		App.My.Options.Run(os.Args)
		if !App.Initialize() {
			panic(INITIALIZEFAIL)
		}
		App.StartUp()

		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, os.Interrupt)
		select {
		case <-sigChan:
			App.ShutDown()
		}
	}).Catch(func() {

	}).End()

}
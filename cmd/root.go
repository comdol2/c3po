package cmd

import (
	"fmt"
	"log"
	"os"
	"runtime"

	c3po "github.com/comdol2/c3po/cmd"
	"github.com/spf13/cobra"
)

var sClient *c3po.Client
var version string
var debug bool

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:     "c3po",
	Version: version,
	Short:   "test code",
	Long:    `This is a test code for Sam's learning`,
	Run: func(cmd *cobra.Command, args []string) {

	},
}

func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	RootCmd.PersistentFlags().BoolVarP(&debug, "debug", "d", false, "To turn-on debugging")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	strMyOS := "Linux"
	if runtime.GOOS == "windows" {
		strMyOS = "Windows"
	}
	if debug {
		fmt.Println("My OS : ", strMyOS)
	}

	c3poAccessToken, err := sClient.getAccessToken()
	if err != nil {
		log.Fatalf("ERROR: %v", err)
	}

	sClient = c3po.NewClient(c3poAccessToken, debug)
	if sClient == nil {
		log.Fatalf("ERROR: Can't create SNOW client")
	}

}

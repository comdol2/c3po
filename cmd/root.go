package cmd

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"runtime"
	"strings"

	c3po "github.com/comdol2/c3po/api"
	"github.com/spf13/cobra"
	"golang.org/x/term"
)

var sClient *c3po.Client
var version, c3poAccessToken string
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

	fmt.Println("My OS : ", strMyOS, "\n")

        reader := bufio.NewReader(os.Stdin)

        fmt.Print("*** Enter Your HubID : ")
        username, _ := reader.ReadString('\n')
        username = strings.TrimSpace(username) // Remove any trailing newline characters

        fmt.Print("*** Enter Password: ")
        bytePassword, _ := term.ReadPassword(int(os.Stdin.Fd()))
        password := string(bytePassword)
	fmt.Println("")

	token := ""

	var err error
	sClient, err = c3po.NewClient(username, password, token, debug)
	if err != nil {
		log.Fatalf("ERROR: Can't create C3PO client: %v", err)
	}
	if sClient == nil {
		log.Fatalf("ERROR: sClient is nil after NewClient call")
	}

	c3poAccessToken, err = sClient.GetAccessToken()
	if err != nil {
		log.Fatalf("ERROR: %v", err)
	}

	fmt.Println("Access Token : ", c3poAccessToken)

}

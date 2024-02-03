package main

import (
        "fmt"
        "github.com/comdol2/c3po/cmd"
        "os"
)

func main() {
        if err := cmd.RootCmd.Execute(); err != nil {
                fmt.Println(err)
                os.Exit(-1)
        }
}


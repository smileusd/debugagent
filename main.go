package main

import (
    "os"
    "fmt"
    "../debugagent/client"
)

const VERSION = "v0.1"

func cleanup() {
    if r := recover(); r != nil {
        e, isErr := r.(error)
        if isErr {
            fmt.Printf("%v", e)
        }
        os.Exit(1)
    }
}

func main() {
    defer cleanup()

    cli := client.NewCli(VERSION)
    err := cli.Run(os.Args)
    if err != nil {
        panic(fmt.Errorf("Error when executing command: %v", err))
    }
}

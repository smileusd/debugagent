package client

import 	"github.com/codegangsta/cli"

const DEFAULT_ADDRESS = "127.0.0.1:10888"

var (
    DaemonFlags = []cli.Flag{
        cli.StringFlag{
            Name:  "log",
            Usage: "specific output log file, otherwise output to stdout by default",
        },
        cli.StringFlag{
            Name:  "addr",
            Value: DEFAULT_ADDRESS,
            Usage: "specific daemon start listen ip address",
        },
    }
)

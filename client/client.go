package client

import 	(
    "github.com/codegangsta/cli"
    "../daemon"
)

// NewCli would generate debug
func NewCli(version string) *cli.App {
    app := cli.NewApp()
    app.Name = "debugagent"
    app.Version = version
    app.Usage = "A debug tool agent to debug local network problem pod"
    app.Flags = []cli.Flag{
        cli.StringFlag{
            Name:  "address, a",
            Value: DEFAULT_ADDRESS,
            Usage: "Specify address for communication between server and client",
        },
    }
    app.Commands = []cli.Command{
        daemonCmd,
    }
    return app
}

var (
    daemonCmd = cli.Command{
        Name:   "daemon",
        Usage:  "start debugagent daemon",
        Flags:  DaemonFlags,
        Action: cmdStartDaemon,
    }
)

func cmdStartDaemon(c *cli.Context) {
    if err := startDaemon(c); err != nil {
        panic(err)
    }
}

func startDaemon(c *cli.Context) error {
    config := &daemon.DebugDaemonConfig{
        Address: c.String("addr"),
    }
    stopCh := make(chan bool)
    return daemon.Start(config, stopCh)
}


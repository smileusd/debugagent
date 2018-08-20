package client

import "io"

type DebugExecResponse struct {
    Stdin io.Reader
    Stdout io.Writer
}

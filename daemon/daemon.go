package daemon

import (
    "fmt"
    "net/http"
    "os/exec"
    "os"
    "log"
    "encoding/json"
    "../client"
    "io"
    "github.com/emicklei/go-restful"
    "net"
)

const NETWORKDEBUGIMAGE = "busybox:networkdebug"

type DebugDaemonConfig struct {
    Address string
}

func IndexHandler(r *restful.Request, w *restful.Response) {
    request := &client.DebugExecRequest{}
    if err := decodeRequest(r, request); err != nil {
        return err
    }
    id, err := findPauseContainer(request)
    if err != nil {
        w.Write(writeError(err))
        return
    }
    stdin, stdout, err := execContainer(id, request.Args)
    if err != nil {
        w.Write(writeError(err))
        return
    }
    // TODO: return exec stream
    response := &client.DebugExecResponse{
        Stdin: stdin,
        Stdout: stdout,
    }
    data, err := responseOutput(response)
    if err != nil {
        w.Write(writeError(err))
        return
    }
    w.Write(data)
}

func writeError(err error) []byte{
    return []byte(fmt.Errorf("failed to exec pod : %v", err).Error())
}

// find puase container according do
func findPauseContainer(request *client.DebugExecRequest) (string, error) {
    pauseName := "k8s_POD_"+request.PodName+"_"+request.Namespace
    grepCommand := fmt.Sprintf("docker ps | grep %s | grep /pause | awk {'print $1'}", pauseName)
    out, err := exec.Command("sh", "-c", grepCommand).CombinedOutput()
    if err != nil {
        return "", err
    }
    return string(out)[:len(string(out))-1], nil
}

func execContainer(containerId string, args []string) (io.Reader, io.Writer, error) {
    cmd := exec.Command("docker", "run", "-it", "--net=container:" + containerId, "--ipc=container:" + containerId, "--pid=container:" + containerId, NETWORKDEBUGIMAGE, args...)
    cmd.Stdin = os.Stdin
    cmd.Stdout = os.Stdout
    err := cmd.Start()
    if err != nil {
        log.Println(err)
        return nil, nil,err
    }
    err = cmd.Wait()
    if err != nil {
        log.Printf("Command finished with error: %v", err)
        return nil, nil, err
    }
    return cmd.Stdin, cmd.Stdout, nil
}

func Start(daemonConfig *DebugDaemonConfig, stopCh chan bool) error{
    http.HandleFunc("/", IndexHandler)
    restContainer := restful.NewContainer()
    ws := new(restful.WebService)
    ws.Path("/")
    ws.Consumes("*/*").Produces("*/*")
    restContainer.Add(ws)

    ws.Route(ws.POST("/").To(IndexHandler).
            Doc("Debug container").
            Reads(*client.DebugExecRequest{}).
            Do(returns200(client.DebugExecResponse{}), returns400))
    l, err := net.Listen("tcp", daemonConfig.Address)
    if err != nil {
        fmt.Printf("listen err", err)
        return err
    }
    defer l.Close()
    go func() {
        err := http.Serve(l, restContainer)
        if err != nil {
            fmt.Printf("serve error: %v", err)
        }
        stopCh <- true
    }()
    <-stopCh
    return nil
}

func returns200(model interface{}) func(*restful.RouteBuilder) {
    return func(b *restful.RouteBuilder) {
        b.Returns(http.StatusOK, "OK", model)
    }
}

func returns400(b *restful.RouteBuilder) {
    b.Returns(http.StatusBadRequest, "Something went wrong", nil)
}

func decodeRequest(r *http.Request, v interface{}) error {
    return json.NewDecoder(r.Body).Decode(v)
}

// ResponseOutput would generate a JSON format byte array of object for output
func responseOutput(v interface{}) ([]byte, error) {
    j, err := json.MarshalIndent(v, "", "\t")
    if err != nil {
        return nil, err
    }
    return j, nil
}

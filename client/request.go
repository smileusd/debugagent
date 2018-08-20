package client

type DebugExecRequest struct{
    Pod int64
    PodName string
    Namespace string
    Args      []string
}


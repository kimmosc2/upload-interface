package main

import (
    "log"
    "net/http"
    "os"
    "upload-interface/heartbeat"
    "upload-interface/locate"
    "upload-interface/objects"
)

func main() {
    // log.SetFlags(log.Llongfile)
    go heartbeat.ListenHeartBeat()
    http.HandleFunc("/objects/", objects.Handler)
    http.HandleFunc("/locate/",locate.Handler)
    log.Fatal(http.ListenAndServe(os.Getenv("LISTEN_ADDRESS"), nil))
}

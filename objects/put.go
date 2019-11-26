package objects

import (
    "fmt"
    "io"
    "log"
    "net/http"
    "strings"
    "upload-interface/heartbeat"
    "upload-interface/util/objectstream"
)

func put(w http.ResponseWriter, r *http.Request) {
    object := strings.Split(r.URL.EscapedPath(), "/")[2]
    c, err := storeObject(r.Body, object)
    if err != nil {
        log.Println(err.Error())
    }
    w.WriteHeader(c)
}

func storeObject(r io.Reader, object string) (int, error) {
    stream, err := putStream(object)
    if err != nil {
        return http.StatusServiceUnavailable, err
    }
    io.Copy(stream, r)
    err = stream.Close()
    if err != nil {
        return http.StatusInternalServerError, err
    }
    return http.StatusOK, nil
}

func putStream(object string) (*objectstream.PutStream, error) {
    server := heartbeat.ChooseRandomDataServer()
    if server == "" {
        return nil, fmt.Errorf("cannot find any dataServer")
    }
    return objectstream.NewPutStream(server, object), nil
}
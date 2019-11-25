package objects

import (
    "fmt"
    "io"
    "log"
    "net/http"
    "strings"
    "upload-interface/locate"
    "upload-interface/util/objectstream"
)

func get(w http.ResponseWriter, r *http.Request) {
    object := strings.Split(r.URL.EscapedPath(), "/")[2]
    stream, err := getStream(object)
    if err != nil {
        log.Println(err)
        w.WriteHeader(http.StatusNotFound)
        return
    }
    io.Copy(w, stream)
}

// 调用locate.Locate()方法定位,返回服务器地址, 再调用
// objectstream.NewGetStream()方法请求目标服务器，
// 将响应体返回给客户端
func getStream(object string) (io.Reader, error) {
    server := locate.Locate(object)
    if server == "" {
        return nil, fmt.Errorf("object %s locate fail", object)
    }
    return objectstream.NewGetStream(server, object)
}

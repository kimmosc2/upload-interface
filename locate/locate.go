package locate

import (
    "encoding/json"
    "net/http"
    "os"
    "strconv"
    "strings"
    "time"
    "upload-interface/util/rabbitmq"
)

func Handler(w http.ResponseWriter, r *http.Request) {
    m := r.Method
    if m != http.MethodGet {
        w.WriteHeader(http.StatusMethodNotAllowed)
        return
    }
    info := Locate(strings.Split(r.URL.EscapedPath(), "/")[2])
    if len(info) == 0 {
        w.WriteHeader(http.StatusNotFound)
        return
    }
    b, _ := json.Marshal(info)
    w.Write(b)
}

// 定位函数，当调用时会向dataServers exchange发送定位消息，并将临时队列
// 告知目标,如果有人持有对象,则将服务器地址发送至临时消息队列,若1S后无响应，
// 则关闭队列
func Locate(name string) string {
    q := rabbitmq.New(os.Getenv("RABBITMQ_SERVER"))
    q.Publish("dataServers", name)
    c := q.Consume()
    go func() {
        time.Sleep(time.Second)
        q.Close()
    }()
    msg := <-c
    s, _ := strconv.Unquote(string(msg.Body))
    return s
}

func Exist(name string) bool {
    return Locate(name) != ""
}

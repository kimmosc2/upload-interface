package heartbeat

import (
    "math/rand"
    "os"
    "strconv"
    "sync"
    "time"
    "upload-interface/util/rabbitmq"
)

var dataServers = make(map[string]time.Time)
var mutex sync.Mutex

// 创建消息队列绑定apiServer exchange,监听心跳消息
func ListenHeartBeat() {
    q := rabbitmq.New(os.Getenv("RABBITMQ_SERVER"))
    defer q.Close()

    q.Bind("apiServers")
    c := q.Consume()
    go removeExpiredDataServer()
    for msg := range c {
        dataServer, err := strconv.Unquote(string(msg.Body))
        if err != nil {
            panic(err)
        }
        mutex.Lock()
        dataServers[dataServer] = time.Now()
        mutex.Unlock()
    }
}

// 清除十秒未收到心跳的服务器列表
func removeExpiredDataServer() {
    for {
        time.Sleep(5 * time.Second)
        mutex.Lock()
        for k, v := range dataServers {
            if v.Add(10 * time.Second).Before(time.Now()) {
                delete(dataServers, k)
            }
        }
        mutex.Unlock()
    }
}

// 将数据服务从map转为slice
// 遍历dataServers并返回当前所有数据服务节点
func GetDataServers() []string {
    mutex.Lock()
    defer mutex.Unlock()
    ds := make([]string,0)
    for s, _ := range dataServers {
        ds = append(ds, s)
    }
    return ds
}

// 随机选择一个节点
func ChooseRandomDataServer() string {
    ds := GetDataServers()
    n := len(ds)
    if n == 0 {
        return ""
    }
    return ds[rand.Intn(n)]
}

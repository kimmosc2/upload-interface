package objectstream

import (
    "fmt"
    "io"
    "net/http"
)

type GetStream struct {
    reader io.Reader
}

func newGetStream(url string) (*GetStream, error) {
    resp, err := http.Get(url)
    if err != nil {
        return nil, err
    }
    if resp.StatusCode != http.StatusOK {
        return nil, fmt.Errorf("dataServer return http code %d", resp.StatusCode)
    }
    return &GetStream{resp.Body}, nil
}

func NewGetStream(server, object string) (*GetStream, error) {
    if server == "" || object == "" {
        return nil, fmt.Errorf("invalid server %s object %s", server, object)
    }
    return newGetStream("http://" + server + "/objects/" + object)
}

func (r *GetStream) Read(p []byte) (n int, err error) {
    return r.reader.Read(p)
}
package gocache

import (
	"fmt"
	"log"
	"net/http"
	"strings"
)

const DEFAULT_PASH = "/gocache/"

type HTTPPool struct {
	self     string
	basePash string
}

func NewHTTPPool(self string) *HTTPPool {
	return &HTTPPool{
		self:     self,
		basePash: DEFAULT_PASH,
	}
}

func (this *HTTPPool) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	//先判断基础路径是否正确
	//然后判断路径是否合法
	//然后尝试取缓存
	path := r.URL.Path
	if !strings.HasPrefix(path, this.basePash) {
		panic("HTTPPool serving unexpected path: " + path)
	}
	this.Log("%s %s", r.Method, path)

	split := strings.Split(path, "/")
	if len(split) != 4 {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	groupName := split[2]
	key := split[3]

	group := GetGroup(groupName)
	if group == nil {
		http.Error(w, "no such group: "+groupName, http.StatusNotFound)
		return
	}

	value, err := group.Get(key)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/octet-stream")
	w.Write(value.ByteSlice())
	return
}

func (this *HTTPPool) Log(format string, v ...interface{}) {
	log.Printf("[Server %s] %s", this.self, fmt.Sprintf(format, v...))
}

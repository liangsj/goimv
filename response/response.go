package response

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
)

type ret struct {
	Msg   string
	Errno int
	Data  interface{}
}

func HttpErrorReturn(w http.ResponseWriter, err error) {
	log.Printf("err:%s\n",err.Error())
	ret := ret{
		Errno: -1,
		Msg:   err.Error(),
	}
	json, _ := json.Marshal(ret)
	io.WriteString(w, string(json))
}

func HttpReturn(w http.ResponseWriter, data interface{}) {
	r := ret{
		Errno: 0,
		Msg:   "",
		Data:  data,
	}
	json, _ := json.Marshal(r)
	io.WriteString(w, string(json))
}

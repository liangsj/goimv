package view

import (
	"io"
	"io/ioutil"
	"net/http"

	"github.com/goimv/response"
)

func Index(w http.ResponseWriter, req *http.Request) {
	index, err := ioutil.ReadFile("views/index.html")

	if err != nil {
		response.HttpErrorReturn(w, err)
		return
	}
	io.WriteString(w, string(index))
}

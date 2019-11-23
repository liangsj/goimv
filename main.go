package main

import (
	"log"
	"net/http"

	"github.com/goimv/goenv"
	"github.com/goimv/problem"
	view "github.com/goimv/views"
)

func main() {
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	//problem
	http.HandleFunc("/goimv/problem/list", problem.List)
	http.HandleFunc("/goimv/problem/content", problem.Content)
	http.HandleFunc("/problem/tips", problem.Tips)
	http.HandleFunc("/", view.Index)

	//go run
	http.HandleFunc("/goimv/goenv/save", goenv.Save)
	http.HandleFunc("/goimv/goenv/build", goenv.Build)
	http.HandleFunc("/goimv/goenv/run", goenv.Run)
	http.HandleFunc("/goimv/goenv/autocomplete", goenv.Autocomplete)

	log.Println("server start at 8888")
	err := http.ListenAndServe(":8888", nil)

	if err != nil {
		log.Fatal("ListenAndServer: ", err.Error())
	}
}

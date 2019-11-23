package problem

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/goimv/response"
)

func List(w http.ResponseWriter, req *http.Request) {
	//	start := req.FormValue("start")
	//	sum := req.FormValue("sum")
	problems := make([]string, 0)
	files, err := ioutil.ReadDir("problems")

	if err != nil {
		response.HttpErrorReturn(w, err)
	}

	for _, f := range files {

		if f.IsDir() && !strings.HasPrefix(f.Name(), ".") {
			problems = append(problems, f.Name())
		}

	}
	response.HttpReturn(w, problems)

}

func Content(w http.ResponseWriter, req *http.Request) {
	title := req.FormValue("title")
	cot, err := getContent(title, "content")

	if err != nil {

		response.HttpErrorReturn(w, err)
	}
	template, _ := getContent(title, "template")
	data := make(map[string]string)
	data["content"] = cot
	data["template"] = template
	response.HttpReturn(w, data)
}

func Tips(w http.ResponseWriter, req *http.Request) {

	title := req.FormValue("title")
	cot, err := getContent(title, "tips")
	if err != nil {

		response.HttpErrorReturn(w, err)
	}

	response.HttpReturn(w, cot)
}

func getContent(title string, fileName string) (string, error) {

	dir := "problems/" + title

	files, err := ioutil.ReadDir(dir)

	if err != nil {
		return "", err
	}

	for _, f := range files {
		if f.Name() == fileName {
			file := "problems/" + title + "/" + f.Name()
			content, err := ioutil.ReadFile(file)
			return string(content), err
		}
	}
	return "", fmt.Errorf(fileName + "not exist")
}

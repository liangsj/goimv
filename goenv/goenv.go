package goenv

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/goimv/response"
)

func Save(w http.ResponseWriter, req *http.Request) {
	code := req.FormValue("code")
	title := req.FormValue("title")
	dir := fmt.Sprintf("problems/%s/", title)
	fileName := "problem.go"

	filePath := filepath.Clean(dir + "/" + fileName)
	fout, err := os.Create(filePath)
	if err != nil {
		response.HttpErrorReturn(w, err)
	}
	fout.WriteString(code)

	cmd := exec.Command("gofmt", "-w", dir+fileName)
	cmd.CombinedOutput()
	cot, err := ioutil.ReadFile(dir + fileName)

	if err != nil {
		response.HttpErrorReturn(w, err)
	}
	if "" != string(cot) {
		code = string(cot)
	}

	data := map[string]interface{}{}

	data["code"] = code
	data["fileName"] = fileName

	if err := fout.Close(); nil != err {
		response.HttpErrorReturn(w, err)
		return
	}
	response.HttpReturn(w, data)
}

func Build(w http.ResponseWriter, req *http.Request) {
	title := req.FormValue("title")
	dir := fmt.Sprintf("problems/%s/", title)

	fileName := "problem.go"
	filePath := filepath.Clean(dir + fileName)

	data := map[string]interface{}{}

	executable := filepath.Clean(dir + strings.Replace(fileName, ".go", "", -1))

	cmd := exec.Command("/bin/bash", "-c", "GOOS=linux", "GOARCH=amd64", "go", "build", "-o", executable, filePath)
	out, err := cmd.CombinedOutput()

	data["output"] = template.HTML(string(out))

	if nil != err {
		response.HttpErrorReturn(w, err)
		return
	}

	data["executable"] = executable
	response.HttpReturn(w, data)
}

func Run(w http.ResponseWriter, req *http.Request) {
	title := req.FormValue("title")
	cmd := req.FormValue("cmd")
	rid := req.FormValue("rid")
	if rid == "" {
		randInt := rand.Int()
		rid = strconv.Itoa(randInt)
	}
	out, err := runCmd(title, cmd, rid)
	if err != nil {
		response.HttpErrorReturn(w, err)
	}
	response.HttpReturn(w, out)

}

func runCmd(title string, cmd string, rid string) (output string, err error) {
	dir := fmt.Sprintf("problems/%s/", title)
	tmpDir := rid + "workspace/"
	err = os.RemoveAll(tmpDir)
	defer os.RemoveAll(tmpDir)
	if err != nil {
		return
	}
	if err = os.Mkdir(tmpDir, os.ModePerm); err != nil {
		return

	}

	absDir, err := filepath.Abs(tmpDir)

	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}
		fContent, _ := ioutil.ReadFile(dir + "/" + file.Name())
		err = ioutil.WriteFile(tmpDir+"/"+file.Name(), fContent, 0644)
		if err != nil {
			return
		}
	}

	var command *exec.Cmd
	switch cmd {
	case "test":
		command = exec.Command("docker", "run", "--rm", "--name", rid, "-v", absDir+":/go", "golang", "go", cmd)
	case "run":
		command = exec.Command("docker", "run", "--rm", "--name", rid, "-v", absDir+":/go", "golang", "go", cmd, "problem.go")
	default:
		return "", fmt.Errorf(cmd + "not support")
	}
	log.Printf("%v", command)
	if err != nil {
		return "", nil
	}
	out := make(chan string)
	after := time.After(5 * time.Second)
	go func() {
		output, err := command.CombinedOutput()
		if err != nil {
			log.Println("err:", err)
		}
		out <- string(output)
		log.Printf("docker out:%s", string(output))
	}()
	select {
	case <-after:
		killCmd := exec.Command("docker", "rm", "-f", rid)
		if err := killCmd.Run(); nil != err {
			log.Println("executes [docker rm -f " + rid + "] failed [" + err.Error() + "], this will cause resource leaking")
		}
		return "your problem is too slow", err
	case msg := <-out:
		return string(msg), err
	}
}

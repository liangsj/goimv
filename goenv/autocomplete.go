package goenv

import (
	"bytes"
	"net/http"
	"os/exec"
	"strconv"
	"strings"

	"github.com/goimv/response"
)

// AutocompleteHandler handles request of code autocompletion.
func Autocomplete(w http.ResponseWriter, r *http.Request) {

	code := r.FormValue("code")
	line, err := strconv.Atoi(r.FormValue("cursorLine"))
	ch, err := strconv.Atoi(r.FormValue("cursorCh"))
	if err != nil {
		response.HttpErrorReturn(w, err)
	}

	offset := getCursorOffset(code, line, ch)

	argv := []string{"-f=json", "autocomplete", strconv.Itoa(offset)}
	cmd := exec.Command("gocode", argv...)

	stdin, _ := cmd.StdinPipe()
	stdin.Write([]byte(code))
	stdin.Close()

	output, err := cmd.CombinedOutput()
	if nil != err {
		response.HttpErrorReturn(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(output)
}

// getCursorOffset calculates the cursor offset.
//
// line is the line number, starts with 0 that means the first line
// ch is the column number, starts with 0 that means the first column
func getCursorOffset(code string, line, ch int) (offset int) {
	lines := strings.Split(code, "\n")

	// calculate sum length of lines before
	for i := 0; i < line; i++ {
		offset += len(lines[i])
	}

	// calculate length of the current line and column
	curLine := lines[line]
	var buffer bytes.Buffer
	r := []rune(curLine)
	for i := 0; i < ch; i++ {
		buffer.WriteString(string(r[i]))
	}

	offset += len(buffer.String()) // append length of current line
	offset += line                 // append number of '\n'

	return offset
}

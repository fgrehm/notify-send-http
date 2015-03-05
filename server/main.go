package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		fmt.Println("PORT has not been set!")
		os.Exit(1)
	}

	// TODO: https://github.com/guard/guard/blob/19351271941a3362a47176c6808ddcb4a675e3ad/lib/guard/notifiers/notifysend.rb#L15
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Println("Handling request")

		// parse request with maximum memory of _5Megabits
		if strings.HasPrefix(r.Header.Get("Content-Type"), "multipart/form-data;") {
			const _5M = (1 << 30) * 5
			if err := r.ParseMultipartForm(_5M); nil != err {
				fmt.Println(err)
				http.Error(w, "Error handling file upload", 500)
				return
			}
		} else {
			r.ParseForm()
		}

		summary := r.Form.Get("summary")
		body := r.Form.Get("body")
		timeout := r.Form.Get("timeout")

		if summary == "" {
			msg := fmt.Sprintf("Invalid request, `summary` has to be specified (got %+v)", r.Form)
			log.Println(msg)
			http.Error(w, msg, 500)
			return
		}

		iconPath := ""
		file, header, _ := r.FormFile("icon")
		if file != nil {
			tmpDir, err := ioutil.TempDir("/tmp", "notify-send-http")

			defer file.Close()
			defer os.RemoveAll(tmpDir)

			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			//create destination file making sure the path is writeable.
			iconPath = tmpDir + "/" + header.Filename
			dst, err := os.Create(iconPath)
			defer dst.Close()
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			//copy the uploaded file to the destination file
			if _, err := io.Copy(dst, file); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}

		msg := "`notify-send` should have been triggered"
		fmt.Fprintln(w, msg)
		log.Println(msg)
		notifySend(summary, body, iconPath, timeout)
	})
	http.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		// Nothing to do here
	})

	log.Println("Starting server for http://0.0.0.0:" + port)
	log.Fatal(http.ListenAndServe(":" + port, nil))
}

func notifySend(summary, body, iconPath, timeout string) {
	var (
		execCmd string
		args []string
	)

	// Try notify-send
	_, err := exec.LookPath("notify-send")
	if err == nil {
		args = []string{summary, body}
		if iconPath != "" {
			args = append([]string{"-i", iconPath}, args...)
		}
		if timeout != "" {
			args = append([]string{"-t", timeout}, args...)
		}
		execCmd = "notify-send"
	} else {
		// Try OSX terminal-notifier
		args = []string{"-title", summary, "-message", body}
		execCmd = "terminal-notifier"
	}
	cmd := exec.Command(execCmd, args...)
	_, err = cmd.CombinedOutput()
	if err != nil {
		log.Println(err)
	}
}

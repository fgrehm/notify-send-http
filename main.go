package main

import (
	"fmt"
	"log"
	"net/http"
	"os/exec"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Println("Handling request")

		r.ParseForm()
		summary := r.Form.Get("summary")
		body := r.Form.Get("body")

		if summary == "" || body == "" {
			msg := fmt.Sprintf("Invalid request, both `summary` and `body` needs to be specified (got %+v)", r.Form)
			log.Println(msg)
			http.Error(w, msg, 500)
		} else {
			msg := "`notify-send` should have been triggered"
			fmt.Fprintln(w, msg)
			log.Println(msg)
			notifySend(summary, body)
		}
	})
	http.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		// Nothing to do here
	})
	log.Println("Starting server for http://0.0.0.0:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func notifySend(summary, body string) {
	cmd := exec.Command("notify-send", summary, body)
	_, err := cmd.CombinedOutput()
	if err != nil {
		log.Println(err)
	}
}

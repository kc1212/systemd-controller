package main

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
)

func hello(name string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		s := "Inactive/NotFound"
		if status(name) == nil {
			s = "Active"
		}
		fmt.Fprintf(w, "<h1>Service "+name+" is "+s+"</h1>"+
			"<p><a href=\"/start/\">Start</a></p>"+
			"<p><a href=\"/stop/\">Stop</a></p>")
	}
}

func systemctl(command string, service string) error {
	var out bytes.Buffer
	cmd := exec.Command("systemctl", command, service)
	cmd.Stdout = &out
	log.Println(cmd)
	return cmd.Run()
}

func start(name string) error {
	return systemctl("start", name)
}

func stop(name string) error {
	return systemctl("stop", name)
}

func status(name string) error {
	return systemctl("status", name)
}

func makeHandler(f func(name string) error, name string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		f(name)
		http.Redirect(w, r, "/", 303)
	}
}

func main() {
	service := os.Args[1] // TODO handle error
	status(service)

	http.HandleFunc("/", hello(service))
	http.HandleFunc("/start/", makeHandler(start, service))
	http.HandleFunc("/stop/", makeHandler(stop, service))

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

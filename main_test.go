package main

import (
	"bytes"
	"log"
	"net/http"
	"os/exec"
	"testing"
)

func TestFirstController(t *testing.T) {
	// Start local webserver:
	cmd := exec.Command("./start")
	if err := cmd.Start(); err != nil {
		log.Fatal(err)
	}

	//test the fine case with standard controller plugin
	res, err := http.Get("http://localhost:8080/first")
	if err != nil {
		log.Fatal(err)
	}

	buf := new(bytes.Buffer)
	buf.ReadFrom(res.Body)
	body := buf.String()
	res.Body.Close()

	if res.StatusCode != http.StatusOK {
		t.Errorf("handler returned unexpected status %v", res.StatusCode)
	}

	// Check the response body is what we expect.
	expected := "Hello FROM CONTROLLER PLUGIN!!! "

	if body != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			body, expected)
	}

	err = cmd.Process.Kill()
	if err != nil {
		panic(err) // panic as can't kill a process.
	}
}

func TestBadPath(t *testing.T) {
	// Start local webserver:
	cmd := exec.Command("./start")
	if err := cmd.Start(); err != nil {
		log.Fatal(err)
	}

	//test the fine case with standard controller plugin
	res, err := http.Get("http://localhost:8080/gigi")
	if err != nil {
		log.Fatal(err)
	}

	if res.StatusCode != http.StatusNotFound {
		t.Errorf("handler returned unexpected status %v", res.StatusCode)
	}

	err = cmd.Process.Kill()
	if err != nil {
		panic(err) // panic as can't kill a process.
	}
}

func TestWrongMethod(t *testing.T) {
	// Start local webserver:
	cmd := exec.Command("./start")
	if err := cmd.Start(); err != nil {
		log.Fatal(err)
	}

	//test the fine case with standard controller plugin but wrong HTTP method
	res, err := http.Head("http://localhost:8080/first")
	if err != nil {
		log.Fatal(err)
	}

	if res.StatusCode != http.StatusBadRequest {
		t.Errorf("handler returned unexpected status %v", res.StatusCode)
	}

	err = cmd.Process.Kill()
	if err != nil {
		panic(err) // panic as can't kill a process.
	}
}

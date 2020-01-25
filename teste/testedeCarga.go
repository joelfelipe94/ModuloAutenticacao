package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"time"
)

func timeTrack(start time.Time, name string) {
	elapsed := time.Since(start)
	log.Printf("Tempo de resposta do %s é de %s", name, elapsed)
}

var done = make(chan bool)

func enviaNPostsHTTPPorSegundo(postsPorSegundo int) {
	intervalo := time.Second / time.Duration(postsPorSegundo)
	requestBody, err := json.Marshal(map[string]string{
		"usuario": "Joel Gaya",
		"senha":   "12345",
	})
	if err != nil {
		log.Fatalln(err)
	}
	for i := 0; i < postsPorSegundo; i++ {
		go enviaPostHTTP(bytes.NewBuffer(requestBody))
		time.Sleep(intervalo)
	}

}

func enviaPostHTTP(requestBody *bytes.Buffer) {
	defer timeTrack(time.Now(), "envio do post")
	resp, err := http.Post("http://127.0.0.1:56666/api/login", "application/json", requestBody)
	if err != nil {
		log.Fatalln(err)
	}
	defer resp.Body.Close()
	log.Println(resp.Status)
	done <- true
}

func main() {
	var numeroThreads int = 2
	var postsPorSegundo int = 4
	totalRequests := numeroThreads * postsPorSegundo
	for indexThread := 0; indexThread < numeroThreads; indexThread++ {
		go enviaNPostsHTTPPorSegundo(postsPorSegundo)
	}
	for indexRequest := 0; indexRequest < totalRequests; indexRequest++ {
		<-done
	}

}

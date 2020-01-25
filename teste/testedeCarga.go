package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"log"
	"net/http"
	"os"
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
	//Cria arquivo de log
	file, err := os.OpenFile("info.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	log.SetOutput(file)

	//Lê parametros do terminal
	numeroThreadsPtr := flag.Int("numeroThreads", 1, "Passe o numero de threads para o test")
	postsPorSegundoPtr := flag.Int("postsPorSegundo", 1,
		"Passe o numero de requests que cada thread faz por segundo")
	flag.Parse()
	//Inica as threads  aguarda o término
	numeroThreads := *numeroThreadsPtr
	postsPorSegundo := *postsPorSegundoPtr
	totalRequests := numeroThreads * postsPorSegundo
	for indexThread := 0; indexThread < numeroThreads; indexThread++ {
		go enviaNPostsHTTPPorSegundo(postsPorSegundo)
	}
	for indexRequest := 0; indexRequest < totalRequests; indexRequest++ {
		<-done
	}

}

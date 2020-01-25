package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"flag"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joelfelipe94/ModuloAutenticacao/gerentedados"
)

//buscaRegistroAleatorio faz uma consulta ao banco e busca uma credencial valida
func buscaRegistroAleatorio() (string, string) {
	db, err := sql.Open("mysql", gerentedados.StringdeConexao)
	if err != nil {
		log.Fatalln(err)
	}
	var senha string
	var nome string
	db.QueryRow("SELECT usuario, senha FROM table ORDER BY RAND() LIMIT 1").Scan(&nome, &senha)
	return nome, senha
}

//timeTrack calcula o tempo entre start e o presente momento e salva em log
func timeTrack(start time.Time, name string) {
	elapsed := time.Since(start)
	log.Printf("Tempo de resposta do %s é de %s", name, elapsed)
}

// cria o channel usado para os resquests avisarem que terminaram
var done = make(chan bool)

//enviaNPostsHTTPPorSegundo envia durante um segundo n posts
//e busca usuarios aleatórios no banco se indicado na flag usuarioVariavel
func enviaNPostsHTTPPorSegundo(postsPorSegundo int, usuarioVariavel bool) {
	intervalo := time.Second / time.Duration(postsPorSegundo)
	var nome, senha string = "Joel Gaya", "12345"
	if usuarioVariavel {
		nome, senha = buscaRegistroAleatorio()
	}
	requestBody, err := json.Marshal(map[string]string{
		"usuario": nome,
		"senha":   senha,
	})
	if err != nil {
		log.Fatalln(err)
	}
	for i := 0; i < postsPorSegundo; i++ {
		go enviaPostHTTP(bytes.NewBuffer(requestBody))
		time.Sleep(intervalo)
	}

}

// enviaPostHTTP envia um post http e avisa quando recebeu a resposta
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
	usarUsuarioVariavelPtr := flag.Bool("usuarioVariavel", false,
		"Indica se o login deve ser fixo ou variável")
	flag.Parse()
	//Inica as threads  aguarda o término
	numeroThreads := *numeroThreadsPtr
	postsPorSegundo := *postsPorSegundoPtr
	usarUsuarioVariavel := *usarUsuarioVariavelPtr
	totalRequests := numeroThreads * postsPorSegundo

	for indexThread := 0; indexThread < numeroThreads; indexThread++ {
		go enviaNPostsHTTPPorSegundo(postsPorSegundo, usarUsuarioVariavel)
	}
	for indexRequest := 0; indexRequest < totalRequests; indexRequest++ {
		<-done
	}

}

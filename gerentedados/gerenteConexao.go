package gerentedados

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	_ "github.com/go-sql-driver/mysql" //não é rerenciado diretamente, mas precisa ser importado
)

type linhaCache struct {
	senha     string
	timestamp time.Time
}

var poolConexoesDB [maxConexoesBanco]*sql.DB //mantém a conexão do banco
var conexaoAtual int = 0
var conexaoAtualMutex sync.RWMutex
var cache map[string]linhaCache = make(map[string]linhaCache, cachSize)
var cacheMutex sync.RWMutex

//PreAlocaPoolConexoesDB cria uma pool de conexões com o banco de dados para
//evitar que a conexão seja criada no momento onde o request é recebido
func PreAlocaPoolConexoesDB() {
	for i := 0; i < maxConexoesBanco; i++ {
		_, err := getConexao()
		if err != nil {
			os.Exit(1)
		}
	}
}

func getConexao() (*sql.DB, error) {
	var err error
	conexaoAtualMutex.Lock()
	db := poolConexoesDB[conexaoAtual]
	if db == nil {
		fmt.Println("criou nova conexão")
		db, err = sql.Open("mysql", stringdeConexao)
	}
	poolConexoesDB[conexaoAtual] = db
	conexaoAtual = (conexaoAtual + 1) % maxConexoesBanco
	conexaoAtualMutex.Unlock()
	return db, err
}

func buscaSenhaUsuario(nome string) (string, error) {
	cacheMutex.Lock()
	if valor, estaNaCache := cache[nome]; estaNaCache {
		valor.timestamp = time.Now()
		cache[nome] = valor
		cacheMutex.Unlock()
		return valor.senha, nil
	}
	log.Println("Busca no banco")
	cacheMutex.Unlock()
	db, err := getConexao()
	if err != nil {
		return "", err
	}
	var senha string
	db.QueryRow("SELECT senha FROM usuarios where nome = ?", nome).Scan(&senha)
	cacheMutex.Lock()
	if len(cache) >= cachSize {
		var chaveRegistroMaisAntigo string
		var menorTempo time.Time
		for chave, valor := range cache {
			if valor.timestamp.Sub(menorTempo) < 0 {
				menorTempo = valor.timestamp
				chaveRegistroMaisAntigo = chave
			}
		}
		delete(cache, chaveRegistroMaisAntigo)
	}
	cache[nome] = linhaCache{senha, time.Now()}
	cacheMutex.Unlock()
	return senha, nil
}

// ChecaUsuarioSenha verifica se o usuario e a senha enviados estão corretos
func ChecaUsuarioSenha(nome string, senha string) (int, error) {
	var status int
	senhaValida, err := buscaSenhaUsuario(nome)
	if err != nil {
		status = 500
	} else {
		if senhaValida == senha {
			status = 200
		} else {
			status = 401
		}
	}
	return status, err
}

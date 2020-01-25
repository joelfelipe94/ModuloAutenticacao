package gerentedados

import (
	"database/sql"
	"log"
	"os"
	"sync"
	"time"

	_ "github.com/go-sql-driver/mysql" //não é rerenciado diretamente, mas precisa ser importado
)

// linhaCache representa o dado armazenado em uma linha da cache
type linhaCache struct {
	senha     string
	timestamp time.Time //armazena o horário do último acesso ao dado
}

var poolConexoesDB [maxConexoesBanco]*sql.DB //armazena a pool de conexões do banco
var conexaoAtual int = 0                     // aponta para a conexão disponível
var conexaoAtualMutex sync.RWMutex           //garante que acessos à conexão atual não serão concorrentes
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

// getConexao busca uma conexão na pool. Se ela não existir é criada
func getConexao() (*sql.DB, error) {
	var err error
	conexaoAtualMutex.Lock()
	defer conexaoAtualMutex.Unlock()
	db := poolConexoesDB[conexaoAtual]
	if db == nil {
		log.Println("criou nova conexão")
		db, err = sql.Open("mysql", StringdeConexao)
	}
	poolConexoesDB[conexaoAtual] = db
	conexaoAtual = (conexaoAtual + 1) % maxConexoesBanco
	return db, err
}

// deletaRegistorMaisAntigoCache apaga o registro na cache a mais tempo sem ser acessado
func deletaRegistorMaisAntigoCache() {

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

// buscaSenhaUsuarioBanco pega uma conexão na pool e usa para buscar senha
func buscaSenhaUsuarioBanco(nome string) (string, error) {
	db, err := getConexao()
	if err != nil {
		return "", err
	}
	var senha string
	db.QueryRow("SELECT senha FROM usuarios where nome = ?", nome).Scan(&senha)
	return senha, nil
}

// buscaSenhaUsuario verfica se o nome buscado está na cache.
// Caso esteja a senha é retornada. Caso contrário a senha é buscada no banco.
// Se o usuario não estiver no banco a senha retornada é vazia
// Não garante que buscas concorrentes pelo mesmo nome resultem em uma única
// consulta ao banco.
func buscaSenhaUsuario(nome string) (string, error) {

	cacheMutex.Lock()
	// Hit, o valor está na cache
	if valor, estaNaCache := cache[nome]; estaNaCache {
		valor.timestamp = time.Now() // atualiza instante de acesso
		cache[nome] = valor
		cacheMutex.Unlock()
		return valor.senha, nil
	}
	// Miss, busca no banco
	cacheMutex.Unlock()
	senha, err := buscaSenhaUsuarioBanco(nome)
	if err != nil {
		return "", err
	}
	cacheMutex.Lock()
	if len(cache) >= cachSize {
		deletaRegistorMaisAntigoCache()
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

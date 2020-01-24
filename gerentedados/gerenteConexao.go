package gerentedados

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql" //não é rerenciado diretamente, mas precisa ser importado
)

var db *sql.DB //mantém a conexão do banco

func ChecaUsuarioSenha(nome string, senha string) (int, error) {
	var senhaValida string
	if db == nil {
		var err error
		db, err = sql.Open("mysql", "joel:12345@/SistemaAutenticacao")
		if err != nil {
			return 500, err
		}
	}
	db.QueryRow("SELECT senha FROM usuarios where nome = ?", nome).Scan(&senhaValida)
	if senhaValida == senha {
		return 200, nil
	}
	return 401, nil
}

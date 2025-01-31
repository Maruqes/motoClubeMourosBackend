package member

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

// Cria a tabela 'socios' caso não exista
func CreateSociosTable(dbV *sql.DB) error {
	createTableSQL := `
	CREATE TABLE IF NOT EXISTS socios (
		id TEXT PRIMARY KEY,
		numero_socio INTEGER,
		junior INTEGER,
		socio_responsavel TEXT,
		data_nascimento TEXT,
		data_adesao TEXT,
		membro_responsavel TEXT,
		nome TEXT,
		email TEXT,
		telefone TEXT,
		tipo_sangue TEXT,
		rua TEXT,
		numero TEXT,
		concelho TEXT,
		distrito TEXT,
		cod_postal TEXT,
		tipo TEXT,
		data_inscricao TEXT
	);
	`
	_, err := dbV.Exec(createTableSQL)
	if err != nil {
		return err
	}

	db = dbV

	return nil
}

func DoesMemberDataExist(id string) bool {
	rows, err := db.Query("SELECT id FROM socios WHERE id = ?", id)
	if err != nil {
		fmt.Println(err)
		return false
	}

	defer rows.Close()

	if rows.Next() {
		return true
	}

	return false
}

// Função para alterar um Member existente na tabela 'socios'
func ChangeMember(member Member) error {
	updateSQL := `
		UPDATE socios
		SET 
			numero_socio = ?,
			junior = ?,
			socio_responsavel = ?,
			data_nascimento = ?,
			data_adesao = ?,
			membro_responsavel = ?,
			nome = ?,
			email = ?,
			telefone = ?,
			tipo_sangue = ?,
			rua = ?,
			numero = ?,
			concelho = ?,
			distrito = ?,
			cod_postal = ?,
			tipo = ?,
			data_inscricao = ?
		WHERE id = ?;
	`

	_, err := db.Exec(
		updateSQL,
		member.NumeroSocio,
		boolToInt(member.Junior),
		member.SocioResponsavel,
		member.DataNascimento,
		member.DataAdesao,
		member.MembroResponsavel,
		member.Nome,
		member.Email,
		member.Telefone,
		member.TipoSangue,
		member.Rua,
		member.Numero,
		member.Concelho,
		member.Distrito,
		member.CodPostal,
		member.Tipo,
		member.DataInscricao,
		member.ID, // Importante para identificar qual registo vamos atualizar
	)

	if err != nil {
		fmt.Println(err)
	}

	return nil
}

// Função para inserir um Member na tabela 'socios'
func InsertMember(member Member) error {

	if DoesMemberDataExist(member.ID) {
		ChangeMember(member)
		return nil
	}

	insertSQL := `
	INSERT INTO socios (
		id, numero_socio, junior, socio_responsavel, data_nascimento, 
		data_adesao, membro_responsavel, nome, email, telefone, tipo_sangue, 
		rua, numero, concelho, distrito, cod_postal, tipo, 
		data_inscricao
	) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?);
	`

	_, err := db.Exec(
		insertSQL,
		member.ID,
		member.NumeroSocio,
		boolToInt(member.Junior),
		member.SocioResponsavel,
		member.DataNascimento,
		member.DataAdesao,
		member.MembroResponsavel,
		member.Nome,
		member.Email,
		member.Telefone,
		member.TipoSangue,
		member.Rua,
		member.Numero,
		member.Concelho,
		member.Distrito,
		member.CodPostal,
		member.Tipo,
		member.DataInscricao,
	)

	if err != nil {
		fmt.Println(err)
	}

	return nil
}

func GetMemberData(id string) (Member, error) {
	var member Member

	rows, err := db.Query("SELECT * FROM socios WHERE id = ?", id)
	if err != nil {
		fmt.Println(err)
		return member, err
	}

	defer rows.Close()

	if rows.Next() {
		err = rows.Scan(
			&member.ID,
			&member.NumeroSocio,
			&member.Junior,
			&member.SocioResponsavel,
			&member.DataNascimento,
			&member.DataAdesao,
			&member.MembroResponsavel,
			&member.Nome,
			&member.Email,
			&member.Telefone,
			&member.TipoSangue,
			&member.Rua,
			&member.Numero,
			&member.Concelho,
			&member.Distrito,
			&member.CodPostal,
			&member.Tipo,
			&member.DataInscricao,
		)
		if err != nil {
			fmt.Println(err)
			return member, err
		}
	}

	return member, nil
}

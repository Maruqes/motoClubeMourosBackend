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
		id  INTEGER PRIMARY KEY,
		numeroSocio INTEGER,
		tipoMembro TEXT,
		dataNascimento TEXT,
		dataAdesao TEXT,
		membroResponsavel TEXT,
		nome TEXT,
		email TEXT,
		telefone TEXT,
		tipoSangue TEXT,
		rua TEXT,
		numeroPorta INTEGER,
		concelho TEXT,
		distrito TEXT,
		codPostal TEXT
	);
	`
	_, err := dbV.Exec(createTableSQL)
	if err != nil {
		return err
	}

	db = dbV

	return nil
}

func DoesMemberDataExist(id int) bool {
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
			numeroSocio = ?,
			tipoMembro = ?,
			dataNascimento = ?,
			dataAdesao = ?,
			membroResponsavel = ?,
			nome = ?,
			email = ?,
			telefone = ?,
			tipoSangue = ?,
			rua = ?,
			numeroPorta INT= ?,
			concelho = ?,
			distrito = ?,
			codPostal = ?
		WHERE id = ?;
	`

	_, err := db.Exec(
		updateSQL,
		member.numeroSocio,
		member.tipoMembro,
		member.dataNascimento,
		member.dataAdesao,
		member.membroResponsavel,
		member.nome,
		member.email,
		member.telefone,
		member.tipoSangue,
		member.rua,
		member.numeroPorta,
		member.concelho,
		member.distrito,
		member.codPostal,
		member.id, // Importante para identificar qual registo vamos atualizar
	)

	if err != nil {
		fmt.Println(err)
	}

	return nil
}

// Função para inserir um Member na tabela 'socios'
func InsertMember(member Member) error {

	if DoesMemberDataExist(member.id) {
		ChangeMember(member)
		return nil
	}

	insertSQL := `
	INSERT INTO socios (
		id, numeroSocio, tipoMembro, dataNascimento, 
		dataAdesao, membroResponsavel, nome, email, telefone, tipoSangue, 
		rua, numeroPorta, concelho, distrito, codigoPostal
	) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?);
	`

	_, err := db.Exec(
		insertSQL,
		member.id,
		member.numeroSocio,
		member.tipoMembro,
		member.dataNascimento,
		member.dataAdesao,
		member.membroResponsavel,
		member.nome,
		member.email,
		member.telefone,
		member.tipoSangue,
		member.rua,
		member.numeroPorta,
		member.concelho,
		member.distrito,
		member.codPostal,
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
			&member.id,
			&member.numeroSocio,
			&member.tipoMembro,
			&member.dataNascimento,
			&member.dataAdesao,
			&member.membroResponsavel,
			&member.nome,
			&member.email,
			&member.telefone,
			&member.tipoSangue,
			&member.rua,
			&member.numeroPorta,
			&member.concelho,
			&member.distrito,
			&member.codPostal,
		)
		if err != nil {
			fmt.Println(err)
			return member, err
		}
	}

	return member, nil
}

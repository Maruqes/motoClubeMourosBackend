package events

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

// Cria a tabela 'eventos' caso não exista
func CreateEventsTable(dbV *sql.DB) error {
	createTableSQL := `
	CREATE TABLE IF NOT EXISTS eventos (
		id INTEGER PRIMARY KEY,
		nome TEXT,
		descricao TEXT,
		dataEvento TEXT,
		dataLimiteInscricao TEXT,
		numInscricoes INTEGER,
		participantes TEXT,
		cartaz TEXT,
		programa TEXT,
	);
	`
	_, err := dbV.Exec(createTableSQL)
	if err != nil {
		return err
	}

	db = dbV

	return nil
}

func DoesEventDataExist(id int) bool {
	rows, err := db.Query("SELECT id FROM eventos WHERE id = ?", id)
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

// Função para alterar um Event existente na tabela 'eventos'
func ChangeEvent(events Event) error {
	updateSQL := `
		UPDATE eventos
		SET 
			nome = ?,
			descricao = ?,
			dataEvento = ?,
			dataLimiteInscricao = ?,
			numInscricoes = ?,
			participantes = ?,
			cartaz = ?,
			programa = ?,
				WHERE id = ?;
	`

	_, err := db.Exec(
		updateSQL,
		events.nome,
		events.descricao,
		events.dataEvento,
		events.dataLimiteInscricao,
		events.numInscricoes,
		events.participantes,
		events.cartaz,
		events.programa,
		events.id, // Importante para identificar qual registo vamos atualizar
	)

	if err != nil {
		fmt.Println(err)
	}

	return nil
}

// Função para inserir um Event na tabela 'events'
func InsertEvent(events Event) error {

	if DoesEventDataExist(events.id) {
		ChangeEvent(events)
		return nil
	}

	insertSQL := `
	INSERT INTO eventos (
		id, nome, descricao, dataEvento, dataLimiteInscricao, numInscricoes, participantes, cartaz, programa
	) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?);
	`

	_, err := db.Exec(
		insertSQL,
		events.id,
		events.nome,
		events.descricao,
		events.dataEvento,
		events.dataLimiteInscricao,
		events.numInscricoes,
		events.participantes,
		events.cartaz,
		events.programa,
	)

	if err != nil {
		fmt.Println(err)
	}

	return nil
}

func GetEventData(id int) (Event, error) {
	var events Event

	rows, err := db.Query("SELECT * FROM eventos WHERE id = ?", id)
	if err != nil {
		fmt.Println(err)
		return events, err
	}

	defer rows.Close()

	if rows.Next() {
		err = rows.Scan(
			&events.id,
			&events.nome,
			&events.descricao,
			&events.dataEvento,
			&events.dataLimiteInscricao,
			&events.numInscricoes,
			&events.participantes,
			&events.cartaz,
			&events.programa,
		)
		if err != nil {
			fmt.Println(err)
			return events, err
		}
	}

	return events, nil
}

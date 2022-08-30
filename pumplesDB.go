package main

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"os"
)

func CreateDbPumples() {
	db, err := os.Create("pumples.db")
	if err != nil {
		log.Fatal(err)
	}
	db.Close()
	pumpleDb, err := sql.Open("sqlite3", "pumples.db")
	if err != nil {
		log.Fatal(err)
	}
	createtables(pumpleDb)
}

func createtables(db *sql.DB) {
	pumlbe_table := `CREATE TABLE pumples(
	id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT ,
	"Cumple" TEXT,
	"ChatID" INT64 , 
	"User_ID" TEXT );`
	query, err := db.Prepare(pumlbe_table)
	if err != nil {
		log.Fatal(err)
	}
	query.Exec()
	fmt.Println("Cosa creada")
	pumlbe_table = `CREATE TABLE saludos(
	id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT ,
	"Saludo" TEXT,
	"ChatID" INT64 , 
	"User_ID" TEXT,
	"Enviado" INTEGER );`
	query, err = db.Prepare(pumlbe_table)
	if err != nil {
		log.Fatal(err)
	}
	query.Exec()
	fmt.Println("Cosa 2 creada ")
}

func insertPumples(db *sql.DB, cumple2 cumple) (int, error) {
	res, err := db.Exec("INSERT INTO pumples VALUES (NULL,?,?,?);", cumple2.cumple, cumple2.ChatID, cumple2.user)
	if err != nil {
		fmt.Println(err)
		return 0, err
	}
	res.LastInsertId()
	return 0, err
}

func insertSaludos(db *sql.DB, saludo2 saludo) (int, error) {
	res, err := db.Exec("INSERT INTO saludos VALUES (NULL,?,?,?,?;", saludo2.saludo, saludo2.ChatID, saludo2.User_ID, 0)
	if err != nil {
		return 0, err
	}
	res.LastInsertId()
	return 0, err
}

func searchPumpleByID(db *sql.DB, ID int64) (cumple, error) {
	row := db.QueryRow("SELECT Cumple,ChatID,User_ID FROM pumples WHERE ChatID=?", ID)
	rowData := cumple{}
	var err error
	if err = row.Scan(&rowData.cumple, &rowData.ChatID, &rowData.user); err == sql.ErrNoRows {
		log.Printf("No encontre el ID")
		return rowData, err
	}
	return rowData, err
}

func searchPumpleByUser(db *sql.DB, user string) (cumple, error) {
	row := db.QueryRow("SELECT Cumple,ChatID,User_ID FROM pumples WHERE User_ID=?", user)
	rowData := cumple{}
	var err error
	if err = row.Scan(&rowData.cumple, &rowData.ChatID, &rowData.user); err == sql.ErrNoRows {
		log.Printf("No encontre el ID")
		return rowData, err
	}
	return rowData, err
}

func searchSaludos(db *sql.DB, username string) ([]saludo, error) {
	rows, err := db.Query("SELECT Saludo,ChatID,User_ID FROM saludos WHERE User_ID=?", username)
	if err != nil {
		return nil, err
	}

	var saludos []saludo
	for rows.Next() {
		i := saludo{}
		err = rows.Scan(&i.saludo, &i.ChatID, &i.User_ID)
		if err != nil {
			return nil, err
		}
		saludos = append(saludos, i)
	}
	return saludos, nil
}

func writeSaludos(db *sql.DB, saludo string, username string) (int, error) {
	res, err := db.Exec("UPDATE saludos SET saludo=? WHERE User_ID=?", saludo, username)
	if err != nil {
		return 0, err
	}
	res.LastInsertId()
	return 0, err

}

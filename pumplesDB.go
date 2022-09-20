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
	"Greatings" TEXT,
	"Receiver" TEXT,
	"Receiver_User_ID" INT64,
	"Sender" TEXT , 
	"Sender_User_ID" INT64);`
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
	res, err := db.Exec("INSERT INTO saludos VALUES (NULL,?,?,?,?,?);", saludo2.saludo, saludo2.Receiver, saludo2.Receiver_User_ID, saludo2.Sender, saludo2.Sender_User_ID)
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

func alreadyAdded(db *sql.DB, receiver_ID int64, sender_id int64) (string, error) {
	row := db.QueryRow("SELECT Greatings FROM saludos WHERE Receiver=? AND Sender=?", receiver_ID, sender_id)
	var greating string
	var err error
	if err = row.Scan(&greating); err == sql.ErrNoRows {
		log.Printf("No encontre el ID")
		return "", err
	}
	return greating, err
}

func searchGreetings(db *sql.DB, Receiver_User_ID int64) ([]saludo, error) {
	rows, err := db.Query("SELECT * FROM saludos WHERE Sender_User_ID=?", Receiver_User_ID)
	if err != nil {
		return nil, err
	}

	var saludos []saludo
	for rows.Next() {
		i := saludo{}
		err = rows.Scan(&i.saludo, &i.Receiver, &i.Receiver_User_ID, &i.Sender, &i.Sender_User_ID)
		if err != nil {
			return nil, err
		}
		saludos = append(saludos, i)
	}
	return saludos, nil
}
func searchLoadedGreetings(db *sql.DB, sender_User_ID int64) ([]saludo, error) {
	rows, err := db.Query("SELECT * FROM saludos WHERE Sender_User_ID=?", sender_User_ID)
	if err != nil {
		return nil, err
	}

	var saludos []saludo
	for rows.Next() {
		i := saludo{}
		var n int
		err = rows.Scan(&n, &i.saludo, &i.Receiver, &i.Receiver_User_ID, &i.Sender, &i.Sender_User_ID)
		if err != nil {
			return nil, err
		}
		saludos = append(saludos, i)
	}
	return saludos, nil
}

func updateGreeting(db *sql.DB, saludo string, username string, chatID int64) (int, error) {
	res, err := db.Exec("UPDATE saludos SET Greatings=? WHERE Receiver=? AND Sender_User_ID=?", saludo, username, chatID)
	if err != nil {
		return 0, err
	}
	res.LastInsertId()
	return 0, err

}

func deleteGreetingDB(db *sql.DB, receiver string, sender_id int64) (int, error) {
	res, err := db.Exec("DELETE FROM saludos WHERE Receiver=? AND Sender_User_ID=?;", receiver, sender_id)
	if err != nil {
		return 0, err
	}
	res.LastInsertId()
	return 0, err
}

package main

import (
	"fmt"

	"./app"
	"./db"
	"log"
	"./models"
)

func main() {

	// load application configurations
	err := app.LoadConfig("./")
	if err != nil {
		panic(fmt.Errorf("Invalid application configuration: %s", err))
	}

	conn := db.OpenConnection()
	defer conn.Close()

	if err != nil {
		log.Fatal("Error: Could not establish a connection with the database")
	}
	fmt.Println("Successfully connected!")


	rows, err := conn.Query("SELECT id, nome FROM public.leis")
	if err != nil {
		log.Fatal(err)
	}

	var leis []models.Lei
	for rows.Next() {
		var lei models.Lei
		rows.Scan(&lei.Id, &lei.Nome)
		fmt.Println(lei)
		leis = append(leis, lei)
	}

	fmt.Println(leis)

	for _, lei := range leis {
		fmt.Println(lei.Nome)
	}


}

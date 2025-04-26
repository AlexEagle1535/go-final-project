package main

import (
	"fmt"
	"log"
	"os"

	"github.com/AlexEagle1535/go-final-project/pkg/db"
	"github.com/AlexEagle1535/go-final-project/pkg/server"
	"github.com/joho/godotenv"
)

func init() {
	if err := godotenv.Load(); err != nil {
		fmt.Print("No .env file found")
	}
}

func main() {
	DBFile := os.Getenv("TODO_DBFILE")
	if len(DBFile) == 0 {
		DBFile = "scheduler.db"
		os.Setenv("TODO_DBFILE", DBFile)
		fmt.Printf("Файл базы данных не найден, установлен файл по умолчанию: %s\n", DBFile)
	}
	err := db.Init(DBFile)
	if err != nil {
		log.Fatalf("ошибка инициализации базы данных: %v", err)
	}
	fmt.Println("Starting server...")
	server.Run()
}

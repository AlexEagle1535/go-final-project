package server

import (
	"fmt"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

const webDir = "/Users/alexeagle1535/Documents/github.com/AlexEagle1535/go_final_project-main/web"

func init() {
	if err := godotenv.Load(); err != nil {
		fmt.Print("No .env file found")
	}
}

func Run() {
	port, exists := os.LookupEnv("TODO_PORT")
	if !exists {
		port = "7450"
		os.Setenv("TODO_PORT", port)
		fmt.Printf("Порт не найден, установлен порт по умолчанию: %s\n", port)
	}
	http.Handle("/", http.FileServer(http.Dir(webDir)))
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		fmt.Printf("ошибка запуска сервера: %s\n", err.Error())
		return
	}
}

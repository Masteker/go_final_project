package main

import (
	"database/sql"
	"fmt"
	"github.com/Masteker/go_final_project/internal/checkers"
	"github.com/Masteker/go_final_project/internal/handlers"
	"github.com/Masteker/go_final_project/internal/modelsAndService"
	"github.com/Masteker/go_final_project/internal/modelsAndService/store"
	mid "github.com/Masteker/go_final_project/middleware"
	"log"
	"net/http"

	"go.uber.org/zap"
	_ "modernc.org/sqlite"
)

// main является точкой входа в приложение. Она инициализирует сервер, настраивает маршрутизацию,
// и запускает прослушивание входящих подключений.
func main() {
	logger := zap.NewExample() // or NewProduction, or NewDevelopment
	defer logger.Sync()
	port := checkers.CheckPort() // Функция для проверки и возврата номера порта
	url := fmt.Sprintf("localhost:%s", port)
	sugar := logger.Sugar()
	webDir := "./web" // Каталог, содержащий статические файлы для обслуживания

	path, install := checkers.CheckDB() // Функция для проверки и возврата пути к базе данных и флага установки

	// Открываем подключение к базе данных
	db, err := sql.Open("sqlite", path)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	store := store.NewTaskStore(db)
	if install {
		store.InitDB()
	}
	service := service.NewTaskService(store, sugar)
	handler := handlers.NewHandler(service, sugar)

	// Создаем новый экземпляр http.Server с указанным портом
	server := &http.Server{
		Addr: ":" + port, // Порт, на котором сервер будет прослушивать
	}

	// Настраиваем маршрутизацию для обслуживания всех файлов в каталоге web и конечных точек API
	http.Handle("/", http.FileServer(http.Dir(webDir)))
	http.HandleFunc("/api/sign", handler.Authentication)
	http.HandleFunc("/api/nextdate", handler.NextDate)
	http.HandleFunc("/api/task", mid.Auth(handler.Task))
	http.HandleFunc("/api/tasks", mid.Auth(handler.GetAllTasks))
	http.HandleFunc("/api/task/done", mid.Auth(handler.TaskDone))

	// Запускаем сервер и прослушиваем входящие подключения
	sugar.Infof("Server started at %s", url)
	err = server.ListenAndServe()
	if err != nil {
		sugar.Panic(err)
	}
}
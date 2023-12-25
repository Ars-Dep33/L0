package main

import (
	"context"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"github.com/sirupsen/logrus"
	"log"
	"net"
	"net/http"
	"time"
)

// Кеш, для хранения заказов
var cache = DataCache{
	cache: make(map[string][]byte),
}

// Основная функция
func main() {
	cfg := GetConfig()
	router := httprouter.New()

	// Подключение к базе данных
	conn, err := connectToDB(cfg.Storage)
	if err != nil {
		log.Fatalf("Ошибка подключения к базе данных: %v", err)
	}
	defer conn.Close(context.Background())
	// Восстановление данных из базы в Кеш
	err = cache.RestoreCacheFromDB(conn)
	if err != nil {
		logrus.Fatalf("Ошибка работы функции restoreCacheFromDB: %v", err)
	}
	// Создаем экземпляр Роутера и регистриуем его
	handler := NewHandler(&cache)
	handler.Register(router)

	//Запускаем HTTP-сервер
	start(router, cfg)

	// Запускаем функцию подписки на Nats
	subscribeNATS(conn, &cache)

	//Вызываем функцию для отправки собщений с заказами в Nats
	go SendNewOrder()

}

// Функция запуска HTTP-сервера
func start(router *httprouter.Router, cfg *Config) {
	logrus.Info("start application")
	//Инициализируем listener для прослушивания порта
	var listener net.Listener
	var listenErr error
	if router == nil {
		logrus.Error("Router is nil")
	}
	if cfg == nil {
		logrus.Error("Config is nil")
	}
	//Передаем параметры из Конфига
	listener, listenErr = net.Listen("tcp", fmt.Sprintf("%s:%s", cfg.Listen.BindIP, cfg.Listen.Port))
	logrus.Infof("Сервер слушает порт port %s:%s", cfg.Listen.BindIP, cfg.Listen.Port)

	if listenErr != nil {
		logrus.Error(listenErr)
	}
	//Создаем экземпляр сервера
	server := &http.Server{
		Handler:      router,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	// запускаем Http сервер, исспользуем корутину, что бы не блокировалась работа программы
	go func() {
		if err := server.Serve(listener); err != nil {
			logrus.Fatalf("HTTP server error: %v", err)
		}
	}()
}

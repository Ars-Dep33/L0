package main

import (
	"encoding/json"
	"github.com/jackc/pgx/v4"
	"github.com/nats-io/stan.go"
	"github.com/sirupsen/logrus"
	"log"
)

// Запуск docker контейнера с сервером Nats-Streaming
// sudo docker run -p 4222:4222 -p 8222:8222 -p 6222:6222 --name nats-server -ti nats-streaming

// Константы для подключения к Nats
const (
	clusterID   = "test-cluster"
	clientID    = "your_client_id_33"
	subject     = "your_subject"
	durableName = "durable-order-sub"
)

// Описываем функцию по подключению и подписки к Nats:
func subscribeNATS(conn *pgx.Conn, cache *DataCache) {
	sc, err := stan.Connect(clusterID, clientID) // Коннектимся
	if err != nil {
		logrus.Fatalf("При подключении к NATS произошла ошибка: %v", err)
	}
	defer sc.Close()
	logrus.Info("Успешное подключение к Nats Streaming")
	logrus.Info("")
	// Выполняем подписку с заданными параметрами
	sub, err := sc.Subscribe(subject, func(msg *stan.Msg) {

		var data SenderOrder
		// десериализуем данные
		if err := json.Unmarshal(msg.Data, &data); err != nil {
			log.Printf("Ошибка десериализации данных заказа: %v", err)
			return
		}
		logrus.Info("Десереализация данных успешна")

		// Проверка валидности JSON
		if !json.Valid(msg.Data) {
			log.Printf("Получен невалидный JSON: %s", msg.Data)
			return
		}
		orderUid := data.OrderUID
		newOrder := NewOrder(orderUid, msg.Data)

		// Сохраняем данные в кеше
		cache.AddOrder(newOrder)
		// Сохраняем данные в базе
		saveOrder(conn, newOrder)
	}, stan.DurableName(durableName))

	if err != nil {
		log.Fatalf("Ошибка установки подписки на NATS: %v", err)
	}
	defer sub.Unsubscribe()
	//цикл для поддержани яподписки
	select {}
}

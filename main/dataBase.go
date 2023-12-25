package main

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4"
	"github.com/sirupsen/logrus"
	"sync"
)

// Мьютекс для обеспечения безопасности работы с базой данных
var dbMutex sync.Mutex

// Функция подключения к базе данных
func connectToDB(cfg StorageConfig) (*pgx.Conn, error) {
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", cfg.UserName, cfg.Password, cfg.Host, cfg.Port, cfg.DataBase)
	var err error
	conn, err := pgx.Connect(context.Background(), connStr)
	if err != nil {
		logrus.Fatalf("Ошибка cоединения с базой: %s", err)
	}
	logrus.Info("Соединение с базой выплнено.")
	return conn, nil
}

// Функция сохранения данных заказа в базу данных
func saveOrder(conn *pgx.Conn, order *Order) {
	// Заблокировать мьютекс перед началом транзакции
	dbMutex.Lock()
	defer dbMutex.Unlock()
	logrus.Info("Сохранение данных в PostgresSQL")
	// Старт транцакции
	tx, err := conn.Begin(context.Background())
	if err != nil {
		logrus.Printf("Ошибка начала транзакции: %v", err)
		return
	}
	logrus.Info("Начало транзацкии...")

	// Вставляем данные в таблицу orders
	_, err = tx.Exec(context.Background(), "INSERT INTO orders VALUES ($1, $2)",
		order.OrderID, order.Data)
	if err != nil {
		logrus.Printf("Ошибка сохранения заказа в PostgresSQL: %v", err)
		tx.Rollback(context.Background()) // Откатить транзакцию при ошибке
		return
	}

	// Подтверждаем, если все успешно
	err = tx.Commit(context.Background())
	if err != nil {
		logrus.Printf("Ошибка подтверждения транзакции: %v", err)
		tx.Rollback(context.Background()) // Откатить транзакцию при ошибке
		return
	}
	logrus.Info("Транзакция прошла успешно! Данные сохранены в PostgresSQL!")
	logrus.Info("")
}

package main

import (
	"context"
	"github.com/jackc/pgx/v4"
	"github.com/sirupsen/logrus"
	"sync"
)

// Структура Кеша

type DataCache struct {
	mu    sync.RWMutex
	cache map[string][]byte
}

// Функция восстановления данных из баны в Кеш

func (oc *DataCache) RestoreCacheFromDB(conn *pgx.Conn) error {
	rows, err := conn.Query(context.Background(), "SELECT order_uid, data FROM orders") //делаем  запрос в базу
	if err != nil {
		logrus.Fatalf("Ошибка получения данных: %v", err)
	}
	defer rows.Close()
	for rows.Next() { // Сканируем строки и передаем значения в поля структуры Order
		var order Order
		err := rows.Scan(&order.OrderID, &order.Data)

		if err != nil {
			logrus.Fatalf("Ошибка сканирования строки: %v", err)
			continue
		}
		// Блочим Мьютекс и добавляем данные из базы в Кеш
		oc.mu.Lock()
		oc.cache[order.OrderID] = order.Data
		oc.mu.Unlock()
	}
	logrus.Info("Восстановение данных из базы в Cache выполнено.")
	return nil
}

// Добавляем новые заказы в кеш

func (oc *DataCache) AddOrder(order *Order) {
	oc.mu.Lock()
	oc.cache[order.OrderID] = order.Data
	oc.mu.Unlock()
	logrus.Infof("Данные сохранены в Кеш! OrderUID: %s", order.OrderID)
}

// возвращает заказ из кеша по его идентификатору

func (oc *DataCache) GetOrder(orderID string) ([]byte, bool) {
	oc.mu.RLock()
	order, found := oc.cache[orderID]
	oc.mu.RUnlock()
	return order, found
}

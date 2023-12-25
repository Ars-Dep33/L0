package main

import (
	"encoding/json"
	"github.com/google/uuid"
	"github.com/nats-io/stan.go"
	"github.com/sirupsen/logrus"
	"log"
	"time"
)

//Структура отправляемых данных

type SenderOrder struct {
	OrderUID    string `json:"order_uid"`
	TrackNumber string `json:"track_number"`
	Entry       string `json:"entry"`
	Delivery    struct {
		Name    string `json:"name"`
		Phone   string `json:"phone"`
		Zip     string `json:"zip"`
		City    string `json:"city"`
		Address string `json:"address"`
		Region  string `json:"region"`
		Email   string `json:"email"`
	} `json:"delivery"`
	Payment struct {
		Transaction  string  `json:"transaction"`
		RequestID    string  `json:"request_id"`
		Currency     string  `json:"currency"`
		Provider     string  `json:"provider"`
		Amount       float64 `json:"amount"`
		PaymentDT    int64   `json:"payment_dt"`
		Bank         string  `json:"bank"`
		DeliveryCost float64 `json:"delivery_cost"`
		GoodsTotal   float64 `json:"goods_total"`
		CustomFee    float64 `json:"custom_fee"`
	} `json:"payment"`
	Items struct {
		ChrtID      int     `json:"chrt_id"`
		TrackNumber string  `json:"track_number"`
		Price       float64 `json:"price"`
		RID         string  `json:"rid"`
		Name        string  `json:"name"`
		Sale        int     `json:"sale"`
		Size        string  `json:"size"`
		TotalPrice  float64 `json:"total_price"`
		NmID        int     `json:"nm_id"`
		Brand       string  `json:"brand"`
		Status      int     `json:"status"`
	} `json:"items"`
	Locale            string `json:"locale"`
	InternalSignature string `json:"internal_signature"`
	CustomerID        string `json:"customer_id"`
	DeliveryService   string `json:"delivery_service"`
	ShardKey          string `json:"shardkey"`
	SmID              int    `json:"sm_id"`
	DateCreated       string `json:"date_created"`
	OOFShard          string `json:"oof_shard"`
}

//Функция отправки данных

func SendNewOrder() {
	//Коннект к Nats
	sc, err := stan.Connect(clusterID, "client_id")
	if err != nil {
		log.Fatal(err)
	}
	defer sc.Close()

	//устанавливаем таймер на отправку сообщений в цикле
	timer := time.NewTicker(10 * time.Second)
	for {
		select {
		case <-timer.C:
			orderData := SenderOrder{
				OrderUID:    uuid.New().String(),
				TrackNumber: "TEST",
				Entry:       "TEST",
				Delivery: struct {
					Name    string `json:"name"`
					Phone   string `json:"phone"`
					Zip     string `json:"zip"`
					City    string `json:"city"`
					Address string `json:"address"`
					Region  string `json:"region"`
					Email   string `json:"email"`
				}{
					Name:    "test_name",
					Phone:   "8-800-555-35-35",
					Zip:     "Проще",
					City:    "Позвонить",
					Address: "Чем у кого-то",
					Region:  "Занимать=)",
					Email:   "test@test.com",
				},
				Payment: struct {
					Transaction  string  `json:"transaction"`
					RequestID    string  `json:"request_id"`
					Currency     string  `json:"currency"`
					Provider     string  `json:"provider"`
					Amount       float64 `json:"amount"`
					PaymentDT    int64   `json:"payment_dt"`
					Bank         string  `json:"bank"`
					DeliveryCost float64 `json:"delivery_cost"`
					GoodsTotal   float64 `json:"goods_total"`
					CustomFee    float64 `json:"custom_fee"`
				}{
					Transaction:  "test_transaction",
					RequestID:    "test_request_id",
					Currency:     "USD",
					Provider:     "test",
					Amount:       1564.5,
					PaymentDT:    time.Now().Unix(),
					Bank:         "test",
					DeliveryCost: 1500,
					GoodsTotal:   5654,
					CustomFee:    0,
				},
				Items: struct {
					ChrtID      int     `json:"chrt_id"`
					TrackNumber string  `json:"track_number"`
					Price       float64 `json:"price"`
					RID         string  `json:"rid"`
					Name        string  `json:"name"`
					Sale        int     `json:"sale"`
					Size        string  `json:"size"`
					TotalPrice  float64 `json:"total_price"`
					NmID        int     `json:"nm_id"`
					Brand       string  `json:"brand"`
					Status      int     `json:"status"`
				}{
					ChrtID:      48946546,
					TrackNumber: "test_track_number",
					Price:       5946,
					RID:         "item_rid",
					Name:        "test_name",
					Sale:        20,
					Size:        "S",
					TotalPrice:  5946,
					NmID:        359465,
					Brand:       "test Brand",
					Status:      202,
				},
				Locale:            "en",
				InternalSignature: "test_signature",
				CustomerID:        "test_id",
				DeliveryService:   "test Service",
				ShardKey:          "test_shard_key",
				SmID:              1,
				DateCreated:       time.Now().Format(time.RFC3339),
				OOFShard:          "test",
			}
			// Сериализуем данные и отправляем
			orderJSON, err := json.Marshal(orderData)
			if err != nil {
				log.Println("Error marshalling order data:", err)
			}
			err = sc.Publish(subject, orderJSON)
			if err != nil {
				log.Println("Error publishing order:", err)
			}
			logrus.Println("Order published successfully. OrderUID:", orderData.OrderUID)
		}
	}
}

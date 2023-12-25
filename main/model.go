package main

type Order struct {
	OrderID string `json:"order_uid"`
	Data    []byte `json:"data"`
}

func NewOrder(orderID string, data []byte) *Order {
	return &Order{
		OrderID: orderID,
		Data:    data,
	}
}

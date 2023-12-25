package main

import (
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"sync"
)

const (
	orderIDUrl = "/:order_uid"
)

type Handler struct {
	cacheMutex sync.RWMutex
	cache      *DataCache
}

// прописываем конструктор New для нашей структуры Handler

func NewHandler(cache *DataCache) *Handler {
	return &Handler{cache: cache}
}

// Регистрируем наш роутер

func (h *Handler) Register(router *httprouter.Router) {
	router.GET(orderIDUrl, h.GetOrderHandler)
}

//Возвращаем данные по order_uid

func (h *Handler) GetOrderHandler(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	// Считываем order_uid из Url-a
	orderID := params.ByName("order_uid")

	// Заблокировать мьютекс для чтения из кеша
	h.cacheMutex.RLock()
	order, found := h.cache.GetOrder(orderID)
	h.cacheMutex.RUnlock()

	if !found {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Заказ не найден"))
		return
	}

	// Отправляем закодированные данные

	json.NewEncoder(w).Encode(order)
}

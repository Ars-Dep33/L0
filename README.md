# Тестовый сервис, отображающий данные о заказе.

- В сервисе реализовано подключение и подписка к nats-streaming.
- Данные, приходящие из nats-streaming записываются в БД и в Кэш.
- В случае падения сервиса, при последующем старте данные восстанавливаются в кэш из БД.
- Так же поднимается Http-сервер для выдачи данных.
- Модель данных в формате JSON.

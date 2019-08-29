package main

import (
	"./consumer"
	"./producer"
	"./shared"
	"./simulator"
	"./storage"
	"./webServer"
	"sync"
)

func main() {
	//настройки из config.json
	cfg := shared.AppConfig();

	//Отправим в очередь все события (Google Hit) пришедшие от Веб-Сервера)
	producer := producer.NewProducer(cfg)

	//Веб-сервер для приема сообщений в формате Google Analytics
	go webServer.NewWebServer(cfg, producer)

	//Хранлищие для хранения метрик
	db := storage.NewRedisClient(cfg.Redis)

	//Временный буффер (на 5 секунд) в памяти программы
	buffer := storage.NewBuffer(db, cfg.Setup);

	//потребители входящих Google Hit'ов
	for i := 1; i <= cfg.Setup.NumberOfConsumers; i++ {
		consumer := consumer.Consumer{Config: cfg, Db: buffer, No: i}
		go consumer.Subscribe()
	}

	//Эмулятор запросов
	go simulator.NewSiumulator(cfg.Client, buffer, db)

	//предотвратить закрытие программы
	notGracefulExit()
}

func notGracefulExit() {
	var wg = &sync.WaitGroup{}
	wg.Add(1)
	wg.Wait()
}

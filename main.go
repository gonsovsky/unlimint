package main

import (
	"./consumer"
	"./producer"
	"./shared"
	"./simulator"
	"./storage"
	"./webServer"
	"context"
	"fmt"
	"os"
	"os/signal"
	"time"
)

var web *webServer.WebServer

func main() {
	//настройки из config.json
	cfg := shared.AppConfig()

	//Отправим в очередь все события (Google Hit) пришедшие от Веб-Сервера)
	producer := producer.NewProducer(cfg)

	//Веб-сервер для приема сообщений в формате Google Analytics
	web = webServer.NewWebServer(cfg, producer)
	go web.Start()

	//Хранлищие для хранения метрик
	db := storage.NewRedisClient(cfg.Redis)

	//Временный буффер (на 5 секунд) в памяти программы
	buffer := storage.NewBuffer(db, cfg.Setup)

	//потребители входящих Google Hit'ов
	for i := 1; i <= cfg.Setup.NumberOfConsumers; i++ {
		consumer := consumer.Consumer{Config: cfg, Db: buffer, No: i}
		go consumer.Subscribe()
	}

	//Эмулятор запросов
	go simulator.NewSiumulator(cfg.Client, buffer, db)

	//закрытие программы
	shutdown()
}

func shutdown() {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	<-stop
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := web.Server.Shutdown(ctx); err != nil {
		panic(err)
	}
	fmt.Println("Good buy.")
}

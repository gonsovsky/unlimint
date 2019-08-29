Проверочная работа Unlimint.

RUN: go run main.go

Схема работы

    - Запускается веб-служба принимающая запросы ф формате Google Analytics
        webServer.go
    
    - Пришедший GoogleHit помещается в очередь AMQP
        producer.go
    
    - Один из N потребителей получает GoogleHit
        consumer.go
        
    - GoogleHit размещается во временном буффере на N секунд
        buffer.go
    
    - Раз в N секунд GoogleHit'ы сохраняются в хранилище
        redis.go
        
    - Эмуляцией запросов и нагрузкой занимается 
        simulator.go
        
    - Настройки в config.json     

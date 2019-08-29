package webServer

import (
	"../producer"
	"../shared"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

//Веб-сервер для приема Google Hit'ов
type WebServer struct {
	config *shared.Cfg
	producer *producer.Producer
	router *mux.Router
}

//Создать новый Веб-Сервер
func NewWebServer(config *shared.Cfg, producer *producer.Producer) *WebServer {
	web := WebServer{config: config, producer: producer}
	web.router = mux.NewRouter().StrictSlash(true)
	web.router.HandleFunc("/", web.index)
	log.Fatal(http.ListenAndServe(config.Web.AddrAndPort(), web.router))
	return &web
}

//Обслужить запросы к ВебСеревру
func (web *WebServer) index(w http.ResponseWriter, r *http.Request) {
	var hit shared.GoogleHit
	hit.FromHTMLForm(r)
	web.producer.Publish(hit)
	w.Header().Set("content-type", "application/json")
	w.Write(hit.ToJSON())
}

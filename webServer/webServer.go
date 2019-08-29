package webServer

import (
	"../producer"
	"../shared"
	"github.com/gorilla/mux"
	"net/http"
)

//Веб-сервер для приема Google Hit'ов
type WebServer struct {
	config   *shared.Cfg
	producer *producer.Producer
	router   *mux.Router
	Server   *http.Server
}

//Создать новый Веб-Сервер
func NewWebServer(config *shared.Cfg, producer *producer.Producer) *WebServer {
	web := WebServer{config: config, producer: producer}
	web.router = mux.NewRouter().StrictSlash(true)
	web.router.HandleFunc("/", web.index)
	web.Server = &http.Server{Addr: config.Web.AddrAndPort(), Handler: web.router}
	return &web
}

func (web *WebServer) Start() {
	if err := web.Server.ListenAndServe(); err != nil {
		panic(err)
	}
}

//Обслужить запросы к ВебСеревру
func (web *WebServer) index(w http.ResponseWriter, r *http.Request) {
	var hit shared.GoogleHit
	hit.FromHTMLForm(r)
	web.producer.Publish(hit)
	w.Header().Set("content-type", "application/json")
	w.Write(hit.ToJSON())
}

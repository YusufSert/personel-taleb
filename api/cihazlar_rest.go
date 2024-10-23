package api

import (
	"encoding/json"
	"log"
	"log/slog"
	"net"
	"net/http"
	"websocketTest/service"
)

type CihazlarRest struct {
	s *service.CihazService
}

func NewCihazlarRest() (*CihazlarRest, error) {
	s, err := service.NewCihazService()
	if err != nil {
		return nil, err
	}
	return &CihazlarRest{s}, err
}

func (cRest CihazlarRest) Run() error {
	net.Interfaces()
	http.HandleFunc("GET /indirim-taleb", cRest.GET)
	http.HandleFunc("POST /indirim-taleb", cRest.POST)
	http.HandleFunc("GET /admin-ekran", cRest.GETAdminEkran)
	slog.Info("server start running on localhost:8080")
	err := http.ListenAndServe("localhost:8080", nil)
	return err
}

func (cRest CihazlarRest) GET(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	cihazlarJTaleb, err := cRest.s.GetIndirimTalebEkraniData()
	if err != nil {
		log.Println(err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	err = json.NewEncoder(w).Encode(cihazlarJTaleb)
	if err != nil {
		log.Println(err)
	}
}

func (cRest CihazlarRest) POST(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	body := struct {
		CihazId int    `json:"cihaz_id"`
		User    string `json:"user"`
	}{}

	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		log.Println(err)
		return
	}

	err = cRest.s.AddTaleb(r.Context(), body.User, body.CihazId)
	if err != nil {
		slog.Error("yeni taleb eklenemedi", "err", err)
	}
}

func (cRest CihazlarRest) GETAdminEkran(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	recs, err := cRest.s.GetAdminEkran(r.Context())
	if err != nil {
		log.Println(err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	err = json.NewEncoder(w).Encode(recs)
	if err != nil {
		log.Println(err)
	}
}

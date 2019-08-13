package application

import (
	"database/sql"
	"encoding/json"
	"github.com/go-chi/chi"
	"github.com/seramirezdev/truora-test/api/config"
	"github.com/seramirezdev/truora-test/api/entities"
	"github.com/seramirezdev/truora-test/api/models"
	"log"
	"net/http"
)

type App struct {
	chi.Router
	*sql.DB
}

func (app *App) Inicialize(user, dbname string) {
	var err error
	app.DB, err = config.GetDB(user, dbname)

	if err != nil {
		log.Fatal(err)
	}

	app.Router = chi.NewRouter()

	app.Router.Get("/consult-domain/{domain}", app.consultDomain)
	app.Router.Get("/domains", app.getDomains)
}

func (app *App) Run(port string) {
	log.Println("Start service")
	log.Fatal(http.ListenAndServe(":"+port, app.Router))
}

func (app *App) consultDomain(w http.ResponseWriter, r *http.Request) {
	domain := chi.URLParam(r, "domain")
	model := models.DomainModel{DB: app.DB}

	infoDomain, err := model.ConsultDomain(domain)

	if err != nil {
		respondJsonWithError(w, http.StatusNotFound, err.Error())
		return
	}

	respondJson(w, http.StatusOK, infoDomain)

}

func (app *App) getDomains(w http.ResponseWriter, r *http.Request) {
	model := models.DomainModel{DB: app.DB}

	domains, err := model.GetDomains()

	if err != nil {
		respondJsonWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	response := map[string][]entities.Domain{}
	response["items"] = domains

	respondJson(w, http.StatusOK, response)
}

func respondJsonWithError(w http.ResponseWriter, code int, message string) {
	respondJson(w, code, map[string]string{"error": message})
}

func respondJson(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(code)
	w.Write(response)
}

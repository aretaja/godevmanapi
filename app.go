package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/aretaja/godevmandb"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v4/pgxpool"
)

type App struct {
	Router *mux.Router
	DB     *pgxpool.Pool
	Ctx    context.Context
}

func (a *App) Initialize(dbURL string) {
	a.Ctx = context.Background()

	pool, err := pgxpool.Connect(a.Ctx, dbURL)
	if err != nil {
		log.Fatal(err)
	}

	a.DB = pool
	a.Router = mux.NewRouter()

	a.initializeRoutes()
}

func (a *App) Run(addr string) {
	log.Fatal(http.ListenAndServe("localhost:48888", a.Router))
}

// Handlers - Basic CRUD
func (a *App) countConProviders(w http.ResponseWriter, r *http.Request) {

	q := godevmandb.New(a.DB)
	res, err := q.CountConProviders(a.Ctx)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]int64{"count": res})
}

func (a *App) getConProviders(w http.ResponseWriter, r *http.Request) {
	var p = godevmandb.GetConProvidersParams{
		Limit:  100,
		Offset: 0,
	}

	l, err := strconv.ParseInt(r.FormValue("count"), 10, 32)
	if err != nil {
		log.Print("Invalid count value. Using default")
	} else {
		if l < 100 || l > 0 {
			p.Limit = int32(l)
		}
	}
	o, err := strconv.ParseInt(r.FormValue("start"), 10, 32)
	if err != nil {
		log.Print("Invalid start value. Using default")
	} else {
		if o > 0 {
			p.Offset = int32(o)
		}
	}

	q := godevmandb.New(a.DB)
	res, err := q.GetConProviders(a.Ctx, p)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, res)
}

func (a *App) getConProvider(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["con_prov_id"], 10, 64)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid provider ID")
		return
	}

	q := godevmandb.New(a.DB)
	res, err := q.GetConProvider(a.Ctx, id)
	if err != nil {
		if err.Error() == "No rows in result set" {
			respondWithError(w, http.StatusNotFound, "Provider not found")
		} else {
			respondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	respondWithJSON(w, http.StatusOK, res)
}

func (a *App) createConProvider(w http.ResponseWriter, r *http.Request) {
	var p godevmandb.CreateConProviderParams
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&p); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	q := godevmandb.New(a.DB)
	res, err := q.CreateConProvider(a.Ctx, p)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusCreated, res)
}

func (a *App) updateConProvider(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["con_prov_id"], 10, 64)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid provider ID")
		return
	}

	var p godevmandb.UpdateConProviderParams
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&p); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()
	p.ConProvID = id

	q := godevmandb.New(a.DB)
	res, err := q.UpdateConProvider(a.Ctx, p)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, res)
}

func (a *App) deleteConProvider(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["con_prov_id"], 10, 64)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid provider ID")
		return
	}

	q := godevmandb.New(a.DB)
	err = q.DeleteConProvider(a.Ctx, id)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
}

// Handlers - Relations
func (a *App) getConProviderConnections(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["con_prov_id"], 10, 64)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid provider ID")
		return
	}

	q := godevmandb.New(a.DB)
	res, err := q.GetConProviderConnections(a.Ctx, id)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, res)
}

// Helpers
func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

// Route definitions
func (a *App) initializeRoutes() {
	a.Router.HandleFunc("/connections/providers/count", a.countConProviders).Methods("GET")
	// Takes parameters: count(100), start(0). Uses default if not set.
	a.Router.HandleFunc("/connections/providers", a.getConProviders).Methods("GET")
	a.Router.HandleFunc("/connections/provider/{con_prov_id:[0-9]+}", a.getConProvider).Methods("GET")
	a.Router.HandleFunc("/connections/provider", a.createConProvider).Methods("POST")
	a.Router.HandleFunc("/connections/provider/{con_prov_id:[0-9]+}", a.updateConProvider).Methods("PUT")
	a.Router.HandleFunc("/connections/provider/{con_prov_id:[0-9]+}", a.deleteConProvider).Methods("DELETE")
	a.Router.HandleFunc("/connections/provider/{con_prov_id:[0-9]+}/connections", a.getConProviderConnections).Methods("GET")
}

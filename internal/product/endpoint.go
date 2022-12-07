package product

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/ncostamagna/gocourse_product/pkg/meta"
)

type (
	Controller func(w http.ResponseWriter, r *http.Request)

	Endpoints struct {
		Create Controller
		Get    Controller
		GetAll Controller
		Update Controller
		Delete Controller
	}

	CreateReq struct {
		Name  string  `json:"name"`
		Price float64 `json:"price"`
	}

	GetAllReq struct {
		Name string
	}

	UpdateReq struct {
		Name *string `json:"name"`
		// Estaba mal escrito el nombre del campo en el tag json
		Price *float64 `json:"price"`
	}

	Response struct {
		Status int         `json:"status"`
		Data   interface{} `json:"data,omitempty"`
		Err    string      `json:"error,omitempty"`
		Meta   *meta.Meta  `json:"meta,omitempty"`
	}
)

func MakeEndpoints(s Service) Endpoints {
	return Endpoints{
		Create: makeCreateEndpoint(s),
		Get:    makeGetEndpoint(s),
		GetAll: makeGetAllEndpoint(s),
		Update: makeUpdateEndpoint(s),
		Delete: makeDeleteEndpoint(s),
	}
}

func makeCreateEndpoint(s Service) Controller {
	return func(w http.ResponseWriter, r *http.Request) {

		var req CreateReq

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			w.WriteHeader(400)
			json.NewEncoder(w).Encode(&Response{Status: 400, Err: fmt.Sprintf("invalid request format: %v", err.Error())})
			return
		}

		if req.Name == "" {
			w.WriteHeader(400)
			json.NewEncoder(w).Encode(&Response{Status: 400, Err: "name is required"})
			return
		}

		if req.Price <= 0 {
			w.WriteHeader(400)
			json.NewEncoder(w).Encode(&Response{Status: 400, Err: "price must be greater than 0"})
			return
		}

		product, err := s.Create(req.Name, req.Price)
		if err != nil {
			w.WriteHeader(400)
			json.NewEncoder(w).Encode(&Response{Status: 400, Err: err.Error()})
			return
		}

		json.NewEncoder(w).Encode(&Response{Status: 200, Data: product})
	}
}

func makeGetEndpoint(s Service) Controller {
	return func(w http.ResponseWriter, r *http.Request) {

		path := mux.Vars(r)
		id := path["id"]

		product, err := s.Get(id)
		if err != nil {
			w.WriteHeader(404)
			json.NewEncoder(w).Encode(&Response{Status: 404, Err: "product doesn't exist"})
			return
		}

		json.NewEncoder(w).Encode(&Response{Status: 200, Data: product})
	}
}

func makeGetAllEndpoint(s Service) Controller {
	return func(w http.ResponseWriter, r *http.Request) {

		v := r.URL.Query()

		// El query string que esperamos es 'name'
		filters := Filters{
			Name: v.Get("name"),
		}

		limit, _ := strconv.Atoi(v.Get("limit"))
		page, _ := strconv.Atoi(v.Get("page"))

		count, err := s.Count(filters)
		if err != nil {
			w.WriteHeader(500)
			json.NewEncoder(w).Encode(&Response{Status: 500, Err: err.Error()})
			return
		}

		meta, err := meta.New(page, limit, count)
		if err != nil {
			w.WriteHeader(500)
			json.NewEncoder(w).Encode(&Response{Status: 500, Err: err.Error()})
			return
		}

		products, err := s.GetAll(filters, meta.Offset(), meta.Limit())
		if err != nil {
			w.WriteHeader(400)
			json.NewEncoder(w).Encode(&Response{Status: 400, Err: err.Error()})
			return
		}

		json.NewEncoder(w).Encode(&Response{Status: 200, Data: products, Meta: meta})
	}
}

func makeUpdateEndpoint(s Service) Controller {
	return func(w http.ResponseWriter, r *http.Request) {
		var req UpdateReq

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			w.WriteHeader(400)
			json.NewEncoder(w).Encode(&Response{Status: 400, Err: fmt.Sprintf("invalid request format: %v", err.Error())})
			return
		}

		if req.Name != nil && *req.Name == "" {
			w.WriteHeader(400)
			json.NewEncoder(w).Encode(&Response{Status: 400, Err: "name is required"})
			return
		}

		// Validamos el campo precio (misma validacion que en el create)
		if req.Price != nil && *req.Price <= 0 {
			w.WriteHeader(400)
			json.NewEncoder(w).Encode(&Response{Status: 400, Err: "price must be greater than 0"})
			return
		}

		path := mux.Vars(r)
		id := path["id"]

		if err := s.Update(id, req.Name, req.Price); err != nil {
			w.WriteHeader(404)
			json.NewEncoder(w).Encode(&Response{Status: 404, Err: "product doesn't exist"})
			return
		}

		json.NewEncoder(w).Encode(&Response{Status: 200, Data: "success"})
	}
}

func makeDeleteEndpoint(s Service) Controller {
	return func(w http.ResponseWriter, r *http.Request) {

		path := mux.Vars(r)
		id := path["id"]

		if err := s.Delete(id); err != nil {
			w.WriteHeader(404)
			json.NewEncoder(w).Encode(&Response{Status: 404, Err: "product doesn't exist"})
			return
		}

		json.NewEncoder(w).Encode(&Response{Status: 200, Data: "success"})
	}
}

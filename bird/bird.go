package bird

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/go-redis/redis"
	"github.com/gorilla/mux"
	"github.com/sankalpjonn/mockingbird/validator"
)

const (
	DB_HOST      = "localhost:6379"
	DB_PASSWORD  = ""
	DB_NUM       = 2
	DB_POOL_SIZE = 50

	EGG_HEADERS_KEY     = "egg:%s:headers"
	EGG_BODY_KEY        = "egg:%s:body"
	EGG_STATUS_CODE_KEY = "egg:%s:statuscode"
)

type Bird struct {
	store *redis.Client
}

func New() *Bird {
	b := new(Bird)
	b.store = redis.NewClient(&redis.Options{
		Addr:     DB_HOST,
		Password: DB_PASSWORD,
		DB:       DB_NUM,
		PoolSize: DB_POOL_SIZE,
	})
	return b
}

func (self *Bird) getEgg(eggid string) (error, *Egg) {
	pipe := self.store.Pipeline()
	headers := pipe.HGetAll(fmt.Sprintf(EGG_HEADERS_KEY, eggid))
	body := pipe.Get(fmt.Sprintf(EGG_BODY_KEY, eggid))
	statuscode := pipe.Get(fmt.Sprintf(EGG_STATUS_CODE_KEY, eggid))
	pipe.Exec()
	if statuscode.Val() == "" {
		return errors.New("No egg found with this id"), nil
	}
	statuscodeint, err := strconv.Atoi(statuscode.Val())
	if err != nil {
		panic(err)
	}

	return nil, &Egg{
		Id:         eggid,
		Headers:    headers.Val(),
		Body:       body.Val(),
		StatusCode: statuscodeint,
	}
}

func (self *Bird) createEgg(egg *Egg) {
	pipe := self.store.Pipeline()
	headerskey := fmt.Sprintf(EGG_HEADERS_KEY, egg.Id)
	bodykey := fmt.Sprintf(EGG_BODY_KEY, egg.Id)
	statuscodekey := fmt.Sprintf(EGG_STATUS_CODE_KEY, egg.Id)
	headers := map[string]interface{}{}
	for k, v := range egg.Headers {
		headers[k] = v
	}

	if egg.Headers != nil {
		pipe.HMSet(headerskey, headers)
		if egg.TTL > 0 {
			pipe.Expire(headerskey, time.Second*time.Duration(egg.TTL))
		}
	}
	if egg.Body != "" {
		pipe.Set(bodykey, egg.Body, time.Second*time.Duration(egg.TTL))
	}
	pipe.Set(statuscodekey, egg.StatusCode, time.Second*time.Duration(egg.TTL))
	_, err := pipe.Exec()
	if err != nil {
		panic(err)
	}
}

func (self *Bird) WriteResponse(w http.ResponseWriter, rawres []byte, err error) {
	if err != nil {
		http.Error(w, fmt.Sprintf("%s", err.Error()), http.StatusBadRequest)
	} else {
		w.Header().Set("Content-Type", "application/json;charset=utf-8")
		w.Write(rawres)
	}
}

func (self *Bird) CreateHandler(w http.ResponseWriter, r *http.Request) {
	// decode request
	egg := &Egg{}
	if err := json.NewDecoder(r.Body).Decode(egg); err != nil {
		self.WriteResponse(w, nil, err)
		return
	}

	// validate request
	if err := validator.ValidateRequest(*egg); err != nil {
		self.WriteResponse(w, nil, err)
		return
	}

	// make egg
	egg.initialize()
	self.createEgg(egg)

	// send response
	response := `
	{
      "egg_id": "` + egg.Id + `"
  	}
  	`
	self.WriteResponse(w, []byte(response), nil)
}

func (self *Bird) GetRawHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	eggid := vars["eggId"]

	err, egg := self.getEgg(eggid)
	if err != nil {
		self.WriteResponse(w, nil, err)
		return
	}
	res, err := json.Marshal(egg)
	if err != nil {
		panic(err)
	}
	w.WriteHeader(200)
	w.Write(res)
}

func (self *Bird) GetHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	eggid := vars["eggId"]

	err, egg := self.getEgg(eggid)
	if err != nil {
		self.WriteResponse(w, nil, err)
		return
	}

	for k, v := range egg.Headers {
		w.Header().Add(k, v)
	}
	w.WriteHeader(egg.StatusCode)
	if egg.Body != "" {
		w.Write([]byte(egg.Body))
	}
}

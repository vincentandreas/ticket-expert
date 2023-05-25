package controller

import (
	"encoding/json"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"io/ioutil"
	"log"
	"net/http"
	"ticket-expert/models"
	"ticket-expert/repo"
	"ticket-expert/utilities"
)

type BaseHandler struct {
	repo repo.AllRepository
}

func HandleRequests(h *BaseHandler) *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/api/user", h.HandleSaveUser).Methods("POST")
	router.HandleFunc("/api/promotor", h.HandleSavePromotor).Methods("POST")
	router.HandleFunc("/api/event", h.HandleSaveEvent).Methods("POST")
	router.HandleFunc("/api/book", h.HandleSaveBooking).Methods("POST")
	router.HandleFunc("/api/purchase", h.HandleSavePurchased).Methods("POST")
	router.HandleFunc("/api/health", h.CheckHealth).Methods("GET")
	return router
}

func NewBaseHandler(repo repo.AllRepository) *BaseHandler {
	return &BaseHandler{
		repo: repo,
	}
}

func (h *BaseHandler) CheckHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	temp := map[string]string{
		"result": "OK",
	}
	json.NewEncoder(w).Encode(temp)
}

func (h *BaseHandler) HandleSaveUser(w http.ResponseWriter, r *http.Request) {
	reqBody, _ := ioutil.ReadAll(r.Body)
	var userRequest models.User
	json.Unmarshal(reqBody, &userRequest)

	if !isValidRequest(w, userRequest) {
		return
	}
	h.repo.SaveUser(userRequest)
	utilities.WriteSuccessResp(w)
}

func (h *BaseHandler) HandleSavePromotor(w http.ResponseWriter, r *http.Request) {
	reqBody, _ := ioutil.ReadAll(r.Body)
	var reqObj models.Promotor
	json.Unmarshal(reqBody, &reqObj)

	if !isValidRequest(w, reqObj) {
		return
	}
	h.repo.SavePromotor(reqObj)
	utilities.WriteSuccessResp(w)
}

func (h *BaseHandler) HandleSaveEvent(w http.ResponseWriter, r *http.Request) {
	reqBody, _ := ioutil.ReadAll(r.Body)
	var reqObj models.Event
	json.Unmarshal(reqBody, &reqObj)

	if !isValidRequest(w, reqObj) {
		return
	}
	err := h.repo.SaveEvent(reqObj)
	if err != nil {
		log.Println(err)
		utilities.WriteErrorResp(w, 403, "Failed to save data")
		return
	}
	utilities.WriteSuccessResp(w)
}

func (h *BaseHandler) HandleSaveBooking(w http.ResponseWriter, r *http.Request) {
	reqBody, _ := ioutil.ReadAll(r.Body)
	var reqObj models.BookingTicket
	json.Unmarshal(reqBody, &reqObj)

	if !isValidRequest(w, reqObj) {
		return
	}
	err := h.repo.SaveBooking(reqObj)
	if err != nil {
		log.Println(err)
		utilities.WriteErrorResp(w, 403, "Failed to save data")
		return
	}
	utilities.WriteSuccessResp(w)
}

func (h *BaseHandler) HandleSavePurchased(w http.ResponseWriter, r *http.Request) {
	reqBody, _ := ioutil.ReadAll(r.Body)
	var reqObj models.PurchasedTicket
	json.Unmarshal(reqBody, &reqObj)

	if !isValidRequest(w, reqObj) {
		return
	}
	err := h.repo.SavePurchase(reqObj)
	if err != nil {
		log.Println(err)
		utilities.WriteErrorResp(w, 403, "Failed to save data")
		return
	}
	utilities.WriteSuccessResp(w)
}

func isValidRequest(w http.ResponseWriter, request interface{}) bool {
	validate := validator.New()
	err := validate.Struct(request)

	if err != nil {
		fmt.Println("Validation failed")
		utilities.WriteErrorResp(w, 400, "Request not valid")
		return false
	}
	return true
}

//func (h *BaseHandler) MainPage(w http.ResponseWriter, r *http.Request) {
//	if r.Method == "GET" {
//		//http.FileServer(http.Dir("static/"))
//		t, _ := template.ParseFiles("static/index.html")
//		t.Execute(w, nil)
//	}
//}

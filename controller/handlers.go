package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"ticket-expert/models"
	"ticket-expert/repo"
	"ticket-expert/utilities"
)

type BaseHandler struct {
	Repo repo.AllRepository
}

func HandleRequests(h *BaseHandler) *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/api/user", h.HandleSaveUser).Methods("POST")
	router.HandleFunc("/api/promotor", h.HandleSavePromotor).Methods("POST")
	//router.HandleFunc("/api/promotor/{id}", h.HandleSavePromotor).Methods("GET")
	router.HandleFunc("/api/event", h.HandleSaveEvent).Methods("POST")
	router.HandleFunc("/api/event", h.HandleSearchEvent).Methods("GET")
	router.HandleFunc("/api/event/{id}", h.HandleSearchEventById).Methods("GET")
	router.HandleFunc("/api/book", h.HandleSaveBooking).Methods("POST")
	router.HandleFunc("/api/book/check", h.HandleCheckBookingPeriod).Methods("GET")
	router.HandleFunc("/api/purchase", h.HandleSavePurchased).Methods("POST")
	router.HandleFunc("/api/health", h.CheckHealth).Methods("GET")
	router.HandleFunc("/api/waitingQueue", h.HandleSaveWaitingQueue).Methods("POST")
	router.HandleFunc("/api/testing/{id}", h.TempHandleTest).Methods("GET")
	router.HandleFunc("/api/checkOrderRoom/{eventId}", h.HandleCheckOrderRoom).Methods("GET")
	return router
}

func NewBaseHandler(repo repo.AllRepository) *BaseHandler {
	return &BaseHandler{
		Repo: repo,
	}
}

func (h *BaseHandler) TempHandleTest(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	idstr := vars["id"]

	h.Repo.SaveUserInOrderRoom(3, "3", "z12", context.TODO())
	h.Repo.SaveUserInOrderRoom(3, "7", "z12", context.TODO())
	if idstr == "yes" {
		h.Repo.PopUserInOrderRoom(3, 3, context.TODO())
		h.Repo.PopUserInOrderRoom(3, 7, context.TODO())
	}

	utilities.WriteSuccessResp(w)
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
	h.Repo.SaveUser(userRequest, r.Context())
	utilities.WriteSuccessResp(w)
}

func (h *BaseHandler) HandleSaveWaitingQueue(w http.ResponseWriter, r *http.Request) {
	reqBody, _ := ioutil.ReadAll(r.Body)
	var userRequest models.NewWaitingUser
	json.Unmarshal(reqBody, &userRequest)

	if !isValidRequest(w, userRequest) {
		return
	}

	//check user
	_, err := h.Repo.FindUserById(userRequest.UserId, r.Context())
	if err != nil {
		log.Println(err)
		if err.Error() == "record not found" {
			utilities.WriteErrorResp(w, 400, "User not found")
			return
		}
		utilities.WriteErrorResp(w, 400, "Error when search the user. ")
		return
	}

	checkRes := h.Repo.GetUserInOrderRoom(userRequest.UserId, userRequest.EventId, r.Context())
	if checkRes != "" {
		utilities.WriteErrorResp(w, 400, "User already in order room")
		return
	}
	h.Repo.SaveWaitingQueue(userRequest, r.Context())
	utilities.WriteSuccessResp(w)
}

func (h *BaseHandler) HandleSaveUserInOrderRoom(w http.ResponseWriter, r *http.Request) {
	reqBody, _ := ioutil.ReadAll(r.Body)
	var userRequest models.NewWaitingUser
	json.Unmarshal(reqBody, &userRequest)

	if !isValidRequest(w, userRequest) {
		return
	}
	h.Repo.SaveWaitingQueue(userRequest, r.Context())
	utilities.WriteSuccessResp(w)
}

func (h *BaseHandler) HandleSavePromotor(w http.ResponseWriter, r *http.Request) {
	reqBody, _ := ioutil.ReadAll(r.Body)
	var reqObj models.Promotor
	json.Unmarshal(reqBody, &reqObj)

	if !isValidRequest(w, reqObj) {
		return
	}
	h.Repo.SavePromotor(reqObj, r.Context())
	utilities.WriteSuccessResp(w)
}

func (h *BaseHandler) HandleSaveEvent(w http.ResponseWriter, r *http.Request) {
	reqBody, _ := ioutil.ReadAll(r.Body)
	var reqObj models.Event
	json.Unmarshal(reqBody, &reqObj)

	if !isValidRequest(w, reqObj) {
		return
	}
	err := h.Repo.SaveEvent(reqObj, r.Context())
	if err != nil {
		log.Println(err)
		utilities.WriteErrorResp(w, 403, err.Error())
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

	err := h.Repo.SaveBooking(reqObj, r.Context())
	if err != nil {
		log.Println(err)
		utilities.WriteErrorResp(w, 403, err.Error())
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
	err := h.Repo.SavePurchase(reqObj, r.Context())
	if err != nil {
		log.Println(err)
		utilities.WriteErrorResp(w, 403, err.Error())
		return
	}
	utilities.WriteSuccessResp(w)
}

func (h *BaseHandler) HandleCheckOrderRoom(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idstr := vars["eventId"]
	numb, err := strconv.ParseUint(idstr, 10, 32)
	if err != nil {
		return
	}
	h.Repo.CheckOrderRoom(uint(numb), r.Context())

	utilities.WriteSuccessResp(w)
}

func (h *BaseHandler) HandleSearchEventById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idstr := vars["id"]
	res, err := h.Repo.FindByEventId(idstr, r.Context())
	if err != nil {
		log.Println(err)
		utilities.WriteErrorResp(w, 403, "Failed to get data")
		return
	}
	utilities.WriteSuccessWithDataResp(w, res)
}

func (h *BaseHandler) HandleCheckBookingPeriod(w http.ResponseWriter, r *http.Request) {
	h.Repo.CheckBookingPeriod(r.Context())

	utilities.WriteSuccessWithDataResp(w, nil)
}

func (h *BaseHandler) HandleSearchEvent(w http.ResponseWriter, r *http.Request) {
	qparams := r.URL.Query()
	city := qparams.Get("city")
	category := qparams.Get("category")
	res, err := h.Repo.FindEventByCondition(city, category, r.Context())
	if err != nil {
		log.Println(err)
		utilities.WriteErrorResp(w, 403, "Failed to get data")
		return
	}
	utilities.WriteSuccessWithDataResp(w, res)
}

func isValidRequest(w http.ResponseWriter, request interface{}) bool {
	validate := validator.New()
	err := validate.Struct(request)

	if err != nil {
		fmt.Println(err)
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

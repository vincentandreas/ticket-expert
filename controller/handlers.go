package controller

import (
	"encoding/json"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/jcuga/golongpoll"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"ticket-expert/models"
	"ticket-expert/repo"
	"ticket-expert/services"
	"ticket-expert/utilities"
)

type BaseHandler struct {
	Repo      repo.AllRepository
	LPManager *golongpoll.LongpollManager
	Store     *sessions.CookieStore
}

func HandleRequests(h *BaseHandler, lp *golongpoll.LongpollManager) *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/api/user", h.HandleSaveUser).Methods("POST")
	router.HandleFunc("/api/user", h.HandleGetUserData).Methods("GET")
	router.HandleFunc("/api/user/login", h.HandleLogin).Methods("POST")
	router.HandleFunc("/api/event", h.HandleSaveEvent).Methods("POST")
	router.HandleFunc("/api/event", h.HandleSearchEvent).Methods("GET")
	router.HandleFunc("/api/event/{id}", h.HandleSearchEventById).Methods("GET")
	router.HandleFunc("/api/book", h.HandleSaveBooking).Methods("POST")
	router.HandleFunc("/api/book", h.HandleFindBookingByUserId).Methods("GET")
	router.HandleFunc("/api/book/check", h.HandleCheckBookingPeriod).Methods("GET")
	router.HandleFunc("/api/purchase", h.HandleSavePurchased).Methods("POST")
	router.HandleFunc("/api/health", h.CheckHealth).Methods("GET")
	router.HandleFunc("/api/waitingQueue/checkTotal/{eventId}", h.HandleWaitingRoomCheckTotal).Methods("GET")
	router.HandleFunc("/api/waitingQueue", h.HandleSaveWaitingQueue).Methods("POST")
	//router.HandleFunc("/api/testing/{id}", h.TempHandleTest).Methods("GET")
	router.HandleFunc("/api/orderRoom/checkAvailable/{eventId}", h.HandleCheckOrderRoom).Methods("GET")
	router.HandleFunc("/api/upload", h.UploadHandler).Methods("POST")
	router.HandleFunc("/api/upload", h.UploadOptHandler).Methods("OPTIONS")
	router.HandleFunc("/api/subQueue", h.WrapSubsHandler).Methods("GET")
	router.HandleFunc("/api/purchase/{id}", h.HandleFindPurchasedEventById).Methods("GET")
	return router
}

func NewBaseHandler(repo repo.AllRepository, lpMngr *golongpoll.LongpollManager, store *sessions.CookieStore) *BaseHandler {
	return &BaseHandler{
		Repo:      repo,
		LPManager: lpMngr,
		Store:     store,
	}
}

func (h *BaseHandler) UploadOptHandler(w http.ResponseWriter, r *http.Request) {
	utilities.WriteSuccessResp(w)
	return
}

func (h *BaseHandler) UploadHandler(w http.ResponseWriter, r *http.Request) {
	// Retrieve the uploaded file
	sessUserId := h.SessionGetUserId(r)
	if sessUserId == 0 {
		utilities.WriteUnauthResp(w)
		return
	}

	file, handler, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Failed to retrieve the file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	photoUrl := services.UploadToBucket(file, handler.Filename)
	if photoUrl == "" {
		utilities.WriteErrorResp(w, 400, "failed to upload photo")
		return
	}
	utilities.WriteSuccessWithDataResp(w, photoUrl)
}

func (h *BaseHandler) HandleFindPurchasedEventById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idstr := vars["id"]
	ticketDetail, err := h.Repo.FindPurchasedEventById(idstr)
	if err != nil {
		log.Println(err.Error())
		utilities.WriteErrorResp(w, 400, "error")
		return
	}
	utilities.WriteSuccessWithDataResp(w, ticketDetail)
}

func (h *BaseHandler) HandleFindBookingByUserId(w http.ResponseWriter, r *http.Request) {
	sessUserId := h.SessionGetUserId(r)
	if sessUserId == 0 {
		utilities.WriteUnauthResp(w)
		return
	}

	ticketDetail, err := h.Repo.FindBookingByUserId(sessUserId, r.Context())
	if err != nil {
		log.Println(err.Error())
		utilities.WriteErrorResp(w, 400, err.Error())
		return
	}
	utilities.WriteSuccessWithDataResp(w, ticketDetail)
}

func (h *BaseHandler) WrapSubsHandler(w http.ResponseWriter, r *http.Request) {
	utilities.SetAllHeaders(w)
	h.LPManager.SubscriptionHandler(w, r)
}

func (h *BaseHandler) CheckHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	temp := map[string]string{
		"result": "OK",
	}
	json.NewEncoder(w).Encode(temp)
}

func (h *BaseHandler) HandleLogin(w http.ResponseWriter, r *http.Request) {
	reqBody, _ := ioutil.ReadAll(r.Body)
	var userRequest models.UserLogin
	json.Unmarshal(reqBody, &userRequest)

	if !isValidRequest(w, userRequest) {
		return
	}
	user, err := h.Repo.Login(userRequest, r.Context())

	if err != nil {
		utilities.WriteErrorResp(w, 401, "Failed login")
		return
	}
	session, err := h.Store.New(r, "te-session")
	if err != nil {
		log.Println(err)
	}
	session.Values["user_name"] = userRequest.UserName
	session.Values["user_id"] = user.ID

	err = session.Save(r, w)
	if err != nil {
		log.Println(err)
	}
	utilities.WriteSuccessWithDataResp(w, user.Role)
}

// return userId, if session invalid, will return ""
func (h *BaseHandler) SessionGetUserId(r *http.Request) uint {
	session, err := h.Store.Get(r, "te-session")
	if err != nil {
		return 0
	}
	var userId uint
	if uid, ok := session.Values["user_id"].(uint); ok {
		userId = uid
	} else {
		userId = 0
	}
	return userId
}

func (h *BaseHandler) HandleGetUserData(w http.ResponseWriter, r *http.Request) {
	sessUserId := h.SessionGetUserId(r)
	if sessUserId == 0 {
		utilities.WriteUnauthResp(w)
		return
	}

	data, err := h.Repo.FindUserById(sessUserId, r.Context())
	if err != nil {
		return
	}
	utilities.WriteSuccessWithDataResp(w, data.Extract())
}

func (h *BaseHandler) HandleSaveUser(w http.ResponseWriter, r *http.Request) {
	reqBody, _ := ioutil.ReadAll(r.Body)
	var userRequest models.User
	json.Unmarshal(reqBody, &userRequest)

	if !isValidRequest(w, userRequest) {
		return
	}
	err := h.Repo.SaveUser(userRequest, r.Context())
	if checkError(w, err) {
		return
	}

	utilities.WriteSuccessResp(w)
}

func checkError(w http.ResponseWriter, err error) bool {
	if err != nil {
		utilities.WriteErrorResp(w, 400, "Error when process the request")
		log.Println(err.Error())
		return true
	}
	return false
}

func (h *BaseHandler) HandleSaveWaitingQueue(w http.ResponseWriter, r *http.Request) {
	sessUserId := h.SessionGetUserId(r)
	if sessUserId == 0 {
		utilities.WriteUnauthResp(w)
		return
	}

	reqBody, _ := ioutil.ReadAll(r.Body)
	var userRequest models.NewWaitingUser
	json.Unmarshal(reqBody, &userRequest)
	userRequest.UserId = sessUserId
	if !isValidRequest(w, userRequest) {
		return
	}
	//todo there's still no validation in queue length.
	//check user
	_, err := h.Repo.FindUserById(userRequest.UserId, r.Context())
	if err != nil {
		log.Println(err)
		if err.Error() == "record not found" {
			utilities.WriteErrorResp(w, 400, "User not found")
			return
		}
		utilities.WriteErrorResp(w, 400, "Error when search the user")
		return
	}

	checkRes := h.Repo.GetUserInOrderRoom(userRequest.UserId, userRequest.EventId, r.Context())
	if checkRes != "" {
		utilities.WriteErrorResp(w, 400, "User already in order room")
		return
	}

	//add unique id
	userRequest.QUniqueCode = uuid.New().String()

	h.Repo.SaveWaitingQueue(userRequest, r.Context())
	utilities.WriteSuccessWithDataResp(w, userRequest)
}

func (h *BaseHandler) HandleSaveUserInOrderRoom(w http.ResponseWriter, r *http.Request) {
	sessUserId := h.SessionGetUserId(r)
	if sessUserId == 0 {
		utilities.WriteUnauthResp(w)
		return
	}

	reqBody, _ := ioutil.ReadAll(r.Body)
	var userRequest models.NewWaitingUser
	json.Unmarshal(reqBody, &userRequest)

	if !isValidRequest(w, userRequest) {
		return
	}
	h.Repo.SaveWaitingQueue(userRequest, r.Context())
	utilities.WriteSuccessResp(w)
}

func (h *BaseHandler) HandleSaveEvent(w http.ResponseWriter, r *http.Request) {
	sessUserId := h.SessionGetUserId(r)
	if sessUserId == 0 {
		utilities.WriteUnauthResp(w)
		return
	}
	user, err := h.Repo.FindUserById(sessUserId, r.Context())
	if err != nil {
		log.Println(err.Error())
		utilities.WriteErrorResp(w, 400, "Error when saving events")
		return
	}
	if user.Role != "PROMOTOR" {
		errDesc := "Only Promotor can add event"
		log.Println(errDesc)
		utilities.WriteErrorResp(w, 400, errDesc)
		return
	}
	reqBody, _ := ioutil.ReadAll(r.Body)
	var reqObj models.Event
	json.Unmarshal(reqBody, &reqObj)

	if !isValidRequest(w, reqObj) {
		return
	}
	reqObj.UserID = sessUserId
	err = h.Repo.SaveEvent(reqObj, r.Context())
	if err != nil {
		log.Println(err)
		utilities.WriteErrorResp(w, 403, err.Error())
		return
	}
	utilities.WriteSuccessResp(w)
}

func (h *BaseHandler) HandleSaveBooking(w http.ResponseWriter, r *http.Request) {
	sessUserId := h.SessionGetUserId(r)
	if sessUserId == 0 {
		utilities.WriteUnauthResp(w)
		return
	}

	reqBody, _ := ioutil.ReadAll(r.Body)
	var reqObj models.BookingTicket
	json.Unmarshal(reqBody, &reqObj)

	reqObj.UserID = sessUserId
	reqObj.BookingStatus = "active"
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
	var reqObj models.PurchaseReq
	json.Unmarshal(reqBody, &reqObj)

	if !isValidRequest(w, reqObj) {
		return
	}
	booking, err := h.Repo.GetBookingByUniqCode(r.Context(), reqObj.BookingUniqCode)
	if checkError(w, err) {
		return
	}

	purchase := models.PurchasedTicket{
		UserID:          booking.UserID,
		BookingTicketID: booking.ID,
		PaymentMethod:   reqObj.PaymentMethod,
		TransRefNo:      reqObj.TransRefNo,
		PaymentStatus:   "success",
	}

	err = h.Repo.SavePurchase(purchase, r.Context())
	if checkError(w, err) {
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
	qUniqueCodes := h.Repo.CheckOrderRoom(uint(numb), r.Context())

	if qUniqueCodes != nil {
		for i := 0; i < len(qUniqueCodes); i++ {
			log.Printf("Allowed queue code %s enter the order room", qUniqueCodes[i])
			h.LPManager.Publish(qUniqueCodes[i], "enter room")
		}
	}

	utilities.WriteSuccessResp(w)
}

func (h *BaseHandler) HandleSearchEventById(w http.ResponseWriter, r *http.Request) {
	sessUserId := h.SessionGetUserId(r)
	if sessUserId == 0 {
		utilities.WriteUnauthResp(w)
		return
	}

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
	h.Repo.CheckBookingPeriodically(r.Context())
	utilities.WriteSuccessWithDataResp(w, nil)
}

func (h *BaseHandler) HandleSearchEvent(w http.ResponseWriter, r *http.Request) {
	sessUserId := h.SessionGetUserId(r)
	if sessUserId == 0 {
		utilities.WriteUnauthResp(w)
		return
	}
	qparams := r.URL.Query()
	city := qparams.Get("city")
	category := qparams.Get("category")
	name := qparams.Get("name")
	res, err := h.Repo.FindEventByCondition(name, city, category, r.Context())
	if err != nil {
		log.Println(err)
		utilities.WriteErrorResp(w, 403, "Failed to get data")
		return
	}
	utilities.WriteSuccessWithDataResp(w, res)
}

func (h *BaseHandler) HandleWaitingRoomCheckTotal(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	idstr := vars["eventId"]
	numb, err := strconv.ParseUint(idstr, 10, 32)
	if err != nil {
		utilities.WriteErrorResp(w, 400, err.Error())
		return
	}
	totalPeople := h.Repo.CountTotalPeopleInWaitingRoom(uint(numb), r.Context())
	utilities.WriteSuccessWithDataResp(w, totalPeople)
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

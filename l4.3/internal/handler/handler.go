package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/AugustSerenity/go-contest-L4/l4.3_Events-calendar/internal/handler/response"
	"github.com/AugustSerenity/go-contest-L4/l4.3_Events-calendar/internal/middleware"
	"github.com/AugustSerenity/go-contest-L4/l4.3_Events-calendar/internal/model"
)

type Handler struct {
	service Service
}

func New(s Service) *Handler {
	return &Handler{
		service: s,
	}
}

func (h *Handler) Route() http.Handler {
	router := http.NewServeMux()

	router.HandleFunc("POST /create_event", middleware.LoggingMiddleware(h.CreateEvent))
	router.HandleFunc("POST /update_event", middleware.LoggingMiddleware(h.UpdateEvent))
	router.HandleFunc("POST /delete_event", middleware.LoggingMiddleware(h.DeleteEvent))
	router.HandleFunc("GET /events_for_day", middleware.LoggingMiddleware(h.EventsForDay))
	router.HandleFunc("GET /events_for_week", middleware.LoggingMiddleware(h.EventsForWeek))
	router.HandleFunc("GET /events_for_month", middleware.LoggingMiddleware(h.EventsForMonth))

	return router
}

func (h *Handler) CreateEvent(w http.ResponseWriter, r *http.Request) {
	requestID := r.Header.Get("X-Request-ID")

	var req model.CreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.SendError(w, http.StatusBadRequest, "bad request")
		return
	}

	eventID, err := h.service.CreateEvent(req, requestID)
	if err != nil {
		response.SendError(w, http.StatusBadRequest, err.Error())
		return
	}

	response.SendSuccess(w, http.StatusOK, map[string]interface{}{
		"event_id": eventID,
		"message":  "successfully created",
	})
}

func (h *Handler) UpdateEvent(w http.ResponseWriter, r *http.Request) {
	requestID := r.Header.Get("X-Request-ID")

	var req model.UpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.SendError(w, http.StatusBadRequest, "bad request")
		return
	}

	if err := h.service.UpdateEvent(req, requestID); err != nil {
		response.SendError(w, http.StatusBadRequest, err.Error())
		return
	}

	response.SendSuccess(w, http.StatusOK, map[string]string{
		"message": "successfully updated",
	})
}

func (h *Handler) DeleteEvent(w http.ResponseWriter, r *http.Request) {
	requestID := r.Header.Get("X-Request-ID")

	var req model.DeleteRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.SendError(w, http.StatusBadRequest, "bad request")
		return
	}

	if err := h.service.DeleteEvent(req, requestID); err != nil {
		response.SendError(w, http.StatusBadRequest, err.Error())
		return
	}

	response.SendSuccess(w, http.StatusOK, map[string]string{
		"message": "successfully deleted",
	})
}

func (h *Handler) EventsForDay(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()

	userID, err := strconv.Atoi(query.Get("user_id"))
	if err != nil {
		response.SendError(w, http.StatusBadRequest, "invalid user_id")
		return
	}

	dateStr := query.Get("date")
	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		response.SendError(w, http.StatusBadRequest, "invalid date format, use YYYY-MM-DD")
		return
	}

	events, err := h.service.EventsForDay(userID, date)
	if err != nil {
		response.SendError(w, http.StatusNotFound, err.Error())
		return
	}

	response.SendSuccess(w, http.StatusOK, events)
}

func (h *Handler) EventsForWeek(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()

	userID, err := strconv.Atoi(query.Get("user_id"))
	if err != nil {
		response.SendError(w, http.StatusBadRequest, "invalid user_id")
		return
	}

	dateStr := query.Get("date")
	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		response.SendError(w, http.StatusBadRequest, "invalid date format, use YYYY-MM-DD")
		return
	}

	events, err := h.service.EventsForWeek(userID, date)
	if err != nil {
		response.SendError(w, http.StatusNotFound, err.Error())
		return
	}

	response.SendSuccess(w, http.StatusOK, events)
}

func (h *Handler) EventsForMonth(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()

	userID, err := strconv.Atoi(query.Get("user_id"))
	if err != nil {
		response.SendError(w, http.StatusBadRequest, "invalid user_id")
		return
	}

	dateStr := query.Get("date")
	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		response.SendError(w, http.StatusBadRequest, "invalid date format, use YYYY-MM-DD")
		return
	}

	events, err := h.service.EventsForMonth(userID, date)
	if err != nil {
		response.SendError(w, http.StatusNotFound, err.Error())
		return
	}

	response.SendSuccess(w, http.StatusOK, events)
}

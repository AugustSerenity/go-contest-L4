package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"l2.18/internal/handler/response"
	"l2.18/internal/middleware"
	"l2.18/internal/model"
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
	if r.Method != http.MethodPost {
		response.SendError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	var eventRequest model.Request
	err := json.NewDecoder(r.Body).Decode(&eventRequest)
	if err != nil {
		response.SendError(w, http.StatusBadRequest, "bad request")
		return
	}

	event, err := model.CastToEvent(eventRequest)
	if err != nil {
		response.SendError(w, http.StatusBadRequest, err.Error())
		return
	}

	err = h.service.CreateEvent(eventRequest.UserID, event)
	if err != nil {
		response.SendError(w, http.StatusServiceUnavailable, err.Error())
		return
	}

	response.SendSuccess(w, http.StatusOK, "successfully created")

}

func (h *Handler) UpdateEvent(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		response.SendError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	var eventRequest model.Request
	err := json.NewDecoder(r.Body).Decode(&eventRequest)
	if err != nil {
		response.SendError(w, http.StatusBadRequest, "bad request")
		return
	}

	originalDate, err := time.Parse("2006-01-02 15:04:05", eventRequest.Date)
	if err != nil {
		response.SendError(w, http.StatusBadRequest, "invalid date format (expected YYYY-MM-DD HH:MM:SS)")
		return
	}

	event, err := model.CastToEvent(eventRequest)
	if err != nil {
		response.SendError(w, http.StatusBadRequest, err.Error())
		return
	}
	err = h.service.UpdateEvent(eventRequest.UserID, originalDate, event)
	if err != nil {
		response.SendError(w, http.StatusServiceUnavailable, err.Error())
		return
	}

	response.SendSuccess(w, http.StatusOK, "successfully updated")
}

func (h *Handler) DeleteEvent(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		response.SendError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	var eventRequest model.Request
	err := json.NewDecoder(r.Body).Decode(&eventRequest)
	if err != nil {
		response.SendError(w, http.StatusBadRequest, "bad request")
		return
	}

	date, err := time.Parse("2006-01-02 15:04:05", eventRequest.Date)
	if err != nil {
		response.SendError(w, http.StatusBadRequest, "invalid date format")
		return
	}

	err = h.service.DeleteEvent(eventRequest.UserID, date, eventRequest.EventName)
	if err != nil {
		response.SendError(w, http.StatusNotFound, err.Error())
		return
	}

	response.SendSuccess(w, http.StatusOK, "successfully deleted")
}

func (h *Handler) EventsForDay(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		response.SendError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	userIDStr := r.URL.Query().Get("user_id")
	dateStr := r.URL.Query().Get("date")

	if userIDStr == "" || dateStr == "" {
		response.SendError(w, http.StatusBadRequest, "missing query parameters")
		return
	}

	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		response.SendError(w, http.StatusBadRequest, "invalid user_id")
		return
	}

	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		response.SendError(w, http.StatusBadRequest, "invalid date format (expected YYYY-MM-DD)")
		return
	}

	events, err := h.service.ShowEventsForDay(userID, date)
	if err != nil {
		response.SendError(w, http.StatusNotFound, err.Error())
		return
	}

	response.SendSuccess(w, http.StatusOK, events)
}

func (h *Handler) EventsForWeek(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		response.SendError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	userIDStr := r.URL.Query().Get("user_id")
	dateStr := r.URL.Query().Get("date")
	if userIDStr == "" || dateStr == "" {
		response.SendError(w, http.StatusBadRequest, "missing query parameters")
		return
	}

	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		response.SendError(w, http.StatusBadRequest, "invalid user_id")
		return
	}

	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		response.SendError(w, http.StatusBadRequest, "invalid date format (expected YYYY-MM-DD)")
		return
	}

	events, err := h.service.ShowEventsForWeek(userID, date)
	if err != nil {
		response.SendError(w, http.StatusNotFound, err.Error())
		return
	}

	response.SendSuccess(w, http.StatusOK, events)
}

func (h *Handler) EventsForMonth(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		response.SendError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	userIDStr := r.URL.Query().Get("user_id")
	dateStr := r.URL.Query().Get("date")
	if userIDStr == "" || dateStr == "" {
		response.SendError(w, http.StatusBadRequest, "missing query parameters")
		return
	}

	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		response.SendError(w, http.StatusBadRequest, "invalid user_id")
		return
	}

	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		response.SendError(w, http.StatusBadRequest, "invalid date format (YYYY-MM-DD)")
		return
	}

	events, err := h.service.ShowEventsForMonth(userID, date)
	if err != nil {
		response.SendError(w, http.StatusNotFound, err.Error())
		return
	}

	response.SendSuccess(w, http.StatusOK, events)
}

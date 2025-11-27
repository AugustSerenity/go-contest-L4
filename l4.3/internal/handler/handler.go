package handler

import (
	"encoding/json"
	"net/http"

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
		response.SendError(w, http.StatusServiceUnavailable, err.Error())
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

	response.SendSuccess(w, http.StatusOK, "successfully updated")
}

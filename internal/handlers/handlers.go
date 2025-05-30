package handlers

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"github.com/NessibeliY/quotes-api/internal/dto"
	"github.com/NessibeliY/quotes-api/internal/models"
	"github.com/NessibeliY/quotes-api/internal/service"
)

type Handler struct {
	Service *service.Service
}

func NewHandler(service *service.Service) *Handler {
	return &Handler{
		Service: service,
	}
}

func (h *Handler) AddQuote(w http.ResponseWriter, r *http.Request) {
	var req dto.AddQuoteRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		slog.Warn("decode request", "error", err)
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	err = req.Validate()
	if err != nil {
		slog.Warn("validate request", "error", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	quote := h.Service.AddQuote(req.Author, req.Quote)
	slog.Info("add quote", "quote", quote)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(quote)
}

func (h *Handler) GetAllQuotes(w http.ResponseWriter, r *http.Request) {
	author := r.URL.Query().Get("author")
	var quotes []models.Quote
	if author != "" {
		quotes = h.Service.GetQuotesByAuthor(author)
		slog.Info("get quotes by author", "author", author, "quotes", quotes)
	} else {
		quotes = h.Service.GetAllQuotes()
		slog.Info("get all quotes", "quotes", quotes)
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(quotes)
}

func (h *Handler) GetRandomQuote(w http.ResponseWriter, r *http.Request) {
	quote, ok := h.Service.GetRandomQuote()
	if !ok {
		http.Error(w, "no quotes", http.StatusNotFound)
		return
	}
	slog.Info("get random quote", "quote", quote)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(quote)
}

func (h *Handler) DeleteQuote(w http.ResponseWriter, r *http.Request) {
	idStr := mux.Vars(r)["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		slog.Warn("parse id", "id", idStr, "error", err)
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	if !h.Service.DeleteQuote(int64(id)) {
		slog.Warn("delete quote", "id", id)
		http.Error(w, "not found", http.StatusNotFound)
		return
	}
	slog.Info("delete quote", "id", id)
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) Health(w http.ResponseWriter, r *http.Request) {
	slog.Info("health check")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

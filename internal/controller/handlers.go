package controller

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/paxaf/BrandScoutTest/internal/entity"
)

func (h *UsecaseHandler) Add(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	quote, err := ParseQuoteFromReq(r)
	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return

	}
	h.service.Set(*quote)
	w.WriteHeader(http.StatusCreated)
}

func (h *UsecaseHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	quotes := h.service.GetAll()
	resp := entity.QuoteResponse{Quotes: quotes}
	data, err := json.Marshal(resp)
	if err != nil {
		http.Error(w, "Internal Error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(data); err != nil {
		log.Printf("Failed to write response: %v", err)
	}
}

func (h *UsecaseHandler) GetRand(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	quote, ok := h.service.Random()
	if !ok {
		http.Error(w, "No content", http.StatusNoContent)
		return
	}
	data, err := json.Marshal(quote)
	if err != nil {
		http.Error(w, "Internal Error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(data); err != nil {
		log.Printf("Failed to write response: %v", err)
	}
}

func (h *UsecaseHandler) ByAutor(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	author := r.URL.Query().Get("author")
	quotes, ok := h.service.GetAllByAuthor(author)
	if !ok {
		http.Error(w, "No content", http.StatusNoContent)
		return
	}
	resp := entity.QuoteResponse{Quotes: quotes}
	data, err := json.Marshal(resp)
	if err != nil {
		http.Error(w, "Internal Error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(data); err != nil {
		log.Printf("Failed to write response: %v", err)
	}
}

func (h *UsecaseHandler) Delete(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	path := strings.Split(r.URL.Path, "/")
	key := path[2]
	err := h.service.Delete(key)
	if err != nil {
		log.Println(err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
	}
	w.WriteHeader(http.StatusOK)
}

func ParseQuoteFromReq(r *http.Request) (*entity.Quote, error) {
	var quote entity.Quote
	contentType := r.Header.Get("Content-Type")
	if contentType != "application/json" {
		return nil, fmt.Errorf("invalid content type: %s", contentType)
	}

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&quote)
	if err != nil {
		return nil, fmt.Errorf("failed to decode request body: %w", err)
	}

	return &quote, nil
}

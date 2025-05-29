package controller_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"sync/atomic"
	"testing"

	"github.com/paxaf/BrandScoutTest/internal/controller"
	"github.com/paxaf/BrandScoutTest/internal/entity"
)

type MockUsecase struct {
	quotes     map[string]entity.Quote
	keyCounter atomic.Uint64
	returnErr  bool
}

func (m *MockUsecase) Set(quote entity.Quote) {
	if m.returnErr {
		return
	}
	key := strconv.FormatUint(m.keyCounter.Add(1), 10)
	quote.Id = key
	m.quotes[key] = quote
}

func (m *MockUsecase) GetAll() []entity.Quote {
	if m.returnErr {
		return nil
	}
	quotes := make([]entity.Quote, 0, len(m.quotes))
	for _, q := range m.quotes {
		quotes = append(quotes, q)
	}
	return quotes
}

func (m *MockUsecase) Random() (entity.Quote, bool) {
	if m.returnErr || len(m.quotes) == 0 {
		return entity.Quote{}, false
	}
	for _, q := range m.quotes {
		return q, true
	}
	return entity.Quote{}, false
}

func (m *MockUsecase) GetAllByAuthor(author string) ([]entity.Quote, bool) {
	if m.returnErr {
		return nil, false
	}
	var result []entity.Quote
	for _, q := range m.quotes {
		if q.Author == author {
			result = append(result, q)
		}
	}
	return result, len(result) > 0
}

func (m *MockUsecase) Delete(id string) error {
	if m.returnErr {
		return errors.New("mock error")
	}
	if _, exists := m.quotes[id]; !exists {
		return errors.New("not found")
	}
	delete(m.quotes, id)
	return nil
}

func TestAddHandler(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		t.Parallel()
		mockUsecase := &MockUsecase{quotes: make(map[string]entity.Quote)}
		h := controller.New(mockUsecase)

		quote := entity.Quote{Author: "Me", Phrase: "Hello"}
		body, _ := json.Marshal(quote)
		req := httptest.NewRequest(http.MethodPost, "/add", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		h.Add(w, req)

		if w.Code != http.StatusCreated {
			t.Errorf("Expected status 201, got %d", w.Code)
		}
		if len(mockUsecase.quotes) != 1 {
			t.Fatalf("Quote not added to service")
		}

		for _, q := range mockUsecase.quotes {
			if q.Author != "Me" || q.Phrase != "Hello" {
				t.Errorf("Unexpected quote content: %+v", q)
			}
			if q.Id == "" {
				t.Error("Quote ID not generated")
			}
		}
	})

	t.Run("wrong method", func(t *testing.T) {
		t.Parallel()
		h := controller.UsecaseHandler{}
		req := httptest.NewRequest(http.MethodGet, "/add", nil)
		w := httptest.NewRecorder()

		h.Add(w, req)

		if w.Code != http.StatusMethodNotAllowed {
			t.Errorf("Expected status 405, got %d", w.Code)
		}
	})

	t.Run("invalid content type", func(t *testing.T) {
		t.Parallel()
		h := controller.UsecaseHandler{}
		req := httptest.NewRequest(http.MethodPost, "/add", nil)
		req.Header.Set("Content-Type", "text/plain")
		w := httptest.NewRecorder()

		h.Add(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("Expected status 400, got %d", w.Code)
		}
	})

	t.Run("malformed JSON", func(t *testing.T) {
		t.Parallel()
		h := controller.UsecaseHandler{}
		req := httptest.NewRequest(http.MethodPost, "/add", strings.NewReader("{invalid}"))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		h.Add(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("Expected status 400, got %d", w.Code)
		}
	})
}

func TestGetAllHandler(t *testing.T) {
	t.Parallel()

	t.Run("success with quotes", func(t *testing.T) {
		t.Parallel()
		mockUsecase := &MockUsecase{
			quotes: map[string]entity.Quote{
				"1": {Id: "1", Author: "Author", Phrase: "Test quote"},
			},
		}
		h := controller.New(mockUsecase)

		req := httptest.NewRequest(http.MethodGet, "/quotes", nil)
		w := httptest.NewRecorder()

		h.GetAll(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d", w.Code)
		}
		if ct := w.Header().Get("Content-Type"); ct != "application/json" {
			t.Errorf("Expected JSON content, got %s", ct)
		}

		var resp entity.QuoteResponse
		if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
			t.Fatalf("Failed to unmarshal response: %v", err)
		}
		if len(resp.Quotes) != 1 {
			t.Fatalf("Expected 1 quote, got %d", len(resp.Quotes))
		}
		if resp.Quotes[0].Id != "1" || resp.Quotes[0].Phrase != "Test quote" {
			t.Errorf("Unexpected quote data: %+v", resp.Quotes[0])
		}
	})

	t.Run("empty storage", func(t *testing.T) {
		t.Parallel()
		mockUsecase := &MockUsecase{quotes: make(map[string]entity.Quote)}
		h := controller.New(mockUsecase)

		req := httptest.NewRequest(http.MethodGet, "/quotes", nil)
		w := httptest.NewRecorder()

		h.GetAll(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d", w.Code)
		}

		var resp entity.QuoteResponse
		if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
			t.Fatalf("Failed to unmarshal response: %v", err)
		}
		if len(resp.Quotes) != 0 {
			t.Errorf("Expected 0 quotes, got %d", len(resp.Quotes))
		}
	})

	t.Run("wrong method", func(t *testing.T) {
		t.Parallel()
		h := controller.UsecaseHandler{}
		req := httptest.NewRequest(http.MethodPost, "/quotes", nil)
		w := httptest.NewRecorder()

		h.GetAll(w, req)

		if w.Code != http.StatusMethodNotAllowed {
			t.Errorf("Expected status 405, got %d", w.Code)
		}
	})
}

func TestGetRandHandler(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		t.Parallel()
		mockUsecase := &MockUsecase{
			quotes: map[string]entity.Quote{
				"1": {Id: "1", Author: "Author", Phrase: "Test quote"},
			},
		}
		h := controller.New(mockUsecase)

		req := httptest.NewRequest(http.MethodGet, "/random", nil)
		w := httptest.NewRecorder()

		h.GetRand(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d", w.Code)
		}

		var quote entity.Quote
		if err := json.Unmarshal(w.Body.Bytes(), &quote); err != nil {
			t.Fatalf("Failed to unmarshal response: %v", err)
		}
		if quote.Id != "1" || quote.Phrase != "Test quote" {
			t.Errorf("Unexpected quote: %+v", quote)
		}
	})

	t.Run("no content", func(t *testing.T) {
		t.Parallel()
		mockUsecase := &MockUsecase{quotes: make(map[string]entity.Quote)}
		h := controller.New(mockUsecase)

		req := httptest.NewRequest(http.MethodGet, "/random", nil)
		w := httptest.NewRecorder()

		h.GetRand(w, req)

		if w.Code != http.StatusNoContent {
			t.Errorf("Expected status 204, got %d", w.Code)
		}
	})
}

func TestByAuthorHandler(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		t.Parallel()
		mockUsecase := &MockUsecase{
			quotes: map[string]entity.Quote{
				"1": {Id: "1", Author: "Tester", Phrase: "Quote 1"},
				"2": {Id: "2", Author: "Tester", Phrase: "Quote 2"},
				"3": {Id: "3", Author: "Other", Phrase: "Other quote"},
			},
		}
		h := controller.New(mockUsecase)

		req := httptest.NewRequest(http.MethodGet, "/byauthor?author=Tester", nil)
		w := httptest.NewRecorder()

		h.ByAutor(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d", w.Code)
		}

		var resp entity.QuoteResponse
		if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
			t.Fatalf("Failed to unmarshal response: %v", err)
		}
		if len(resp.Quotes) != 2 {
			t.Fatalf("Expected 2 quotes, got %d", len(resp.Quotes))
		}
		for _, q := range resp.Quotes {
			if q.Author != "Tester" {
				t.Errorf("Unexpected author in quote: %s", q.Author)
			}
		}
	})

	t.Run("no content", func(t *testing.T) {
		t.Parallel()
		mockUsecase := &MockUsecase{quotes: make(map[string]entity.Quote)}
		h := controller.New(mockUsecase)

		req := httptest.NewRequest(http.MethodGet, "/byauthor?author=Unknown", nil)
		w := httptest.NewRecorder()

		h.ByAutor(w, req)

		if w.Code != http.StatusNoContent {
			t.Errorf("Expected status 204, got %d", w.Code)
		}
	})

	t.Run("missing author parameter", func(t *testing.T) {
		t.Parallel()
		mockUsecase := &MockUsecase{
			quotes: map[string]entity.Quote{
				"1": {Id: "1", Author: "Tester", Phrase: "Quote"},
			},
		}
		h := controller.New(mockUsecase)

		req := httptest.NewRequest(http.MethodGet, "/byauthor", nil)
		w := httptest.NewRecorder()

		h.ByAutor(w, req)

		if w.Code != http.StatusNoContent {
			t.Errorf("Expected status 204 for missing author, got %d", w.Code)
		}
	})
}

func TestDeleteHandler(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		t.Parallel()
		mockUsecase := &MockUsecase{
			quotes: map[string]entity.Quote{
				"test-id": {Id: "test-id", Author: "Author", Phrase: "Test quote"},
			},
		}
		h := controller.New(mockUsecase)

		req := httptest.NewRequest(http.MethodDelete, "/delete/test-id", nil)
		w := httptest.NewRecorder()

		h.Delete(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d", w.Code)
		}
		if _, exists := mockUsecase.quotes["test-id"]; exists {
			t.Errorf("Quote was not deleted")
		}
	})

	t.Run("not found", func(t *testing.T) {
		t.Parallel()
		mockUsecase := &MockUsecase{quotes: make(map[string]entity.Quote)}
		h := controller.New(mockUsecase)

		req := httptest.NewRequest(http.MethodDelete, "/delete/missing-id", nil)
		w := httptest.NewRecorder()

		h.Delete(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("Expected status 400, got %d", w.Code)
		}
	})

	t.Run("service error", func(t *testing.T) {
		t.Parallel()
		mockUsecase := &MockUsecase{returnErr: true}
		h := controller.New(mockUsecase)

		req := httptest.NewRequest(http.MethodDelete, "/delete/id", nil)
		w := httptest.NewRecorder()

		h.Delete(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("Expected status 400, got %d", w.Code)
		}
	})

	t.Run("missing id in path", func(t *testing.T) {
		t.Parallel()
		mockUsecase := &MockUsecase{quotes: make(map[string]entity.Quote)}
		h := controller.New(mockUsecase)

		req := httptest.NewRequest(http.MethodDelete, "/delete/", nil)
		w := httptest.NewRecorder()

		h.Delete(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("Expected status 400 for missing id, got %d", w.Code)
		}
	})
}

func TestParseQuoteFromReq(t *testing.T) {
	t.Parallel()

	t.Run("valid", func(t *testing.T) {
		t.Parallel()
		quote := entity.Quote{Author: "Me", Phrase: "Hello"}
		body, _ := json.Marshal(quote)
		req := httptest.NewRequest("POST", "/", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		result, err := controller.ParseQuoteFromReq(req)
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
		if result.Author != quote.Author || result.Phrase != quote.Phrase {
			t.Errorf("Parsed quote doesn't match original")
		}
	})

	t.Run("invalid content type", func(t *testing.T) {
		t.Parallel()
		req := httptest.NewRequest("POST", "/", strings.NewReader("{}"))
		req.Header.Set("Content-Type", "text/plain")

		_, err := controller.ParseQuoteFromReq(req)
		if err == nil {
			t.Fatal("Expected error but got none")
		}
	})

	t.Run("malformed json", func(t *testing.T) {
		t.Parallel()
		req := httptest.NewRequest("POST", "/", strings.NewReader("{invalid}"))
		req.Header.Set("Content-Type", "application/json")

		_, err := controller.ParseQuoteFromReq(req)
		if err == nil {
			t.Fatal("Expected error but got none")
		}
	})

	t.Run("empty body", func(t *testing.T) {
		t.Parallel()
		req := httptest.NewRequest("POST", "/", nil)
		req.Header.Set("Content-Type", "application/json")

		_, err := controller.ParseQuoteFromReq(req)
		if err == nil {
			t.Fatal("Expected error but got none")
		}
	})
}

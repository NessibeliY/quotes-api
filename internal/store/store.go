package store

import (
	"sync"

	"github.com/NessibeliY/quotes-api/internal/models"
)

type Store struct {
	mu     sync.RWMutex
	quotes []models.Quote
	nextID int64
}

func NewStore() *Store {
	return &Store{
		quotes: []models.Quote{},
		nextID: 1,
	}
}

func (s *Store) AddQuote(author, quote string) models.Quote {
	s.mu.Lock()
	defer s.mu.Unlock()
	q := models.Quote{
		ID:     s.nextID,
		Author: author,
		Quote:  quote,
	}
	s.quotes = append(s.quotes, q)
	s.nextID++
	return q
}

func (s *Store) GetAllQuotes() []models.Quote {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return append([]models.Quote{}, s.quotes...)
}

func (s *Store) GetQuotesByAuthor(author string) []models.Quote {
	s.mu.RLock()
	defer s.mu.RUnlock()
	var res []models.Quote
	for _, q := range s.quotes {
		if q.Author == author {
			res = append(res, q)
		}
	}
	return res
}

func (s *Store) DeleteQuoteByID(id int64) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	for i, q := range s.quotes {
		if q.ID == id {
			s.quotes = append(s.quotes[:i], s.quotes[i+1:]...)
			return true
		}
	}
	return false
}

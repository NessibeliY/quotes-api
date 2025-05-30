package service

import (
	"math/rand"
	"time"

	"github.com/NessibeliY/quotes-api/internal/models"
	"github.com/NessibeliY/quotes-api/internal/store"
)

type QuoteService interface {
	AddQuote(author, quote string) models.Quote
	GetQuotesByAuthor(author string) []models.Quote
	GetAllQuotes() []models.Quote
	GetRandomQuote() (models.Quote, bool)
	DeleteQuote(id int64) bool
}

type Service struct {
	store store.QuoteStore
	rand  *rand.Rand
}

func NewService(store *store.Store) *Service {
	return &Service{
		store: store,
		rand:  rand.New(rand.NewSource(time.Now().UnixNano())), //nolint:gosec
	}
}

func (s *Service) AddQuote(author, quote string) models.Quote {
	return s.store.AddQuote(author, quote)
}

func (s *Service) GetAllQuotes() []models.Quote {
	return s.store.GetAllQuotes()
}

func (s *Service) GetQuotesByAuthor(author string) []models.Quote {
	return s.store.GetQuotesByAuthor(author)
}

func (s *Service) GetRandomQuote() (models.Quote, bool) {
	all := s.store.GetAllQuotes()
	if len(all) == 0 {
		return models.Quote{}, false
	}
	return all[rand.Intn(len(all))], true //nolint:gosec
}

func (s *Service) DeleteQuote(id int64) bool {
	return s.store.DeleteQuoteByID(id)
}

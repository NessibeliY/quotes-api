package service

import (
	"testing"

	"github.com/NessibeliY/quotes-api/internal/store"
)

func TestService_AddAndGetAllQuotes(t *testing.T) {
	st := store.NewStore()
	srv := NewService(st)

	q1 := srv.AddQuote("Author1", "Quote1")
	q2 := srv.AddQuote("Author2", "Quote2")

	all := srv.GetAllQuotes()
	if len(all) != 2 {
		t.Errorf("expected 2 quotes, got %d", len(all))
	}

	if all[0] != q1 || all[1] != q2 {
		t.Errorf("expected quotes %v and %v, got %v", q1, q2, all)
	}
}

func TestService_GetQuotesByAuthor(t *testing.T) {
	st := store.NewStore()
	srv := NewService(st)

	srv.AddQuote("Author1", "Quote1")
	srv.AddQuote("Author2", "Quote2")
	srv.AddQuote("Author1", "Quote3")

	quotes := srv.GetQuotesByAuthor("Author1")
	if len(quotes) != 2 {
		t.Errorf("expected 2 quotes by Author1, got %d", len(quotes))
	}
}

func TestService_GetRandomQuote(t *testing.T) {
	st := store.NewStore()
	srv := NewService(st)

	q, ok := srv.GetRandomQuote()
	if ok {
		t.Errorf("expected false when there are no quotes, got true with %+v", q)
	}

	srv.AddQuote("Author1", "Quote1")
	q, ok = srv.GetRandomQuote()
	if !ok {
		t.Errorf("expected true when there is a quote, got false")
	}
	if q.Author != "Author1" || q.Quote != "Quote1" {
		t.Errorf("expected author=Author1, quote=Quote1, got author=%s, quote=%s", q.Author, q.Quote)
	}

	srv.AddQuote("Author2", "Quote2")
	counts := make(map[int64]int)
	for i := 0; i < 100; i++ {
		q, ok = srv.GetRandomQuote()
		if !ok {
			t.Errorf("expected true when there is a quote, got false")
		}
		counts[q.ID]++
	}

	if len(counts) != 2 {
		t.Errorf("expected 2 quotes, got %d", len(counts))
	}
}

func TestService_DeleteQuote(t *testing.T) {
	st := store.NewStore()
	srv := NewService(st)

	q1 := srv.AddQuote("Author1", "Quote1")
	ok := srv.DeleteQuote(q1.ID)
	if !ok {
		t.Errorf("expected to delete quote with ID %d", q1.ID)
	}

	if len(srv.GetAllQuotes()) != 0 {
		t.Errorf("expected 0 quotes left, got %d", len(srv.GetAllQuotes()))
	}
}

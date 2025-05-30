package store

import "testing"

func TestAddQuote(t *testing.T) {
	s := NewStore()
	q := s.AddQuote("Author1", "Quote1")

	if q.ID != 1 {
		t.Errorf("expected id=1, got id=%d", q.ID)
	}
	if q.Author != "Author1" || q.Quote != "Quote1" {
		t.Errorf("expected author=author, quote=quote, got author=%s, quote=%s", q.Author, q.Quote)
	}

	all := s.GetAllQuotes()
	if len(all) != 1 {
		t.Errorf("expected 1 quote, got %d", len(all))
	}
}

func TestGetAllQuotes(t *testing.T) {
	s := NewStore()
	s.AddQuote("Author1", "Quote1")
	s.AddQuote("Author2", "Quote2")

	all := s.GetAllQuotes()
	if len(all) != 2 {
		t.Errorf("expected 2 quotes, got %d", len(all))
	}

	if all[0].Author != "Author1" || all[1].Author != "Author2" {
		t.Errorf("expected author=author1, author=author2, got author=%s, author=%s", all[0].Author, all[1].Author)
	}
}

func TestGetQuotesByAuthor(t *testing.T) {
	s := NewStore()
	s.AddQuote("Author1", "Quote1")
	s.AddQuote("Author2", "Quote2")
	s.AddQuote("Author1", "Quote3")

	quotes := s.GetQuotesByAuthor("Author1")
	if len(quotes) != 2 {
		t.Errorf("expected 2 quotes by Author1, got %d", len(quotes))
	}

	for _, q := range quotes {
		if q.Author != "Author1" {
			t.Errorf("expected author=Author1, got author=%s", q.Author)
		}
	}
}

func TestDeleteQuoteByID(t *testing.T) {
	s := NewStore()
	q1 := s.AddQuote("Author1", "Quote1")
	q2 := s.AddQuote("Author2", "Quote2")

	ok := s.DeleteQuoteByID(q1.ID)
	if !ok {
		t.Errorf("expected to delete quote with ID %d", q1.ID)
	}

	all := s.GetAllQuotes()
	if len(all) != 1 {
		t.Errorf("expected 1 quote left, got %d", len(all))
	}
	if all[0].ID != q2.ID {
		t.Errorf("expected quote with ID %d, got quote with ID %d", q2.ID, all[0].ID)
	}

	ok = s.DeleteQuoteByID(999)
	if ok {
		t.Errorf("expected false when deleting non-existent quote")
	}
}

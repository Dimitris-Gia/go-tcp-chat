package server

import (
	"sync"
	"testing"
)

func TestHistory_EmptyOnCreate(t *testing.T) {
	h := NewHistory()
	if entries := h.GetAll(); len(entries) != 0 {
		t.Fatalf("expected empty history, got %d entries", len(entries))
	}
}

func TestHistory_AddAndGetAll(t *testing.T) {
	h := NewHistory()
	h.Add("msg1\n")
	h.Add("msg2\n")

	entries := h.GetAll()
	if len(entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(entries))
	}
	if entries[0] != "msg1\n" || entries[1] != "msg2\n" {
		t.Fatalf("unexpected entries: %v", entries)
	}
}

func TestHistory_GetAllReturnsCopy(t *testing.T) {
	h := NewHistory()
	h.Add("original\n")

	entries := h.GetAll()
	entries[0] = "mutated\n"

	fresh := h.GetAll()
	if fresh[0] != "original\n" {
		t.Fatal("GetAll should return a copy, not a reference to internal slice")
	}
}

func TestHistory_ConcurrentAdd(t *testing.T) {
	h := NewHistory()
	var wg sync.WaitGroup
	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			h.Add("msg\n")
		}()
	}
	wg.Wait()

	if got := len(h.GetAll()); got != 50 {
		t.Fatalf("expected 50 entries, got %d", got)
	}
}

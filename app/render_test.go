package main

import (
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestRenderNotesListHTML(t *testing.T) {
	rec := httptest.NewRecorder()
	notes := []noteRow{{ID: 1, Title: "First"}, {ID: 2, Title: "Second"}}
	renderNotesListHTML(rec, notes)

	assert.Equal(t, "text/html; charset=utf-8", rec.Header().Get("Content-Type"))
	body := rec.Body.String()
	assert.Contains(t, body, "<table")
	assert.Contains(t, body, "First")
	assert.Contains(t, body, "Second")
}

func TestRenderNoteHTML(t *testing.T) {
	rec := httptest.NewRecorder()
	n := Note{
		ID:        7,
		Title:     "Test",
		Content:   "Body text",
		CreatedAt: time.Date(2024, 1, 2, 3, 4, 5, 0, time.UTC),
	}
	renderNoteHTML(rec, n)

	assert.Equal(t, "text/html; charset=utf-8", rec.Header().Get("Content-Type"))
	body := rec.Body.String()
	assert.Contains(t, body, "Note #7")
	assert.Contains(t, body, "Test")
	assert.Contains(t, body, "Body text")
}

func TestRenderJSON(t *testing.T) {
	rec := httptest.NewRecorder()
	renderJSON(rec, map[string]string{"status": "ok"})

	assert.Equal(t, "application/json", rec.Header().Get("Content-Type"))
	assert.True(t, strings.Contains(rec.Body.String(), `"status":"ok"`))
}

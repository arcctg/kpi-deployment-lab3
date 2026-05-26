package main

import (
	"database/sql"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newTestApp(t *testing.T) (*App, sqlmock.Sqlmock, func()) {
	t.Helper()
	db, mock, err := sqlmock.New(sqlmock.MonitorPingsOption(true))
	require.NoError(t, err)
	return &App{db: db}, mock, func() { db.Close() }
}

func TestHandleAlive(t *testing.T) {
	app, _, cleanup := newTestApp(t)
	defer cleanup()

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/health/alive", nil)
	app.handleAlive(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "OK", rec.Body.String())
}

func TestHandleReadySuccess(t *testing.T) {
	app, mock, cleanup := newTestApp(t)
	defer cleanup()

	mock.ExpectPing()
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/health/ready", nil)
	app.handleReady(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "OK", rec.Body.String())
}

func TestHandleReadyFailure(t *testing.T) {
	app, mock, cleanup := newTestApp(t)
	defer cleanup()

	mock.ExpectPing().WillReturnError(sql.ErrConnDone)
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/health/ready", nil)
	app.handleReady(rec, req)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
}

func TestHandleRoot(t *testing.T) {
	app, _, cleanup := newTestApp(t)
	defer cleanup()

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	app.handleRoot(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Body.String(), "GET /notes")
}

func TestHandleRootNotFound(t *testing.T) {
	app, _, cleanup := newTestApp(t)
	defer cleanup()

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/missing", nil)
	app.handleRoot(rec, req)

	assert.Equal(t, http.StatusNotFound, rec.Code)
}

func TestListNotesJSON(t *testing.T) {
	app, mock, cleanup := newTestApp(t)
	defer cleanup()

	rows := sqlmock.NewRows([]string{"id", "title"}).
		AddRow(1, "One").
		AddRow(2, "Two")
	mock.ExpectQuery(`SELECT id, title FROM notes ORDER BY id`).
		WillReturnRows(rows)

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/notes", nil)
	req.Header.Set("Accept", "application/json")
	app.handleNotes(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Body.String(), `"title":"One"`)
}

func TestCreateNoteJSON(t *testing.T) {
	app, mock, cleanup := newTestApp(t)
	defer cleanup()

	created := time.Date(2024, 5, 1, 12, 0, 0, 0, time.UTC)
	mock.ExpectQuery(`INSERT INTO notes`).
		WithArgs("Title", "Content").
		WillReturnRows(sqlmock.NewRows([]string{"id", "title", "content", "created_at"}).
			AddRow(1, "Title", "Content", created))

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/notes", strings.NewReader(`{"title":"Title","content":"Content"}`))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	app.handleNotes(rec, req)

	assert.Equal(t, http.StatusCreated, rec.Code)
	assert.Contains(t, rec.Body.String(), `"title":"Title"`)
}

func TestHandleNoteByIDNotFound(t *testing.T) {
	app, mock, cleanup := newTestApp(t)
	defer cleanup()

	mock.ExpectQuery(`SELECT id, title, content, created_at FROM notes WHERE id = \$1`).
		WithArgs(99).
		WillReturnError(sql.ErrNoRows)

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/notes/99", nil)
	app.handleNoteByID(rec, req)

	assert.Equal(t, http.StatusNotFound, rec.Code)
}

func TestDemoFailingCheck(t *testing.T) {
	t.Fatal("intentional failure for Lab 3 blocked PR demonstration")
}

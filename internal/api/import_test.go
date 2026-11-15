package api

import (
	"bytes"
	"mime/multipart"
	"net/http/httptest"
	"testing"

	"github.com/tanq16/expenseowl/internal/storage"
)

func TestCSVPreservesSourceIDsAndSkipsReimport(t *testing.T) {
	store, err := storage.InitializeJsonStore(storage.SystemConfig{StorageURL: t.TempDir()})
	if err != nil {
		t.Fatal(err)
	}
	handler := NewHandler(store)
	for attempt := 0; attempt < 2; attempt++ {
		var body bytes.Buffer
		writer := multipart.NewWriter(&body)
		file, err := writer.CreateFormFile("file", "statement.csv")
		if err != nil {
			t.Fatal(err)
		}
		file.Write([]byte("ID,Name,Category,Amount,Date,Tags,Currency\n bank_001 ,Groceries,Food,-160.10,2026-09-16,mat,nok\n"))
		writer.Close()
		request := httptest.NewRequest("POST", "/import/csv", &body)
		request.Header.Set("Content-Type", writer.FormDataContentType())
		response := httptest.NewRecorder()
		handler.ImportCSV(response, request)
		if response.Code != 200 {
			t.Fatal(response.Code, response.Body.String())
		}
	}
	rows, err := store.GetAllExpenses()
	if err != nil || len(rows) != 1 {
		t.Fatalf("rows=%v err=%v", rows, err)
	}
	if rows[0].ID != "bank_001" || rows[0].Currency != "nok" || rows[0].Amount != -160.10 {
		t.Fatal(rows[0])
	}
}

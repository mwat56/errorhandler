/*
Copyright © 2019, 2024 M.Watermann, 10247 Berlin, Germany

	    All rights reserved
	EMail : <support@mwat.de>
*/
package errorhandler

//lint:file-ignore ST1017 - I prefer Yoda conditions

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

type tMockErrorPager struct{}

func (m *tMockErrorPager) GetErrorPage(aData []byte, aStatus int) []byte {
	return []byte("<html><body>Custom Error Page</body></html>")
}

func Test_WriteHeader_NoContentType_StatusOK(t *testing.T) {
	// Arrange
	rec := httptest.NewRecorder()
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ew := &tErrorWriter{
			ResponseWriter: w,
			errPager:       nil,
			status:         0,
		}
		ew.WriteHeader(http.StatusOK)
	})

	// Act
	handler.ServeHTTP(rec, nil)

	// Assert
	contentType := rec.Header().Get("Content-Type")
	if "" != contentType {
		t.Errorf("Expected no Content-Type header, but got %s", contentType)
	}
} // Test_WriteHeader_NoContentType_StatusOK()

func Test_WriteHeader_ContentType_StatusNotOK_WithPager(t *testing.T) {
	// Arrange
	rec := httptest.NewRecorder()
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ew := &tErrorWriter{
			ResponseWriter: w,
			errPager:       &tMockErrorPager{},
			status:         0,
		}
		ew.WriteHeader(http.StatusNotFound)
	})

	// Act
	handler.ServeHTTP(rec, nil)

	// Assert
	contentType := rec.Header().Get("Content-Type")
	if contentType != "text/html; charset=utf-8" {
		t.Errorf("Expected Content-Type 'text/html; charset=utf-8', but got %q", contentType)
	}
} // Test_WriteHeader_ContentType_StatusNotOK_WithPager()

func Test_WriteHeader_NoContentType_StatusOK_WithPager(t *testing.T) {
	// Arrange
	rec := httptest.NewRecorder()
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ew := &tErrorWriter{
			ResponseWriter: w,
			errPager:       &tMockErrorPager{},
			status:         0,
		}
		ew.WriteHeader(http.StatusOK)
	})

	// Act
	handler.ServeHTTP(rec, nil)

	// Assert
	contentType := rec.Header().Get("Content-Type")
	if "" != contentType {
		t.Errorf("Expected no Content-Type header, but got %s", contentType)
	}
} // Test_WriteHeader_NoContentType_StatusOK_WithPager()

func Test_tErrorWriter_WriteHeader(t *testing.T) {
	// Create a mock ResponseWriter
	w := httptest.ResponseRecorder{}
	ew := tErrorWriter{
		/*http.ResponseWriter:*/ &w,
		/*errPager:*/ &tMockErrorPager{},
		/*status:*/ 0,
	}

	tests := []struct {
		name   string
		fields tErrorWriter
		args   int
	}{
		{"01", ew, 0},
		{"02", ew, 100},
		{"03", ew, 200},
		{"04", ew, 300},
		{"05", ew, 400},
		{"06", ew, 500},
		{"07", ew, 600},
		{"08", ew, 700},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name,
			func(t *testing.T) {
				ew := &tErrorWriter{
					ResponseWriter: tt.fields.ResponseWriter,
					errPager:       tt.fields.errPager,
					status:         tt.fields.status,
				}
				ew.WriteHeader(tt.args)
				if ew.status != tt.args {
					t.Errorf("Expected status code %d, got %d", tt.args, ew.status)
				}
			})
	}
} // Test_tErrorWriter_WriteHeader()

func Test_tErrorWriter_Write(t *testing.T) {
	w := httptest.ResponseRecorder{}
	ew := tErrorWriter{
		/*http.ResponseWriter:*/ &w,
		/*errPager:*/ &tMockErrorPager{},
		/*status:*/ 0,
	}

	tests := []struct {
		name    string
		fields  tErrorWriter
		args    []byte
		want    int
		wantErr bool
	}{
		{"01", ew, []byte("test 1"), 6, false},
		{"02", ew, []byte("test 02"), 7, false},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ew := &tErrorWriter{
				ResponseWriter: tt.fields.ResponseWriter,
				errPager:       tt.fields.errPager,
				status:         tt.fields.status,
			}
			got, err := ew.Write(tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("tErrorWriter.Write() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("tErrorWriter.Write() = %v, want %v", got, tt.want)
			}
		})
	}
} // Test_tErrorWriter_Write()

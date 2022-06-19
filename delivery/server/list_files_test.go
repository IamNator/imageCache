package server_test

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"imageCache/delivery/server"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestListUploadedFile(t *testing.T) {
	r := server.SetUpRouter()
	r.GET("/list", server.ListFilesHandler)

	req, _ := http.NewRequest("GET", "/list", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var files []string
	er := json.Unmarshal(w.Body.Bytes(), &files)
	if er != nil {
		t.Error(er.Error())
	}

	assert.NotEmpty(t, files)
}

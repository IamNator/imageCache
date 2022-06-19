package server_test

import (
	"github.com/stretchr/testify/assert"
	"imageCache/delivery/server"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestDownloadFile(t *testing.T) {
	r := server.SetUpRouter()

	r.GET("/data/files/:fileName", server.DownloadFileHandler)

	req := httptest.NewRequest(http.MethodGet, "/data/files/sample2.png", nil)

	//var u = url.Values{}
	//u.Set("fileName", "sample2.png")
	//req.URL.RawQuery = u.Encode()

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

}

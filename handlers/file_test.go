package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFileHandler(t *testing.T) {

	r, err := http.NewRequest("GET", "public/index.html", nil)
	assert.NoError(t, err, "error making a request")

	w := httptest.NewRecorder()

	fileHandler("public/index.html")(w, r)

	responseHeader := w.Header()

	assert.Equal(t, []string{"no-cache, no-store, must-revalidate"}, responseHeader["Cache-Control"])
	assert.Equal(t, "no-cache", responseHeader["Pragma"][0])
	assert.Equal(t, "0", responseHeader["Expires"][0])
}

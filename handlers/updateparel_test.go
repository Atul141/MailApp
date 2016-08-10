package handlers

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	tu "git.mailbox.com/mailbox/testutils"
	"fmt"
)


func TestParcelUpdateSuccess(t *testing.T) {
	r, err := http.NewRequest("PATCH", "/parcels/23",
		strings.NewReader(`{"status":false}`))
	require.NoError(t, err, "failed to create a request")
	w := httptest.NewRecorder()

	mockDbObj := new(tu.MockDB)
	mockDbObj.On("UpdateParcelStatusById", "23", false).Return(nil)

	updateParcelHandler(mockDbObj)(w, r)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestParcelUpdateFailedForNoRequest(t *testing.T) {
	r, err := http.NewRequest("PATCH", "/parcels/23", nil)
	require.NoError(t, err, "failed to create a request")
	w := httptest.NewRecorder()

	mockDbObj := new(tu.MockDB)
	mockDbObj.On("UpdateParcelStatusById", "23", false).Return(nil)

	updateParcelHandler(mockDbObj)(w, r)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Equal(t, w.Body.String(), "request body is empty\n")
}

func TestParcelUpdateFailedForEmptyRequestBody(t *testing.T) {
	r, err := http.NewRequest("PATCH", "/parcels/23", strings.NewReader(``))
	require.NoError(t, err, "failed to create a request")
	w := httptest.NewRecorder()

	mockDbObj := new(tu.MockDB)
	mockDbObj.On("UpdateParcelStatusById", "23", false).Return(nil)

	updateParcelHandler(mockDbObj)(w, r)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Equal(t, w.Body.String(), "request body parsing failed\n")
}

func TestParcelUpdateFailedForMalformedJsonRequest(t *testing.T) {
	r, err := http.NewRequest("PATCH", "/parcels/23",
		strings.NewReader(`malformed json request`))
	require.NoError(t, err, "failed to create a request")
	w := httptest.NewRecorder()

	mockDbObj := new(tu.MockDB)
	mockDbObj.On("UpdateParcelStatusById", "23", false).Return(nil)

	updateParcelHandler(mockDbObj)(w, r)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Equal(t, w.Body.String(), "request body parsing failed\n")
}

func TestParcelUpdationFailed(t *testing.T) {
	r, err := http.NewRequest("PATCH", "/parcels/23",
		strings.NewReader(`{"status":false}`))
	require.NoError(t, err, "failed to create a request")
	w := httptest.NewRecorder()

	mockDbObj := new(tu.MockDB)
	mockDbObj.On("UpdateParcelStatusById", "23", false).Return(fmt.Errorf("Failed to update the parcel"))

	updateParcelHandler(mockDbObj)(w, r)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Equal(t, w.Body.String(), fmt.Sprintf("{\"code\":%d,\"message\":\"Internal server error\"}\n", w.Code))
}

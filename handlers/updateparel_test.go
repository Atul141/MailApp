package handlers

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	tu "git.mailbox.com/mailbox/testutils"
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
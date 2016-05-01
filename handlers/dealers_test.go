package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
    "fmt"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"


	tu "git.mailbox.com/mailbox/testutils"
	m "git.mailbox.com/mailbox/models"
	u "git.mailbox.com/mailbox/utils"
)

var dealers = []*m.Dealer{
	&m.Dealer{
		ID:   u.SPtr("ada1103c-4024-4ea4-b955-58c1c2c702b7"),
		Name: u.SPtr("Flipkart"),
	},
	&m.Dealer{
		ID:   u.SPtr("884a4cf3-2399-4481-b872-13233eaa3d6f"),
		Name: u.SPtr("Amazon"),
	},
}

func TestGetDealersSuccess(t *testing.T) {
	r, err := http.NewRequest("GET", "/dealers", nil)
	require.NoError(t, err, "failed to create a request: dealers")
	w := httptest.NewRecorder()

	mockDbObj := new(tu.MockDB)
	mockDbObj.On("GetDealers").Return(dealers, nil)

	dealersHandler(mockDbObj)(w, r)

	var actualDealers []*m.Dealer
	err = json.Unmarshal(w.Body.Bytes(), &actualDealers)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, *dealers[0].ID, *actualDealers[0].ID)
	assert.Equal(t, *dealers[1].ID, *actualDealers[1].ID)
	assert.Equal(t, *dealers[0].Name, *actualDealers[0].Name)
	assert.Equal(t, *dealers[1].Name, *actualDealers[1].Name)

	mockDbObj.AssertExpectations(t)
}

func TestGetDealersDBError(t *testing.T) {
	r, err := http.NewRequest("GET", "/dealers", nil)
	require.NoError(t, err, "failed to create a request: dealers")
	w := httptest.NewRecorder()

	expectedErrorMessage := "some-db-specific-error"
	mockDbObj := new(tu.MockDB)
	mockDbObj.On("GetDealers").Return(nil, fmt.Errorf(expectedErrorMessage))

	dealersHandler(mockDbObj)(w, r)

	var actualError *m.Error
	err = json.Unmarshal(w.Body.Bytes(), &actualError)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	require.NotNil(t, actualError)
	assert.Equal(t, int32(500), *actualError.Code)
	assert.Equal(t, "Internal server error", *actualError.Message)
	require.Nil(t, actualError.Fields)

	mockDbObj.AssertExpectations(t)
}


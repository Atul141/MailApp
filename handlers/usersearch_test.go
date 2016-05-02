package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	strfmt "github.com/go-openapi/strfmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	m "git.mailbox.com/mailbox/models"
	tu "git.mailbox.com/mailbox/testutils"
	u "git.mailbox.com/mailbox/utils"
)

var users = []*m.User{
	&m.User{
		Email:   strfmt.Email("hello@hello.com"),
		EmpID:   "11111",
		Name:    u.SPtr("Hello"),
		PhoneNo: u.SPtr("9900099900"),
	},
	&m.User{
		Email:   strfmt.Email("mello@mello.com"),
		Name:    u.SPtr("Mello"),
		EmpID:   "11113",
		PhoneNo: u.SPtr("9910399900"),
	},
	&m.User{
		Email:   strfmt.Email("jello@jello.com"),
		Name:    u.SPtr("Jello"),
		EmpID:   "11114",
		PhoneNo: u.SPtr("9900299900"),
	},
}

func TestUserSearchSuccess(t *testing.T) {
	r, err := http.NewRequest("GET", "/users/search?q=ello", nil)
	require.NoError(t, err, "failed to create a request: dealers")
	w := httptest.NewRecorder()

	mockDbObj := new(tu.MockDB)
	mockDbObj.On("GetUsersWith", "ello").Return(users, nil)

	userSearchHandler(mockDbObj)(w, r)

	var actualUsers []*m.User
	err = json.Unmarshal(w.Body.Bytes(), &actualUsers)
	require.NoError(t, err, "failed to unmarshal the response")

	assert.Equal(t, http.StatusOK, w.Code)

	assert.Equal(t, *users[0].Name, *actualUsers[0].Name)
	assert.Equal(t, *users[1].Name, *actualUsers[1].Name)
	assert.Equal(t, *users[2].Name, *actualUsers[2].Name)

	assert.Equal(t, users[0].Email, actualUsers[0].Email)
	assert.Equal(t, users[1].Email, actualUsers[1].Email)
	assert.Equal(t, users[2].Email, actualUsers[2].Email)

	assert.Equal(t, users[0].EmpID, actualUsers[0].EmpID)
	assert.Equal(t, users[1].EmpID, actualUsers[1].EmpID)
	assert.Equal(t, users[2].EmpID, actualUsers[2].EmpID)

	assert.Equal(t, *users[0].PhoneNo, *actualUsers[0].PhoneNo)
	assert.Equal(t, *users[1].PhoneNo, *actualUsers[1].PhoneNo)
	assert.Equal(t, *users[2].PhoneNo, *actualUsers[2].PhoneNo)

	mockDbObj.AssertExpectations(t)
}

func TestUserSearchSuccessNoQueryParam(t *testing.T) {
	r, err := http.NewRequest("GET", "/users/search", nil)
	require.NoError(t, err, "failed to create a request: dealers")
	w := httptest.NewRecorder()

	mockDbObj := new(tu.MockDB)

	userSearchHandler(mockDbObj)(w, r)

	var actualUsers []*m.User
	err = json.Unmarshal(w.Body.Bytes(), &actualUsers)
	require.NoError(t, err, "failed to unmarshal the response")

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, 0, len(actualUsers))
}

func TestUserSearchSuccessWhenLessThan3QueryParam(t *testing.T) {
	r, err := http.NewRequest("GET", "/users/search?q=el", nil)
	require.NoError(t, err, "failed to create a request: dealers")
	w := httptest.NewRecorder()

	mockDbObj := new(tu.MockDB)

	userSearchHandler(mockDbObj)(w, r)

	var actualUsers []*m.User
	err = json.Unmarshal(w.Body.Bytes(), &actualUsers)
	require.NoError(t, err, "failed to unmarshal the response")

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, 0, len(actualUsers))
}

func TestUserSearchDBFailure(t *testing.T) {
	r, err := http.NewRequest("GET", "/users/search?q=ello", nil)
	require.NoError(t, err, "failed to create a request: dealers")
	w := httptest.NewRecorder()

	mockDbObj := new(tu.MockDB)
	mockDbObj.On("GetUsersWith", "ello").Return(nil, fmt.Errorf("failure to connect to db"))

	userSearchHandler(mockDbObj)(w, r)

	var actualError *m.Error
	err = json.Unmarshal(w.Body.Bytes(), &actualError)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	require.NotNil(t, actualError)
	assert.Equal(t, int32(500), *actualError.Code)
	assert.Equal(t, "Internal server error", *actualError.Message)
	require.Nil(t, actualError.Fields)

	mockDbObj.AssertExpectations(t)
}

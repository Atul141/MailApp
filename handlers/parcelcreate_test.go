package handlers

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	strfmt "github.com/go-openapi/strfmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	m "git.mailbox.com/mailbox/models"
	tu "git.mailbox.com/mailbox/testutils"
	u "git.mailbox.com/mailbox/utils"
)

var dealer = &m.Dealer{
	ID:   u.SPtr("a8f4e46c-8295-4f53-ab0a-7fc2d2f47d28"),
	Name: u.SPtr("Flipkart"),
}

var user = &m.User{
	ID:      u.SPtr("23abce0a-ceb7-4127-8d98-b0bb5df4cce7"),
	Email:   strfmt.Email("hello@hello.com"),
	EmpID:   "11111",
	Name:    u.SPtr("Hello"),
	PhoneNo: u.SPtr("9900099900"),
}

var parcel = &m.Parcel{
	ID:       "23abce0a-ceb7-4127-8d98-b0bb5df4cce8",
	Dealer:   dealer,
	Owner:    user,
	DealerID: dealer.ID,
	OwnerID:  user.ID,
}

func TestParcelCreateSuccess(t *testing.T) {
	r, err := http.NewRequest("POST", "/parcels",
		strings.NewReader(`{ "dealerId": "a8f4e46c-8295-4f53-ab0a-7fc2d2f47d28", "ownerId": "23abce0a-ceb7-4127-8d98-b0bb5df4cce7" }`))
	require.NoError(t, err, "failed to create a request")
	w := httptest.NewRecorder()

	mockDbObj := new(tu.MockDB)
	mockDbObj.On("GetDealerByID", "a8f4e46c-8295-4f53-ab0a-7fc2d2f47d28").Return(dealer, nil)
	mockDbObj.On("GetUserByID", "23abce0a-ceb7-4127-8d98-b0bb5df4cce7").Return(user, nil)
	mockDbObj.On("CreateParcel", *dealer.ID, *user.ID).Return(parcel, nil)

	parcelCreateHandler(mockDbObj)(w, r)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "application/json", w.Header().Get("Content-Type"))
	assert.Equal(t, fmt.Sprintf(`{"id":"%s"}`, parcel.ID), w.Body.String())
}

func TestParcelCreateWhenNoRequestBody(t *testing.T) {
	r, err := http.NewRequest("POST", "/parcels", nil)
	require.NoError(t, err, "failed to create a request")
	w := httptest.NewRecorder()

	mockDbObj := new(tu.MockDB)
	mockDbObj.On("GetDealerByID", "a8f4e46c-8295-4f53-ab0a-7fc2d2f47d28").Return(dealer, nil)
	mockDbObj.On("GetUserByID", "23abce0a-ceb7-4127-8d98-b0bb5df4cce7").Return(user, nil)
	mockDbObj.On("CreateParcel", *dealer.ID, *user.ID).Return(parcel, nil)

	parcelCreateHandler(mockDbObj)(w, r)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestParcelCreateMalformedRequest(t *testing.T) {
	r, err := http.NewRequest("POST", "/parcels", strings.NewReader(`malformed json`))
	require.NoError(t, err, "failed to create a request")
	w := httptest.NewRecorder()

	mockDbObj := new(tu.MockDB)
	mockDbObj.On("GetDealerByID", "a8f4e46c-8295-4f53-ab0a-7fc2d2f47d28").Return(dealer, nil)
	mockDbObj.On("GetUserByID", "23abce0a-ceb7-4127-8d98-b0bb5df4cce7").Return(user, nil)
	mockDbObj.On("CreateParcel", *dealer.ID, *user.ID).Return(parcel, nil)

	parcelCreateHandler(mockDbObj)(w, r)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestParcelCreateValidationFail(t *testing.T) {
	r, err := http.NewRequest("POST", "/parcels", strings.NewReader(`{ }`))
	require.NoError(t, err, "failed to create a request")
	w := httptest.NewRecorder()

	mockDbObj := new(tu.MockDB)
	mockDbObj.On("GetDealerByID", "a8f4e46c-8295-4f53-ab0a-7fc2d2f47d28").Return(dealer, nil)
	mockDbObj.On("GetUserByID", "23abce0a-ceb7-4127-8d98-b0bb5df4cce7").Return(user, nil)
	mockDbObj.On("CreateParcel", *dealer.ID, *user.ID).Return(parcel, nil)

	parcelCreateHandler(mockDbObj)(w, r)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	validationFailMessage := fmt.Sprintf(`{"code":%d,"message":"dealerId should not be empty\ndealerId should be a valid UUID V4 string\ndealerId should not be empty\ndealerId should be a valid UUID V4 string"}`, w.Code)
	assert.Equal(t, fmt.Sprintf("%s\n", validationFailMessage), w.Body.String())
}

func TestParcelCreateGetDealerFail(t *testing.T) {
	r, err := http.NewRequest("POST", "/parcels",
		strings.NewReader(`{ "dealerId": "a8f4e46c-8295-4f53-ab0a-7fc2d2f47d28", "ownerId": "23abce0a-ceb7-4127-8d98-b0bb5df4cce7" }`))
	require.NoError(t, err, "failed to create a request")
	w := httptest.NewRecorder()

	mockDbObj := new(tu.MockDB)
	mockDbObj.On("GetDealerByID", "a8f4e46c-8295-4f53-ab0a-7fc2d2f47d28").Return(nil, fmt.Errorf("failed to get dealer with id"))
	mockDbObj.On("GetUserByID", "23abce0a-ceb7-4127-8d98-b0bb5df4cce7").Return(user, nil)
	mockDbObj.On("CreateParcel", *dealer.ID, *user.ID).Return(parcel, nil)

	parcelCreateHandler(mockDbObj)(w, r)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestParcelCreateGetUserFail(t *testing.T) {
	r, err := http.NewRequest("POST", "/parcels",
		strings.NewReader(`{ "dealerId": "a8f4e46c-8295-4f53-ab0a-7fc2d2f47d28", "ownerId": "23abce0a-ceb7-4127-8d98-b0bb5df4cce7" }`))
	require.NoError(t, err, "failed to create a request")
	w := httptest.NewRecorder()

	mockDbObj := new(tu.MockDB)
	mockDbObj.On("GetDealerByID", "a8f4e46c-8295-4f53-ab0a-7fc2d2f47d28").Return(dealer, nil)
	mockDbObj.On("GetUserByID", "23abce0a-ceb7-4127-8d98-b0bb5df4cce7").Return(nil, fmt.Errorf("failed to get user with id"))
	mockDbObj.On("CreateParcel", *dealer.ID, *user.ID).Return(parcel, nil)

	parcelCreateHandler(mockDbObj)(w, r)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestParcelCreateFail(t *testing.T) {
	r, err := http.NewRequest("POST", "/parcels",
		strings.NewReader(`{ "dealerId": "a8f4e46c-8295-4f53-ab0a-7fc2d2f47d28", "ownerId": "23abce0a-ceb7-4127-8d98-b0bb5df4cce7" }`))
	require.NoError(t, err, "failed to create a request")
	w := httptest.NewRecorder()

	mockDbObj := new(tu.MockDB)
	mockDbObj.On("GetDealerByID", "a8f4e46c-8295-4f53-ab0a-7fc2d2f47d28").Return(dealer, nil)
	mockDbObj.On("GetUserByID", "23abce0a-ceb7-4127-8d98-b0bb5df4cce7").Return(user, nil)
	mockDbObj.On("CreateParcel", *dealer.ID, *user.ID).Return(nil, fmt.Errorf("failed to create a parcel"))

	parcelCreateHandler(mockDbObj)(w, r)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

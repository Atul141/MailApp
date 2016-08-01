package handlers

import (
	"github.com/stretchr/testify/require"
	"net/http"
	"testing"

	"encoding/json"
	"fmt"
	m "git.mailbox.com/mailbox/models"
	tu "git.mailbox.com/mailbox/testutils"
	u "git.mailbox.com/mailbox/utils"
	"net/http/httptest"

	strfmt "github.com/go-openapi/strfmt"

	"github.com/stretchr/testify/assert"
)

var parcelWithUserAndDealer = []*m.ParcelUserDetails{
	&m.ParcelUserDetails{
		ID:          "ada1103c-4024-4ea4-b955-58c1c2c702b7",
		UserName:    u.SPtr("foobar"),
		Status:      false,
		OwnerID:     u.SPtr("123e4567-e89b-12d3-a456-426655440000"),
		RecieverID:  u.SPtr("123e4567-e89b-12d3-a456-426655440000"),
		UserEmail:   strfmt.Email("foobar@gmail.com"),
		UserEmpID:   "112312",
		UserPhoneNo: u.SPtr("12312312312"),
		DealerName:  u.SPtr("flipkart"),
	},
}

func TestParcelsHandlerForClosedParcelsSuccess(t *testing.T) {
	r, err := http.NewRequest("GET", "/parcels/close", nil)
	require.NoError(t, err, "failed to get a request: parcels")

	fmt.Println(fmt.Sprint(r))
	w := httptest.NewRecorder()

	mockDbObj := new(tu.MockDB)
	mockDbObj.On("GetCloseParcels").Return(parcelWithUserAndDealer, nil)

	parcelsHandler(mockDbObj)(w, r)

	var actualParcel []*m.ParcelUserDetails
	err = json.Unmarshal(w.Body.Bytes(), &actualParcel)
	require.NoError(t, err, "failed to unmarshal the response")

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, parcelWithUserAndDealer[0].ID, actualParcel[0].ID)
	assert.Equal(t, parcelWithUserAndDealer[0].UserName, actualParcel[0].UserName)
	assert.Equal(t, parcelWithUserAndDealer[0].UserPhoneNo, actualParcel[0].UserPhoneNo)
	assert.Equal(t, parcelWithUserAndDealer[0].UserEmpID, actualParcel[0].UserEmpID)
	assert.Equal(t, parcelWithUserAndDealer[0].UserEmail, actualParcel[0].UserEmail)
	assert.Equal(t, parcelWithUserAndDealer[0].DealerIcon, actualParcel[0].DealerIcon)

	mockDbObj.AssertExpectations(t)
}

func TestParcelsHandlerForOpenParcelsSuccess(t *testing.T) {
	r, err := http.NewRequest("GET", "/parcels/open", nil)
	require.NoError(t, err, "failed to get a request: parcels")

	fmt.Println(fmt.Sprint(r))
	w := httptest.NewRecorder()

	mockDbObj := new(tu.MockDB)
	mockDbObj.On("GetOpenParcels").Return(parcelWithUserAndDealer, nil)

	parcelsHandler(mockDbObj)(w, r)

	var actualParcel []*m.ParcelUserDetails
	err = json.Unmarshal(w.Body.Bytes(), &actualParcel)
	require.NoError(t, err, "failed to unmarshal the response")

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, parcelWithUserAndDealer[0].ID, actualParcel[0].ID)
	assert.Equal(t, parcelWithUserAndDealer[0].UserName, actualParcel[0].UserName)
	assert.Equal(t, parcelWithUserAndDealer[0].UserPhoneNo, actualParcel[0].UserPhoneNo)
	assert.Equal(t, parcelWithUserAndDealer[0].UserEmpID, actualParcel[0].UserEmpID)
	assert.Equal(t, parcelWithUserAndDealer[0].UserEmail, actualParcel[0].UserEmail)
	assert.Equal(t, parcelWithUserAndDealer[0].DealerIcon, actualParcel[0].DealerIcon)

	mockDbObj.AssertExpectations(t)
}

func TestParcelsHandlerForDBError(t *testing.T) {
	r, err := http.NewRequest("GET", "/parcels/close", nil)
	require.NoError(t, err, "failed to create a request: Parcels")
	w := httptest.NewRecorder()

	expectedErrorMessage := "some-db-specific-error"
	mockDbObj := new(tu.MockDB)
	mockDbObj.On("GetCloseParcels").Return(nil, fmt.Errorf(expectedErrorMessage))

	parcelsHandler(mockDbObj)(w, r)

	var actualError *m.Error
	err = json.Unmarshal(w.Body.Bytes(), &actualError)
	require.NoError(t, err, "failed to unmarshal the response")

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	require.NotNil(t, actualError)
	assert.Equal(t, int32(500), *actualError.Code)
	assert.Equal(t, "Internal server error", *actualError.Message)
	require.Nil(t, actualError.Fields)

	mockDbObj.AssertExpectations(t)
}

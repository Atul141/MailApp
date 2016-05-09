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

	r, err := http.NewRequest("GET", "/parcels",
		strings.NewReader(`{ "dealerId": "a8f4e46c-8295-4f53-ab0a-7fc2d2f47d28", "ownerId": "23abce0a-ceb7-4127-8d98-b0bb5df4cce7" }`))
	require.NoError(t, err, "failed to create a request")
	w := httptest.NewRecorder()

	mockDbObj := new(tu.MockDB)
	mockDbObj.On("GetDealerByID", "a8f4e46c-8295-4f53-ab0a-7fc2d2f47d28").Return(dealer, nil)
	mockDbObj.On("GetUserByID", "23abce0a-ceb7-4127-8d98-b0bb5df4cce7").Return(user, nil)
	mockDbObj.On("CreateParcel", *dealer.ID, *user.ID).Return(parcel, nil)

	parcelCreateHandler(mockDbObj)(w, r)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, fmt.Sprintf(`{"id":"%s"}`, parcel.ID), w.Body.String())
}

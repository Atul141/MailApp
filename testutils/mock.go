package testutils

import(
  "github.com/stretchr/testify/mock"

	m "git.mailbox.com/mailbox/models"
)

type MockDB struct {
	mock.Mock
}

func (db *MockDB) GetDealers() ([]*m.Dealer, error) {
	args := db.Called()
	if args.Get(0) != nil {
		return args.Get(0).([]*m.Dealer), nil
	}
	return nil, args.Error(1)
}

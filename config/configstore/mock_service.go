// +build unit

package configstore

import (
	"github.com/centrifuge/go-centrifuge/config"
	"github.com/stretchr/testify/mock"
)

type MockService struct {
	mock.Mock
}

func (m MockService) GenerateAccount() (config.Account, error) {
	args := m.Called()
	return args.Get(0).(config.Account), args.Error(1)
}

func (m MockService) GetConfig() (config.Configuration, error) {
	args := m.Called()
	return args.Get(0).(*NodeConfig), args.Error(1)
}

func (m MockService) GetAccount(identifier []byte) (config.Account, error) {
	args := m.Called(identifier)
	return args.Get(0).(config.Account), args.Error(1)
}

func (m MockService) GetAllAccounts() ([]config.Account, error) {
	args := m.Called()
	v, _ := args.Get(0).([]config.Account)
	return v, nil
}

func (m MockService) CreateConfig(data config.Configuration) (config.Configuration, error) {
	args := m.Called(data)
	return args.Get(0).(*NodeConfig), args.Error(0)
}

func (m MockService) CreateAccount(data config.Account) (config.Account, error) {
	args := m.Called(data)
	return args.Get(0).(*Account), args.Error(0)
}

func (m MockService) UpdateAccount(data config.Account) (config.Account, error) {
	args := m.Called(data)
	return args.Get(0).(*Account), args.Error(0)
}

func (m MockService) DeleteAccount(identifier []byte) error {
	args := m.Called(identifier)
	return args.Error(0)
}

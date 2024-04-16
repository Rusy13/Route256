package postgresql

import (
	"go.uber.org/mock/gomock"
	"testing"

	mock_database "Homework/internal/storage/db/mocks"
	"Homework/internal/storage/repository"
)

type pvzRepoFixtures struct {
	ctrl   *gomock.Controller
	repo   repository.PvzRepo
	mockDB *mock_database.MockDBops
}

func SetUp(t *testing.T) pvzRepoFixtures {
	ctrl := gomock.NewController(t)
	mockDB := mock_database.NewMockDBops(ctrl)
	repo := NewPvzRepo(mockDB)
	return pvzRepoFixtures{
		ctrl:   ctrl,
		repo:   repo,
		mockDB: mockDB,
	}
}

func (a *pvzRepoFixtures) TearDown() {
	a.ctrl.Finish()
}

package postgresql

import (
	mock_database "HW1/internal/storage/db/mocks"
	"HW1/internal/storage/repository"
	"go.uber.org/mock/gomock"
	"testing"
)

type pvzRepoFixtures struct {
	ctrl   *gomock.Controller
	repo   repository.PvzRepo
	mockDB *mock_database.MockDBops
}

func setUp(t *testing.T) pvzRepoFixtures {
	ctrl := gomock.NewController(t)
	mockDB := mock_database.NewMockDBops(ctrl)
	repo := NewArticles(mockDB)
	return pvzRepoFixtures{
		ctrl:   ctrl,
		repo:   repo,
		mockDB: mockDB,
	}
}

func (a *pvzRepoFixtures) tearDown() {
	a.ctrl.Finish()
}

package postgresql

import (
	"go.uber.org/mock/gomock"
	"testing"

	mock_repository "Homework/internal/storage/repository/mocks"
)

type pvzRepoFixtures struct {
	ctrl    *gomock.Controller
	srv     Server1
	mockPvz *mock_repository.MockPvzRepo
}

func setUp(t *testing.T) pvzRepoFixtures {
	ctrl := gomock.NewController(t)
	mockPvz := mock_repository.NewMockPvzRepo(ctrl)
	srv := Server1{mockPvz}
	return pvzRepoFixtures{
		ctrl:    ctrl,
		mockPvz: mockPvz,
		srv:     srv,
	}
}

func (a *pvzRepoFixtures) tearDown() {
	a.ctrl.Finish()
}

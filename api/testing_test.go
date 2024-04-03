package postgresql

import (
	mock_repository "HW1/internal/storage/repository/mocks"
	"go.uber.org/mock/gomock"
	"testing"
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

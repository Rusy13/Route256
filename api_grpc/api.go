package postgresql

import (
	"Homework/internal/storage/repository"
	pb "Homework/protos/gen/go/app"
	"context"
	"errors"
	"fmt"
)

type Server struct {
	Repo repository.PvzRepo
	pb.UnimplementedPvzServiceServer
}

func (s *Server) CreatePvz(ctx context.Context, req *pb.CreatePvzRequest) (*pb.CreatePvzResponse, error) {
	pvzRepo := &repository.Pvz{
		PvzName: req.Pvzname,
		Address: req.Address,
		Email:   req.Email,
	}
	id, err := s.Repo.Add(ctx, pvzRepo)
	if err != nil {
		return nil, fmt.Errorf("failed to add pvz: %v", err)
	}

	return &pb.CreatePvzResponse{
		Id:      id,
		Pvzname: pvzRepo.PvzName,
		Address: pvzRepo.Address,
		Email:   pvzRepo.Email,
	}, nil
}

func (s *Server) GetPvzByID(ctx context.Context, req *pb.GetPvzByIDRequest) (*pb.GetPvzByIDResponse, error) {
	key := req.Key

	pvz, err := s.Repo.GetByID(ctx, key)
	if err != nil {
		if errors.Is(err, repository.ErrObjectNotFound) {
			return nil, fmt.Errorf("pvz with key %v not found", key)
		}
		return nil, fmt.Errorf("failed to get pvz: %v", err)
	}

	return &pb.GetPvzByIDResponse{
		Pvzname: pvz.PvzName,
		Address: pvz.Address,
		Email:   pvz.Email,
	}, nil
}

func (s *Server) UpdatePvz(ctx context.Context, req *pb.UpdatePvzRequest) (*pb.UpdatePvzResponse, error) {
	key := req.Key

	updatedPvz := &repository.Pvz{
		ID:      key,
		PvzName: req.Pvzname,
		Address: req.Address,
		Email:   req.Email,
	}

	if err := s.Repo.Update(ctx, key, updatedPvz); err != nil {
		return nil, fmt.Errorf("failed to update pvz: %v", err)
	}

	return &pb.UpdatePvzResponse{
		Id:      key,
		Pvzname: req.Pvzname,
		Address: req.Address,
		Email:   req.Email,
	}, nil
}

func (s *Server) DeletePvz(ctx context.Context, req *pb.DeletePvzRequest) (*pb.DeletePvzResponse, error) {
	key := req.Key

	if err := s.Repo.Delete(ctx, key); err != nil {
		return nil, fmt.Errorf("failed to delete pvz: %v", err)
	}

	return &pb.DeletePvzResponse{
		Message: "Successfully deleted",
	}, nil
}

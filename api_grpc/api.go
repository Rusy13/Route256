package postgresql

import (
	"context"
	"errors"
	"fmt"
	"github.com/opentracing/opentracing-go"
	"time"

	"Homework/internal/storage/repository"
	metrics "Homework/metrics/metrics"
	pb "Homework/protos/gen/go/app"
)

type Server struct {
	Repo repository.PvzRepo
	pb.UnimplementedPvzServiceServer
}

func (s *Server) CreatePvz(ctx context.Context, req *pb.CreatePvzRequest) (*pb.CreatePvzResponse, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "CreatePvz")
	defer span.Finish() // Завершаем спан после выполнения функции

	metrics.OrdersCounter.Inc()
	metrics.OrdersInProgress.Inc()

	start := time.Now() // Начало измерения времени обработки

	pvzRepo := &repository.Pvz{
		PvzName: req.Pvzname,
		Address: req.Address,
		Email:   req.Email,
	}
	id, err := s.Repo.Add(ctx, pvzRepo)
	if err != nil {
		metrics.OrdersInProgress.Dec()
		return nil, fmt.Errorf("failed to add pvz: %v", err)
	}

	metrics.ProcessingHistogram.Observe(time.Since(start).Seconds())
	metrics.OrdersInProgress.Dec()

	return &pb.CreatePvzResponse{
		Id:      id,
		Pvzname: pvzRepo.PvzName,
		Address: pvzRepo.Address,
		Email:   pvzRepo.Email,
	}, nil
}

func (s *Server) GetPvzByID(ctx context.Context, req *pb.GetPvzByIDRequest) (*pb.GetPvzByIDResponse, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "GetPvzByID")
	defer span.Finish() // Завершаем спан после выполнения функции

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
	span, ctx := opentracing.StartSpanFromContext(ctx, "UpdatePvz")
	defer span.Finish() // Завершаем спан после выполнения функции

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
	span, ctx := opentracing.StartSpanFromContext(ctx, "DeletePvz")
	defer span.Finish() // Завершаем спан после выполнения функции

	key := req.Key

	if err := s.Repo.Delete(ctx, key); err != nil {
		return nil, fmt.Errorf("failed to delete pvz: %v", err)
	}

	return &pb.DeletePvzResponse{
		Message: "Successfully deleted",
	}, nil
}

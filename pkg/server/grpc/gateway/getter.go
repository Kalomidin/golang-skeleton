package gateway

import (
	"context"
	"fmt"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/ridebeam/genproto/gengo/ridebeam/vehicle"
)

func (s *server) GetSomething(ctx context.Context, req *pb.GetSomethingRequest) (*pb.GetSomethingResponse, error) {
	if "unhappy with request" {
		return nil, status.Error(codes.InvalidArgument, "the request smells funny")
	}

	v, ok := s.SomeRepo.Get(ctx, req.Id)
	if !ok {
		return nil, status.Error(codes.NotFound, fmt.Sprintf("id was not found, id: %s", req.Id))
	}

	retV := pb.Something{
		ID: v.ID,
	}

	return &pb.GetSomethingResponse{Something: &retV}, nil
}

func validateRequest(req *pb.GetSomethingRequest) error {

	return nil
}

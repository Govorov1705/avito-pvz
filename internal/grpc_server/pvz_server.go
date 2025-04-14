package server

import (
	"context"

	"github.com/Govorov1705/avito-pvz/internal/repositories"
	pvz_v1 "github.com/Govorov1705/avito-pvz/pkg/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type PvzServer struct {
	pvz_v1.UnimplementedPVZServiceServer
	PvzRepository repositories.PvzRepository
}

func NewPvzServer(r repositories.PvzRepository) *PvzServer {
	return &PvzServer{PvzRepository: r}
}

func (s *PvzServer) GetPVZList(ctx context.Context, req *pvz_v1.GetPVZListRequest) (*pvz_v1.GetPVZListResponse, error) {
	pvzs, err := s.PvzRepository.List(ctx)
	if err != nil {
		return nil, err
	}

	pbPvzs := make([]*pvz_v1.PVZ, len(pvzs))
	for i, pvz := range pvzs {
		pbPvzs[i] = &pvz_v1.PVZ{
			Id:               pvz.ID.String(),
			RegistrationDate: timestamppb.New(pvz.RegistrationDate),
			City:             string(pvz.City),
		}
	}

	return &pvz_v1.GetPVZListResponse{
		Pvzs: pbPvzs,
	}, nil
}

package server

import (
	"context"
	"fmt"
	"net"

	"github.com/go-pg/pg/v9"
	"github.com/go-redis/redis/v7"
	pricingpb "github.com/weslenng/petssenger/protos"
	"github.com/weslenng/petssenger/services/pricing/models"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type pricingServer struct {
	pg    *pg.DB
	redis *redis.Client
}

const addr = "0.0.0.0:50051"

func (ps *pricingServer) GetPricingFeesByCity(ctx context.Context, req *pricingpb.GetPricingFeesByCityRequest) (*pricingpb.GetPricingFeesByCityResponse, error) {
	city := req.GetCity()
	fees, err := models.GetPricingFees(city, ps.pg, ps.redis)
	if err != nil {
		if err == pg.ErrNoRows {
			return nil, status.Errorf(
				codes.NotFound,
				fmt.Sprintf("The city %v was not found", city),
			)
		}

		panic(err)
	}

	proto := models.ProtoPricingFees(fees)
	return proto, nil
}

// PricingServerListen is a helper function to lis and gRPC server
func PricingServerListen(pg *pg.DB, redis *redis.Client) (net.Listener, *grpc.Server, error) {
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, nil, err
	}

	ser := grpc.NewServer()
	ps := &pricingServer{
		pg:    pg,
		redis: redis,
	}

	pricingpb.RegisterPricingServer(ser, ps)
	if err := ser.Serve(lis); err != nil {
		lis.Close()
		return nil, nil, err
	}

	return lis, ser, nil
}

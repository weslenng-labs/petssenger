package models

import (
	"time"

	"github.com/go-pg/pg/v9"
	"github.com/go-redis/redis/v7"
	"github.com/vmihailenco/msgpack/v4"
	pricingpb "github.com/weslenng/petssenger/protos"
)

// Fees represents a city-fees payload structure
type Fees struct {
	ID       string
	Base     float32
	Distance float32
	Minute   float32
	Service  float32
}

// ProtoPricingFees transforms the type Fees in a protobuf message
func ProtoPricingFees(fees *Fees) *pricingpb.GetPricingFeesByCityResponse {
	return &pricingpb.GetPricingFeesByCityResponse{
		Id:       fees.ID,
		Base:     fees.Base,
		Distance: fees.Distance,
		Minute:   fees.Minute,
		Service:  fees.Service,
	}
}

// GetPricingFees retrieve the ride fees by a given city
func GetPricingFees(ID string, pg *pg.DB, redis *redis.Client) (*Fees, error) {
	fees := &Fees{ID: ID}

	val, err := redis.Get(ID).Bytes()
	if err == nil {
		err = msgpack.Unmarshal(val, fees)
		if err == nil {
			return fees, nil
		}
	}

	err = pg.Select(fees)
	if err != nil {
		return nil, err
	}

	val, err = msgpack.Marshal(fees)
	if err != nil {
		return nil, err
	}

	err = redis.Set(ID, val, 5*time.Minute).Err()
	if err != nil {
		return nil, err
	}

	return fees, nil
}

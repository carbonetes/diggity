package types

import (
	"github.com/golistic/urn"
	"github.com/google/uuid"
)

type Payload struct {
	Address *urn.URN
	Layer   string
	Body    interface{}
}

func NewAddress() (*urn.URN, error) {
	addr, err := urn.New("bom", "diggity", urn.WithQuery(uuid.NewString()))
	if err != nil {
		return nil, err
	}

	return addr, nil
}

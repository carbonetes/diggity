package types

import (
	"fmt"

	"github.com/golistic/urn"
	"github.com/google/uuid"
)

type Payload struct {
	Address *urn.URN
	Layer   string
	Body    interface{}
}

func NewAddress(input string) (*urn.URN, error) {

	if input == "" {
		return nil, fmt.Errorf("input is empty")
	}

	addr, err := urn.New("bom", input, urn.WithQuery(uuid.NewString()))
	if err != nil {
		return nil, err
	}

	return addr, nil
}

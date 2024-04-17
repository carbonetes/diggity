package types

import (
	"fmt"

	"github.com/golistic/urn"
	"github.com/google/uuid"
)

type Payload struct {
	Address Address
	Body    interface{}
}

type Address urn.URN

func NewAddress(input string) (Address, error) {
	addr, err := urn.New("bom", input, urn.WithQuery(uuid.NewString()))
	if err != nil {
		return Address(*addr), err
	}

	return Address(*addr), nil
}

func (a Address) ToString() string {
	return fmt.Sprintf("%s", a)
}

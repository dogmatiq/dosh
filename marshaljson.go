package dosh

import (
	"encoding/json"

	"google.golang.org/genproto/googleapis/type/money"
)

// MarshalJSON mashals an amount to its JSON representation.
func (a Amount) MarshalJSON() ([]byte, error) {
	pb, err := a.MarshalProto()
	if err != nil {
		return nil, err
	}

	return json.Marshal(pb)
}

// UnmarshalJSON unmarshals an amount from its protocol buffers representation.
func (a *Amount) UnmarshalJSON(data []byte) error {
	var pb money.Money

	if err := json.Unmarshal(data, &pb); err != nil {
		return err
	}

	return a.UnmarshalProto(&pb)
}

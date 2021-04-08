package dosh

import (
	"fmt"

	"google.golang.org/genproto/googleapis/type/money"
	"google.golang.org/protobuf/encoding/protojson"
)

// jsonMarshaler is the marshaler used to encode a money.Money to JSON.
var jsonMarshaler = protojson.MarshalOptions{
	UseProtoNames: true,
}

// jsonUnmarshaler is the text marshaler used by JSONCodec if none is
// provided.
var jsonUnmarshaler = protojson.UnmarshalOptions{
	DiscardUnknown: false,
}

// MarshalJSON mashals an amount to its JSON representation.
//
// It uses the canonical JSON format of the protocol buffers money.Money type.
func (a Amount) MarshalJSON() ([]byte, error) {
	pb, err := a.marshalProto()
	if err != nil {
		return nil, fmt.Errorf("cannot marshal amount to JSON representation: %w", err)
	}

	data, err := jsonMarshaler.Marshal(pb)
	if err != nil {
		// CODE COVERAGE: It does not appear this branch can currently be
		// reached as it's not clear how to make jsonMarshaler fail with this
		// configuration.
		return nil, fmt.Errorf("cannot marshal amount to JSON representation: %w", err)
	}

	return data, nil
}

// UnmarshalJSON unmarshals an amount from its protocol buffers representation.
//
// NOTE: In order to comply with Go's json.Unmarshaler interface, this method
// mutates the internals of a, violating Amount's immutability guarantee.
func (a *Amount) UnmarshalJSON(data []byte) error {
	var pb money.Money

	if err := jsonUnmarshaler.Unmarshal(data, &pb); err != nil {
		return fmt.Errorf("cannot unmarshal amount from JSON representation: %w", err)
	}

	if err := a.unmarshalProto(&pb); err != nil {
		return fmt.Errorf("cannot unmarshal amount from JSON representation: %w", err)
	}

	return nil
}

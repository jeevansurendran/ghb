package ghb

import (
	"fmt"
	"strconv"
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"
)

type timeFormat struct {
	unmarshal func(any) (*timestamppb.Timestamp, error)
	marshal   func(*timestamppb.Timestamp) (any, error)
}

var EpochTimeFormat = &timeFormat{
	unmarshal: func(input any) (*timestamppb.Timestamp, error) {
		s, ok := input.(string)
		if !ok {
			return nil, fmt.Errorf("expected string for timestamp, got %T", s)
		}
		seconds, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			return nil, err
		}
		return timestamppb.New(time.Unix(seconds, 0)), nil
	},
	marshal: func(t *timestamppb.Timestamp) (any, error) {
		return strconv.FormatInt(t.Seconds, 10), nil
	},
}

var ISOTimeFormat = &timeFormat{
	unmarshal: func(input any) (*timestamppb.Timestamp, error) {
		s, ok := input.(string)
		if !ok {
			return nil, fmt.Errorf("expected string for timestamp, got %T", s)
		}
		t, err := time.Parse(time.RFC3339, s)
		if err != nil {
			return nil, err
		}
		return &timestamppb.Timestamp{Seconds: t.Unix()}, nil
	},
	marshal: func(t *timestamppb.Timestamp) (any, error) {
		return t.AsTime().Format(time.RFC3339), nil
	},
}

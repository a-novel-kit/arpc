package arpcdata

import (
	"fmt"

	"google.golang.org/protobuf/types/known/structpb"
)

func StructOptional(src map[string]interface{}) (*structpb.Struct, error) {
	if src == nil {
		return nil, nil //nolint:nilnil
	}

	res, err := structpb.NewStruct(src)
	if err != nil {
		return nil, fmt.Errorf("convert map to struct: %w", err)
	}

	return res, nil
}

func StructOptionalProto(src *structpb.Struct) map[string]interface{} {
	if src == nil {
		return nil
	}

	return src.AsMap()
}

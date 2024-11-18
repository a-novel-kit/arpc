package arpcdata_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	arpcdata "github.com/a-novel-kit/arpc/data"
)

func TestStructOptional(t *testing.T) {
	structure, err := arpcdata.StructOptional(map[string]interface{}{
		"key": "value",
	})
	require.Nil(t, err)
	require.NotNil(t, structure)
	require.Equal(t, "value", structure.Fields["key"].GetStringValue())

	structure, err = arpcdata.StructOptional(nil)
	require.Nil(t, err)
	require.Nil(t, structure)
}

func TestStructOptionalProto(t *testing.T) {
	initial, err := arpcdata.StructOptional(map[string]interface{}{
		"key": "value",
	})
	require.Nil(t, err)

	structure := arpcdata.StructOptionalProto(initial)
	require.NotNil(t, structure)
	require.Equal(t, "value", structure["key"])

	structure = arpcdata.StructOptionalProto(nil)
	require.Nil(t, structure)
}

package arpcdata_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	arpcdata "github.com/a-novel-kit/arpc/data"
)

func TestProtoConverterFromProto(t *testing.T) {
	testCases := []struct {
		name string

		source string

		mapper        arpcdata.ProtoMapper[string, int]
		protoDefault  string
		entityDefault int

		expect int
	}{
		{
			name: "OK",

			source: "one",

			mapper: arpcdata.ProtoMapper[string, int]{
				"one":   1,
				"two":   2,
				"three": 3,
			},

			protoDefault:  "zero",
			entityDefault: 0,

			expect: 1,
		},
		{
			name: "Default",

			source: "four",

			mapper: arpcdata.ProtoMapper[string, int]{
				"one":   1,
				"two":   2,
				"three": 3,
			},

			protoDefault:  "zero",
			entityDefault: 0,

			expect: 0,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			converter := arpcdata.NewProtoConverter(testCase.mapper, testCase.protoDefault, testCase.entityDefault)
			require.Equal(t, testCase.expect, converter.FromProto(testCase.source))
		})
	}
}

func TestProtoConverterToProto(t *testing.T) {
	testCases := []struct {
		name string

		source int

		mapper        arpcdata.ProtoMapper[string, int]
		protoDefault  string
		entityDefault int

		expect string
	}{
		{
			name: "OK",

			source: 1,

			mapper: arpcdata.ProtoMapper[string, int]{
				"one":   1,
				"two":   2,
				"three": 3,
			},

			protoDefault:  "zero",
			entityDefault: 0,

			expect: "one",
		},
		{
			name: "Default",

			source: 4,

			mapper: arpcdata.ProtoMapper[string, int]{
				"one":   1,
				"two":   2,
				"three": 3,
			},

			protoDefault:  "zero",
			entityDefault: 0,

			expect: "zero",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			converter := arpcdata.NewProtoConverter(testCase.mapper, testCase.protoDefault, testCase.entityDefault)
			require.Equal(t, testCase.expect, converter.ToProto(testCase.source))
		})
	}
}

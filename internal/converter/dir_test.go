package converter_test

import (
	"errors"
	"testing"

	"github.com/convert/internal/converter"
	"github.com/convert/pkg/types"
)

var errFake = errors.New("fake error")

func TestSummaryLine(t *testing.T) {
	cases := []struct {
		name     string
		results  []types.ConvertResult
		expected string
	}{
		{
			"all ok",
			[]types.ConvertResult{{}, {}, {}},
			"conversion complete: 3 succeeded, 0 failed",
		},
		{
			"mixed",
			[]types.ConvertResult{{}, {Error: errFake}, {}},
			"conversion complete: 2 succeeded, 1 failed",
		},
		{
			"all fail",
			[]types.ConvertResult{{Error: errFake}, {Error: errFake}},
			"conversion complete: 0 succeeded, 2 failed",
		},
		{
			"empty",
			[]types.ConvertResult{},
			"conversion complete: 0 succeeded, 0 failed",
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := converter.SummaryLine(c.results)
			if got != c.expected {
				t.Errorf("got %q, want %q", got, c.expected)
			}
		})
	}
}

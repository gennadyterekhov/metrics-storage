package errcheck

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/tools/go/analysis/analysistest"
)

func TestErrcheck(t *testing.T) {
	reportsUnchecked(t)
	doesNotReportChecked(t)
}

func reportsUnchecked(t *testing.T) {
	var result analysistest.Result
	results := analysistest.Run(t, analysistest.TestData(), New(), "./unchecked")

	assert.Equal(t, 1, len(results))
	result = *results[0]
	assert.NoError(t, result.Err)

	assert.Equal(t, 4, len(result.Diagnostics))
}

func doesNotReportChecked(t *testing.T) {
	var result analysistest.Result
	results := analysistest.Run(t, analysistest.TestData(), New(), "./checked")

	assert.Equal(t, 1, len(results))
	result = *results[0]
	assert.NoError(t, result.Err)

	assert.Equal(t, 0, len(result.Diagnostics))
}

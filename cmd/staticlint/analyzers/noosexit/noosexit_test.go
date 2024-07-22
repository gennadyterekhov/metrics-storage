package noosexit

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/tools/go/analysis/analysistest"
)

func TestNoOsExit(t *testing.T) {
	noOsExitInMain(t)
	canUseOsExitInAnotherPackage(t)
}

func noOsExitInMain(t *testing.T) {
	var result analysistest.Result
	results := analysistest.Run(t, analysistest.TestData(), New(), "./pkg1")

	assert.Equal(t, 1, len(results))
	result = *results[0]
	assert.NoError(t, result.Err)

	assert.Equal(t, 1, len(result.Diagnostics))
	assert.Equal(t, "cannot use os Exit in main function of package main", result.Diagnostics[0].Message)
}

func canUseOsExitInAnotherPackage(t *testing.T) {
	var result analysistest.Result
	results := analysistest.Run(t, analysistest.TestData(), New(), "./notmainpackage")

	assert.Equal(t, 1, len(results))
	result = *results[0]
	assert.NoError(t, result.Err)

	assert.Equal(t, 0, len(result.Diagnostics))
}

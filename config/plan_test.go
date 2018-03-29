// +build unit

package config_test

import (
	"testing"

	"github.com/stefanprodan/mgob/config"
)

func assertError(t *testing.T, err error) {
	t.Log(err)
	if err == nil {
		t.Error(err)
	}
}

func assertNoError(t *testing.T, err error) {
	t.Log(err)
	if err != nil {
		t.Error(err)
	}
}

func TestPlanReturnErrorOnInvalidPath(t *testing.T) {
	planDir := "./"
	planName := "test.yaml"
	_, err := config.LoadPlan(planDir, planName)
	assertError(t, err)
}

func TestPlanReturnErrorOnInvalidYAML(t *testing.T) {
	planDir := "./resources"
	planName := "wrong_example"
	_, err := config.LoadPlan(planDir, planName)
	assertError(t, err)
}

func TestShouldValidatePlanYAMLCorrectly(t *testing.T) {
	planDir := "./resources"
	planName := "correctly_example"
	_, err := config.LoadPlan(planDir, planName)
	assertNoError(t, err)
}

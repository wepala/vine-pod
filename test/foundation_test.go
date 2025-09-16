package test

import (
	"context"
	"os"
	"testing"

	"github.com/cucumber/godog"

	"github.com/wepala/vine-pod/test/steps"
)

func TestFoundationFeatures(t *testing.T) {
	suite := godog.TestSuite{
		ScenarioInitializer: func(sc *godog.ScenarioContext) {
			// Initialize foundation steps
			foundationSteps := steps.NewFoundationSteps(t)
			steps.RegisterFoundationSteps(sc, foundationSteps)
		},
		Options: &godog.Options{
			Format:    "pretty",
			Paths:     []string{"features/foundation.feature"},
			Output:    os.Stdout,
			Randomize: 0, // Consistent test ordering
			Strict:    true,
			TestingT:  t,
		},
	}

	if suite.Run() != 0 {
		t.Fatal("Non-zero status returned, failed to run foundation BDD tests")
	}
}

func TestMain(m *testing.M) {
	// Setup and teardown for all tests can go here
	// For now, we'll just run the tests
	ctx := context.Background()
	_ = ctx // Use context if needed for test setup

	// Run tests
	m.Run()
}

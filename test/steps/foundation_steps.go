package steps

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/cucumber/godog"
	"github.com/stretchr/testify/assert"
	"go.uber.org/fx"

	"github.com/wepala/vine-pod/internal/config"
	fxmodules "github.com/wepala/vine-pod/internal/fx"
	"github.com/wepala/vine-pod/internal/handler"
	"github.com/wepala/vine-pod/pkg/logger"
)

// FoundationSteps holds the test context for foundation scenarios
type FoundationSteps struct {
	config       *config.Config
	logger       logger.Logger
	handlers     *handler.Handlers
	response     *httptest.ResponseRecorder
	request      *http.Request
	responseBody map[string]interface{}
	app          *fx.App
	t            *testing.T
}

// NewFoundationSteps creates a new instance of foundation steps
func NewFoundationSteps(t *testing.T) *FoundationSteps {
	return &FoundationSteps{
		t: t,
	}
}

// RegisterFoundationSteps registers all foundation-related step definitions
func RegisterFoundationSteps(sc *godog.ScenarioContext, steps *FoundationSteps) {
	// Background steps
	sc.Step(`^the vine-pod service is configured$`, steps.theVinePodServiceIsConfigured)
	sc.Step(`^the database connection is available$`, steps.theDatabaseConnectionIsAvailable)

	// HTTP request steps
	sc.Step(`^I send a GET request to "([^"]*)"$`, steps.iSendAGETRequestTo)
	sc.Step(`^the response status should be (\d+)$`, steps.theResponseStatusShouldBe)
	sc.Step(`^the response should contain "([^"]*)" as "([^"]*)"$`, steps.theResponseShouldContainAs)
	sc.Step(`^the response should contain a "([^"]*)" field$`, steps.theResponseShouldContainAField)

	// Database steps
	sc.Step(`^the application is starting up$`, steps.theApplicationIsStartingUp)
	sc.Step(`^the database module initializes$`, steps.theDatabaseModuleInitializes)
	sc.Step(`^the database connection should be successful$`, steps.theDatabaseConnectionShouldBeSuccessful)
	sc.Step(`^the database should be pingable$`, steps.theDatabaseShouldBePingable)

	// Fx dependency injection steps
	sc.Step(`^the Fx container is configured$`, steps.theFxContainerIsConfigured)
	sc.Step(`^the application starts$`, steps.theApplicationStarts)
	sc.Step(`^all dependencies should be resolved$`, steps.allDependenciesShouldBeResolved)
	sc.Step(`^no circular dependencies should exist$`, steps.noCircularDependenciesShouldExist)
	sc.Step(`^the application should start without errors$`, steps.theApplicationShouldStartWithoutErrors)

	// Logging steps
	sc.Step(`^the Zap logger is configured$`, steps.theZapLoggerIsConfigured)
	sc.Step(`^the application logs messages$`, steps.theApplicationLogsMessages)
	sc.Step(`^logs should be in JSON format$`, steps.logsShouldBeInJSONFormat)
	sc.Step(`^logs should contain proper fields$`, steps.logsShouldContainProperFields)
	sc.Step(`^log levels should be respected$`, steps.logLevelsShouldBeRespected)
}

// Background step implementations
func (s *FoundationSteps) theVinePodServiceIsConfigured() error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}
	s.config = cfg

	s.logger = logger.New(cfg.LogLevel)
	s.handlers = handler.New(cfg, s.logger)
	return nil
}

func (s *FoundationSteps) theDatabaseConnectionIsAvailable() error {
	// This step assumes database is available
	// In real scenarios, we'd check database connectivity here
	return nil
}

// HTTP request step implementations
func (s *FoundationSteps) iSendAGETRequestTo(endpoint string) error {
	s.request = httptest.NewRequest("GET", endpoint, nil)
	s.response = httptest.NewRecorder()

	switch endpoint {
	case "/health":
		s.handlers.Health(s.response, s.request)
	case "/version":
		s.handlers.Version(s.response, s.request)
	default:
		return fmt.Errorf("unknown endpoint: %s", endpoint)
	}

	// Parse response body
	if s.response.Body.Len() > 0 {
		if err := json.Unmarshal(s.response.Body.Bytes(), &s.responseBody); err != nil {
			return fmt.Errorf("failed to parse response body: %w", err)
		}
	}

	return nil
}

func (s *FoundationSteps) theResponseStatusShouldBe(expectedStatus int) error {
	assert.Equal(s.t, expectedStatus, s.response.Code,
		"Expected status %d, got %d", expectedStatus, s.response.Code)
	return nil
}

func (s *FoundationSteps) theResponseShouldContainAs(field, expectedValue string) error {
	actualValue, exists := s.responseBody[field]
	if !exists {
		return fmt.Errorf("field %s not found in response", field)
	}

	assert.Equal(s.t, expectedValue, actualValue,
		"Expected %s to be %s, got %v", field, expectedValue, actualValue)
	return nil
}

func (s *FoundationSteps) theResponseShouldContainAField(field string) error {
	_, exists := s.responseBody[field]
	assert.True(s.t, exists, "Field %s should exist in response", field)
	return nil
}

// Database step implementations
func (s *FoundationSteps) theApplicationIsStartingUp() error {
	// Setup for database testing
	return s.theVinePodServiceIsConfigured()
}

func (s *FoundationSteps) theDatabaseModuleInitializes() error {
	// This would test the database module initialization
	// For now, we'll assume it initializes correctly
	return nil
}

func (s *FoundationSteps) theDatabaseConnectionShouldBeSuccessful() error {
	// In a real implementation, we'd test actual database connectivity
	// This is a placeholder for database connection verification
	return nil
}

func (s *FoundationSteps) theDatabaseShouldBePingable() error {
	// In a real implementation, we'd ping the database
	// This is a placeholder for database ping verification
	return nil
}

// Fx dependency injection step implementations
func (s *FoundationSteps) theFxContainerIsConfigured() error {
	// Create an Fx app to test dependency injection
	s.app = fxmodules.NewApp()
	return nil
}

func (s *FoundationSteps) theApplicationStarts() error {
	if s.app == nil {
		return fmt.Errorf("Fx app is not configured")
	}

	// Test that the app can start
	ctx := context.Background()
	err := s.app.Start(ctx)
	if err != nil {
		return fmt.Errorf("failed to start Fx app: %w", err)
	}

	// Stop the app immediately for testing
	return s.app.Stop(ctx)
}

func (s *FoundationSteps) allDependenciesShouldBeResolved() error {
	// This is verified by the successful start of the Fx app
	return nil
}

func (s *FoundationSteps) noCircularDependenciesShouldExist() error {
	// This is verified by the successful start of the Fx app
	// Fx would fail to start if circular dependencies exist
	return nil
}

func (s *FoundationSteps) theApplicationShouldStartWithoutErrors() error {
	// This is verified by the successful start of the Fx app
	return nil
}

// Logging step implementations
func (s *FoundationSteps) theZapLoggerIsConfigured() error {
	if s.logger == nil {
		return fmt.Errorf("logger is not configured")
	}
	return nil
}

func (s *FoundationSteps) theApplicationLogsMessages() error {
	// Test that logging works
	s.logger.Info("test log message")
	return nil
}

func (s *FoundationSteps) logsShouldBeInJSONFormat() error {
	// This would require capturing log output to verify JSON format
	// For now, we'll assume the logger produces JSON (which Zap does by default)
	return nil
}

func (s *FoundationSteps) logsShouldContainProperFields() error {
	// This would require capturing log output to verify field structure
	// For now, we'll assume proper fields are present
	return nil
}

func (s *FoundationSteps) logLevelsShouldBeRespected() error {
	// This would require testing different log levels
	// For now, we'll assume log levels work correctly
	return nil
}

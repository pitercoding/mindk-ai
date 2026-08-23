package config

import "testing"

// fakeEnv builds a getenv func backed by a map, returning "" for unset keys
// (matching os.Getenv's behavior).
func fakeEnv(values map[string]string) func(string) string {
	return func(key string) string {
		return values[key]
	}
}

func validProductionEnv() map[string]string {
	return map[string]string{
		"APP_ENV":          "production",
		"OPENAI_API_KEY":   "sk-test",
		"CLERK_SECRET_KEY": "clerk-test",
		"FRONTEND_ORIGIN":  "https://app.example.com",
		"DATABASE_PATH":    "/var/lib/mindk/mindk.db",
	}
}

func TestFromEnv_DevelopmentDefaults(t *testing.T) {
	cfg, err := FromEnv(fakeEnv(map[string]string{
		"OPENAI_API_KEY":   "sk-test",
		"CLERK_SECRET_KEY": "clerk-test",
	}))
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if cfg.Environment != EnvDevelopment {
		t.Errorf("expected default environment %q, got %q", EnvDevelopment, cfg.Environment)
	}
	if cfg.FrontendOrigin != defaultFrontendOrigin {
		t.Errorf("expected default frontend origin %q, got %q", defaultFrontendOrigin, cfg.FrontendOrigin)
	}
	if cfg.DatabasePath != defaultDatabasePath {
		t.Errorf("expected default database path %q, got %q", defaultDatabasePath, cfg.DatabasePath)
	}
	if cfg.IsProduction() {
		t.Error("expected IsProduction() to be false in development")
	}
}

func TestFromEnv_ValidProduction(t *testing.T) {
	cfg, err := FromEnv(fakeEnv(validProductionEnv()))
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if !cfg.IsProduction() {
		t.Error("expected IsProduction() to be true")
	}
	if cfg.FrontendOrigin != "https://app.example.com" {
		t.Errorf("unexpected frontend origin %q", cfg.FrontendOrigin)
	}
	if cfg.DatabasePath != "/var/lib/mindk/mindk.db" {
		t.Errorf("unexpected database path %q", cfg.DatabasePath)
	}
}

func TestFromEnv_ProductionRequiresFrontendOrigin(t *testing.T) {
	values := validProductionEnv()
	delete(values, "FRONTEND_ORIGIN")

	if _, err := FromEnv(fakeEnv(values)); err == nil {
		t.Fatal("expected an error when FRONTEND_ORIGIN is missing in production")
	}
}

func TestFromEnv_ProductionRequiresDatabasePath(t *testing.T) {
	values := validProductionEnv()
	delete(values, "DATABASE_PATH")

	if _, err := FromEnv(fakeEnv(values)); err == nil {
		t.Fatal("expected an error when DATABASE_PATH is missing in production")
	}
}

func TestFromEnv_MissingOpenAIKey(t *testing.T) {
	values := validProductionEnv()
	delete(values, "OPENAI_API_KEY")

	if _, err := FromEnv(fakeEnv(values)); err == nil {
		t.Fatal("expected an error when OPENAI_API_KEY is missing")
	}
}

func TestFromEnv_MissingClerkSecret(t *testing.T) {
	values := validProductionEnv()
	delete(values, "CLERK_SECRET_KEY")

	if _, err := FromEnv(fakeEnv(values)); err == nil {
		t.Fatal("expected an error when CLERK_SECRET_KEY is missing")
	}
}

func TestFromEnv_InvalidAppEnv(t *testing.T) {
	values := validProductionEnv()
	values["APP_ENV"] = "staging"

	if _, err := FromEnv(fakeEnv(values)); err == nil {
		t.Fatal("expected an error for an unrecognized APP_ENV value")
	}
}

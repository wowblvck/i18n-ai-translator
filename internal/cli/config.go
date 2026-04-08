package cli

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

func loadConfig(configPath string) (appConfig, string, error) {
	configPath = strings.TrimSpace(configPath)
	paths := []string{".i18n-translator.yaml", ".i18n-translator.yml"}
	if configPath != "" {
		paths = []string{configPath}
	}

	for _, path := range paths {
		data, err := os.ReadFile(path)
		if err != nil {
			if os.IsNotExist(err) {
				if configPath != "" {
					return appConfig{}, "", fmt.Errorf("config file not found: %s", path)
				}
				continue
			}
			return appConfig{}, "", fmt.Errorf("failed to read config file %s: %w", path, err)
		}

		var cfg appConfig
		if err := yaml.Unmarshal(data, &cfg); err != nil {
			return appConfig{}, "", fmt.Errorf("failed to parse config file %s: %w", path, err)
		}
		return cfg, path, nil
	}

	return appConfig{}, "", nil
}

func collectExplicitFlags(fs *flag.FlagSet) map[string]bool {
	explicit := map[string]bool{}
	fs.Visit(func(f *flag.Flag) {
		explicit[f.Name] = true
	})
	return explicit
}

func resolveStringOption(flagName, flagValue, envValue, configValue, defaultValue string, explicitFlags map[string]bool) string {
	if explicitFlags[flagName] {
		return flagValue
	}
	if env := strings.TrimSpace(envValue); env != "" {
		return env
	}
	if cfg := strings.TrimSpace(configValue); cfg != "" {
		return cfg
	}
	return defaultValue
}

func resolveIntOption(flagName, envName string, flagValue int, envValue string, configValue, defaultValue int, explicitFlags map[string]bool) (int, error) {
	if explicitFlags[flagName] {
		return flagValue, nil
	}

	if env := strings.TrimSpace(envValue); env != "" {
		v, err := strconv.Atoi(env)
		if err != nil {
			return 0, fmt.Errorf("%s must be a valid integer: %w", envName, err)
		}
		return v, nil
	}

	if configValue != 0 {
		return configValue, nil
	}

	return defaultValue, nil
}

func resolveDurationOption(flagName, envName string, flagValue time.Duration, envValue, configValue string, defaultValue time.Duration, explicitFlags map[string]bool) (time.Duration, error) {
	if explicitFlags[flagName] {
		return flagValue, nil
	}

	if env := strings.TrimSpace(envValue); env != "" {
		v, err := time.ParseDuration(env)
		if err != nil {
			return 0, fmt.Errorf("%s must be a valid duration: %w", envName, err)
		}
		return v, nil
	}

	if cfg := strings.TrimSpace(configValue); cfg != "" {
		v, err := time.ParseDuration(cfg)
		if err != nil {
			return 0, fmt.Errorf("config field for %s must be a valid duration: %w", flagName, err)
		}
		return v, nil
	}

	return defaultValue, nil
}

func resolveBoolOption(flagName string, flagValue bool, envValue string, configValue, defaultValue bool, explicitFlags map[string]bool) (bool, error) {
	if explicitFlags[flagName] {
		return flagValue, nil
	}

	if env := strings.TrimSpace(envValue); env != "" {
		v, err := strconv.ParseBool(env)
		if err != nil {
			return false, fmt.Errorf("%s must be a valid boolean: %w", flagName, err)
		}
		return v, nil
	}

	if configValue {
		return true, nil
	}

	return defaultValue, nil
}

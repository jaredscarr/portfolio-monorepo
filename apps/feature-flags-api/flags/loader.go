package flags

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
)

var (
	flagsCache = make(map[string]map[string]bool) // env -> map[flagKey]bool
	flagsLock  sync.RWMutex
)

// LoadFlagsFromDisk loads the JSON flag file for a given environment.
func LoadFlagsFromDisk(env string) error {
	file := filepath.Join("flags", fmt.Sprintf("%s.json", env))
	data, err := os.ReadFile(file)
	if err != nil {
		return fmt.Errorf("failed to read flag file for env=%s: %w", env, err)
	}

	var parsed map[string]bool
	if err := json.Unmarshal(data, &parsed); err != nil {
		return fmt.Errorf("invalid JSON in %s: %w", file, err)
	}

	flagsLock.Lock()
	defer flagsLock.Unlock()
	flagsCache[env] = parsed
	return nil
}

// GetAllFlags returns all flags for a given environment.
func GetAllFlags(env string) (map[string]bool, error) {
	flagsLock.RLock()
	defer flagsLock.RUnlock()

	flags, ok := flagsCache[env]
	if !ok {
		return nil, fmt.Errorf("flags not loaded for env: %s", env)
	}
	return flags, nil
}

// GetSingleFlag returns a specific flag's value by key.
func GetSingleFlag(env, key string) (bool, bool, error) {
	allFlags, err := GetAllFlags(env)
	if err != nil {
		return false, false, err
	}

	val, exists := allFlags[key]
	return val, exists, nil
}

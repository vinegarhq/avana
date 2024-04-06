package binary

import (
	"fmt"

	"golang.org/x/sys/windows/registry"
)

func (b *binary) verKeyName() string {
	return "LastKnown" + b.d.Type.Short() + "Version"
}

func (b *binary) key() (registry.Key, error) {
	k, err := registry.OpenKey(registry.CURRENT_USER, stateRegistryPath, registry.ALL_ACCESS)
	if err == nil {
		return k, nil
	} else if err != registry.ErrNotExist {
		return 0, fmt.Errorf("open: %w", err)
	}

	k, _, err = registry.CreateKey(registry.CURRENT_USER, stateRegistryPath, registry.ALL_ACCESS)
	if err != nil {
		return 0, fmt.Errorf("create: %w", err)
	}
	return k, nil
}

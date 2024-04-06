package binary

import (
	"errors"
	"fmt"
	"io"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"
	"golang.org/x/sys/windows/registry"
	"time"

	"github.com/apprehensions/rbxbin"
	cs "github.com/apprehensions/rbxweb/clientsettings"

	"github.com/vinegarhq/avana/internal/dirs"
)

const logTimeout = 6 * time.Second
const stateRegistryPath = `Software\Avana\State`

type binary struct {
	dir string
	d   rbxbin.Deployment
	m   rbxbin.Mirror
	lf  *os.File
}

func New() *binary {
	return &binary{
		m: rbxbin.DefaultMirror,
	}
}

func (b *binary) Run(args ...string) error {
	if b.dir == "" {
		return errors.New("deployment directory unknown")
	}

	exe := "Roblox" + b.d.Type.Short() + "Beta.exe"
	cmd := exec.Command(filepath.Join(b.dir, exe), args...)
	if cmd.Err != nil {
		return cmd.Err
	}
	cmd.Stderr = io.MultiWriter(os.Stderr, b.lf)

	slog.Info("Running!", "name", b.d.Type, "cmd", cmd)

	go func() {
		for {
			if cmd.Process != nil {
				break
			}
		}

		lf, err := robloxLogFile()
		if err != nil {
			slog.Error("Failed to find Roblox log file", "error", err.Error())
			return
		}

		tailFile(lf)
	}()

	return cmd.Run()
}

func (b *binary) Setup(bt cs.BinaryType) error {
	if err := b.setupLogging(bt.Short()); err != nil {
		return fmt.Errorf("logging: %w", err)
	}

	k, err := b.key()
	if err != nil {
		return fmt.Errorf("state reg key: %w", err)
	}
	defer k.Close()

	d, err := rbxbin.GetDeployment(bt, "")
	if err != nil {
		return fmt.Errorf("get deployment: %w", err)
	}
	b.d = d
	b.dir = filepath.Join(dirs.Versions, b.d.GUID)

	ver, _, err := k.GetStringValue(b.verKeyName())
	if err != nil && err != registry.ErrNotExist {
		return fmt.Errorf("retrieve last known state guid: %w", err)
	}

	if ver == b.d.GUID {
		slog.Info("Up to date", "guid", b.d.GUID)
		return nil
	}

	slog.Info("Installing!", "old_guid", ver, "new_guid", b.d.GUID)

	if err := b.install(); err != nil {
		return fmt.Errorf("install: %w", err)
	}

	if err := k.SetStringValue(b.verKeyName(), b.d.GUID); err != nil {
		return fmt.Errorf("set last known state guid: %w", err)
	}

	return nil
}

package main

import (
	"fmt"
	"path/filepath"
	
	"github.com/vinegarhq/avana/internal/dirs"
	"github.com/vinegarhq/vinegar/roblox/bootstrapper"
	"github.com/vinegarhq/vinegar/util"
)

func SetupManifest(m *bootstrapper.Manifest, dir string) error {
	lim := len(m.Packages)-1

	for idx, pkg := range m.Packages {
		url := m.Version.DeployURL + "-" + pkg.Name
		dest := filepath.Join(dirs.Downloads, pkg.Checksum)

		dialog.Text("Downloading package: " + pkg.Name)

		if err := util.Download(url, dest); err != nil {
			return err
		}

		if err := util.VerifyFileMD5(dest, pkg.Checksum); err != nil {
			return err
		}

		dialog.Value(int((float64(idx) / float64(lim)) * 100))
	}

	for idx, pkg := range m.Packages {
		src := filepath.Join(dirs.Downloads, pkg.Checksum)
		dest, ok := bootstrapper.PlayerDirectories[pkg.Name]

		dialog.Text("Extracting package: " + pkg.Name)

		if !ok {
			return fmt.Errorf("unhandled package: %s", pkg.Name)
		}

		if err := util.Extract(src, filepath.Join(dir, dest)); err != nil {
			return err
		}

		dialog.Value(int((float64(idx) / float64(lim)) * 100))
	}
	dialog.Complete()

	dialog.Text("Writing AppSettings")

	return bootstrapper.WriteAppSettings(dir)
}
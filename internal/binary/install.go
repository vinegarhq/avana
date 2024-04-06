package binary

import (
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"

	"github.com/apprehensions/rbxbin"
	"github.com/vinegarhq/avana/internal/dirs"
	"golang.org/x/sync/errgroup"
)

func (b *binary) install() error {
	if b.dir == "" {
		return errors.New("deployment directory unknown")
	}

	ps, err := b.m.GetPackages(b.d)
	if err != nil {
		return err
	}

	if err := os.MkdirAll(b.dir, 0o755); err != nil {
		return err
	}

	dsts := rbxbin.BinaryDirectories(b.d.Type)
	eg := new(errgroup.Group)

	for _, p := range ps {
		if p.Name == "RobloxPlayerLauncher.exe" {
			continue
		}

		eg.Go(func() error {
			return b.setupPackage(&p, dsts[p.Name])
		})
	}

	if err := eg.Wait(); err != nil {
		return err
	}

	if err := rbxbin.WriteAppSettings(b.dir); err != nil {
		return err
	}

	return nil
}

func (b *binary) setupPackage(p *rbxbin.Package, dir string) error {
	src := filepath.Join(dirs.Downloads, p.Checksum)

	if err := p.Verify(src); err != nil {
		slog.Info("Downloading package", "name", p.Name, "size", p.Size)

		if err := download(b.m.PackageURL(b.d, p.Name), src); err != nil {
			return err
		}

		if err := p.Verify(src); err != nil {
			return err
		}
	}

	if err := p.Extract(src, b.dir); err != nil {
		return err
	}

	return nil
}

func download(url, file string) error {
	out, err := os.Create(file)
	if err != nil {
		return err
	}
	defer out.Close()

	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %s", resp.Status)
	}

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	return nil
}

package binary

import (
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/sewnie/rbxbin"
	"github.com/vinegarhq/avana/internal/dirs"
	"golang.org/x/sync/errgroup"
)

func (b *binary) install() error {
	if b.dir == "" {
		return errors.New("deployment directory unknown")
	}

	ps, err := b.m.GetPackages(b.c, b.d)
	if err != nil {
		return fmt.Errorf("packages: %w", err)
	}

	dsts, err := b.m.BinaryDirectories(b.c, b.d)
	if err != nil {
		return fmt.Errorf("dirs: %w", err)
	}

	if err := os.MkdirAll(b.dir, 0o755); err != nil {
		return fmt.Errorf("mkdir: %w", err)
	}

	eg := new(errgroup.Group)

	for _, p := range ps {
		if !strings.HasSuffix(p.Name, ".zip") {
			continue
		}

		eg.Go(func() error {
			if err := b.setupPackage(&p, dsts[p.Name]); err != nil {
				return fmt.Errorf("%s: %w", p.Name, err)
			}
			return nil
		})
	}

	if err := eg.Wait(); err != nil {
		return err
	}

	if err := rbxbin.WriteAppSettings(b.dir); err != nil {
		return fmt.Errorf("appsettings: %w", err)
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

	if err := p.Extract(src, filepath.Join(b.dir, dir)); err != nil {
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

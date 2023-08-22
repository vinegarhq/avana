package main

import (
	"log"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/ncruces/zenity"
	"github.com/vinegarhq/vinegar/roblox"
	"github.com/vinegarhq/vinegar/roblox/bootstrapper"
	"github.com/vinegarhq/avana/internal/dirs"
)

var dialog zenity.ProgressDialog

// this is horrible i hate it
func dialogerror(err error) {
	zenity.Error(err.Error(),
		zenity.Title("Avana"),
		zenity.ErrorIcon,
	)

	panic(err)
}

func SetupBinary(ver roblox.Version, dir string) {
	if err := dirs.Mkdirs(dir, dirs.Downloads); err != nil {
		dialogerror(err)
	}

	manifest, err := bootstrapper.FetchManifest(ver, dirs.Downloads)
	if err != nil {
		dialogerror(err)
	}

	if err := SetupManifest(&manifest, dir); err != nil {
		dialogerror(err)
	}
}

func main() {
	var err error

	dialog, err = zenity.Progress(
		zenity.Title("Avana"),
	)
	if err != nil {
		panic(err)
	}
	defer dialog.Close()

	log.Println(dirs.Data)

	dialog.Text("Fetching Version")

	ver, err := roblox.LatestVersion(roblox.Player, "LIVE")
	if err != nil {
		log.Fatal(err)
	}

	verDir := filepath.Join(dirs.Versions, ver.GUID)

	_, err = os.Stat(filepath.Join(verDir, "AppSettings.xml"))
	if err != nil {
		dialog.Text("Updating/Installing Player")
		SetupBinary(ver, verDir)
	}
	
	dialog.Text("Running Player")

	cmd := exec.Command(filepath.Join(verDir, roblox.Player.Executable()))
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout

	dialog.Close()

	err = cmd.Run()
	if err != nil {
		dialogerror(err)
	}
}
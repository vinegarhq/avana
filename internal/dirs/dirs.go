package dirs

import (
	"os"
	"path/filepath"
)

var (
	Data      string
	Downloads string
	Versions  string
	Logs      string
)

func init() {
	// this is %LOCALAPPDATA% but who cares its data
	cache, err := os.UserCacheDir()
	if err != nil {
		panic(err)
	}

	Data = filepath.Join(cache, "Avana")
	Downloads = filepath.Join(Data, "Downloads")
	Versions = filepath.Join(Data, "Versions")
	Logs = filepath.Join(Data, "Logs")

	for _, dir := range []string{Data, Downloads, Versions, Logs} {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			panic(err)
		}
	}
}

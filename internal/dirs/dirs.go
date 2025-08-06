package dirs

import (
	"os"
	"path/filepath"
)

var (
	Data      string
	Downloads string
	Versions  string
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

	for _, dir := range []string{Data, Downloads, Versions} {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			panic(err)
		}
	}
}

package dirs

import (
	"os"
	"path/filepath"
)

var (
	Data string
	Downloads string
	Versions string
)

func init() {
	// this is %LOCALAPPDATA% but who cares its data
	userCache, err := os.UserCacheDir()
	if err != nil {
		panic(err)
	}

	Data = filepath.Join(userCache, "Avana")
	Downloads = filepath.Join(Data, "Downloads")
	Versions = filepath.Join(Data, "Versions")

	err = Mkdirs(Data, Downloads, Versions)
	if err != nil {
		panic(err)
	}
}

func Mkdirs(dirs ...string) error {
	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return err
		}
	}

	return nil
}
package provider_v2

import (
	"io/fs"
	"os"
)

type pathModeMap map[string]fs.FileMode

func combinePaths(v ...[]string) pathModeMap {

	paths := make(pathModeMap)

	for _, list := range v {

		for _, path := range list {

			info, err := os.Stat(path)
			if err != nil {
				paths[path] = fs.ModeIrregular
				continue
			}

			paths[path] = info.Mode()

		}

	}

	return paths

}

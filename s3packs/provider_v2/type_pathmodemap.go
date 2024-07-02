package provider_v2

import (
	"io/fs"
	"log"
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

	// TODO: Remove this output
	for name, mode := range paths {
		log.Printf("Added Path: %s [%v]", name, mode.String())
	}

	return paths

}

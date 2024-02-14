package objectify

import (
	"fmt"

	"github.com/orme292/s3packer/conf"
)

func NewRootList(ac *conf.AppConfig, paths []string) (rl RootList, err error) {
	rl = make(RootList, len(paths))
	for _, p := range paths {
		dirs, err := getSubDirs(p)
		if err != nil {
			return nil, err
		}
		for _, d := range dirs {
			fmt.Printf("Processing directory: %q\n", d)
			files, err := getFiles(ac, d)
			if err != nil {
				return nil, err
			}
			fol, err := NewFileObjList(ac, files, p)
			if err != nil {
				return nil, err
			}
			rl = append(rl, fol)
		}
	}
	rl.repairRedundantKeys()
	return rl, nil
}

func (rl RootList) repairRedundantKeys() {
	// Having to traverse the nested references can be confusing, so here's a breakdown:
	// rli = RootList index
	// doli = DirObjList index
	// foli = FileObjList index
	// RootList => DirObjList => DirObj        => FileObjList        => FileObj
	// rl       => rl[rli]    => rl[rli][doli] => rl[rli][doli].fol => rl[rli][doli].fol[foli]
	seen := make(map[string]int)

	for i := range rl {
		for file := range rl[i] {
			key := rl[i][file].FKey()
			if _, ok := seen[key]; ok {
				seen[key]++
			} else {
				seen[key] = 1
			}
		}
	}

	for k, num := range seen {
		if num > 1 {
			count := 0
			for i := range rl {
				for file := range rl[i] {
					if rl[i][file].FKey() == k {
						old := rl[i][file].FKey()
						rl[i][file].FName = s("%s_%d", rl[i][file].FName, count)
						rl[i][file].ac.Log.Debug("Key %q => %q", old, rl[i][file].FKey())
					}
					count++
				}
			}
		}
	}

}

func (rl RootList) GetStats() (stats *Stats) {
	stats = &Stats{}
	for i := range rl {
		stats.Add(rl[i].GetStats())
	}
	return stats
}

/* DEBUG */

func (rl RootList) Values() {
	for i := range rl {
		rl[i].Values()
	}
}

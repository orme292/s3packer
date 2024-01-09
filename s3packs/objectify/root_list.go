package objectify

import (
	"fmt"

	"github.com/orme292/s3packer/conf"
)

/*
A RootList, DirObjList, DirObj, FileObjList and FileObj are all components of a nested,
tree-like file structure. Here's a simple representation of their relationships:

 RootList ........................................ -> RootList is the list of directories specified
 | ...............................................    in the profile.
 +-> DirObjList, DirObjList, DirObjList, ......... -> Each DirObjList is built from a directory in
 | ...............................................    and contains all subdirectories, with infinite
 | ...............................................    depth. Each subdirectory is represented by a
 | ...............................................    DirObj.
 +-> DirObj, DirObj, DirObj, ..................... -> Each DirObj is a single directory and contains
 | ...............................................    a single FileObjList.
 +-> FileObjList, FileObjList, FileObjList, ...... -> Each FileObjList is built from the directory
 | ...............................................    represented by the DirObj and contains all
 | ...............................................    files in that directory. Each file is
 | ...............................................    represented by a FileObj.
 +-> FileObj, FileObj, FileObj, .................. -> Each FileObj is a single file and contains
 | ...............................................    all the information needed to upload that file.

  Echo branch of the tree has its own type and methods. Uploads are intended to happen at either
  the RootList or FileObjList level. Traversing the nested references can be confusing when using
  the range keyword. See RootList.repairRedundantKeys() for an example of how to do this.

  An alternative way of executing a method would be to code the same method in each type and
  call it down the tree. For example, RootList.Count() calls DirObjList.Count() which calls
  DirObj.Count() which calls FileObjList.Count() which calls FileObj.Count(). The results are
  passed back up the tree.
*/

type RootList []DirObjList

func NewRootList(ac *conf.AppConfig, paths []string) (rl RootList, err error) {
	rl = make(RootList, len(paths))
	for i, p := range paths {
		fmt.Printf("Processing Root: %q\n", p)
		dol, err := NewDirObjList(ac, p)
		if err != nil {
			return nil, err
		}
		rl[i] = dol
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
	for rli := range rl {
		for doli := range rl[rli] {
			for foli := range rl[rli][doli].Fol {
				key := rl[rli][doli].Fol[foli].FKey()
				if _, ok := seen[key]; ok {
					seen[key]++
				} else {
					seen[key] = 1
				}
			}
		}
	}

	for k, num := range seen {
		if num > 1 {
			count := 0
			for rli := range rl {
				for doli := range rl[rli] {
					for foli := range rl[rli][doli].Fol {
						if rl[rli][doli].Fol[foli].FKey() == k {
							old := rl[rli][doli].Fol[foli].FKey()
							rl[rli][doli].Fol[foli].FName = s("%s_%d",
								rl[rli][doli].Fol[foli].FName, count)
							rl[rli][doli].Ac.Log.Debug("Key %q => %q",
								old, rl[rli][doli].Fol[foli].FKey())
						}
						count++
					}
				}
			}
		}
	}
}

func (rl RootList) GetStats() (stats *Stats) {
	stats = &Stats{}
	for _, dol := range rl {
		for _, do := range dol {
			stats.Add(do.Fol.GetStats())
		}
	}
	return stats
}

/* DEBUG */

func (rl RootList) Values() {
	for _, dol := range rl {
		for _, do := range dol {
			do.Values()
			for _, fo := range do.Fol {
				fo.Values()
			}
		}
	}
}

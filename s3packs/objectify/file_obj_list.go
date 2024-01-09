package objectify

import (
	"github.com/orme292/s3packer/conf"
)

type FileObjList []*FileObj

func NewFileObjList(ac *conf.AppConfig, files []string, rel string) (fol FileObjList, err error) {
	fol = make(FileObjList, len(files))
	for i, f := range files {
		grp := i % ac.Opts.MaxUploads
		fol[i], err = NewFileObj(ac, f, rel, grp)
		if err != nil {
			return nil, err
		}
	}
	if rel != EmptyString {
		fol.repairRedundantKeys()
	}
	return fol, nil
}

func (fol FileObjList) repairRedundantKeys() {
	if len(fol) <= 1 {
		return
	}

	fol[0].ac.Log.Debug("Checking for redundant keys in %q...", fol[0].FPseudoP)
	seen := make(map[string]int)
	for i := range fol {
		if _, ok := seen[fol[i].FKey()]; ok {
			seen[fol[i].FKey()]++
		} else {
			seen[fol[i].FKey()] = 1
		}
	}

	for k, num := range seen {
		if num > 1 {
			count := 0
			for i := range fol {
				if fol[i].FKey() == k {
					old := fol[i].FKey()
					fol[i].FName = s("%s_%d", fol[i].FName, count)
					fol[i].ac.Log.Debug("Key %q => %q", old, fol[i].FKey())
				}
				count++
			}
		}
	}
}

func (fol FileObjList) MaxGroup() (max int) {
	if len(fol) == 0 {
		return 0
	}

	max = fol[0].Group
	for _, fo := range fol {
		if fo.Group > max {
			max = fo.Group
		}
	}

	return
}
func (fol FileObjList) GetStats() (stats Stats) {
	for _, fo := range fol {
		stats.Objects++
		if fo.IsUploaded {
			stats.Bytes += fo.FileSize
			stats.Uploaded++
		}
		if fo.Ignore {
			stats.Ignored++
		}
	}
	return stats
}

/* DEBUG */

func (fol FileObjList) Values() {
	for _, fo := range fol {
		fo.Values()
	}
}

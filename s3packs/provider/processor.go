package provider

import (
	"fmt"

	"github.com/orme292/s3packer/conf"
	"github.com/orme292/s3packer/s3packs/objectify"
)

func NewProcessor(ac *conf.AppConfig, ops Operator, iterFn IteratorFunc) (p *Processor, err error) {
	var (
		rl  objectify.RootList
		fol objectify.FileObjList
	)

	fmt.Printf("Starting Processor...\n")
	if len(ac.Directories) > 0 {
		rl, err = objectify.NewRootList(ac, ac.Directories)
		if err != nil {
			return nil, err
		}
	}
	if len(ac.Files) > 0 {
		fol, err = objectify.NewFileObjList(ac, ac.Files, EmptyPath)
		if err != nil {
			return nil, err
		}
	}

	p = &Processor{
		ac:     ac,
		rl:     rl,
		fol:    fol,
		ops:    ops,
		iterFn: iterFn,
	}

	// Single Bucket Check
	p.ac.Log.Info("Checking if bucket %q exists...", p.ac.Bucket.Name)
	exists, _ := p.ops.BucketExists()
	if exists == false {
		p.ac.Log.Info("Bucket %q does not exist.", p.ac.Bucket.Name)
		if p.ac.Bucket.Create == true {
			err := p.ops.CreateBucket()
			if err != nil {
				return nil, err
			} else {
				p.ac.Log.Info("Created bucket %q", p.ac.Bucket.Name)
			}
		} else {
			return nil, err
		}
	}

	return p, nil
}

func (p *Processor) Run() (errs Errs) {
	if p.rl != nil {
		if len(p.rl) > 0 {
			for i := range p.rl {
				iterErrs := p.RunIterator(p.rl[i], DisregardGroups)
				errs.Append(iterErrs)
			}
		}
	}
	if p.fol != nil {
		if len(p.fol) > 0 {
			iterErrs := p.RunIterator(p.fol, DisregardGroups)
			errs.Append(iterErrs)
		}
	}
	p.populateStats()
	p.outputIgnored()
	return
}

func (p *Processor) RunIterator(fol objectify.FileObjList, grp int) (errs Errs) {
	if len(fol) == 0 {
		return
	}

	iter, err := p.iterFn(p.ac, fol, grp)
	if err != nil {
		errs.Add(err)
		return
	}

	if len(fol) > 0 {
		fmt.Printf("Uploading %q...\n", fol[0].OriginDir)
	}
	if err := iter.First(); err != nil {
		errs.Add(err)
	}
	for iter.Next() {
		object := iter.Prepare()
		overwrite, msg := p.Overwrite(object)
		if !overwrite {
			iter.MarkIgnore(msg)
			continue
		}
		if object.Before != nil {
			if err := object.Before(); err != nil {
				errs.Add(err)
			}
		}
		if object.Fo().FileSize > MultipartThreshold && p.ops.SupportsMultipartUploads() {
			err = p.ops.UploadMultipart(*object)
		} else {
			err = p.ops.Upload(*object)
		}
		if err != nil {
			p.ac.Log.Error("Error uploading %q: %q", object.Output().Key, err.Error())
			iter.MarkFailed(err.Error())
		}

		if object.After == nil {
			continue
		}
		if err = object.After(); err != nil {
			errs.Add(err)
		}
	}
	if err := iter.Final(); err != nil {
		errs.Add(err)
	}

	return
}

func (p *Processor) Overwrite(object *PutObject) (exists bool, msg string) {
	switch p.ac.Opts.Overwrite {
	case conf.OverwriteNever:
		if exists, _ := p.ops.ObjectExists(object.Output().Key); exists {
			return false, ObjectExists
		} else {
			return true, EmptyString
		}
	case conf.OverwriteAlways:
		return true, EmptyString
	default:
		return false, fmt.Sprintf("unknown overwrite mode: %q", p.ac.Opts.Overwrite.String())
	}
}

func (p *Processor) populateStats() {
	p.Stats = p.rl.GetStats()
	p.Stats.Add(p.fol.GetStats())
	p.Stats.Discrep = p.Stats.Objects - p.Stats.Uploaded - p.Stats.Failed - p.Stats.Ignored
}

func (p *Processor) outputIgnored() {
	if len(p.rl) > 0 {
		for i := range p.rl {
			for file := range p.rl[i] {
				if p.rl[i][file].Ignore {
					p.ac.Log.Warn("Ignored %q: %q", p.rl[i][file].FKey(), p.rl[i][file].IgnoreString)
				}
			}
		}
	}

	if len(p.fol) > 0 {
		for i := range p.fol {
			if p.fol[i].Ignore {
				p.ac.Log.Warn("Ignored %q: %q", p.fol[i].FKey(), p.fol[i].IgnoreString)
			}
		}
	}
}

/* DEBUG */

func (p *Processor) DebugStats() {
	fmt.Printf("Stats: %+v\n", p.Stats)
}

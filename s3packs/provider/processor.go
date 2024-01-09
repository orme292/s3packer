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
	return &Processor{
		ac:     ac,
		rl:     rl,
		fol:    fol,
		ops:    ops,
		iterFn: iterFn,
	}, nil
}

func (p *Processor) Run() (errs Errs) {
	exists, errs := p.ops.BucketExists()
	if !exists {
		fmt.Printf("Bucket does not exist, should create? %t\n", p.ac.Bucket.Create)
		if p.ac.Bucket.Create == true {
			err := p.ops.CreateBucket()
			if err != nil {
				errs.Add(err)
			} else {
				errs.Release()
			}
		} else {
			errs.Add(fmt.Errorf("bucket %q does not exist", p.ac.Bucket.Name))
			return
		}
	}
	if len(errs.Each) > 0 {
		return errs
	}

	if len(p.rl) > 0 {
		for rli := range p.rl {
			for doli := range p.rl[rli] {
				iterErrs := p.RunIterator(p.rl[rli][doli].Fol, DisregardGroups)
				errs.Append(iterErrs)
			}
		}
	}
	if len(p.fol) > 0 {
		iterErrs := p.RunIterator(p.fol, DisregardGroups)
		errs.Append(iterErrs)
	}
	p.populateStats()
	return errs
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
}

/* DEBUG */

func (p *Processor) DebugStats() {
	fmt.Printf("Stats: %+v\n", p.Stats)
}

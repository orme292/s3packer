package conf

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type Builder struct {
	Filename string
	inc      *ProfileIncoming
	ac       *AppConfig
}

func NewBuilder(path string) *Builder {

	fpath, err := filepath.Abs(expandHome(path))
	if err != nil {
		fpath = path
	}

	return &Builder{
		Filename: fpath,
		inc:      NewProfile(),
	}

}

func (b *Builder) FromYaml() (*AppConfig, error) {

	err := b.inc.LoadFromYaml(b.Filename)
	if err != nil {
		return b.ac, err
	}

	b.ac = NewAppConfig()
	err = b.ac.ImportFromProfile(b.inc)
	if err != nil {
		return b.ac, err
	}

	return b.ac, nil

}

func (b *Builder) YamlOut() error {

	profile := NewProfile()
	profile.loadSampleData()

	output, err := yaml.Marshal(&profile)
	if err != nil {
		return err
	}

	_, err = canCreate(b.Filename)
	if err != nil {
		return err
	}

	f, err := os.Create(b.Filename)
	defer f.Close()
	if err != nil {
		return err
	}

	n, err := f.WriteString("---\n")
	if err != nil || n == 0 {
		return fmt.Errorf("bad write: %v", err)
	}

	n, err = f.Write(output)
	if err != nil || n == 0 {
		return fmt.Errorf("bad write: %v", err)
	}

	return nil

}

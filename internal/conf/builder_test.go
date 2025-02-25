package conf

import (
	"fmt"
	"os"
	"testing"
)

func TestNewBuilder(t *testing.T) {
	filename, err := createTestFile()
	if err != nil {
		t.Fatal(err)
	}
	defer func(name string) {
		err := os.Remove(name)
		if err != nil {
			t.Fatal("Could not remove temp file ", name)
		}
	}(filename)

	builder := NewBuilder(filename)
	app, err := builder.FromYaml()
	if err != nil {
		t.Fatal(err)
	}

	fmt.Printf("%+v\n", *app)
}

func TestNewBuilderSample(t *testing.T) {
	filename := os.TempDir() + "/builder_test_sample.yaml"
	defer func(name string) {
		err := os.Remove(name)
		if err != nil {
			t.Fatal("Could not remove temp file ", name)
		}
	}(filename)

	builder := NewBuilder(filename)

	err := builder.YamlOut()
	if err != nil {
		t.Fatal("Could not create sample profile: ", err)
	}
}

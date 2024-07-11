package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/orme292/s3packer/conf"
	"github.com/orme292/s3packer/s3packs"
	flag "github.com/spf13/pflag"
)

// Partial: Overwrite options -- only-if checksum changes (overwrite: always, on-change, never)
// Done: Generate checksums concurrently
// Done: Remove FileObject attribute ShouldMultiPart, not used.
// Done: LogBot, support sprintf style formatting
// Done: Config, add naming section for KeyNamingMethod, pathPrefix, etc
// Done: Config, rename indexes from camel case to dashed "pathPrefix" to "path-prefix"
// Done: Upload/Ignore function return args can be removed -- they can be counted on the fly
// Done: Add option to create sample profile YAML
// Done: Upgrade to AWS SDK v2
// Done: Modular Provider Support (AWS, OCI, GCP, Azure, etc)
// TODO: Test.
// TODO: Concurrent FileObjList creation, it's a drag.
// TODO: Add profile support for Ignoring directories / files.
// TODO: ^^ support blobs.
// TODO: More debug messages
// TODO: Concurrent checksum tagging
// TODO: Support checksum-only overwrite mode
// TODO: Add some console styling, maybe a progress bar.
// TODO: Add silent option
// TODO: Update all comments for each function/method.
// TODO: Add more readable log output, check log levels make sense
// TODO: Consider ErrorAs implementation and hard coding error messages in Const

/*
getFlags uses the flag package to configure and get command line arguments. It returns:
-- profile: The filename of the profile to load.
*/
func getFlags() (profile, create string, err error) {

	flag.StringVar(&profile, "profile", "", "The filename of the profile you want to open.")
	flag.StringVar(&create, "create", "", "Create a new profile with the specified filename.")
	flag.Parse()

	if create == "" && profile == "" {
		err = errors.New("must specify --create=\"filename\" or --profile=\"filename\"")
		return
	}
	if create != "" && profile != "" {
		err = errors.New("use either --create or --profile, not both")
	}
	return

}

/*
main is the entry point of the program. It does the following:
 1. Calls getFlags to get the command line arguments. (See above)
 2. If the --create flag is specified, it calls conf.Create to create a new profile file.
 2. Creates a new configuration object (See conf/appconfig.go) using the provided profile filename
 3. Checks whether the loaded profile has specified directories to process (c.Dirs)
 4. Checks whether the loaded profile has specified individual files to process (c.Files)
 5. Any returned errors from either of the above are printed as warnings and the program terminates with a 0.
*/
func main() {

	profile, create, err := getFlags()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	if create != "" {

		builder := conf.NewBuilder(create)
		err = builder.YamlOut()
		if err != nil {
			log.Fatalf("Unable to write profile: %v", err)
		}

		log.Printf("File written: %s", create)
		os.Exit(0)

	}

	builder := conf.NewBuilder(profile)
	app, err := builder.FromYaml()
	if err != nil {
		log.Fatalf("Error loading profile: %v", err)
	}

	if app.Tui.Output.Screen {
		startWithScreen(app)
	} else {
		startWithoutScreen(app)
	}

	os.Exit(0)

}

func startPacker(app *conf.AppConfig) {
	_, err := s3packs.Do(app)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}

func startWithoutScreen(app *conf.AppConfig) {
	startPacker(app)
}

func startWithScreen(app *conf.AppConfig) {

	go func() {
		time.Sleep(time.Second * 5)
		startPacker(app)
	}()

	_, err := app.Tui.Screen.Run()
	if err != nil {
		fmt.Println("Couldn't start TUI.")
		os.Exit(1)
	}

}

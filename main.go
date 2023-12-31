package main

import (
	"errors"
	"fmt"
	"os"

	pal "github.com/abusomani/go-palette/palette"
	"github.com/orme292/s3packer/conf"
	"github.com/orme292/s3packer/s3pack"
	flag "github.com/spf13/pflag"
)

// Done: Option to turn off checksum tagging (big bottleneck)
// TODO: More debug messages
// TODO: Concurrent checksum tagging
// Done: Remove FileObject attribute ShouldMultiPart, not used.
// TODO: Overwrite options -- only-if checksum changes (overwrite: always, on-change, never)
// TODO: Upload/Ignore function return args can be removed -- they can be counted on the fly
// Done: LogBot, support sprintf style formatting
// Done: Config, add naming section for KeyNamingMethod, pathPrefix, etc
// Done: Config, rename indexes from camel case to dashed "pathPrefix" to "path-prefix"
// TODO: Add some console styling, maybe a progress bar.
// TODO: Add silent option
// TODO: Add option to create sample profile YAML
// TODO: Update all comments for each function/method.
// TODO: Update CHANGELOG.md
// TODO: Update README.md
// TODO: Update VERSION
// TODO: Add more readable log output, check log levels make sense
// TODO: Consider ErrorAs implementation and hard coding error messages in Const
// Done: Upgrade to aws SDK v2
// TODO: GCP Support
// TODO: Azure Support
// TODO: OCI support

/*
getFlags uses the flag package to configure and get command line arguments. It returns:
-- profile: The filename of the profile to load.
*/
func getFlags() (profile string, create string, max int, err error) {
	flag.StringVar(&profile, "profile", "", "The profile filename you want to use.")
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
	p := pal.New(pal.WithBackground(pal.Color(21)), pal.WithForeground(pal.BrightWhite), pal.WithSpecialEffects([]pal.Special{pal.Bold}))
	_, _ = p.Println("s3packer v", s3pack.Version)
	p.SetOptions(pal.WithDefaults(), pal.WithForeground(pal.BrightWhite))
	_, _ = p.Println("https://github.com/orme292/s3packer\n")

	profileF, createF, _, err := getFlags()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	if createF != "" {
		err = conf.Create(createF)
		if err != nil {
			fmt.Printf("An error occurred: %q\n\n", err.Error())
			os.Exit(1)
		} else {
			os.Exit(0)
		}
	}

	a, err := conf.New(profileF)
	if err != nil {
		a.Log.Fatal(err.Error())
	}

	var dirFilesUploaded, filesIgnored, dirFilesIgnored, filesUploaded int
	var dirBytes, fileBytes int64

	fmt.Println("Processing objects...")

	if len(a.Directories) != 0 {
		err, dirBytes, dirFilesUploaded, dirFilesIgnored = s3pack.DirectoryUploader(a, a.Directories)
		if err != nil {
			a.Log.Error(err.Error())
		}
	}

	if len(a.Files) != 0 {
		err, fileBytes, filesUploaded, filesIgnored = s3pack.IndividualFileUploader(a, a.Files)
		if err != nil {
			a.Log.Error(err.Error())
		}
	}

	fmt.Printf("s3packer Finished, Uploaded %d objects, %s, Ignored %d objects.\n", dirFilesUploaded+filesUploaded, s3pack.FileSizeString(dirBytes+fileBytes), dirFilesIgnored+filesIgnored)
	os.Exit(0)
}

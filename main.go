package main

import (
	"errors"
	"flag"
	"fmt"
	"os"

	pal "github.com/abusomani/go-palette/palette"
	"github.com/orme292/s3packer/config"
	"github.com/orme292/s3packer/s3pack"
)

/*
getFlags uses the flag package to configure and get command line arguments. It returns:
-- profile: The filename of the profile to load.
-- max: The maximum number of files to upload at once (not supported yet)
*/
func getFlags() (profile string, max int, err error) {
	flag.StringVar(&profile, "profile", "", "Filename of the YAML profile to load.")
	flag.IntVar(&max, "max", 5, "Maximum number of files to upload at once.")
	flag.Parse()

	if profile == "" {
		err = errors.New("must specify a profile with -profile \"filename\"")
		return
	}
	return
}

/*
main is the entry point of the program. It does the following:
 1. Creates a new configuration object (See config/config.go). This will be passed around to all the modules. It initially
    contains default values and a logbot instance. (See logbot/logbot.go)
 2. Calls getFlags to get the command line arguments. (See above)
 3. Loads the YAML profile specified in the command line arguments. (See config/config.go#Load)
 4. Checks whether the loaded profile has specified directories to process (c.Dirs)
 5. Checks whether the loaded profile has specified individual files to process (c.Files)
 6. Any returned errors from either of the above are printed as warnings and the program terminates with a 0.
*/
func main() {

	p := pal.New(pal.WithBackground(pal.Color(21)), pal.WithForeground(pal.BrightWhite), pal.WithSpecialEffects([]pal.Special{pal.Bold}))
	_, _ = p.Println("s3packer v", s3pack.Version)
	p.SetOptions(pal.WithDefaults(), pal.WithForeground(pal.BrightWhite))
	_, _ = p.Println("https://github.com/orme292/s3packer\n")

	var dirsIgnored, filesIgnored, dirsUploaded, filesUploaded int
	c := config.New()

	filename, _, err := getFlags()
	if err != nil {
		c.Logger.Fatal(err.Error())
	}

	if err = c.Load(filename); err != nil {
		c.Logger.Fatal("Problem loading profile: " + err.Error())
	}

	fmt.Println("Processing objects...")

	if len(c.Dirs) != 0 {
		err, dirsUploaded, dirsIgnored = s3pack.DirectoryUploader(&c, c.Dirs)
		if err != nil {
			c.Logger.Error(err.Error())
		}
	}

	if len(c.Files) != 0 {
		err, filesUploaded, filesIgnored = s3pack.IndividualFileUploader(&c, c.Files)
		if err != nil {
			c.Logger.Error(err.Error())
		}
	}

	fmt.Printf("s3packer Finished, Uploaded %d objects, Ignored %d objects.\n", dirsUploaded+filesUploaded, dirsIgnored+filesIgnored)
	os.Exit(0)
}

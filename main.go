package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/orme292/s3packer/conf"
	"github.com/orme292/s3packer/s3packs"
	"github.com/orme292/s3packer/tuipack"
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

type appFlags struct {
	profile  string
	create   string
	noscreen bool
}

/*
getFlags uses the flag package to configure and get command line arguments. It returns:
-- profile: The filename of the profile to load.
*/
func getFlags() (appFlags, error) {

	var err error
	flags := appFlags{}

	flag.StringVar(&flags.profile, "profile", "", "The filename of the profile you want to open.")
	flag.StringVar(&flags.create, "create", "", "Create a new profile with the specified filename.")
	flag.BoolVar(&flags.noscreen, "noscreen", false, "No fancy text effects, output logs as configured in the profile.")
	flag.Parse()

	if flags.create == "" && flags.profile == "" {
		err = errors.New("must specify --create=\"filename\" or --profile=\"filename\"")
	}
	if flags.create != "" && flags.profile != "" {
		err = errors.New("use either --create or --profile, not both")
	}

	return flags, err

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

	flags, err := getFlags()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	if flags.create != "" {

		builder := conf.NewBuilder(flags.create)
		err = builder.YamlOut()
		if err != nil {
			log.Fatalf("Unable to write profile: %v", err)
		}

		log.Printf("File written: %s", flags.create)
		os.Exit(0)

	}

	builder := conf.NewBuilder(flags.profile)
	app, err := builder.FromYaml()
	if err != nil {
		log.Fatalf("Error loading profile: %v", err)
	}

	if flags.noscreen == true {
		app.LogOpts.Screen = false
		app.Tui.Output.Screen = false
	}

	startMessage(app)

	if app.Tui.Output.Screen {
		startWithScreen(app)
	} else {
		startWithoutScreen(app)
	}

	os.Exit(0)

}

func startMessage(app *conf.AppConfig) {

	fmt.Printf("\ns3packer\n\n")
	fmt.Printf("Logging [screen:%v] [file:%v] [console:%v]\n", app.LogOpts.Screen, app.LogOpts.File, app.LogOpts.Console)
	time.Sleep(1 * time.Second)

}

func startWithoutScreen(app *conf.AppConfig) {

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT)
	go func() {
		<-sig
		os.Exit(1)
	}()

	startPacker(app)

}

func startWithScreen(app *conf.AppConfig) {

	go func() {
		for {
			if app.Tui.Screen != nil {
				break
			}
		}
		startPacker(app)
	}()

	_, err := app.Tui.Screen.Run()
	if err != nil {
		if app.Tui.Screen != nil {
			app.Tui.Screen.ExitAltScreen()
			app.Tui.ScreenQuit()
		}
		log.Fatalf("Couldn't start TUI.\n")
	}

}

func startPacker(app *conf.AppConfig) {

	stats, err := s3packs.Init(app)
	if err != nil {
		if app.Tui.Screen != nil {
			app.Tui.Screen.ExitAltScreen()
		}
		log.Printf("s3packer exited with error: %s\n\n", err.Error())
		os.Exit(1)
	}

	hrb := stats.ReadableString()
	msg := fmt.Sprintf("%s uploaded, %s skipped", hrb[stats.ObjectsBytes],
		hrb[stats.SkippedBytes])
	app.Tui.SendOutput(tuipack.NewLogMsg(msg, tuipack.ScrnLfDefault, tuipack.INFO, msg))
	app.Tui.SendOutput(tuipack.NewLogMsg(stats.String(), tuipack.ScrnLfDefault, tuipack.INFO, stats.String()))

	if app.Tui.Screen != nil {
		app.Tui.Screen.ExitAltScreen()
	}

	os.Exit(0)

}

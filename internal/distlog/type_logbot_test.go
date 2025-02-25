package distlog

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"
)

var lb = &LogBot{
	Level: zerolog.DebugLevel,
	Output: &LogOutput{
		Console: false,
		File:    false,
	},
}

func restoreExitFunc() {
	exitFunc = os.Exit
}

// mockExit is a helper that replaces os.Exit. Instead of exiting the process,
// it will panic. This allows tests to catch the panic and confirm the exit call.
func mockExit(code int) {
	panic(fmt.Sprintf("os.Exit called with code %d", code))
}

// captureOutput helps us redirect output to a buffer so we can inspect it.
func captureOutput(f func()) string {
	oldStderr := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	f()

	w.Close()
	os.Stderr = oldStderr

	var buf bytes.Buffer
	_, _ = io.Copy(&buf, r)
	return buf.String()
}

func TestSetLogLevel(t *testing.T) {
	lb.Level = zerolog.InfoLevel
	lb.Output.Console = false
	lb.Output.File = false

	require.Equal(t, zerolog.InfoLevel, lb.Level)

	lb.SetLogLevel(zerolog.DebugLevel)
	require.Equal(t, zerolog.DebugLevel, lb.Level)
	require.Equal(t, zerolog.GlobalLevel(), zerolog.DebugLevel)

	lb.SetLogLevel(zerolog.WarnLevel)
	require.Equal(t, zerolog.WarnLevel, lb.Level)
	require.Equal(t, zerolog.GlobalLevel(), zerolog.WarnLevel)
}

// blConsoleOnly tests that the logger is built for console output only.
func blConsoleOnly(t *testing.T) {
	lb.Output.Console = true
	lb.Output.File = false

	// Capture console output
	out := captureOutput(func() {
		lb.SetLogLevel(zerolog.InfoLevel)
		logger := lb.BuildLogger(zerolog.InfoLevel)
		logger.Info().Msg("Testing console only")
	})

	require.Contains(t, out, "Testing console only", "Expected console output to contain log message.")
}

// blFileOnly tests that the logger is built for file output only.
func blFileOnly(t *testing.T) {
	tmpFile := filepath.Join(t.TempDir(), fmt.Sprintf("test-%d.log", time.Now().Unix()))
	lb.Logfile = tmpFile
	lb.Output.Console = false
	lb.Output.File = true

	lb.SetLogLevel(zerolog.InfoLevel)
	logger := lb.BuildLogger(zerolog.InfoLevel)
	logger.Info().Msg("Testing file only logger")

	// Read the file and check content
	data, err := os.ReadFile(tmpFile)
	require.NoError(t, err)
	require.Contains(t, string(data), "Testing file only logger", "Expected file output to contain log message.")
}

// blConsoleAndFile tests that the logger is built for both console and file output.
func blConsoleAndFile(t *testing.T) {
	tmpFile := filepath.Join(t.TempDir(), fmt.Sprintf("test-%d.log", time.Now().Unix()))
	lb.Logfile = tmpFile
	lb.Output.Console = true
	lb.Output.File = true

	var consoleOut string
	// Capture console output
	consoleOut = captureOutput(func() {
		lb.SetLogLevel(zerolog.InfoLevel)
		logger := lb.BuildLogger(zerolog.InfoLevel)
		logger.Info().Msg("Test console + file output")
	})

	data, err := os.ReadFile(tmpFile)
	require.NoError(t, err)

	require.Contains(t, consoleOut, "Test console + file output", "Expected console output to contain log message.")
	require.Contains(t, string(data), "Test console + file output", "Expected file output to contain log message.")
}

func TestBuildLogger(t *testing.T) {
	t.Run("BuildLogger-ConsoleOnly", blConsoleOnly)
	t.Run("BuildLogger-FileOnly", blFileOnly)
	t.Run("BuildLogger-ConsoleAndFile", blConsoleAndFile)
}

// TestRouteLogMsg_Console tests that RouteLogMsg routes a message to console only.
func TestRouteLogMsg_Console(t *testing.T) {
	lb.Output.Console = true
	lb.Output.File = false

	msg := "route console test"
	out := captureOutput(func() {
		lb.SetLogLevel(zerolog.DebugLevel)
		lb.RouteLogMsg(zerolog.DebugLevel, msg)
	})

	require.Contains(t, out, msg, "Expected console output to contain the routed message.")
}

// TestRouteLogMsg_File tests that RouteLogMsg routes a message to file only.
func TestRouteLogMsg_File(t *testing.T) {
	tmpFile := filepath.Join(t.TempDir(), fmt.Sprintf("test-%d.log", time.Now().Unix()))
	lb.Logfile = tmpFile
	lb.Output.Console = false
	lb.Output.File = true

	msg := "route file test"
	lb.SetLogLevel(zerolog.InfoLevel)
	lb.RouteLogMsg(zerolog.InfoLevel, msg)

	data, err := os.ReadFile(tmpFile)
	require.NoError(t, err)
	require.Contains(t, string(data), msg, "Expected file output to contain the routed message.")
}

// TestPanic verifies that calling Panic routes the message at panic level and calls os.Exit(1).
// We override exitFunc so the test does not terminate the process.
func TestPanic(t *testing.T) {
	defer restoreExitFunc()
	exitFunc = mockExit // override to catch the call

	lb := &LogBot{
		Level: zerolog.InfoLevel,
		Output: &LogOutput{
			Console: true,
		},
	}

	msg := "panic-level test"

	defer func() {
		if r := recover(); r != nil {
			log.Println(r.(string))
			require.True(t, strings.Contains(r.(string), msg),
				"Expected os.Exit(1) to be called on Panic")
		} else {
			t.Errorf("Expected os.Exit(1) to be called on Panic but did not catch panic.")
		}
	}()

	lb.SetLogLevel(zerolog.PanicLevel)
	lb.Panic("%s", msg)
}

// TestFatal verifies that calling Fatal routes the message at fatal level and calls os.Exit(1).
func TestFatal(t *testing.T) {
	defer restoreExitFunc()
	exitFunc = mockExit

	lb := &LogBot{
		Level: zerolog.InfoLevel,
		Output: &LogOutput{
			Console: true,
		},
	}

	defer func() {
		if r := recover(); r != nil {
			require.True(t, strings.Contains(r.(string), "os.Exit called with code 1"),
				"Expected os.Exit(1) to be called on Fatal")
		} else {
			t.Errorf("Expected os.Exit(1) to be called on Fatal but did not catch panic.")
		}
	}()

	lb.SetLogLevel(zerolog.FatalLevel)
	lb.Fatal("fatal-level test")
}

func TestMsg(t *testing.T) {
	lb := &LogBot{
		Level: zerolog.DebugLevel,
		Output: &LogOutput{
			Console: true,
		},
	}

	levels := map[zerolog.Level]func(m string, a ...any){
		zerolog.DebugLevel: lb.Debug,
		zerolog.InfoLevel:  lb.Info,
		zerolog.WarnLevel:  lb.Warn,
		zerolog.ErrorLevel: lb.Error,
		zerolog.TraceLevel: lb.Blast,
	}
	for l, f := range levels {
		lb.Level = l
		lb.SetLogLevel(l)
		msg := fmt.Sprintf("%s-level test", l.String())
		out := captureOutput(func() {
			f("%s", msg)
		})

		require.Contains(t, out, msg, "Expected to log error message to console.")
	}
}

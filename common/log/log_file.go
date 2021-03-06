// File I/O for logs.

package log

import (
	"errors"
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"strings"
	"time"
)

// MaxSize is the maximum size of a log file in bytes.
var MaxSize uint64 = 1024 * 1024 * 100

var LogDir = "./log"

var (
	pid      = os.Getpid()
	program  = filepath.Base(os.Args[0])
	host     = "unknownhost"
	userName = "unknownuser"
)

func init() {
	h, err := os.Hostname()
	if err == nil {
		host = shortHostname(h)
	}

	current, err := user.Current()
	if err == nil {
		userName = current.Username
	}

	// Sanitize userName since it may contain filepath separators on Windows.
	userName = strings.Replace(userName, `\`, "_", -1)
}

// shortHostname returns its argument, truncating at the first period.
// For instance, given "www.google.com" it returns "www".
func shortHostname(hostname string) string {
	if i := strings.Index(hostname, "."); i >= 0 {
		return hostname[:i]
	}
	return hostname
}

// logName returns a new log file name containing tag, with start time t, and
// the name for the symlink for tag.
func logName(tag string, t time.Time) (name, link string) {
	name = fmt.Sprintf("%s.%s.%s.log.%s.%04d%02d%02d-%02d%02d%02d.%d",
		program,
		host,
		userName,
		tag,
		t.Year(),
		t.Month(),
		t.Day(),
		t.Hour(),
		t.Minute(),
		t.Second(),
		pid)
	return name, program + "." + tag
}

// create creates a new log file and returns the file and its filename, which
// contains tag ("INFO", "FATAL", etc.) and t.  If the file is created
// successfully, create also attempts to update the symlink for that tag, ignoring
// errors.
func create(tag string, t time.Time) (f *os.File, filename string, err error) {
	if LogDir == "" {
		return nil, "", errors.New("log: no log dirs")
	}
	err = os.MkdirAll(LogDir, os.ModePerm)
	if err != nil {
		return nil, "", fmt.Errorf("log: cannot mkdir: %v", err)
	}
	name, link := logName(tag, t)

	fname := filepath.Join(LogDir, name)
	f, err = os.OpenFile(fname, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		return nil, "", fmt.Errorf("log: cannot mkdir: %v", err)
	}

	symlink := filepath.Join(LogDir, link)
	os.Remove(symlink)        // ignore err
	os.Symlink(name, symlink) // ignore err
	return f, fname, nil

}

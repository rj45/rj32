// Copyright (c) 2018-2021 TinyGo Authors. All rights reserved.
// Licensed under a 3 clause BSD license. See LICENSE.tinygo.
//
// Copyright (c) 2021 rj45 (github.com/rj45), MIT Licensed, see LICENSE.

// Package goenv returns environment variables that are used in various parts of
// the compiler. You can query it manually with the `gorj env` subcommand.
package goenv

import (
	"fmt"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"runtime"
)

// Keys is a slice of all available environment variable keys.
var Keys = []string{
	"GOOS",
	"GOARCH",
	"GOROOT",
	"GOPATH",
	"GOCACHE",
	"CGO_ENABLED",
	"GORJROOT",
}

func init() {
	if Get("GOARCH") == "arm" {
		Keys = append(Keys, "GOARM")
	}
}

// GORJROOT is the path to the final location for checking gorj files. If
// unset (by a -X ldflag), then sourceDir() will fallback to the original build
// directory.
var GORJROOT string

// Get returns a single environment variable, possibly calculating it on-demand.
// The empty string is returned for unknown environment variables.
func Get(name string) string {
	switch name {
	case "GOOS":
		if goos := os.Getenv("GOOS"); goos != "" {
			return goos
		}
		return runtime.GOOS
	case "GOARCH":
		if arch := os.Getenv("GOARCH"); arch != "" {
			return arch
		}
		return runtime.GOARCH
	case "GOROOT":
		return getGoroot()
	case "GOPATH":
		if dir := os.Getenv("GOPATH"); dir != "" {
			return dir
		}

		// fallback
		home := getHomeDir()
		return filepath.Join(home, "go")
	case "GOCACHE":
		// Get the cache directory, usually ~/.cache/gorj
		dir, err := os.UserCacheDir()
		if err != nil {
			panic("could not find cache dir: " + err.Error())
		}
		return filepath.Join(dir, "gorj")
	case "GORJROOT":
		return sourceDir()
	default:
		return ""
	}
}

// Return the GORJROOT, or exit with an error.
func sourceDir() string {
	// Use $GORJROOT as root, if available.
	root := os.Getenv("GORJROOT")
	if root != "" {
		if !isSourceDir(root) {
			fmt.Fprintln(os.Stderr, "error: $GORJROOT was not set to the correct root")
			os.Exit(1)
		}
		return root
	}

	if GORJROOT != "" {
		if !isSourceDir(GORJROOT) {
			fmt.Fprintln(os.Stderr, "error: GORJROOT was not set to the correct root")
			os.Exit(1)
		}
		return GORJROOT
	}

	// Find root from executable path.
	path, err := os.Executable()
	if err != nil {
		// Very unlikely. Bail out if it happens.
		panic("could not get executable path: " + err.Error())
	}
	root = filepath.Dir(filepath.Dir(path))
	if isSourceDir(root) {
		return root
	}

	// Fallback: use the original directory from where it was built
	// https://stackoverflow.com/a/32163888/559350
	_, path, _, _ = runtime.Caller(0)
	root = filepath.Dir(filepath.Dir(path))
	if isSourceDir(root) {
		return root
	}

	fmt.Fprintln(os.Stderr, "error: could not autodetect root directory, set the GORJROOT environment variable to override")
	os.Exit(1)
	panic("unreachable")
}

// isSourceDir returns true if the directory looks like a TinyGo source directory.
func isSourceDir(root string) bool {
	// _, err := os.Stat(filepath.Join(root, "src/runtime/internal/sys/zversion.go"))
	// if err != nil {
	// 	return false
	// }
	// todo: make this more unique to gorj
	_, err := os.Stat(filepath.Join(root, "src/runtime/string.go"))
	return err == nil
}

func getHomeDir() string {
	u, err := user.Current()
	if err != nil {
		panic("cannot get current user: " + err.Error())
	}
	if u.HomeDir == "" {
		// This is very unlikely, so panic here.
		// Not the nicest solution, however.
		panic("could not find home directory")
	}
	return u.HomeDir
}

// getGoroot returns an appropriate GOROOT from various sources. If it can't be
// found, it returns an empty string.
func getGoroot() string {
	// An explicitly set GOROOT always has preference.
	goroot := os.Getenv("GOROOT")
	if goroot != "" {
		// Convert to the standard GOROOT being referenced, if it's a TinyGo cache.
		return getStandardGoroot(goroot)
	}

	// Check for the location of the 'go' binary and base GOROOT on that.
	binpath, err := exec.LookPath("go")
	if err == nil {
		binpath, err = filepath.EvalSymlinks(binpath)
		if err == nil {
			goroot := filepath.Dir(filepath.Dir(binpath))
			if isGoroot(goroot) {
				return goroot
			}
		}
	}

	// Check what GOROOT was at compile time.
	if isGoroot(runtime.GOROOT()) {
		return runtime.GOROOT()
	}

	// Check for some standard locations, as a last resort.
	var candidates []string
	switch runtime.GOOS {
	case "linux":
		candidates = []string{
			"/usr/local/go",     // manually installed
			"/usr/lib/go",       // from the distribution
			"/snap/go/current/", // installed using snap
		}
	case "darwin":
		candidates = []string{
			"/usr/local/go",             // manually installed
			"/usr/local/opt/go/libexec", // from Homebrew
		}
	}

	for _, candidate := range candidates {
		if isGoroot(candidate) {
			return candidate
		}
	}

	// Can't find GOROOT...
	return ""
}

// isGoroot checks whether the given path looks like a GOROOT.
func isGoroot(goroot string) bool {
	_, err := os.Stat(filepath.Join(goroot, "src", "runtime", "internal", "sys", "zversion.go"))
	return err == nil
}

// getStandardGoroot returns the physical path to a real, standard Go GOROOT
// implied by the given path.
// If the given path appears to be a TinyGo cached GOROOT, it returns the path
// referenced by symlinks contained in the cache. Otherwise, it returns the
// given path as-is.
func getStandardGoroot(path string) string {
	// Check if the "bin" subdirectory of our given GOROOT is a symlink, and then
	// return the _parent_ directory of its destination.
	if dest, err := os.Readlink(filepath.Join(path, "bin")); nil == err {
		// Clean the destination to remove any trailing slashes, so that
		// filepath.Dir will always return the parent.
		//   (because both "/foo" and "/foo/" are valid symlink destinations,
		//   but filepath.Dir would return "/" and "/foo", respectively)
		return filepath.Dir(filepath.Clean(dest))
	}
	return path
}

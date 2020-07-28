// Package testlicense provides the function to check a license file.
//
// usage: In your `some_test.go`, write a test code like below.
//
//   package yourpackage_test
//
//   import (
//   	"testing"
//
//   	"github.com/google/licensecheck"
//   	"github.com/takumakei/go-testlicense"
//   )
//
//   func TestLicense(t *testing.T) {
//   	testlicense.Test(t, licensecheck.MIT)
//   }
//
package testlicense

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"testing"

	"github.com/google/licensecheck"
)

// Filenames for license.
//
// https://pkg.go.dev/license-policy
//
//   > We currently use github.com/google/licensecheck for license detection, and
//   > look for licenses in files with the following names: COPYING, COPYING.md,
//   > COPYING.markdown, COPYING.txt, LICENCE, LICENCE.md, LICENCE.markdown,
//   > LICENCE.txt, LICENSE, LICENSE.md, LICENSE.markdown, LICENSE.txt,
//   > LICENSE-2.0.txt, LICENCE-2.0.txt, LICENSE-APACHE, LICENCE-APACHE,
//   > LICENSE-APACHE-2.0.txt, LICENCE-APACHE-2.0.txt, LICENSE-MIT, LICENCE-MIT,
//   > LICENSE.MIT, LICENCE.MIT, LICENSE.code, LICENCE.code, LICENSE.docs,
//   > LICENCE.docs, LICENSE.rst, LICENCE.rst, MIT-LICENSE, MIT-LICENCE,
//   > MIT-LICENSE.md, MIT-LICENCE.md, MIT-LICENSE.markdown, MIT-LICENCE.markdown,
//   > MIT-LICENSE.txt, MIT-LICENCE.txt, MIT_LICENSE, MIT_LICENCE, UNLICENSE,
//   > UNLICENCE. The match is case-insensitive.
//
var Filenames = []string{
	"COPYING",
	"COPYING.md",
	"COPYING.markdown",
	"COPYING.txt",
	"LICENCE",
	"LICENCE.md",
	"LICENCE.markdown",
	"LICENCE.txt",
	"LICENSE",
	"LICENSE.md",
	"LICENSE.markdown",
	"LICENSE.txt",
	"LICENSE-2.0.txt",
	"LICENCE-2.0.txt",
	"LICENSE-APACHE",
	"LICENCE-APACHE",
	"LICENSE-APACHE-2.0.txt",
	"LICENCE-APACHE-2.0.txt",
	"LICENSE-MIT",
	"LICENCE-MIT",
	"LICENSE.MIT",
	"LICENCE.MIT",
	"LICENSE.code",
	"LICENCE.code",
	"LICENSE.docs",
	"LICENCE.docs",
	"LICENSE.rst",
	"LICENCE.rst",
	"MIT-LICENSE",
	"MIT-LICENCE",
	"MIT-LICENSE.md",
	"MIT-LICENCE.md",
	"MIT-LICENSE.markdown",
	"MIT-LICENCE.markdown",
	"MIT-LICENSE.txt",
	"MIT-LICENCE.txt",
	"MIT_LICENSE",
	"MIT_LICENCE",
	"UNLICENSE",
	"UNLICENCE",
}

var filenames = make(map[string]struct{})

func init() {
	for _, v := range Filenames {
		filenames[strings.ToLower(v)] = struct{}{}
	}
}

// IsLicenseFilename returns true if the string s matches any one of license
// filenames.
func IsLicenseFilename(s string) bool {
	_, ok := filenames[strings.ToLower(s)]
	return ok
}

// TestPercent calls t.Fatal(err) if no license file is found or percentage is
// less than 90%.
func Test(t *testing.T, want licensecheck.Type) {
	TestPercent(t, want, 90)
}

// TestPercent calls t.Fatal(err) if no license file is found or percentage is
// less than percent.
func TestPercent(t *testing.T, want licensecheck.Type, percent float64) {
	if err := AssertLicense(want, percent); err != nil {
		t.Fatal(err)
	}
}

// AssertLicense returns err != nil if the license in the current directory
// does not match want or percentage is less than percent.
func AssertLicense(want licensecheck.Type, percent float64) error {
	_, b, err := ReadLicense()
	if err != nil {
		return err
	}
	return assertLicense(b, want, percent)
}

// AssertLicenseDir returns err != nil if the license in the dir does not match
// want or percentage is less than percent.
func AssertLicenseDir(dir DirnamesReader, want licensecheck.Type, percent float64) error {
	_, b, err := ReadLicenseDir(dir)
	if err != nil {
		return err
	}
	return assertLicense(b, want, percent)
}

func assertLicense(b []byte, want licensecheck.Type, percent float64) error {
	cov, ok := licensecheck.Cover(b, licensecheck.Options{})
	if !ok {
		return errors.New("license not found")
	}
	if cov.Percent < percent {
		return fmt.Errorf("percentage %f is less than wanted %f", cov.Percent, percent)
	}
	var list []string
	for _, m := range cov.Match {
		if m.Type == want {
			return nil
		}
		list = append(list, fmt.Sprintf("%s(%3.1f%%)", m.Name, m.Percent))
	}
	return fmt.Errorf("license does not match. found: %s", strings.Join(list, ","))
}

// ReadLicenseDir searchs the license file in the current dir and returns the
// filename and the contents.
func ReadLicense() (filename string, contents []byte, err error) {
	f, err := os.Open(".")
	if err != nil {
		return "", nil, err
	}
	defer f.Close()
	return ReadLicenseDir(f)
}

type DirnamesReader interface {
	Readdirnames(int) ([]string, error)
}

// ReadLicenseDir searchs the license file in the dir and returns the filename
// and the contents.
func ReadLicenseDir(dir DirnamesReader) (filename string, contents []byte, err error) {
	names, err := dir.Readdirnames(0)
	if err != nil {
		return "", nil, err
	}
	for _, filename = range names {
		if IsLicenseFilename(filename) {
			b, err := ioutil.ReadFile(filename)
			return filename, b, err
		}
	}
	return "", nil, os.ErrNotExist
}

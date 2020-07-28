package testlicense_test

import (
	"testing"

	"github.com/google/licensecheck"
	"github.com/takumakei/go-testlicense"
)

func TestLicense(t *testing.T) {
	testlicense.Test(t, licensecheck.MIT)
}

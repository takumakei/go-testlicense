# go-testlicense

Package testlicense provides the function to check a license file.

## Usage

In your `some_test.go`, write a test code like below.

```go
package yourpackage_test

import (
  "testing"

  "github.com/google/licensecheck"
  "github.com/takumakei/go-testlicense"
)

func TestLicense(t *testing.T) {
  testlicense.Test(t, licensecheck.MIT)
}
```
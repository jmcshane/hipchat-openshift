package util_test

import (
	"testing"

	"github.com/jmcshane/hipchat-openshift/util"
)

func TestValidateSliceLocator(t *testing.T) {
	expectedIndex := 4
	a := []string{"a", "b", "c", "d", "e", "f"}
	result := util.SliceIndex(len(a), func(i int) bool { return a[i] == "e" })
	if expectedIndex != result {
		t.Errorf("SliceIndex located incorrect index: got %d want %d",
			result, expectedIndex)
	}
}

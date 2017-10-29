// +build !"go1.9"

package main

import "testing"

func runTest(t *testing.T, tc testCase) {
    runFileTest(t, tc)
}

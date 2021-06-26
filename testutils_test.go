package main

import "testing"

func runTest(t *testing.T, tc testCase) {
	t.Helper()
	runFileTest(t, tc)
}

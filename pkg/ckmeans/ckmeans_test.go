package ckmeans

import (
	"reflect"
	"testing"
)

type testCase struct {
	vals      []float64
	nClusters int
	expected  [][]float64
}

func testCkmeansOnce(t *testing.T, tc testCase) bool {
	clusters, err := Ckmeans(tc.vals, tc.nClusters)
	if err != nil {
		t.Errorf("Ckmeans() returned an error on input (%v, %v)", tc.vals, tc.nClusters)
	}
	if !reflect.DeepEqual(clusters, tc.expected) {
		t.Errorf("Ckmeans() gave incorrect result on input (%v, %v) (result: %v; expected: %v)", tc.vals, tc.nClusters, clusters, tc.expected)
	}

	return true
}

func TestCkmeans(t *testing.T) {
	// These tests are based on Simple Statistics:
	// https://github.com/simple-statistics/simple-statistics/blob/master/test/ckmeans.test.js
	_, err := Ckmeans([]float64{}, 10)
	if err == nil {
		t.Errorf("Cannot generate more values than input")
	}
	_, err = Ckmeans([]float64{1}, 2)
	if err == nil {
		t.Errorf("Cannot generate more values than input")
	}

	testCases := []testCase{
		{[]float64{1}, 1, [][]float64{{1}}},
		{[]float64{1, 1, 1, 1}, 1, [][]float64{{1, 1, 1, 1}}},
		{[]float64{1, 1, 1, 1}, 4, [][]float64{{1, 1, 1, 1}}}, // return one cluster for consistency with Simple Statistics
		{[]float64{-1, 2, -1, 2, 4, 5, 6, -1, 2, -1}, 3, [][]float64{
			{-1, -1, -1, -1}, {2, 2, 2}, {4, 5, 6}}},
		{[]float64{1, 1, 1, 1}, 1, [][]float64{{1, 1, 1, 1}}},
		{[]float64{1, 1, 2}, 3, [][]float64{{1}, {1}, {2}}},
		{[]float64{1, 2, 3}, 3, [][]float64{{1}, {2}, {3}}},
		{[]float64{0, 3, 4}, 2, [][]float64{{0}, {3, 4}}},
		{[]float64{-3, 0, 4}, 2, [][]float64{{-3, 0}, {4}}},
		{[]float64{1, 2, 2, 3}, 3, [][]float64{{1}, {2, 2}, {3}}},
		{[]float64{1, 2, 2, 3, 3}, 3, [][]float64{{1}, {2, 2}, {3, 3}}},
		{[]float64{1, 2, 3, 2, 3}, 3, [][]float64{{1}, {2, 2}, {3, 3}}},
		{[]float64{3, 2, 3, 2, 1}, 3, [][]float64{{1}, {2, 2}, {3, 3}}},
		{[]float64{3, 2, 3, 5, 2, 1}, 3, [][]float64{{1, 2, 2}, {3, 3}, {5}}},
		{[]float64{0, 1, 2, 100, 101, 103}, 2, [][]float64{{0, 1, 2}, {100, 101, 103}}},
		{[]float64{0, 1, 2, 50, 100, 101, 103}, 3, [][]float64{{0, 1, 2}, {50}, {100, 101, 103}}},
		{[]float64{64.64249127327881, 64.64249127328245, 57.79216426169771}, 2, [][]float64{{57.79216426169771}, {64.64249127327881, 64.64249127328245}}},
	}
	for _, c := range testCases {
		testCkmeansOnce(t, c)
	}
}

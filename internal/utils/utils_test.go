package utils

import (
	"bytes"
	"testing"
)

func TestInsertCatalogToFile(t *testing.T) {
	testcases := []struct {
		data    []byte
		toc     []byte
		catalog []byte
		expect  []byte
	}{
		{
			data:    []byte("test\n"),
			toc:     nil,
			catalog: []byte("catalog\n"),
			expect:  []byte("catalog\ntest\n"),
		},
		{
			data:    []byte("test\n[TOC]\nend"),
			toc:     []byte("[TOC]"),
			catalog: []byte("catalog\n"),
			expect:  []byte("test\ncatalog\nend"),
		},
	}
	for _, tc := range testcases {
		if result := insertCatalogToFile(tc.data, tc.catalog, tc.toc); !bytes.Equal(result, tc.expect) {
			t.Errorf("want %v, got %v\n", string(tc.expect), string(result))
		}
	}
}

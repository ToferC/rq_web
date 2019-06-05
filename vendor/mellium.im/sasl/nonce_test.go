// Copyright 2016 The Mellium Contributors.
// Use of this source code is governed by the BSD 2-clause license that can be
// found in the LICENSE file.

package sasl

import (
	"crypto/rand"
	"errors"
	"testing"
)

type zeroReader struct{}

func (zeroReader) Read(p []byte) (int, error) {
	for i := 0; i < cap(p); i++ {
		p[i] = 0
	}
	return cap(p), nil
}

func TestNoncePanicsIfTooShort(t *testing.T) {
	func() {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("Expected nonce() to panic if given zero length")
			}
		}()

		nonce(0, rand.Reader)
	}()
	func() {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("Expected nonce() to panic if given negative length")
			}
		}()

		nonce(-1, rand.Reader)
	}()
}

type errReader struct{}

func (errReader) Read(p []byte) (int, error) {
	return 0, errors.New("Expected errReader error")
}

type nopReader struct{}

func (nopReader) Read(p []byte) (int, error) {
	return 0, nil
}

func TestNoncePanicsOnErrorReadingRand(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected nonce to panic if generating randomness fails")
		}
	}()

	nonce(1, errReader{})
}

func TestNoncePanicsOnIncompleteReadingRand(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected nonce to panic if too little randomness was generated")
		}
	}()

	nonce(1, nopReader{})
}

func BenchmarkNonce(b *testing.B) {
	cr := rand.Reader
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		nonce(16, cr)
	}
}

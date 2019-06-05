// Copyright 2016 The Mellium Contributors.
// Use of this source code is governed by the BSD 2-clause license that can be
// found in the LICENSE file.

package sasl_test

import (
	"bytes"
	"fmt"

	"mellium.im/sasl"
)

// A custom SASL Mechanism that implements XOAUTH2:
// https://developers.google.com/gmail/xoauth2_protocol
var xoauth2 = sasl.Mechanism{
	Name: "XOAUTH2",
	Start: func(m *sasl.Negotiator) (bool, []byte, interface{}, error) {
		// Start is called only by clients and returns the client first message.

		username, password, _ := m.Credentials()

		payload := []byte(`user=`)
		payload = append(payload, username...)
		payload = append(payload, '\x01')
		payload = append(payload, []byte(`auth=Bearer `)...)
		payload = append(payload, password...)
		payload = append(payload, '\x01', '\x01')

		return false, payload, nil, nil
	},
	Next: func(m *sasl.Negotiator, challenge []byte, _ interface{}) (bool, []byte, interface{}, error) {
		// Next is called by both clients and servers and must be able to generate
		// and handle every challenge except for the client first message which is
		// generated (but not handled by) by Start.

		state := m.State()

		// If we're a client or a server that's past the AuthTextSent step, we
		// should never actually hit this step for the XOAUTH2 mechanism so return
		// an error.
		if state&sasl.Receiving != sasl.Receiving || state&sasl.StepMask != sasl.AuthTextSent {
			return false, nil, nil, sasl.ErrTooManySteps
		}

		parts := bytes.Split(challenge, []byte{1})
		if len(parts) != 3 {
			return false, nil, nil, sasl.ErrInvalidChallenge
		}
		user := bytes.TrimPrefix([]byte("user="), parts[0])
		if len(user) == len(parts[0]) {
			return false, nil, nil, sasl.ErrInvalidChallenge
		}
		pass := bytes.TrimPrefix([]byte("Auth=Bearer "), parts[1])
		if len(pass) == len(parts[1]) {
			return false, nil, nil, sasl.ErrInvalidChallenge
		}
		if len(parts[2]) > 0 {
			return false, nil, nil, sasl.ErrInvalidChallenge
		}

		if m.Permissions(sasl.Credentials(func() ([]byte, []byte, []byte) {
			return user, pass, nil
		})) {
			return false, nil, nil, nil
		}
		return false, nil, nil, sasl.ErrAuthn
	},
}

func Example_xOAUTH2() {
	c := sasl.NewClient(
		xoauth2,
		sasl.Credentials(func() ([]byte, []byte, []byte) {
			return []byte("someuser@example.com"), []byte("vF9dft4qmTc2Nvb3RlckBhdHRhdmlzdGEuY29tCg=="), []byte{}
		}),
	)

	// This is the first step and we haven't received any challenge from the
	// server yet.
	more, resp, _ := c.Step(nil)
	fmt.Printf("%v %s", more, bytes.Replace(resp, []byte{1}, []byte{' '}, -1))

	// Output: false user=someuser@example.com auth=Bearer vF9dft4qmTc2Nvb3RlckBhdHRhdmlzdGEuY29tCg==
}

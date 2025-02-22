// This file is part of DiceDB.
// Copyright (C) 2024 DiceDB (dicedb.io).
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program. If not, see <https://www.gnu.org/licenses/>.

package http

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestHSetNX(t *testing.T) {
	exec := NewHTTPCommandExecutor()

	testCases := []struct {
		name     string
		commands []HTTPCommand
		expected []interface{}
		delays   []time.Duration
	}{
		{
			name: "HSetNX returns 0 when field is already set",
			commands: []HTTPCommand{
				{Command: "HSETNX", Body: map[string]interface{}{"key": "key_nx_t1", "field": "field", "value": "value"}},
				{Command: "HSETNX", Body: map[string]interface{}{"key": "key_nx_t1", "field": "field", "value": "value_new"}},
			},
			expected: []interface{}{float64(1), float64(0)},
			delays:   []time.Duration{0, 0},
		},
		{
			name: "HSetNX with new field",
			commands: []HTTPCommand{
				{Command: "HSETNX", Body: map[string]interface{}{"key": "key_nx_t2", "field": "field", "value": "value"}},
			},
			expected: []interface{}{float64(1)},
			delays:   []time.Duration{0},
		},
		{
			name: "HSetNX with wrong number of arguments",
			commands: []HTTPCommand{
				{Command: "HSETNX", Body: map[string]interface{}{"key": "key_nx_t3", "field": "field", "value": "value"}},
				{Command: "HSETNX", Body: map[string]interface{}{"key": "key_nx_t3", "field": "field", "value": "value_new"}},
				{Command: "HSETNX", Body: map[string]interface{}{"key": "key_nx_t3"}},
			},
			expected: []interface{}{float64(1), float64(0), "ERR wrong number of arguments for 'hsetnx' command"},
			delays:   []time.Duration{0, 0, 0},
		},
		{
			name: "HSetNX with wrong type",
			commands: []HTTPCommand{
				{Command: "SET", Body: map[string]interface{}{"key": "key_nx_t4", "value": "v"}},
				{Command: "HSETNX", Body: map[string]interface{}{"key": "key_nx_t4", "field": "f", "value": "v_new"}},
			},
			expected: []interface{}{"OK", "WRONGTYPE Operation against a key holding the wrong kind of value"},
			delays:   []time.Duration{0, 0},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			for i, cmd := range tc.commands {
				if tc.delays[i] > 0 {
					time.Sleep(tc.delays[i])
				}
				result, err := exec.FireCommand(cmd)
				if err != nil {
					// Check if the error message matches the expected result
					assert.Equal(t, tc.expected[i], err.Error(), "Error message mismatch for cmd %s", cmd)
				} else {
					assert.Equal(t, tc.expected[i], result, "Value mismatch for cmd %s, expected %v, got %v", cmd, tc.expected[i], result)
				}
			}
		})
	}
}

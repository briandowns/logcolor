// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"reflect"
	"strings"
)

type HTTP struct {
	ID    int   `json:"id"`
	Match Match `json:"match"`
}

func (h *HTTP) GoodWords() []string {
	return h.Match.GoodWords
}

func (h *HTTP) GoodLines() []string {
	return h.Match.GoodLines
}

func (h *HTTP) WarnWords() []string {
	return h.Match.WarnWords
}

func (h *HTTP) BadLines() []string {
	return h.Match.BadLines
}

type FTP struct {
	ID    int   `json:"id"`
	Match Match `json:"match"`
}

type Match struct {
	GoodWords []string `json:"good_words"`
	GoodLines []string `json:"good_lines"`
	WarnWords []string `json:"warn_words"`
	BadLines  []string `json:"bad_lines"`
}

func WordExists(word string, words []string) bool {
	for _, i := range words {
		if strings.Contains(word, i) {
			return true
		}
	}
	return false
}

// FieldCount will return the number of fields on a given struct
func FieldCount(i interface{}) int {
	return reflect.TypeOf(i).NumField()
}

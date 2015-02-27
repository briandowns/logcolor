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

package renderer

import (
	"reflect"
	"strings"
)

type Log interface {
	ExistsInGoodWords(word string) bool
	ExistsInWarnWords(word string) bool
	ExistsInGoodLines(word string) bool
	ExistsInBadLines(word string) bool
}

type HTTP struct {
	GoodWords []string `json:"good_words"`
	GoodLines []string `json:"good_lines"`
	WarnWords []string `json:"warn_words"`
	BadLines  []string `json:"bad_lines"`
}

type FTP struct {
	GoodWords []string `json:"good_words"`
	GoodLines []string `json:"good_lines"`
	WarnWords []string `json:"warn_words"`
	BadLines  []string `json:"bad_lines"`
}

func (h *HTTP) ExistsInGoodWords(word string) bool {
	for _, i := range h.GoodWords {
		if strings.Contains(word, i) {
			return true
		}
	}
	return false
}

func (h *HTTP) ExistsInWarnWords(word string) bool {
	for _, i := range h.WarnWords {
		if strings.Contains(word, i) {
			return true
		}
	}
	return false
}

func (h *HTTP) ExistsInGoodWords() bool {
	for _, i := range h.GoodWords {
		if strings.Contains(word, i) {
			return true
		}
	}
	return false
}

func (h *HTTP) ExistsInBadLines() bool {
	for _, i := range h.BadLines {
		if strings.Contains(word, i) {
			return true
		}
	}
	return false
}

func FieldCount(i interface{}) int {
	return reflect.TypeOf(i).NumField()
}

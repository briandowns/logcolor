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

type HTTP struct {
	GoodWords map[string][]string `json:"good_words"`
	GoodLines map[string][]string `json:"good_lines"`
	WarnWords map[string][]string `json:"warn_words"`
	BadLines  map[string][]string `json:"bad_lines"`
}

type FTP struct {
	GoodWords map[string]string   `json:"good_words"`
	WarnWords map[string][]string `json:"warn_words"`
	BadLines  map[string][]string `json:"bad_lines"`
}

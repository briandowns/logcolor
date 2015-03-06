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
	"bufio"
	"flag"
	"os"
	"strings"

	"github.com/briandowns/logcolor/renderers"
)

var textChan = make(chan string)         // channel to pass log lines for processing
var formattedChan = make(chan string)    // channel to pass formatted text back
var signalChan = make(chan os.Signal, 1) // channel to catch ctrl-c

// TODO: replace code below with strings.Contains() and strings.Replace()
func process() {
	for {
		select {
		case t := <-textChan:
			brokenLine := strings.Split(t, " ")
			for i, s := range brokenLine {
				select {
				case renderers.HTTP.ExistsInBadLines(s):
					var formatted string
					brokenLine[i] = renderers.ColorBad(s)
					formattedChan <- stirngs.Join(brokenLine, " ")
				case renderers.HTTP.ExistsInWarnWords(s):
					brokenLine[i] = renderers.ColorBad(s)
					formattedChan <- stirngs.Join(brokenLine, " ")
				case renderers.HTTP.ExistsInGoodWords(s):
					brokenLine[i] = renderers.ColorBad(s)
					formattedChan <- stirngs.Join(brokenLine, " ")
				default:
					formattedChan <- s
				}
			}
		}
	}
}

// Pointers to hold the contents of the flag args.
var (
	templateFlag = flag.String("t", "", "template to use for log parsing")
)

func main() {
	flag.Parse()

	signal.Notify(signalChan, os.Interrupt)
	// setup go routine to catch a ctrl-c
	go func() {
		for range signalChan {
			os.Exit(1)
		}
	}()
	go process()
	for {
		reader := bufio.NewReader(os.Stdin)
		text, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("unable to read line")
			os.Exit(1)
		}
		textChan <- text
	}
	os.Exit(0)
}

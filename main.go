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

	"github.com/fatih/color"
)

var textChan = make(chan string)
var signalChan = make(chan os.Signal, 1) // channel to catch ctrl-c

func colorize() {
	for {
		select {
		case t := <-textChan:
			brokenLine := strings.Split(t, " ")
			for i, s := range brokenLine {
				brokenLine[i] = color.Red(s)
			}
		}
	}
}

func main() {
	signal.Notify(signalChan, os.Interrupt)
	// setup go routine to catch a ctrl-c
	go func() {
		for range signalChan {
			os.Exit(1)
		}
	}()
	go colorize()
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

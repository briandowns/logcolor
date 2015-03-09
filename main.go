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
	"fmt"
	"os"
	"os/signal"
	"strings"

	"github.com/briandowns/logcolor/renderers"
)

var textChan = make(chan string)         // channel to pass log lines for processing
var formattedChan = make(chan string)    // channel to pass formatted text back
var signalChan = make(chan os.Signal, 1) // channel to catch ctrl-c
var stopChan = make(chan struct{})       // channel to kill the process thread

// TODO: replace code below with strings.Contains() and strings.Replace()
func processLine(log renderers.Log) {
	for {
		select {
		case t := <-textChan:
			brokenLine := strings.Split(t, " ")
			for i, s := range brokenLine {
				switch {
				case log.ExistsInBadLines(s):
					var formatted string
					brokenLine[i] = renderers.ColorBad(s)
					formattedChan <- strings.Join(brokenLine, " ")
				case log.ExistsInWarnWords(s):
					brokenLine[i] = renderers.ColorBad(s)
					formattedChan <- strings.Join(brokenLine, " ")
				case log.ExistsInGoodWords(s):
					brokenLine[i] = renderers.ColorBad(s)
					formattedChan <- strings.Join(brokenLine, " ")
				default:
					formattedChan <- s
				}
			}
		case <-stopChan:
			return
		}
	}
}

// Pointers to hold the contents of the flag args.
var (
	templateFlag = flag.String("t", "", "template to use for log parsing")
)

const USAGE = `Usage: logcolor -t template [-h]`

func main() {
	flag.Parse()

	signal.Notify(signalChan, os.Interrupt)
	// setup go routine to catch a ctrl-c
	go func() {
		for range signalChan {
			stopChan <- struct{}{} // clean up
			os.Exit(1)
		}
	}()

	if len(*templateFlag) != 1 {
		fmt.Println(USAGE)
		os.Exit(1)
	}

	var log renderers.Log
	switch *templateFlag {
	case "http":
		h := &renderers.HTTP{}
		log = renderers.Log(h)
	case "ftp":
		f := &renderers.FTP{}
		log = renderers.Log(f)
	case "sip":
		s := &renderers.SIP{}
		log = renderers.Log(s)
	case "mysql":
		m := &renderers.MySQL{}
		log = renderers.Log(m)
	case "rsync":
		r := &renderers.Rsync{}
		log = renderers.Log(r)
	case "postgresql":
		p := &renderers.Postgresql{}
		log = renderers.Log(p)
	case "openstack":
		o := &renderers.Openstack{}
		log = renderers.Log(o)
	}

	go processLine(log)

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

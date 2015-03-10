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
)

var textChan = make(chan string)         // channel to pass log lines for processing
var formattedChan = make(chan string)    // channel to pass formatted text back
var signalChan = make(chan os.Signal, 1) // channel to catch ctrl-c
var stopChan = make(chan struct{})       // channel to kill the process thread

type logger interface{}

func processLine(log logger) {
	for {
		select {
		case t := <-textChan:
			brokenLine := strings.Split(t, " ")
			for _, s := range brokenLine {
				switch {
				case WordExists(s, log.Match.GoodWords):
					formattedChan <- strings.Join(brokenLine, " ")
				case WordExists(s, log.Match.GoodLines):
					formattedChan <- strings.Join(brokenLine, " ")
				case WordExists(s, log.Match.WarnWords):
					formattedChan <- strings.Join(brokenLine, " ")
				case WordExists(s, log.Match.BadWords):
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

	/*
		if len(*templateFlag) != 1 {
			fmt.Println(USAGE)
			os.Exit(1)
		}
	*/

	var log interface{}
	switch *templateFlag {
	case "http":
		log := logger(&HTTP{})
	case "ftp":
		//f := &FTP{}
		/*
			case "sip":
				s := &SIP{}
			case "mysql":
				m := &MySQL{}
			case "rsync":
				r := &Rsync{}
			case "postgresql":
				p := &Postgresql{}
			case "openstack":
				o := &Openstack{}
		*/
	}

	go processLine(log)

	go func() {
		for {
			reader := bufio.NewReader(os.Stdin)
			text, err := reader.ReadString('\n')
			if err != nil {
				fmt.Println("unable to read line")
				os.Exit(1)
			}
			textChan <- text
		}
	}()

	for {
		select {
		case f := <-formattedChan:
			fmt.Println(f)
		default:
			continue
		}
	}
	os.Exit(0)
}

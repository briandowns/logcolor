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
	"io"
	"os"
	"os/signal"
	"strings"

	"github.com/ActiveState/tail"
)

var textChan = make(chan string)         // channel to pass log lines for processing
var formattedChan = make(chan string)    // channel to pass formatted text back
var signalChan = make(chan os.Signal, 1) // channel to catch ctrl-c
var stopChan = make(chan struct{})       // channel to kill the process thread

type logger interface {
	GoodWords() []string
	GoodLines() []string
	WarnWords() []string
	BadLines() []string
}

func processLine(log logger) {
	for {
		select {
		case t := <-textChan:
			if len(t) == 0 {
				formattedChan <- t
			}
			brokenLine := strings.Split(t, " ")
			for _, s := range brokenLine {
				switch {
				case WordExists(s, log.GoodWords()):
					formattedChan <- strings.Join(brokenLine, " ")
				case WordExists(s, log.GoodLines()):
					formattedChan <- strings.Join(brokenLine, " ")
				case WordExists(s, log.WarnWords()):
					formattedChan <- strings.Join(brokenLine, " ")
				case WordExists(s, log.BadLines()):
					formattedChan <- strings.Join(brokenLine, " ")
				default:
					formattedChan <- s + "\n"
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
	logFileFlag  = flag.String("l", "", "log file to colorize")
)

const USAGE = `Usage: logcolor -t template -l logfile [-h]`

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

	switch *templateFlag {
	case "http":
		go processLine(logger(&HTTP{}))
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

	t, err := tail.TailFile("tail_test.log", tail.Config{Follow: true})
	if err != nil {
		os.Exit(1)
	}
	for line := range t.Lines {
		fmt.Println(line.Text)
	}
	go func() {
		for {
			reader := bufio.NewReader(os.Stdin)
			text, err := reader.ReadString('\n')
			if err == io.EOF {
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

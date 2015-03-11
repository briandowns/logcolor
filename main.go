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
	"flag"
	"fmt"
	"os"
	"os/signal"
	"strings"

	"github.com/ActiveState/tail"
)

var signalChan = make(chan os.Signal, 1) // channel to catch ctrl-c

type logger interface {
	GoodWords() []string
	GoodLines() []string
	WarnWords() []string
	BadLines() []string
}

func processLine(line string, log logger) string {
	if len(line) == 0 {
		return line
	}
	brokenLine := strings.Split(line, " ")
	for _, s := range brokenLine {
		switch {
		case WordExists(s, log.GoodWords()):
			return strings.Join(brokenLine, " ")
		case WordExists(s, log.GoodLines()):
			return strings.Join(brokenLine, " ")
		case WordExists(s, log.WarnWords()):
			return strings.Join(brokenLine, " ")
		case WordExists(s, log.BadLines()):
			return strings.Join(brokenLine, " ")
		default:
			return line
		}
	}
	return ""
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
			os.Exit(1)
		}
	}()

	/*
		if len(*templateFlag) != 1 {
			fmt.Println(USAGE)
			os.Exit(1)
		}
	*/

	var log logger
	switch *templateFlag {
	case "http":
		log = logger(&HTTP{})
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

	t, err := tail.TailFile(*logFileFlag, tail.Config{Follow: true})
	if err != nil {
		os.Exit(1)
	}
	for line := range t.Lines {
		fmt.Println(processLine(line.Text, log))
	}
	os.Exit(0)
}

package debug

import "os"
import "time"
import "strings"
import "io/ioutil"
import "runtime/pprof"

func SetupTestUtil() {
	go func() {
		for {
			input, err := ioutil.ReadFile("test.cmd")

			if err == nil && len(input) > 0 {
				ioutil.WriteFile("test.cmd", []byte(""), 0744)

				cmd := strings.Trim(string(input), " \n\r\t")

				var p *pprof.Profile

				switch cmd {
				case "lookup goroutine":
					p = pprof.Lookup("goroutine")
				case "lookup heap":
					p = pprof.Lookup("heap")
				case "lookup threadcreate":
					p = pprof.Lookup("threadcreate")
				default:
					println("unknow command: '" + cmd + "'")
				}

				if p != nil {
					file, err := os.Create("test.out")
					if err != nil {
						println("couldn't create test.out")
					} else {
						p.WriteTo(file, 2)
					}
				}
			}

			time.Sleep(2 * time.Second)
		}
	}()
}

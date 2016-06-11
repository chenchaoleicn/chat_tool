package log

import (
	"fmt"
)

func AddTestFlag(funcName string) {
	fmt.Println("---------------------------------------")
	if funcName != "" {
		fmt.Println("(" + funcName + ")")
	}
}

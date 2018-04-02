package common

import (
	"fmt"
)

const (
	MTYPE_NORMAL  = iota
	MTYPE_WARNING = iota
	MTYPE_ERROR   = iota
)

func Print(base string, mType int) {
	var message string = ""
	switch mType {

	case MTYPE_NORMAL:
		message = fmt.Sprintf("[+] %s\n", base)
		break
	case MTYPE_WARNING:
		message = fmt.Sprintf("[*] %s\n", base)
		break
	case MTYPE_ERROR:
		message = fmt.Sprintf("[!] %s\n", base)
		break
	}
	fmt.Printf("%s", message)
}

func ErrorCheck(location string, err error) bool {
	if err != nil {
		fmt.Printf("[!] Error thrown @ %s: %s.\n", location, err)
		return true
	}

	return false
}

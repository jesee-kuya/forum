package util

import (
	"fmt"
	"os"
	"strconv"
)

/*
ValidatePort check for the validity of the choice of port used to run and access the program's server.
*/
func ValidatePort() (string, error) {
	var port string

	switch len(os.Args) {
	case 1:
		port = ":9000"
	case 2:
		portNum, err := strconv.Atoi(os.Args[1])
		if err != nil {
			return "", fmt.Errorf("error converting %v to int: %v", os.Args[1], err)
		}

		if portNum < 1024 || portNum > 65535 {
			return "", fmt.Errorf("use a range between 1024 and 65535")
		}
		port = ":" + strconv.Itoa(portNum)
	default:
		return "", fmt.Errorf("usage: 'go run main.go' OR 'go run main.go [PORT]'")

	}
	return port, nil
}

package parser

import (
	"errors"
	"fmt"
)

// 1 to 65635 is the range of ports
const (
	defaultPort = 8989
	minPort     = 1
	maxPort     = 65535
)

func GetPortNumber(args []string) (int, error) {
	if len(args) == 2 {
		var num int

		_, err := fmt.Sscanf(args[1], "%d", &num)
		if err != nil {
			return 0, errors.New("[USAGE]: ./TCPChat $port\nInvalid port: must be a number")
		}

		if num >= minPort && num <= maxPort {
			return num, nil
		} else {
			return 0, errors.New("[USAGE]: ./TCPChat $port\nPort must be between 1 and 65535")
		}

	} else if len(args) > 2 {
		return 0, fmt.Errorf("[USAGE]: ./TCPChat $port\n")
	}
	return defaultPort, nil
}

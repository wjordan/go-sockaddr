package sockaddr

import (
	"errors"
	"os/exec"
)

// defaultWindowsIfNameCmd is the comamnd to run on Windows to get the default
// interface.
func defaultWindowsIfNameCmd() []string {
	return []string{"netstat", "-rn"}
}

// getDefaultIfName is a Windows-specific function for extracting the name of
// the interface from `netstat -rn`.
func getDefaultIfName() (string, error) {
	var cmd []string = defaultWindowsIfNameCmd()
	out, err := exec.Command(cmd[0], cmd[1:]...).Output()
	if err != nil {
		return "", err
	}

	var ifName string
	if ifName, err = parseDefaultIfNameFromWindowsNetstatRN(string(out)); err != nil {
		return "", errors.New("No default interface found")
	}
	return ifName, nil
}

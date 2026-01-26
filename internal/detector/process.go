package detector

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/shirou/gopsutil/v3/process"
)

var (
	ErrProcessNotFound = errors.New("antigravity language server process not found")
	ErrPortNotFound    = errors.New("port not found in process arguments or connections")
	portArgRegex       = regexp.MustCompile(`--port[=\s](\d+)`)
)

type ProcessInfo struct {
	Pid     int32
	Cmdline string
	Port    int
}

// FindAntigravityProcess scans for the Antigravity Language Server process
func FindAntigravityProcess() (*ProcessInfo, error) {
	processes, err := process.Processes()
	if err != nil {
		return nil, fmt.Errorf("failed to list processes: %w", err)
	}

	for _, p := range processes {
		// Optimization: Check name first if possible, but cmdline is more reliable for node processes
		cmdline, err := p.Cmdline()
		if err != nil {
			continue
		}

		// Look for the language server process
		// Common patterns: "antigravity-language-server", or specific binary
		if strings.Contains(cmdline, "antigravity-language-server") {
			info := &ProcessInfo{
				Pid:     p.Pid,
				Cmdline: cmdline,
			}
			
			// Try to extract port
			port, err := extractPort(p, cmdline)
			if err == nil {
				info.Port = port
			}
			
			return info, nil
		}
	}

	return nil, ErrProcessNotFound
}

func extractPort(p *process.Process, cmdline string) (int, error) {
	// 1. Try regex on cmdline
	matches := portArgRegex.FindStringSubmatch(cmdline)
	if len(matches) > 1 {
		return strconv.Atoi(matches[1])
	}

	// 2. Try scanning listening ports (Linux specific mainly)
	// This requires permissions usually, but might work if owned by same user
	connections, err := p.Connections()
	if err != nil {
		return 0, err
	}

	for _, conn := range connections {
		if conn.Status == "LISTEN" && conn.Laddr.Port > 0 {
			// If there are multiple listening ports, this might be ambiguous.
			// But usually LS has one main port.
			return int(conn.Laddr.Port), nil
		}
	}
	
	// Fallback: check net connections via net package if process methods fail?
	// process.Connections() wraps net.ConnectionsPid

	return 0, ErrPortNotFound
}

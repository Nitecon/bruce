package exe

import (
	"fmt"
	"github.com/rs/zerolog/log"
	"os/exec"
	"regexp"
	"runtime"
	"strings"
)

type Execution struct {
	input       string
	fields      []string
	useSudo     bool
	outputStr   string
	errorStr    string
	isError     bool
	cmnd        string
	args        []string
	regex       *regexp.Regexp
	regexString string
	err         error
}

func Run(c string, useSudo bool) *Execution {
	e := &Execution{}
	e.input = c
	e.fields = strings.Fields(c)
	if useSudo {
		e.useSudo = true
		e.cmnd = "sudo"
		e.args = e.fields
	} else {
		e.cmnd = e.fields[0]
		e.args = e.fields[1:]
	}
	cmd := exec.Command(e.cmnd, e.args...)
	d, err := cmd.CombinedOutput()
	if err != nil {
		e.isError = true
	}
	e.outputStr = strings.TrimSuffix(strings.TrimLeft(strings.TrimRight(string(d), " "), " "), "\n")
	if err != nil {
		e.errorStr = strings.TrimSuffix(strings.TrimLeft(strings.TrimRight(err.Error(), " "), " "), "\n")
	}
	return e
}

// Failed will return true if the command returned an error.
func (e *Execution) Failed() bool {
	return e.isError
}

// ContainsLC will check if either the output or error strings contain a value all lower case.
func (e *Execution) ContainsLC(c string) bool {
	if strings.Contains(strings.ToLower(e.outputStr), c) {
		return true
	}
	if strings.Contains(strings.ToLower(e.errorStr), c) {
		return true
	}
	return false
}

// Get will return the currently populated Output string even if it's empty
func (e *Execution) Get() string {
	return e.outputStr
}

// GetErrStr will return the currently populated error output string even if it's empty
func (e *Execution) GetErrStr() string {
	return e.errorStr
}

// GetErr will return the actual error
func (e *Execution) GetErr() error {
	return e.err
}

// SetRegex will compile a regex for RegexMatch to run.
func (e *Execution) SetRegex(re string) (*regexp.Regexp, error) {
	rc, err := regexp.Compile(re)
	if err != nil {
		return nil, err
	}
	e.regex = rc
	return rc, err
}

// RegexMatch will check if either the output or error strings match the previous regex.
func (e *Execution) RegexMatch() bool {
	if e.regex == nil {
		log.Error().Err(fmt.Errorf("chain this after SetRegex(re string)")).Msg("use SetRegex first")
		return false
	}
	if e.regex.MatchString(e.outputStr) {
		return true
	}
	if e.regex.MatchString(e.errorStr) {
		return true
	}
	return false
}

func HasExecInPath(name string) string {
	if runtime.GOOS == "linux" {
		hasPkg := Run(fmt.Sprintf("which %s", name), false)
		log.Debug().Msgf("Output of HasExec: %s", hasPkg.outputStr)
		log.Debug().Msgf("Error of HasExec: %s", hasPkg.errorStr)
		if hasPkg.Get() != "" {
			return hasPkg.Get()
		}
	}
	return ""
}

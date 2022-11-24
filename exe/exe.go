package exe

import (
	"bruce/config"
	"bruce/random"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/rs/zerolog/log"
	"io"
	"os"
	"os/exec"
	"os/user"
	"path"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

type Execution struct {
	input       string
	fields      []string
	useSudo     bool
	outputStr   string
	isError     bool
	cmnd        string
	args        []string
	regex       *regexp.Regexp
	regexString string
	err         error
}

func GetFileChecksum(fname string) (string, error) {
	hasher := sha256.New()
	s, err := os.ReadFile(fname)
	hasher.Write(s)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(hasher.Sum(nil)), nil
}

func DeleteFile(src string) error {
	err := os.Remove(src)
	if err != nil {
		log.Error().Err(err).Msgf("delete failure with: %s", src)
		return err
	}
	return nil
}

func FileExists(src string) bool {
	if _, err := os.Stat(src); !os.IsNotExist(err) {
		return true
	}
	return false
}

func CopyFile(src, dst string, makedirs bool) error {
	source, err := os.Open(src)
	if err != nil {
		log.Debug().Msgf("copy fail, src does not exist: %s", src)
		return err
	}
	if makedirs {
		err = os.MkdirAll(path.Dir(dst), 0775)
		if err != nil {
			log.Error().Err(err).Msgf("cannot create directories for %s", dst)
		}
	}
	destination, err := os.Create(dst)
	if err != nil {
		log.Error().Err(err).Msgf("copy fail, cannot create destination file: %s", dst)
		return err
	}
	defer destination.Close()
	buf := make([]byte, 4096)
	for {
		n, err := source.Read(buf)
		if err != nil && err != io.EOF {
			return err
		}
		if n == 0 {
			break
		}
		if _, err := destination.Write(buf[:n]); err != nil {
			return err
		}
	}
	log.Info().Msgf("copy complete %s to: %s", src, dst)
	return nil
}

func EchoToFile(cmd string) string {
	randFileName := fmt.Sprintf("%s%c%s.sh", config.Get().Template.TempDir, os.PathSeparator, random.String(16))
	// Create the directory not just temp
	err := os.MkdirAll(path.Dir(randFileName), 0775)
	if err != nil {
		log.Error().Err(err).Msgf("cannot create directories for %s", path.Dir(randFileName))
	}

	tempFile, err := os.Create(randFileName)
	if err != nil {
		log.Error().Err(err).Msgf("temp file creation failed for: %s", randFileName)
		return ""
	}
	size, err := tempFile.WriteString("#!/bin/sh\n" + cmd + "\n")
	if err != nil {
		log.Error().Err(err).Msgf("could not write cmd: %s to file %s", cmd, randFileName)
		return ""
	}
	log.Debug().Msgf("wrote %db in %s", size, randFileName)
	tempFile.Close()
	return randFileName
}

func SetOwnership(item config.OwnerShip) error {
	grp, err := user.LookupGroupId(item.Group)
	if err != nil {
		log.Error().Msgf("cannot lookup group for %s", item.Group)
		return err
	}
	gid, err := strconv.Atoi(grp.Gid)
	if err != nil {
		log.Error().Msgf("not a valid group id number to convert to int")
		return err
	}
	usr, err := user.Lookup(item.Owner)
	if err != nil {
		log.Error().Msgf("cannot lookup user for %s", item.Owner)
		return err
	}
	uid, err := strconv.Atoi(usr.Uid)
	if err != nil {
		log.Error().Msgf("not a valid user id number to convert to int")
		return err
	}
	if item.Recursive {
		err := filepath.Walk(item.Path, func(path string, f os.FileInfo, err error) error {
			return os.Chown(path, uid, gid)
		})
		if err != nil {
			log.Error().Err(err).Msg("could not recursively set ownership")
			return err
		}
	} else {
		return os.Chown(item.Path, uid, gid)
	}
	return nil
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
		e.err = fmt.Errorf("%s", strings.TrimSuffix(strings.TrimLeft(strings.TrimRight(err.Error(), " "), " "), "\n"))
	}
	return e
}

// Failed will return true if the command returned an error.
func (e *Execution) Failed() bool {
	return e.isError
}

// ContainsLC will check if either the output or error strings contain a value all lower case.
func (e *Execution) ContainsLC(c string) bool {
	if strings.Contains(strings.ToLower(e.Get()), c) {
		return true
	}
	if strings.Contains(strings.ToLower(e.GetErrStr()), c) {
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
	if e.err != nil {
		return e.err.Error()
	}
	return ""
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
	if e.regex.MatchString(e.Get()) {
		return true
	}
	if e.regex.MatchString(e.GetErrStr()) {
		return true
	}
	return false
}

func HasExecInPath(name string) string {
	path, err := exec.LookPath(name)
	if errors.Is(err, exec.ErrDot) {
		err = nil
	}
	if err != nil {
		log.Error().Err(err).Msgf("error searching for %s in path", name)
		return ""
	}
	return path
}

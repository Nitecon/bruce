package templates

import (
	"bruce/loader"
	"bruce/system"
	"bytes"
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"github.com/rs/zerolog/log"
	"io"
	"io/fs"
	"os"
	"os/exec"
	"path"
	"strings"
	"text/template"
)

var (
	// may want to re-use this later but tbd
	templateFuncs = template.FuncMap{
		"contains": strings.Contains,
		"dump":     func(field interface{}) string { return dump(field) },
	}
)

func dump(field interface{}) string {
	buf := &bytes.Buffer{}
	spew.Fdump(buf, field)
	return buf.String()
}

func doFileBackup(backupDir, srcName string) (string, error) {
	source, err := os.Open(srcName)
	if err != nil {
		log.Debug().Msgf("source file doesn't exist, cannot backup: %s", srcName)
		return "", err
	}
	backupFileName := fmt.Sprintf("%s%c%s", backupDir, os.PathSeparator, strings.TrimLeft(srcName, string(os.PathSeparator)))
	err = os.MkdirAll(path.Dir(backupFileName), 0775)
	if err != nil {
		log.Fatal().Err(err).Msgf("cannot create temporary storage dirs")
	}
	destination, err := os.Create(backupFileName)
	if err != nil {
		log.Error().Err(err).Msgf("could not create a backup file for: %s", backupFileName)
		return "", err
	}
	defer destination.Close()
	buf := make([]byte, 4096)
	for {
		n, err := source.Read(buf)
		if err != nil && err != io.EOF {
			return "", err
		}
		if n == 0 {
			break
		}
		if _, err := destination.Write(buf[:n]); err != nil {
			return "", err
		}
	}
	log.Info().Msgf("completed backup of %s to: %s", srcName, backupFileName)
	return srcName, nil
}

func RestoreBackupFile(backupDir, srcName string) error {
	log.Debug().Msgf("preparing to restore: %s", srcName)
	backupFileName := fmt.Sprintf("%s%c%s", backupDir, os.PathSeparator, strings.TrimLeft(srcName, string(os.PathSeparator)))
	_, err := os.Stat(backupFileName)
	if os.IsNotExist(err) {
		log.Info().Msgf("local backup file does not exist removing source to replicate original state for: %s", backupFileName)
		_, exErr := os.Stat(srcName)
		if os.IsNotExist(exErr) {
			log.Info().Msgf("template that should exist does not, user removal?: %s", srcName)
			return exErr
		}
		// we return early here as the original state has been restored already.
		return os.Remove(srcName)

	}
	source, err := os.Open(backupFileName)
	if err != nil {
		log.Info().Msgf("backup source file doesn't exist: %s", srcName)
		return err
	}
	defer source.Close()
	srcInfo, err := source.Stat()
	if err != nil {
		log.Fatal().Msg("should not be an error to read the file info")

	}
	destination, err := os.OpenFile(srcName, os.O_RDWR|os.O_CREATE, srcInfo.Mode().Perm())
	if err != nil {
		log.Fatal().Err(err).Msgf("could not open backup file to write it to: %s", srcName)
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
	log.Debug().Msgf("completed restore of [%s] to: %s", backupFileName, srcName)
	return nil
}

// BackupLocal will first create a backup of all existing templates so we can revert if need be
func BackupLocal(backupDir string, tpls []system.ActionTemplate) error {
	// Backup dir should already exist so we can just check if file exists and make a backup

	for _, t := range tpls {
		log.Debug().Msgf("local backup started for: %s", t.LocalLocation)
		_, err := doFileBackup(backupDir, t.LocalLocation)
		if err != nil {
			log.Debug().Msgf("file backup failed as the local file doesn't exist yet: %s", err.Error())
		}
	}

	return nil
}

func loadTemplateFromRemote(remoteLoc string) (*template.Template, error) {
	d, err := loader.ReadRemoteFile(remoteLoc)
	if err != nil {
		log.Error().Err(err).Msgf("could not read remote template file: %s", remoteLoc)
	}
	log.Debug().Msgf("remote template read completed for: %s", remoteLoc)
	return template.New(path.Base(remoteLoc)).Parse(string(d))
}

func loadTemplateValue(v system.Vars) string {
	if v.Type == "value" {
		return v.Output
	}
	if v.Type == "command" {
		var outb, errb bytes.Buffer
		cText := strings.Fields(v.Output)
		if len(cText) > 1 {
			cmd := exec.Command(cText[0], cText[1:]...)
			cmd.Stdout = &outb
			cmd.Stderr = &errb
			err := cmd.Run()
			if err != nil {
				log.Err(err).Msg("error executing command returning error statement")
				// we don't want to put crazy errors in our templates anyway...
				return "ERROR_IN_CMD"
			}
		} else {
			cmd := exec.Command(cText[0])
			cmd.Stdout = &outb
			cmd.Stderr = &errb
			err := cmd.Run()
			if err != nil {
				log.Err(err).Msg("error executing command returning error statement")
				// we don't want to put crazy errors in our templates anyway...
				return "ERROR_IN_CMD"
			}
		}
		fmt.Println(cText[0])

		return outb.String()
	}
	// sometimes we will actually want an empty string so this is okay
	return ""
}

func doTemplateExec(local, remote string, vars []system.Vars, perms fs.FileMode) error {
	// we have the backup so now we can delete the file if it exists
	_, err := os.Stat(local)
	if os.IsNotExist(err) {
		log.Debug().Msg("local file does not exist yet no need to delete")
	} else {
		rerr := os.Remove(local)
		if rerr != nil {
			log.Error().Err(rerr).Msg("could not delete the existing file, writes may have issues...")
		}
	}

	log.Debug().Msgf("template exec starting on: %s", local)
	t, err := loadTemplateFromRemote(remote)
	if err != nil {
		log.Err(err).Msgf("cannot render %s", local)
		return err
	}
	var content = make(map[string]string)
	for _, v := range vars {
		content[v.Variable] = loadTemplateValue(v)
	}
	destination, err := os.OpenFile(local, os.O_RDWR|os.O_CREATE, perms)
	if err != nil {
		log.Error().Err(err).Msgf("could not open backup file to write it to: %s", local)
		return err
	}
	defer destination.Close()
	err = t.Execute(destination, content)
	if err != nil {
		log.Err(err).Msgf("could not write template for: %s", local)
		return err
	}
	log.Info().Msgf("completed template update: %s", local)
	return nil
}

// RenderTemplates post backup this renders the templates that have been loaded.
func RenderTemplates() {
	//wg := sync.WaitGroup{}
	for _, tpl := range system.GetSysInfo().Configuration.Templates {
		err := doTemplateExec(tpl.LocalLocation, tpl.RemoteLocation, tpl.Variables, tpl.Permissions)
		if err != nil {
			log.Debug().Err(err).Msgf("could not execute template: %s", tpl.LocalLocation)
		} /*
			wg.Add(1)
			log.Debug().Msgf("template selection for: %s", tpl.LocalLocation)
			go func() {
				defer wg.Done()
				// more here
			}()*/
	}
	//wg.Wait()
}

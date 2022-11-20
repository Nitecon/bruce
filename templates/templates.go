package templates

import (
	"bruce/config"
	"bruce/exe"
	"bruce/loader"
	"bytes"
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"github.com/rs/zerolog/log"
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
	modifiedTemplates []string
)

func dump(field interface{}) string {
	buf := &bytes.Buffer{}
	spew.Fdump(buf, field)
	return buf.String()
}

func doFileBackup(backupDir, srcName string) (string, error) {
	if exe.FileExists(srcName) {
		backupFileName := fmt.Sprintf("%s%c%s", backupDir, os.PathSeparator, strings.TrimLeft(srcName, string(os.PathSeparator)))
		cerr := exe.CopyFile(srcName, backupFileName, true)
		if cerr != nil {
			log.Fatal().Err(cerr).Msgf("cannot create backup file: %s", backupFileName)
		}
	}
	return srcName, nil
}

func RestoreBackupFile(srcName string) error {
	log.Debug().Msgf("preparing to restore: %s", srcName)
	backupDir := config.Get().Template.BackupDir
	backupFileName := fmt.Sprintf("%s%c%s", backupDir, os.PathSeparator, strings.TrimLeft(srcName, string(os.PathSeparator)))
	if exe.FileExists(backupFileName) {
		err := exe.CopyFile(backupFileName, srcName, false)
		if err != nil {
			log.Fatal().Err(err).Msgf("could not restore backup file")
			return err
		}
	}
	return nil
}

func GetBackupFileChecksum(src string) (string, error) {
	backupDir := config.Get().Template.BackupDir
	backupFileName := fmt.Sprintf("%s%c%s", backupDir, os.PathSeparator, strings.TrimLeft(src, string(os.PathSeparator)))
	return exe.GetFileChecksum(backupFileName)
}

// BackupLocal will first create a backup of all existing templates so we can revert if need be
func BackupLocal(backupDir string, tpls []config.ActionTemplate) error {
	// Backup dir should already exist so we can just check if file exists and make a backup
	for _, t := range tpls {
		log.Debug().Msgf("local backup started for: %s", t.LocalLocation)
		if exe.FileExists(t.LocalLocation) {
			_, err := doFileBackup(backupDir, t.LocalLocation)
			if err != nil {
				log.Debug().Msgf("file backup failed : %s", err.Error())
			}
		} else {
			log.Debug().Msgf("local file does not exist yet... skipping")
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

func loadTemplateValue(v config.Vars) string {
	if v.Type == "value" {
		return config.GetValueForOSHandler(v.Output)
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

func doTemplateExec(local, remote string, vars []config.Vars, perms fs.FileMode) error {
	// we have the backup so now we can delete the file if it exists
	if exe.FileExists(local) {
		exe.DeleteFile(local)
	} else {
		// check if the directories exist to render the file
		if !exe.FileExists(path.Dir(local)) {
			os.MkdirAll(path.Dir(local), 0775)
		}
	}

	log.Debug().Msgf("template exec starting on: %s", local)
	t, err := loadTemplateFromRemote(remote)
	if err != nil {
		log.Err(err).Msgf("cannot read template source %s", local)
		return err
	}
	var content = make(map[string]string)
	for _, v := range vars {
		content[v.Variable] = loadTemplateValue(v)
	}

	destination, err := os.OpenFile(local, os.O_RDWR|os.O_CREATE, 0664)
	if err != nil {
		log.Error().Err(err).Msgf("could not open backup file to write it to: %s", local)
		return err
	}
	defer destination.Close()
	err = t.Execute(destination, content)
	if err != nil {
		log.Err(err).Msgf("could not write template: %s", local)
		return err
	}
	log.Info().Msgf("template written: %s", local)
	localHash, err := exe.GetFileChecksum(local)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to get new file checksum")
	}
	backupHash, err := GetBackupFileChecksum(local)
	if err != nil {
		// no backup exists so lets add it as a changed template as it should be net new.
		modifiedTemplates = append(modifiedTemplates, local)
	}
	if localHash != backupHash {
		modifiedTemplates = append(modifiedTemplates, local)
	}
	return nil
}

// RenderTemplates post backup this renders the templates that have been loaded.
func RenderTemplates(templates []config.ActionTemplate) {
	//wg := sync.WaitGroup{}
	for _, tpl := range templates {
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

func GetModifiedTemplates() []string {
	log.Debug().Msgf("modified templates: %s", modifiedTemplates)
	return modifiedTemplates
}

package cmd

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"

	"github.com/manifoldco/promptui"
	"github.com/meinto/glow/pkg/cli/cmd/internal/command"
	"github.com/meinto/glow/pkg/cli/cmd/internal/util"
	"github.com/meinto/promter"
	"github.com/spf13/cobra"
)

type InstallCommand struct {
	command.Service
}

func (cmd *InstallCommand) PostSetup(parent command.Service) command.Service {
	parent.Add(cmd)
	return cmd
}

var installCmd = SetupInstallCommand(RootCmd)

func SetupInstallCommand(parent command.Service) command.Service {
	return command.Setup(&InstallCommand{
		&command.Command{
			Command: &cobra.Command{
				Use:   "install",
				Short: "install glow",
			},
			Run: func(cmd command.Service, args []string) {
				p := promter.NewPromter()

				flist, err := fileList(".")
				util.ExitOnErrorWithMessage("cannot get file list")(err)

				index, _, err := p.Select(
					"Select your downloaded glow file",
					flist,
				)
				util.ExitOnErrorWithMessage(fmt.Sprintf("cannot get path to glow file: %s", err))(err)

				filePath, err := filepath.Abs(flist[index])
				util.ExitOnErrorWithMessage("cannot get absolute file path")(err)

				index, err = usageOptions()
				util.ExitOnErrorWithMessage("cannot get usage option")(err)

				var newFileName string
				switch index {
				case 0:
					newFileName = "/usr/local/bin/glow"
				case 1:
					newFileName = "/usr/local/bin/git-glow"
				}

				if _, err := os.Stat(newFileName); !os.IsNotExist(err) {
					replace, err := promtReplaceFile(newFileName)
					util.ExitOnErrorWithMessage(err.Error())(err)
					if !replace {
						log.Fatal("file not replaced")
					}
				}

				err = os.Rename(filePath, newFileName)
				if err != nil {
					log.Fatal(err)
				}

				fmt.Println("successfully moved glow")
			},
		},
	}, parent)
}

func pathToGlowFile() (string, error) {
	validate := func(input string) error {
		filePath, err := filepath.Abs(input)
		if err != nil {
			return fmt.Errorf("error while creating absolute path: %s", err.Error())
		}
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			return errors.New("file does not exist")
		}
		return nil
	}

	getFileName := promptui.Prompt{
		Label:    "Name of binary: ",
		Validate: validate,
	}

	fileName, err := getFileName.Run()
	if err != nil {
		return "", err
	}

	filePath, err := filepath.Abs(fileName)
	if err != nil {
		return "", err
	}

	return filePath, nil
}

func usageOptions() (int, error) {
	prompt := promptui.Select{
		Label: "How do you want to use glow?",
		Items: []string{"global", "git plugin"},
	}

	index, _, err := prompt.Run()
	if err != nil {
		return -1, err
	}

	return index, nil
}

func promtReplaceFile(filePath string) (bool, error) {
	if _, err := os.Stat(filePath); !os.IsNotExist(err) {
		prompt := promptui.Select{
			Label: "File exists. Do you want to replace it?",
			Items: []string{"yes", "no"},
		}

		index, _, err := prompt.Run()
		if err != nil {
			return false, err
		}

		if index == 0 {
			return true, nil
		}
		return false, nil
	}
	return false, nil
}

func fileList(rootPath string) ([]string, error) {
	var files []string

	err := filepath.Walk(rootPath, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() && info.Name() != "." {
			return filepath.SkipDir
		}
		if !strings.HasPrefix(path, "glow_") {
			return nil
		}
		files = append(files, path)
		return nil
	})
	if err != nil {
		return files, errors.Wrap(err, "error creating file list")
	}

	return files, nil
}

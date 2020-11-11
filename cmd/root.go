package cmd

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/mitchellh/go-homedir"
	"github.com/murpg/gsp/pkg/lib"
	"github.com/murpg/gsp/pkg/release"
	"github.com/murpg/gsp/pkg/tpl"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	config lib.Configuration

	cfgFile string
	dryRun  bool

	rootCmd = &cobra.Command{
		Use:   "gsp",
		Short: "gsp - git simple packager",

		Run: func(cmd *cobra.Command, args []string) {
			changedFilesWithPossibleDuplicates := getChangedFileNames(config)
			changedFileNames := removeDuplicates(changedFilesWithPossibleDuplicates)

			// list changed files in console
			notes := release.New()
			notes.ChangedFiles = changedFileNames
			releaseNotesTemplate := template.Must(template.New("release-notes").Parse(tpl.ReleaseNotesTextTemplate))
			_ = releaseNotesTemplate.Execute(os.Stdout, notes)

			if dryRun {
				fmt.Println()
				fmt.Println("=> Not creating a release archive as --dry-run.")
				return
			}

			if len(changedFileNames) > 0 {
				createReleaseArchive(config, notes)
			}
		},

		DisableSuggestions: true,
		SilenceErrors:      true,
		SilenceUsage:       true,
	}
)

// Execute is an main entry point that executes rootCmd.Run function.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "Config file (default is $HOME/.gsp-config.json)")
	rootCmd.PersistentFlags().BoolVarP(&dryRun, "dry-run", "n", false, "Do a trial run without creating a release archive.")
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		gspDir, err := filepath.Abs(filepath.Dir(os.Args[0]))
		if err != nil {
			fmt.Println(err)
			os.Exit(2)
		}

		viper.SetConfigName(".gsp-config")
		viper.AddConfigPath(home)
		viper.AddConfigPath(gspDir)
		viper.AddConfigPath(".")
	}

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil {
		config = lib.LoadConfiguration(viper.ConfigFileUsed())
	}
}

func getChangedFileNames(config lib.Configuration) []string {
	var cmdOut []byte
	var cmdArgs = []string{"diff", "--name-only"}
	var err error

	cmdArgs = append(cmdArgs, fmt.Sprintf("--diff-filter=%v", config.DiffFilter))

	if len(config.GitHashNewest) > 0 || len(config.GitHashOldest) > 0 {
		cmdArgs = append(cmdArgs, config.GitHashNewest)
		if len(config.GitHashOldest) == 0 {
			cmdArgs = append(cmdArgs, "HEAD")
		} else {
			cmdArgs = append(cmdArgs, config.GitHashOldest)
		}
	} else {
		cmdArgs = append(cmdArgs, fmt.Sprintf("HEAD~%v", config.CommitsCount))
		cmdArgs = append(cmdArgs, "HEAD")
	}
	cmdArgs = append(cmdArgs, "--")

	if os.Getenv("MODE") == "DEV" {
		fmt.Printf("[DEBUG] command arguments: %v\n", cmdArgs)
	}

	gitCommand := exec.Command("git", cmdArgs...)
	gitCommand.Dir = config.RepositoryPath
	if cmdOut, err = gitCommand.Output(); err != nil {
		fmt.Printf("[ERROR] %v\n", err.Error())
		fmt.Printf("[ERROR] command output: %v\n", string(cmdOut))
		os.Exit(-1)
	}

	var files []string
	fileSlices := strings.Split(string(cmdOut), "\n")
	for _, fileSlice := range fileSlices {
		if len(fileSlice) == 0 {
			continue
		}

		if len(config.DirectoryNames) > 0 {
			for _, filteredDirectoryName := range config.DirectoryNames {
				if strings.Contains(fileSlice, filteredDirectoryName) {
					files = append(files, fileSlice)
				}
			}
		} else {
			files = append(files, fileSlice)
		}
	}

	return files
}

func removeDuplicates(s []string) []string {
	encountered := map[string]bool{}
	for v := range s {
		encountered[s[v]] = true
	}
	var result []string
	for key := range encountered {
		result = append(result, key)
	}
	return result
}

func createReleaseArchive(config lib.Configuration, notes *release.Notes) {
	if _, err := os.Stat(config.OutputPath); os.IsNotExist(err) {
		if err := os.MkdirAll(config.OutputPath, os.ModePerm); err != nil {
			fmt.Printf("[ERROR] %v\n", err.Error())
			fmt.Printf("[ERROR] Could not create folder for outputs: %v\n", config.OutputPath)
			os.Exit(-1)
		}
	}

	releaseOutputFolder := notes.ReleaseDate.Format("20060102_150405")
	releaseOutputFolderAbsPath := filepath.Join(config.OutputPath, releaseOutputFolder)
	if err := os.Mkdir(releaseOutputFolderAbsPath, os.ModePerm); err != nil {
		fmt.Printf("[ERROR] %v\n", err.Error())
		fmt.Printf("[ERROR] Could not create folder for release output: %v\n", releaseOutputFolderAbsPath)
		os.Exit(-1)
	}

	releaseFile, err := os.Create(path.Join(releaseOutputFolderAbsPath, "RELEASE.txt"))
	if err != nil {
		fmt.Printf("[ERROR] %v\n", err.Error())
		fmt.Printf("[ERROR] Could not create release text file: %v\n", releaseFile)
		os.Exit(-1)
	}
	releaseNotesTemplate := template.Must(template.New("release-notes").Parse(tpl.ReleaseNotesTextTemplate))
	_ = releaseNotesTemplate.Execute(releaseFile, notes)

	// fix line endings
	data, err := ioutil.ReadFile(path.Join(releaseOutputFolderAbsPath, "RELEASE.txt"))
	if err != nil {
		fmt.Printf("[ERROR] %v\n", err.Error())
		fmt.Printf("[ERROR] Could not fix line endings in release file: %v\n", releaseFile)
		os.Exit(-1)
	}
	data = bytes.Replace(data, []byte{10}, []byte{13, 10}, -1)
	err = ioutil.WriteFile(path.Join(releaseOutputFolderAbsPath, "RELEASE.txt"), data, 0644)
	if err != nil {
		fmt.Printf("[ERROR] %v\n", err.Error())
		fmt.Printf("[ERROR] Could not write release file with fixed line endings: %v\n", releaseFile)
		os.Exit(-1)
	}

	releaseArchive := path.Join(releaseOutputFolderAbsPath, "archive.zip")
	var changedFiles []string
	for _, changedFile := range notes.ChangedFiles {
		changedFiles = append(changedFiles, changedFile)
	}

	err = zipFiles(releaseArchive, config, changedFiles)
	if err != nil {
		fmt.Printf("[ERROR] %v\n", err.Error())
		fmt.Printf("[ERROR] Could not create release archive\n")
		os.Exit(-1)
	}
}

func zipFiles(archive string, config lib.Configuration, files []string) error {
	newZipFile, err := os.Create(archive)
	if err != nil {
		return err
	}
	defer newZipFile.Close()

	zipWriter := zip.NewWriter(newZipFile)
	defer zipWriter.Close()

	for _, file := range files {
		zipFile, err := os.Open(path.Join(config.RepositoryPath, file))
		if err != nil {
			fmt.Printf("[ERROR] %v\n", err.Error())
			fmt.Printf("[ERROR] zip file: %v\n", zipFile)
			return err
		}
		defer zipFile.Close()

		info, err := zipFile.Stat()
		if err != nil {
			return err
		}

		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}
		header.Name = file
		header.Method = zip.Deflate

		writer, err := zipWriter.CreateHeader(header)
		if err != nil {
			return err
		}
		if _, err = io.Copy(writer, zipFile); err != nil {
			return err
		}
	}
	return nil
}

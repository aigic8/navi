package main

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/alexflint/go-arg"
	"gopkg.in/yaml.v3"
)

// TODO: add placeholders
// TODO: add helps

var CONFIG_NAMES = []string{
	"navi.yaml",
	"navi.yml",
}

// TODO: add short veresions of subcommands
var (
	args struct {
		Init *CmdInit `arg:"subcommand:init"`
		Find *CmdFind `arg:"subcommand:find"`
		Add  *CmdAdd  `arg:"subcommand:add"`
	}
)

type (
	CmdInit struct {
		Path string `arg:"positional"`
	}

	CmdFind struct {
		Tags string `arg:"--tags,-t"`
		Path string `arg:"positional"`
	}

	CmdAdd struct {
		Path string `arg:"positional"`
		Tags string `arg:"--tags,-t"`
	}
)

type (
	Config struct {
		Paths map[string][]string `yaml:"paths"`
	}
)

func main() {
	// FIXME: error handling for argument parsing
	arg.MustParse(&args)

	wd, err := os.Getwd()
	if err != nil {
		// FIXME: remove panic
		panic(err)
	}

	if args.Init != nil {
		configPath := wd
		if args.Init.Path != "" {
			configPath = path.Join(wd, args.Init.Path)
		}

		// TODO: handle the case where the path does not exist
		// TODO: handle the case where the file already exists

		pathHasConfigFile := strHasOneOfTheSuffixes(configPath, CONFIG_NAMES)

		if !pathHasConfigFile {
			configPath = path.Join(configPath, CONFIG_NAMES[0])
		}

		relativePath := args.Init.Path
		if !pathHasConfigFile {
			relativePath = path.Join(relativePath, CONFIG_NAMES[0])
		}

		baseConfig := Config{Paths: map[string][]string{}}
		baseConfigBytes, err := yaml.Marshal(&baseConfig)
		if err != nil {
			// FIXME: remove panic
			panic(err)
		}
		if err := os.WriteFile(configPath, baseConfigBytes, 0644); err != nil {
			// FIXME: remove panic
			panic(err)
		}
		fmt.Printf("File '%s' was generated.\n", relativePath)
	} else if args.Find != nil {
		configFilePath, err := findCurrentConfig()
		if err != nil {
			// FIXME: remove panic
			panic(err)
		}
		if configFilePath == "" {
			// FIXME: remove panic
			panic("config not found")
		}

		fileData, err := os.ReadFile(configFilePath)
		if err != nil {
			// FIXME: remove panic
			panic(err)
		}

		config := Config{}
		if err = yaml.Unmarshal(fileData, &config); err != nil {
			// FIXME: remove panic
			panic(err)
		}

		// TODO: add config validation

		resultPathsMap := map[string]bool{}
		// TODO: add logical operators for tags
		argsTags := strings.Split(args.Find.Tags, ",")
		argsTagsMap := ArrToMap(argsTags)
		for configPath, configTags := range config.Paths {
			for _, configTag := range configTags {
				if _, ok := argsTagsMap[configTag]; ok {
					resultPathsMap[configPath] = true
				}
			}
		}

		// TODO: add support for tagging dirs

		resultPaths := []string{}
		configFileDir := path.Dir(configFilePath)
		for resultPath := range resultPathsMap {
			if !filepath.IsAbs(resultPath) {
				resultPath = path.Join(configFileDir, resultPath)
			}
			// FIXME:  change the base path to the current wd
			// For example, if the in config file is './this/path/is/mine.md'
			// and wd is './this/path'
			// then, the result path should become './is/mine.md'

			// TODO: check the path filter (positional argument)
			resultPaths = append(resultPaths, resultPath)
		}
		for _, resultPath := range resultPaths {
			fmt.Println(resultPath)
		}
	} else if args.Add != nil {
		configFilePath, err := findCurrentConfig()
		if err != nil {
			// FIXME: remove panic
			panic(err)
		}

		configData, err := os.ReadFile(configFilePath)
		if err != nil {
			// FIXME: remove panic
			panic(err)
		}

		var config Config
		if err := yaml.Unmarshal(configData, &config); err != nil {
			// FIXME: remove panic
			panic(err)
		}

		cwd, err := os.Getwd()
		if err != nil {
			// FIXME: remove panic
			panic(err)
		}

		// TODO: check if the adding path is file and exists (maybe?)
		addingPathAbs := makePathAbsolute(cwd, args.Add.Path)
		configFileDirEndingSep := path.Dir(configFilePath) + string(os.PathSeparator)
		if !strings.HasPrefix(addingPathAbs, configFileDirEndingSep) {
			// FIXME: remove panic
			panic("out of scope path")
		}
		addingPathConfigRel := strings.TrimPrefix(addingPathAbs, configFileDirEndingSep)
		fmt.Printf("adding path rel: %s", addingPathConfigRel)

		// TODO: validate and trim tags
		addingTagsArr := strings.Split(args.Add.Tags, ",")
		// FIXME: all of the paths in config should be converted to a singular format to be compared with addingPathConfigRel
		if configTags, ok := config.Paths[addingPathConfigRel]; ok {
			configTagsMap := ArrToMap(configTags)
			for _, tag := range addingTagsArr {
				configTagsMap[tag] = true
			}
			config.Paths[addingPathConfigRel] = MapKeys(configTagsMap)
		} else {
			config.Paths[addingPathConfigRel] = addingTagsArr
		}

		configBytes, err := yaml.Marshal(&config)
		if err != nil {
			// FIXME: remove panic
			panic(err)
		}
		// TODO: create a backup for the last file in case of breaking anything
		if err := os.WriteFile(configFilePath, configBytes, 0644); err != nil {
			// FIXME: remove panic
			panic(err)
		}
	}
	// TODO: add remove tags

}

func strHasOneOfTheSuffixes(in string, suffixes []string) bool {
	for _, suffix := range suffixes {
		if strings.HasSuffix(in, suffix) {
			return true
		}
	}
	return false
}

func findCurrentConfig() (string, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("finding current config: %w", err)
	}

	currPath := cwd
	for {
		for _, configName := range CONFIG_NAMES {
			configPath := path.Join(cwd, configName)
			isPathFile, err := IsFileAndExists(configPath)
			if err != nil {
				return "", fmt.Errorf("finding current config: %w", err)
			}
			if isPathFile {
				return configPath, nil
			}
		}
		currPath = path.Dir(currPath)
		// NOTE: check on windows if the separator is the result when the path reaches the OS
		// TODO: fix when traverse so high that we don't have permission, for linux it would be '/root/' for example
		if len(currPath) == 1 && (currPath[0] == os.PathSeparator || currPath[0] == '.') {
			return "", nil
		}
	}
}

func IsFileAndExists(inputPath string) (bool, error) {
	pathStat, err := os.Stat(inputPath)
	if err != nil {
		if !os.IsNotExist(err) {
			return false, err
		}
		return false, nil
	}
	return !pathStat.IsDir(), nil
}

func ArrToMap[T comparable](arr []T) map[T]bool {
	res := map[T]bool{}
	for _, item := range arr {
		res[item] = true
	}
	return res
}

func MapKeys[T comparable, Y any](inMap map[T]Y) []T {
	res := []T{}
	for item := range inMap {
		res = append(res, item)
	}
	return res
}

func makePathAbsolute(cwd, in string) string {
	if filepath.IsAbs(in) {
		return in
	}
	return path.Join(cwd, in)
}

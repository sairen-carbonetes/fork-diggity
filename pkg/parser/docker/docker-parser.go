package docker

import (
	"encoding/json"
	"errors"
	"os"
	"strings"

	"github.com/carbonetes/diggity/pkg/model"
	"github.com/carbonetes/diggity/pkg/parser/bom"
)

var (
	// ImageInfo docker image information
	ImageInfo model.ImageInfo
	// DockerManifest docker manifest json file
	dockerManifest []model.DockerManifest = make([]model.DockerManifest, 0)
	// DockerConfig docker config json file
	dockerConfig model.DockerConfig = model.DockerConfig{}
)

// ParseDockerProperties appends docker json files to parser.Result
func ParseDockerProperties(req *bom.ParserRequirements) {
	tarDirectory, err := os.Open(*req.Dir)
	if err != nil {
		if len(*req.Arguments.Dir) > 0 {
			tarDirectory, err = os.Open(*req.Arguments.Dir)
			if err != nil {
				err = errors.New("docker-parser: " + err.Error())
				*req.Errors = append(*req.Errors, err)
			}
		} else {
			err = errors.New("docker-parser: " + err.Error())
			*req.Errors = append(*req.Errors, err)
		}
	}
	files, err := getJSONFilesFromDir(tarDirectory.Name())
	if err != nil {
		err = errors.New("docker-parser: " + err.Error())
		*req.Errors = append(*req.Errors, err)
	}

	for _, jsonFile := range files {
		jsonFile, err := os.Open(jsonFile)
		if err != nil {
			err = errors.New("docker-parser: " + err.Error())
			*req.Errors = append(*req.Errors, err)
		}
		jsonparser := json.NewDecoder(jsonFile)
		if strings.Contains(jsonFile.Name(), "manifest") {
			if err := jsonparser.Decode(&dockerManifest); err != nil {
				err = errors.New("docker-parser: " + err.Error())
				*req.Errors = append(*req.Errors, err)
			}
		} else {
			if err := jsonparser.Decode(&dockerConfig); err != nil {
				err = errors.New("docker-parser: " + err.Error())
				*req.Errors = append(*req.Errors, err)
			}
		}
	}

	ImageInfo = model.ImageInfo{
		DockerConfig:   dockerConfig,
		DockerManifest: dockerManifest,
	}

	req.SBOM.ImageInfo = ImageInfo

	defer req.WG.Done()
}

// Get JSON files from extracted image
func getJSONFilesFromDir(root string) ([]string, error) {
	var files []string
	fileInfo, err := os.ReadDir(root)
	if err != nil {
		return files, err
	}

	for _, file := range fileInfo {
		if strings.HasSuffix(file.Name(), ".json") {
			files = append(files, root+string(os.PathSeparator)+file.Name())
		}
	}
	return files, nil
}

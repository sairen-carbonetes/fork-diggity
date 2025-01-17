package conan

import (
	"bufio"
	"encoding/json"
	"errors"
	"os"
	"strings"

	"github.com/carbonetes/diggity/internal/cpe"
	"github.com/carbonetes/diggity/pkg/model"
	"github.com/carbonetes/diggity/pkg/model/metadata"
	"github.com/carbonetes/diggity/pkg/parser/bom"
	"github.com/carbonetes/diggity/pkg/parser/util"

	"github.com/google/uuid"
)

const (
	conan       = "conan"
	conanFile   = "conanfile.txt"
	conanLock   = "conan.lock"
	requiresTag = "[requires]"
)

// ConanLockMetadata conan.lock metadata type
var conanLockMetadata metadata.ConanLockMetadata

// FindConanPackagesFromContent Find Conan packages in the file content
func FindConanPackagesFromContent(req *bom.ParserRequirements) {
	if util.ParserEnabled(conan, req.Arguments.EnabledParsers) {
		for _, content := range *req.Contents {
			if strings.Contains(content.Path, conanFile) {
				if err := readConanFileContent(&content, req.SBOM.Packages); err != nil {
					err = errors.New("conan-parser: " + err.Error())
					*req.Errors = append(*req.Errors, err)
				}
			}
			if strings.Contains(content.Path, conanLock) {
				if err := readConanLockContent(&content, req.SBOM.Packages); err != nil {
					err = errors.New("conan-parser: " + err.Error())
					*req.Errors = append(*req.Errors, err)
				}
			}
		}
	}
	defer req.WG.Done()
}

// Read conanfile.txt contents
func readConanFileContent(location *model.Location, pkgs *[]model.Package) error {
	file, err := os.Open(location.Path)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	var requires bool

	for scanner.Scan() {
		values := scanner.Text()

		// Detect requires section
		if strings.Contains(values, requiresTag) {
			requires = true
		}

		// Parse conan package metadata
		if requires && strings.Contains(values, "/") {
			*pkgs = append(*pkgs, *initConanPackage(location, values))
		}

		// End of requires section
		if !strings.Contains(values, requiresTag) && !strings.Contains(values, "/") {
			requires = false
		}

	}
	return nil
}

// Parse conan.lock contents
func readConanLockContent(location *model.Location, pkgs *[]model.Package) error {
	file, err := os.ReadFile(location.Path)
	if err != nil {
		return err
	}

	if err = json.Unmarshal(file, &conanLockMetadata); err != nil {
		return err
	}

	if len(conanLockMetadata.GraphLock.Nodes) > 0 {
		for _, conanPkg := range conanLockMetadata.GraphLock.Nodes {
			if conanPkg.Ref != "" {
				*pkgs = append(*pkgs, *initConanPackage(location, conanPkg))
			}
		}
	}

	return nil
}

// Init Conan Package
func initConanPackage(location *model.Location, conanMetadata interface{}) *model.Package {
	pkg := new(model.Package)
	pkg.ID = uuid.NewString()

	// Get conan package name, version, and metadata based on parsed metadata type
	var name, version string
	switch md := conanMetadata.(type) {
	case string:
		name, version = conanNameVersion(md)
		pkg.Metadata = metadata.ConanMetadata{
			Name:    name,
			Version: version,
		}
	case metadata.ConanLockNode:
		name, version = conanNameVersion(md.Ref)
		pkg.Metadata = md
	}

	pkg.Name = name
	pkg.Version = version
	pkg.Path = pkg.Name
	pkg.Type = conan
	pkg.Locations = append(pkg.Locations, model.Location{
		Path:      util.TrimUntilLayer(*location),
		LayerHash: location.LayerHash,
	})
	pkg.Licenses = []string{}

	// get purl
	parseConanPackageURL(pkg)

	// get CPEs
	cpe.NewCPE23(pkg, "", pkg.Name, pkg.Version)

	return pkg
}

// Parse PURL
func parseConanPackageURL(pkg *model.Package) {
	pkg.PURL = model.PURL("pkg" + ":" + conan + "/" + pkg.Name + "@" + pkg.Version)
}

// Get Name and Version from package or node ref metadata
func conanNameVersion(ref string) (name string, version string) {
	var nv []string
	if strings.Contains(ref, "@") {
		nv = strings.Split(ref, "@")
	} else if strings.Contains(ref, "#") {
		nv = strings.Split(ref, "#")
	} else {
		nv = append(nv, ref)
	}

	result := strings.Split(nv[0], "/")
	name = result[0]
	version = result[1]

	// cleanup version
	if strings.Contains(version, "[") && strings.Contains(version, "]") {
		version = strings.Replace(version, "[", "", -1)
		version = strings.Replace(version, "]", "", -1)
	}

	return name, version
}

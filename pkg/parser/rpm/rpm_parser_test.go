package rpm

import (
	"path/filepath"
	"testing"

	"github.com/carbonetes/diggity/pkg/model"
	"github.com/carbonetes/diggity/pkg/model/metadata"

	rpmdb "github.com/knqyf263/go-rpmdb/pkg"
)

type (
	FormatVendorResult struct {
		input    string
		expected string
	}

	FormatLicensesResult struct {
		pkg      *model.Package
		licenses string
		expected []string
	}

	RpmMetadataResult struct {
		pkg      *model.Package
		rpmPkg   *rpmdb.PackageInfo
		expected metadata.RPMMetadata
	}

	RpmPurlResult struct {
		pkg      *model.Package
		arch     string
		expected model.PURL
	}

	InitRpmPackageResult struct {
		pkg      *model.Package
		location *model.Location
		rpmdb    *rpmdb.PackageInfo
		expected *model.Package
	}
)

var (
	rpmLocation1 = model.Location{
		Path:      filepath.Join("AppData", "Local", "Temp", "3175519915", "diggity-tmp-614678a1-5579-42fb-8e8f-0d8e2101c803", "69a15d957a7a6f77e3fe31f330da5f4b6b582f228917a713a7a9e59449a3f413", "var", "lib", "rpm", "Packages"),
		LayerHash: "69a15d957a7a6f77e3fe31f330da5f4b6b582f228917a713a7a9e59449a3f413",
	}
	rpmLocation2 = model.Location{
		Path:      filepath.Join("AppData", "Local", "Temp", "2256324509", "diggity-tmp-32db5a22-f1f3-4603-8b53-25a16418dfed", "99d5c4b75475235491986963958036ff26733f018bcafa2758534f235cefeaa2", "var", "lib", "rpm", "Packages"),
		LayerHash: "99d5c4b75475235491986963958036ff26733f018bcafa2758534f235cefeaa2",
	}
	rpmLocation3 = model.Location{
		Path:      filepath.Join("AppData", "Local", "Temp", "632174759", "diggity-tmp-535a9533-57f9-4562-b092-f982ccfeab3c", "d1fd2cca7a7751ca9786b088cf639e65088fa0bda34492bb5ba292c32195461a", "var", "lib", "rpm", "Packages"),
		LayerHash: "d1fd2cca7a7751ca9786b088cf639e65088fa0bda34492bb5ba292c32195461a",
	}

	epoch = 3 // Epoch Num

	rpmdb1 = rpmdb.PackageInfo{
		Release:   "14.el8",
		Arch:      "x86_64",
		SourceRpm: "lzo-2.08-14.el8.src.rpm",
		License:   "GPLv2+",
		Size:      198757,
		Name:      "lzo",
		PGP:       "RSA/SHA256, Tue Jul  2 00:01:31 2019, Key ID 05b555b38483c65d",
		Summary:   "Data compression library with very fast (de)compression",
		Vendor:    "CentOS",
		Version:   "2.08",
	}

	rpmdb2 = rpmdb.PackageInfo{
		Release:   "2.fc29",
		Arch:      "x86_64",
		SourceRpm: "p11-kit-0.23.15-2.fc29.src.rpm",
		License:   "BSD",
		Size:      506497,
		Name:      "p11-kit-trust",
		PGP:       "RSA/SHA256, Tue Feb 19 02:39:25 2019, Key ID a20aa56b429476b4",
		Summary:   "System trust module from p11-kit",
		Vendor:    "Fedora Project",
		Version:   "0.23.15",
	}

	rpmdb3 = rpmdb.PackageInfo{
		Release:   "19.el7",
		Arch:      "x86_64",
		SourceRpm: "hardlink-1.0-19.el7.src.rpm",
		License:   "GPL+",
		Size:      16545,
		Name:      "hardlink",
		PGP:       "RSA/SHA256, Tue Apr  1 17:48:32 2014, Key ID 199e2f91fd431d51",
		Summary:   "Create a tree of hardlinks",
		Vendor:    "Red Hat, Inc.",
		Version:   "1.0",
		Epoch:     &epoch,
	}

	rpmPackage1 = model.Package{
		Name:    "lzo",
		Type:    rpmType,
		Version: "2.08-14.el8",
		Path:    rpmPackagesPath,
		Locations: []model.Location{
			{
				Path:      filepath.Join("var", "lib", "rpm", "Packages"),
				LayerHash: "69a15d957a7a6f77e3fe31f330da5f4b6b582f228917a713a7a9e59449a3f413",
			},
		},
		Description: "Data compression library with very fast (de)compression",
		Licenses: []string{
			"GPLv2+",
		},
		CPEs: []string{
			"cpe:2.3:a:centos:lzo:2.08-14.el8:*:*:*:*:*:*:*",
			"cpe:2.3:a:lzo:lzo:2.08-14.el8:*:*:*:*:*:*:*",
		},
		PURL: model.PURL("pkg:rpm/lzo@2.08-14.el8?arch=x86_64"),
		Metadata: metadata.RPMMetadata{
			Release:      "14.el8",
			Architecture: "x86_64",
			SourceRpm:    "lzo-2.08-14.el8.src.rpm",
			License:      "GPLv2+",
			Size:         198757,
			Name:         "lzo",
			PGP:          "RSA/SHA256, Tue Jul  2 00:01:31 2019, Key ID 05b555b38483c65d",
			Summary:      "Data compression library with very fast (de)compression",
			Vendor:       "CentOS",
			Version:      "2.08",
		},
	}
	rpmPackage2 = model.Package{
		Name:    "p11-kit-trust",
		Type:    rpmType,
		Version: "0.23.15-2.fc29",
		Path:    rpmPackagesPath,
		Locations: []model.Location{
			{
				Path:      filepath.Join("var", "lib", "rpm", "Packages"),
				LayerHash: "99d5c4b75475235491986963958036ff26733f018bcafa2758534f235cefeaa2",
			},
		},
		Description: "System trust module from p11-kit",
		Licenses: []string{
			"BSD",
		},
		CPEs: []string{
			"cpe:2.3:a:fedoraproject:p11-kit-trust:0.23.15-2.fc29:*:*:*:*:*:*:*",
			"cpe:2.3:a:fedoraproject:p11_kit-trust:0.23.15-2.fc29:*:*:*:*:*:*:*",
			"cpe:2.3:a:fedoraproject:p11_kit_trust:0.23.15-2.fc29:*:*:*:*:*:*:*",
			"cpe:2.3:a:p11-kit-trust:p11-kit-trust:0.23.15-2.fc29:*:*:*:*:*:*:*",
			"cpe:2.3:a:p11-kit-trust:p11_kit-trust:0.23.15-2.fc29:*:*:*:*:*:*:*",
			"cpe:2.3:a:p11-kit-trust:p11_kit_trust:0.23.15-2.fc29:*:*:*:*:*:*:*",
			"cpe:2.3:a:p11_kit-trust:p11_kit_trust:0.23.15-2.fc29:*:*:*:*:*:*:*",
			"cpe:2.3:a:p11_kit-trust:p11-kit_trust:0.23.15-2.fc29:*:*:*:*:*:*:*",
			"cpe:2.3:a:p11_kit-trust:p11-kit-trust:0.23.15-2.fc29:*:*:*:*:*:*:*",
			"cpe:2.3:a:p11_kit_trust:p11-kit-trust:0.23.15-2.fc29:*:*:*:*:*:*:*",
			"cpe:2.3:a:p11_kit_trust:p11_kit-trust:0.23.15-2.fc29:*:*:*:*:*:*:*",
			"cpe:2.3:a:p11_kit_trust:p11_kit_trust:0.23.15-2.fc29:*:*:*:*:*:*:*",
			"cpe:2.3:a:p11:p11_kit_trust:0.23.15-2.fc29:*:*:*:*:*:*:*",
			"cpe:2.3:a:p11:p11-kit_trust:0.23.15-2.fc29:*:*:*:*:*:*:*",
			"cpe:2.3:a:p11:p11-kit-trust:0.23.15-2.fc29:*:*:*:*:*:*:*",
		},
		PURL: model.PURL("pkg:rpm/p11-kit-trust@0.23.15-2.fc29?arch=x86_64"),
		Metadata: metadata.RPMMetadata{
			Release:      "2.fc29",
			Architecture: "x86_64",
			SourceRpm:    "p11-kit-0.23.15-2.fc29.src.rpm",
			License:      "BSD",
			Size:         506497,
			Name:         "p11-kit-trust",
			PGP:          "RSA/SHA256, Tue Feb 19 02:39:25 2019, Key ID a20aa56b429476b4",
			Summary:      "System trust module from p11-kit",
			Vendor:       "Fedora Project",
			Version:      "0.23.15",
		},
	}
	rpmPackage3 = model.Package{
		Name:    "hardlink",
		Type:    rpmType,
		Version: "3:1.0-19.el7",
		Path:    rpmPackagesPath,
		Locations: []model.Location{
			{
				Path:      filepath.Join("var", "lib", "rpm", "Packages"),
				LayerHash: "d1fd2cca7a7751ca9786b088cf639e65088fa0bda34492bb5ba292c32195461a",
			},
		},
		Description: "Create a tree of hardlinks",
		Licenses: []string{
			"GPL+",
		},
		CPEs: []string{
			"cpe:2.3:a:redhat:hardlink:3\\:1.0-19.el7:*:*:*:*:*:*:*",
			"cpe:2.3:a:hardlink:hardlink:3\\:1.0-19.el7:*:*:*:*:*:*:*",
		},
		PURL: model.PURL("pkg:rpm/hardlink@1.0-19.el7?arch=x86_64"),
		Metadata: metadata.RPMMetadata{
			Release:      "19.el7",
			Architecture: "x86_64",
			SourceRpm:    "hardlink-1.0-19.el7.src.rpm",
			License:      "GPL+",
			Size:         16545,
			Name:         "hardlink",
			PGP:          "RSA/SHA256, Tue Apr  1 17:48:32 2014, Key ID 199e2f91fd431d51",
			Summary:      "Create a tree of hardlinks",
			Vendor:       "Red Hat, Inc.",
			Version:      "1.0",
			Epoch:        3,
		},
	}
)

func TestReadRpmContent(t *testing.T) {
	rpmPath := filepath.Join("..", "..", "..", "docs", "references", "rpm", "Packages")
	testLocation := model.Location{Path: rpmPath}
	pkgs := new([]model.Package)
	err := readRpmContent(&testLocation, pkgs)
	if err != nil {
		t.Error("Test Failed: Error occurred while reading RPM content.")
	}
}

func TestInitRpmPackage(t *testing.T) {
	var pkg1, pkg2, pkg3 model.Package

	tests := []InitRpmPackageResult{
		{&pkg1, &rpmLocation1, &rpmdb1, &rpmPackage1},
		{&pkg2, &rpmLocation2, &rpmdb2, &rpmPackage2},
		{&pkg3, &rpmLocation3, &rpmdb3, &rpmPackage3},
	}

	for _, test := range tests {
		output := initRpmPackage(test.pkg, test.location, test.rpmdb)
		outputMetadata := output.Metadata.(metadata.RPMMetadata)
		expectedMetadata := test.expected.Metadata.(metadata.RPMMetadata)

		if output.Type != test.expected.Type ||
			output.Path != test.expected.Path ||
			output.Name != test.expected.Name ||
			output.Version != test.expected.Version ||
			output.Description != test.expected.Description ||
			len(output.Licenses) != len(test.expected.Licenses) ||
			len(output.Locations) != len(test.expected.Locations) ||
			len(output.CPEs) != len(test.expected.CPEs) ||
			string(output.PURL) != string(test.expected.PURL) ||
			outputMetadata.Release != expectedMetadata.Release ||
			outputMetadata.Architecture != expectedMetadata.Architecture ||
			outputMetadata.SourceRpm != expectedMetadata.SourceRpm ||
			outputMetadata.License != expectedMetadata.License ||
			outputMetadata.Size != expectedMetadata.Size ||
			outputMetadata.Name != expectedMetadata.Name ||
			outputMetadata.PGP != expectedMetadata.PGP ||
			outputMetadata.ModularityLabel != expectedMetadata.ModularityLabel ||
			outputMetadata.Summary != expectedMetadata.Summary ||
			outputMetadata.Vendor != expectedMetadata.Vendor ||
			outputMetadata.Version != expectedMetadata.Version ||
			outputMetadata.Epoch != expectedMetadata.Epoch {
			t.Errorf("Test Failed:\n Expected output of %v \n, Received: %v \n", test.expected, output)
		}

		for i := range output.Licenses {
			if output.Licenses[i] != test.expected.Licenses[i] {
				t.Errorf("Test Failed:\n Expected output of %v \n, Received: %v \n", test.expected.Licenses[i], output.Licenses[i])
			}
		}
		for i := range output.Locations {
			if output.Locations[i] != test.expected.Locations[i] {
				t.Errorf("Test Failed:\n Expected output of %v \n, Received: %v \n", test.expected.Locations[i], output.Locations[i])
			}
		}
		for i := range output.CPEs {
			if output.CPEs[i] != test.expected.CPEs[i] {
				t.Errorf("Test Failed:\n Expected output of %v \n, Received: %v \n", test.expected.CPEs[i], output.CPEs[i])
			}
		}
	}
}

func TestParseRpmPackageURL(t *testing.T) {
	pkg1 := model.Package{
		Name:    rpmPackage1.Name,
		Version: "2.08-14.el8",
	}
	pkg2 := model.Package{
		Name:    rpmPackage2.Name,
		Version: "0.23.15-2.fc29",
	}
	pkg3 := model.Package{
		Name:    rpmPackage3.Name,
		Version: "1.0-19.el7",
	}

	tests := []RpmPurlResult{
		{&pkg1, rpmdb1.Arch, model.PURL("pkg:rpm/lzo@2.08-14.el8?arch=x86_64")},
		{&pkg2, rpmdb2.Arch, model.PURL("pkg:rpm/p11-kit-trust@0.23.15-2.fc29?arch=x86_64")},
		{&pkg3, rpmdb3.Arch, model.PURL("pkg:rpm/hardlink@1.0-19.el7?arch=x86_64")},
	}

	for _, test := range tests {
		parseRpmPackageURL(test.pkg, test.arch)
		if test.pkg.PURL != test.expected {
			t.Errorf("Test Failed: Expected an output of %v, received: %v", test.expected, test.pkg.PURL)
		}
	}
}

func TestInitFinalRpmMetadata(t *testing.T) {
	var pkg1, pkg2, pkg3 model.Package

	tests := []RpmMetadataResult{
		{&pkg1, &rpmdb1, metadata.RPMMetadata{
			Release:      "14.el8",
			Architecture: "x86_64",
			SourceRpm:    "lzo-2.08-14.el8.src.rpm",
			License:      "GPLv2+",
			Size:         198757,
			Name:         "lzo",
			PGP:          "RSA/SHA256, Tue Jul  2 00:01:31 2019, Key ID 05b555b38483c65d",
			Summary:      "Data compression library with very fast (de)compression",
			Vendor:       "CentOS",
			Version:      "2.08",
		}},
		{&pkg2, &rpmdb2, metadata.RPMMetadata{
			Release:      "2.fc29",
			Architecture: "x86_64",
			SourceRpm:    "p11-kit-0.23.15-2.fc29.src.rpm",
			License:      "BSD",
			Size:         506497,
			Name:         "p11-kit-trust",
			PGP:          "RSA/SHA256, Tue Feb 19 02:39:25 2019, Key ID a20aa56b429476b4",
			Summary:      "System trust module from p11-kit",
			Vendor:       "Fedora Project",
			Version:      "0.23.15",
		}},
		{&pkg3, &rpmdb3, metadata.RPMMetadata{
			Release:      "19.el7",
			Architecture: "x86_64",
			SourceRpm:    "hardlink-1.0-19.el7.src.rpm",
			License:      "GPL+",
			Size:         16545,
			Name:         "hardlink",
			PGP:          "RSA/SHA256, Tue Apr  1 17:48:32 2014, Key ID 199e2f91fd431d51",
			Summary:      "Create a tree of hardlinks",
			Vendor:       "Red Hat, Inc.",
			Version:      "1.0",
			Epoch:        3,
		}},
	}

	for _, test := range tests {
		initFinalRpmMetadata(test.pkg, test.rpmPkg)
		outputMetadata := test.pkg.Metadata.(metadata.RPMMetadata)
		expectedMetadata := test.expected
		if outputMetadata.Release != expectedMetadata.Release ||
			outputMetadata.Architecture != expectedMetadata.Architecture ||
			outputMetadata.SourceRpm != expectedMetadata.SourceRpm ||
			outputMetadata.License != expectedMetadata.License ||
			outputMetadata.Size != expectedMetadata.Size ||
			outputMetadata.Name != expectedMetadata.Name ||
			outputMetadata.PGP != expectedMetadata.PGP ||
			outputMetadata.ModularityLabel != expectedMetadata.ModularityLabel ||
			outputMetadata.Summary != expectedMetadata.Summary ||
			outputMetadata.Vendor != expectedMetadata.Vendor ||
			outputMetadata.Version != expectedMetadata.Version ||
			outputMetadata.Epoch != expectedMetadata.Epoch {
			t.Errorf("Test Failed: Expected output of %v, received: %v", test.expected, test.pkg.Metadata)
		}
	}
}

func TestFormatLicenses(t *testing.T) {
	var pkg1, pkg2, pkg3, pkg4 model.Package

	tests := []FormatLicensesResult{
		{&pkg1, "license01 and license02", []string{"license01", "license02"}},
		{&pkg2, "license01 or license02", []string{"license01", "license02"}},
		{&pkg3, " ", []string{}},
		{&pkg4, "", []string{}},
	}

	for _, test := range tests {
		formatLicenses(test.pkg, test.licenses)
		if len(test.expected) == 0 && len(test.pkg.Licenses) != 0 {
			t.Errorf("Test Failed: Slice length must be equal with the expected result. Expected: %v, Received: %v", len(test.expected), len(test.pkg.Licenses))
		}
		if len(test.pkg.Licenses) != len(test.expected) {
			t.Errorf("Test Failed: Slice length must be equal with the expected result. Expected: %v, Received: %v", len(test.expected), len(test.pkg.Licenses))
		}
		for i := range test.pkg.Licenses {
			if test.pkg.Licenses[i] != test.expected[i] {
				t.Errorf("Test Failed:\n Expected output of %v \n, Received: %v \n", test.expected[i], test.pkg.Licenses[i])
			}
		}
	}
}

func TestFormatVendor(t *testing.T) {
	tests := []FormatVendorResult{
		{"CentOS", "centos"},
		{"Red Hat, Inc.", "redhat"},
		{"fedoraproject", "fedoraproject"},
		{"test", "test"},
		{"   testWithSpace   ", "testwithspace"},
		{"", ""},
		{"   ", ""},
	}

	for _, test := range tests {
		if output := formatVendor(test.input); output != test.expected {
			t.Errorf("Test Failed: Input %v must have output of %v, received: %v", test.input, test.expected, output)
		}
	}
}

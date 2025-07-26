package build

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/layer5io/meshery-adapter-library/adapter"

	"github.com/layer5io/meshkit/utils"
	"github.com/layer5io/meshkit/utils/manifests"
	smp "github.com/layer5io/service-mesh-performance/spec"
)

var DefaultGenerationMethod string
var DefaultGenerationURL string
var LatestVersion string
var MeshModelPath string
var AllVersions []string

const Component = "Istio"

var MeshModelConfig = adapter.MeshModelConfig{ //Move to build/config.go
	Category: "Cloud Native Network",
	Metadata: map[string]interface{}{},
}

// NewConfig creates the configuration for creating components
func NewConfig(version string) manifests.Config {
	return manifests.Config{
		Name:        smp.ServiceMesh_Type_name[int32(smp.ServiceMesh_ISTIO)],
		Type:        Component,
		MeshVersion: version,
		CrdFilter: manifests.NewCueCrdFilter(manifests.ExtractorPaths{
			NamePath:    "spec.names.kind",
			IdPath:      "spec.names.kind",
			VersionPath: "spec.versions[0].name",
			GroupPath:   "spec.group",
			SpecPath:    "spec.versions[0].schema.openAPIV3Schema"}, false),
		ExtractCrds: func(manifest string) []string {
			crds := strings.Split(manifest, "---")
			return crds
		},
	}
}
func init() {
	//Initialize Metadata including logo svgs
	f, _ := os.Open("./build/meshmodel_metadata.json")
	defer func() {
		if err := f.Close(); err != nil {
			fmt.Printf("Error closing file: %s\n", err)
		}
	}()
	byt, _ := io.ReadAll(f)

	_ = json.Unmarshal(byt, &MeshModelConfig.Metadata)
	wd, _ := os.Getwd()
	MeshModelPath = filepath.Join(wd, "templates", "meshmodel", "components")
	AllVersions, _ = utils.GetLatestReleaseTagsSorted("istio", "istio")
	if len(AllVersions) == 0 {
		return
	}
	LatestVersion = AllVersions[len(AllVersions)-1]
	DefaultGenerationMethod = adapter.Manifests
	DefaultGenerationURL = "https://raw.githubusercontent.com/istio/istio/" + LatestVersion + "/manifests/charts/base/files/crd-all.gen.yaml"
}

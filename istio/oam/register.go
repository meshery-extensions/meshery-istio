package oam

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"sync"

	"github.com/layer5io/meshery-adapter-library/adapter"
	"github.com/layer5io/meshkit/models/meshmodel/core/types"
)

var (
	basePath, _         = os.Getwd()
	MeshmodelComponents = filepath.Join(basePath, "templates", "meshmodel", "components")
)

// AvailableVersions denote the component versions available statically
var AvailableVersions = map[string]bool{}

type meshmodelDefinitionPathSet struct {
	meshmodelDefinitionPath string
}

func RegisterMeshModelComponents(uuid, runtime, host, port string) error {
	meshmodelRDP := []adapter.MeshModelRegistrantDefinitionPath{}
	pathSets, err := loadMeshmodelComponents(MeshmodelComponents)
	if err != nil {
		return err
	}
	portint, _ := strconv.Atoi(port)
	for _, pathSet := range pathSets {
		meshmodelRDP = append(meshmodelRDP, adapter.MeshModelRegistrantDefinitionPath{
			EntityDefintionPath: pathSet.meshmodelDefinitionPath,
			Host:                host,
			Port:                portint,
			Type:                types.ComponentDefinition,
		})
	}

	return adapter.
		NewMeshModelRegistrant(meshmodelRDP, fmt.Sprintf("%s/api/meshmodel/components/register", runtime)).
		Register(uuid)
}

var versionLock sync.Mutex

func loadMeshmodelComponents(basepath string) ([]meshmodelDefinitionPathSet, error) {
	res := []meshmodelDefinitionPathSet{}
	if err := filepath.Walk(basepath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		res = append(res, meshmodelDefinitionPathSet{
			meshmodelDefinitionPath: path,
		})
		versionLock.Lock()
		AvailableVersions[filepath.Base(filepath.Dir(path))] = true // Getting available versions already existing on file system
		versionLock.Unlock()
		return nil
	}); err != nil {
		return nil, err
	}

	return res, nil
}

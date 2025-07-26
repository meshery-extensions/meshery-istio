package istio

import (
	"context"
	"fmt"
	"net/url"
	"strings"
	"sync"

	"github.com/layer5io/meshery-adapter-library/adapter"
	"github.com/layer5io/meshery-adapter-library/status"
	"github.com/layer5io/meshkit/utils"
	mesherykube "github.com/layer5io/meshkit/utils/kubernetes"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

// installAddon installs/uninstalls an addon in the given namespace
//
// the template defines the manifest's link/location which needs to be used to
// install the addon
func (istio *Istio) installAddon(namespace string, del bool, service string, patches []string, templates []adapter.Template, kubeconfigs []string) (string, error) {
	st := status.Installing

	if del {
		st = status.Removing
	}

	istio.Log.Debug(fmt.Sprintf("Overidden namespace: %s", namespace))
	namespace = "istio-system"
	var wg sync.WaitGroup
	var errMx sync.Mutex
	var errs []error
	for _, k8sconfig := range kubeconfigs {
		wg.Add(1)
		go func(k8sconfig string) {
			defer wg.Done()
			mclient, err := mesherykube.New([]byte(k8sconfig))
			if err != nil {
				errMx.Lock()
				errs = append(errs, err)
				errMx.Unlock()
				return
			}
			for _, template := range templates {
				err := istio.applyManifestOnSingleCluster([]byte(template.String()), del, namespace, mclient)
				// Specifically choosing to ignore kiali dashboard's error.
				// Referring to: https://github.com/kiali/kiali/issues/3112
				if err != nil && !strings.Contains(err.Error(), "no matches for kind \"MonitoringDashboard\" in version \"monitoring.kiali.io/v1alpha1\"") {
					if !strings.Contains(err.Error(), "clusterIP") {
						errMx.Lock()
						errs = append(errs, err)
						errMx.Unlock()
					}
				}
			}

			for _, patch := range patches {
				if patch == "" {
					continue //avoid throwing error when a given patch key didn't exist for a specific addon type in operations
				}
				if !del {
					_, err := url.ParseRequestURI(patch)
					if err != nil {
						errMx.Lock()
						errs = append(errs, err)
						errMx.Unlock()
						return
					}

					content, err := utils.ReadFileSource(patch)
					if err != nil {
						errMx.Lock()
						errs = append(errs, err)
						errMx.Unlock()
						return
					}

					_, err = mclient.KubeClient.CoreV1().Services(namespace).Patch(context.TODO(), service, types.MergePatchType, []byte(content), metav1.PatchOptions{})
					if err != nil {
						errMx.Lock()
						errs = append(errs, err)
						errMx.Unlock()
						return
					}
				}
			}
		}(k8sconfig)
	}
	wg.Wait()
	if len(errs) == 0 {
		return status.Installed, nil
	}
	return st, ErrAddonFromTemplate(mergeErrors(errs))
}

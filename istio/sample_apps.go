package istio

import (
	"context"
	"sync"

	"github.com/layer5io/meshery-adapter-library/adapter"
	"github.com/layer5io/meshery-adapter-library/status"
	"github.com/layer5io/meshkit/utils"
	mesherykube "github.com/layer5io/meshkit/utils/kubernetes"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
)

func (istio *Istio) installSampleApp(namespace string, del bool, templates []adapter.Template, kubeconfigs []string) (string, error) {
	st := status.Installing

	if del {
		st = status.Removing
	}

	for _, template := range templates {
		err := istio.applyManifest([]byte(template.String()), del, namespace, kubeconfigs)
		if err != nil {
			return st, ErrSampleApp(err)
		}
	}

	return status.Installed, nil
}

func (istio *Istio) patchWithEnvoyFilter(namespace string, del bool, app string, templates []adapter.Template, patchObject string, kubeconfigs []string) (string, error) {
	st := status.Deploying

	if del {
		st = status.Removing
	}

	jsonContents, err := utils.ReadFileSource(patchObject)
	if err != nil {
		return st, ErrEnvoyFilter(err)
	}
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
			_, err = mclient.KubeClient.AppsV1().Deployments(namespace).Patch(context.TODO(), app, types.MergePatchType, []byte(jsonContents), metav1.PatchOptions{})
			if err != nil {
				errMx.Lock()
				errs = append(errs, err)
				errMx.Unlock()
				return
			}

			for _, template := range templates {
				contents, err := utils.ReadFileSource(string(template))
				if err != nil {
					errMx.Lock()
					errs = append(errs, err)
					errMx.Unlock()
					continue
				}

				err = istio.applyManifestOnSingleCluster([]byte(contents), del, namespace, mclient)
				if err != nil {
					errMx.Lock()
					errs = append(errs, err)
					errMx.Unlock()
					return
				}
			}
		}(k8sconfig)
	}
	wg.Wait()
	if len(errs) == 0 {
		return status.Deployed, nil
	}
	return st, ErrEnvoyFilter(mergeErrors(errs))
}
func (istio *Istio) applyPolicy(namespace string, del bool, templates []adapter.Template, kubeconfigs []string) (string, error) {
	st := status.Deploying

	if del {
		st = status.Removing
	}

	for _, template := range templates {
		contents, err := utils.ReadFileSource(string(template))
		if err != nil {
			return st, ErrApplyPolicy(err)
		}

		err = istio.applyManifest([]byte(contents), del, namespace, kubeconfigs)
		if err != nil {
			return st, ErrApplyPolicy(err)
		}
	}
	return status.Deployed, nil
}

// LoadToMesh is used to mark deployment for automatic sidecar injection (or not)
func (istio *Istio) LoadToMesh(namespace string, service string, remove bool, kubeconfigs []string) error {
	var wg sync.WaitGroup
	var errMx sync.Mutex
	var errs []error
	for _, k8sconfig := range kubeconfigs {
		wg.Add(1)
		go func(k8sconfig string) {
			defer wg.Done()
			kclient, err := mesherykube.New([]byte(k8sconfig))
			if err != nil {
				errMx.Lock()
				errs = append(errs, err)
				errMx.Unlock()
				return
			}

			deploy, err := kclient.KubeClient.AppsV1().Deployments(namespace).Get(context.TODO(), service, metav1.GetOptions{})
			if err != nil {
				errMx.Lock()
				errs = append(errs, err)
				errMx.Unlock()
				return
			}

			if deploy.ObjectMeta.Labels == nil {
				deploy.ObjectMeta.Labels = map[string]string{}
			}
			deploy.ObjectMeta.Labels["istio-injection"] = "enabled"

			if remove {
				delete(deploy.ObjectMeta.Labels, "istio-injection")
			}

			_, err = kclient.KubeClient.AppsV1().Deployments(namespace).Update(context.TODO(), deploy, metav1.UpdateOptions{})
			if err != nil {
				errMx.Lock()
				errs = append(errs, err)
				errMx.Unlock()
				return
			}
		}(k8sconfig)
	}
	wg.Wait()
	if len(errs) == 0 {
		return nil
	}
	return mergeErrors(errs)
}

// LoadNamespaceToMesh is used to mark namespaces for automatic sidecar injection (or not)
func (istio *Istio) LoadNamespaceToMesh(namespace string, remove bool, kubeconfigs []string) error {
	var wg sync.WaitGroup
	var errMx sync.Mutex
	var errs []error
	for _, k8sconfig := range kubeconfigs {
		wg.Add(1)
		go func(k8sconfig string) {
			defer wg.Done()
			kclient, err := mesherykube.New([]byte(k8sconfig))
			if err != nil {
				errMx.Lock()
				errs = append(errs, err)
				errMx.Unlock()
				// return ErrLoadNamespace(err, namespace)
				return
			}

			ns, err := kclient.KubeClient.CoreV1().Namespaces().Get(context.TODO(), namespace, metav1.GetOptions{})
			if err != nil {
				errMx.Lock()
				errs = append(errs, err)
				errMx.Unlock()
				// return ErrLoadNamespace(err, namespace)
				return
			}
			if ns.ObjectMeta.Labels == nil {
				ns.ObjectMeta.Labels = map[string]string{}
			}
			ns.ObjectMeta.Labels["istio-injection"] = "enabled"

			if remove {
				delete(ns.ObjectMeta.Labels, "istio-injection")
			}

			_, err = kclient.KubeClient.CoreV1().Namespaces().Update(context.TODO(), ns, metav1.UpdateOptions{})
			if err != nil {
				errMx.Lock()
				errs = append(errs, err)
				errMx.Unlock()
				// return ErrLoadNamespace(err, namespace)
				return
			}
		}(k8sconfig)
	}
	wg.Wait()
	if len(errs) == 0 {
		return nil
	}
	return ErrLoadNamespace(mergeErrors(errs), namespace)
}

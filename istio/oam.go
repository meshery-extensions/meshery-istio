package istio

import (
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/layer5io/meshery-adapter-library/common"
	"github.com/layer5io/meshery-adapter-library/meshes"
	"github.com/layer5io/meshery-istio/internal/config"
	"github.com/layer5io/meshkit/errors"
	"github.com/layer5io/meshkit/models/oam/core/v1alpha1"
	"gopkg.in/yaml.v2"
)

// CompHandler is the type for functions which can handle OAM components
type CompHandler func(*Istio, v1alpha1.Component, bool, []string) (string, error)

// HandleComponents handles the processing of OAM components
func (istio *Istio) HandleComponents(comps []v1alpha1.Component, isDel bool, kubeconfigs []string) (string, error) {
	var errs []error
	var msgs []string

	compFuncMap := map[string]CompHandler{
		"IstioMesh":            handleComponentIstioMesh,
		"GrafanaIstioAddon":    handleComponentIstioAddon,
		"PrometheusIstioAddon": handleComponentIstioAddon,
		"ZipkinIstioAddon":     handleComponentIstioAddon,
		"JaegerIstioAddon":     handleComponentIstioAddon,
	}
	stat1 := "deploying"
	stat2 := "deployed"
	if isDel {
		stat1 = "removing"
		stat2 = "removed"
	}
	for _, comp := range comps {
		ee := &meshes.EventsResponse{
			OperationId:   uuid.New().String(),
			Component:     config.ServerConfig["type"],
			ComponentName: config.ServerConfig["name"],
		}
		fnc, ok := compFuncMap[comp.Spec.Type]
		if !ok {
			msg, err := handleIstioCoreComponent(istio, comp, isDel, "", "", kubeconfigs)
			if err != nil {
				ee.Summary = fmt.Sprintf("Error while %s %s", stat1, comp.Spec.Type)
				ee.Details = err.Error()
				ee.ErrorCode = errors.GetCode(err)
				ee.ProbableCause = errors.GetCause(err)
				ee.SuggestedRemediation = errors.GetRemedy(err)
				istio.StreamErr(ee, err)
				errs = append(errs, err)
				continue
			}
			ee.Summary = fmt.Sprintf("%s %s successfully", comp.Spec.Type, stat2)
			ee.Details = fmt.Sprintf("The %s is now %s.", comp.Spec.Type, stat2)
			istio.StreamInfo(ee)
			msgs = append(msgs, msg)
			continue
		}

		msg, err := fnc(istio, comp, isDel, kubeconfigs)
		if err != nil {
			ee.Summary = fmt.Sprintf("Error while %s %s", stat1, comp.Spec.Type)
			ee.Details = err.Error()
			ee.ErrorCode = errors.GetCode(err)
			ee.ProbableCause = errors.GetCause(err)
			ee.SuggestedRemediation = errors.GetRemedy(err)
			istio.StreamErr(ee, err)
			errs = append(errs, err)
			continue
		}
		ee.Summary = fmt.Sprintf("%s %s %s successfully", comp.Name, comp.Spec.Type, stat2)
		ee.Details = fmt.Sprintf("The %s %s is now %s.", comp.Name, comp.Spec.Type, stat2)
		istio.StreamInfo(ee)
		msgs = append(msgs, msg)
	}

	if err := mergeErrors(errs); err != nil {
		return mergeMsgs(msgs), err
	}

	return mergeMsgs(msgs), nil
}

// HandleApplicationConfiguration handles the processing of OAM application configuration
func (istio *Istio) HandleApplicationConfiguration(config v1alpha1.Configuration, isDel bool, kubeconfigs []string) (string, error) {
	var errs []error
	var msgs []string
	for _, comp := range config.Spec.Components {
		for _, trait := range comp.Traits {
			if trait.Name == "mTLS" {
				namespaces := castSliceInterfaceToSliceString(trait.Properties["namespaces"].([]interface{}))
				policy := trait.Properties["policy"].(string)

				if err := handleMTLS(istio, namespaces, policy, isDel, kubeconfigs); err != nil {
					errs = append(errs, err)
				}
			}

			if trait.Name == "automaticSidecarInjection" {
				namespaces := castSliceInterfaceToSliceString(trait.Properties["namespaces"].([]interface{}))
				if err := handleNamespaceLabel(istio, namespaces, isDel, kubeconfigs); err != nil {
					errs = append(errs, err)
				}
			}

			msgs = append(msgs, fmt.Sprintf("applied trait \"%s\" on service \"%s\"", trait.Name, comp.ComponentName))
		}
	}

	if err := mergeErrors(errs); err != nil {
		return mergeMsgs(msgs), err
	}

	return mergeMsgs(msgs), nil
}

func handleMTLS(istio *Istio, namespaces []string, policy string, isDel bool, kubeconfigs []string) error {
	var errs []error
	for _, ns := range namespaces {
		policyName := fmt.Sprintf("%s-mtls-policy-operation", policy)

		if _, err := istio.applyPolicy(ns, isDel, config.GetOperations(common.Operations, "master")[policyName].Templates, kubeconfigs); err != nil {
			errs = append(errs, err)
		}
	}

	return mergeErrors(errs)
}

func handleNamespaceLabel(istio *Istio, namespaces []string, isDel bool, kubeconfigs []string) error {
	var errs []error
	for _, ns := range namespaces {
		if err := istio.LoadNamespaceToMesh(ns, isDel, kubeconfigs); err != nil {
			errs = append(errs, err)
		}
	}

	return mergeErrors(errs)
}

func handleComponentIstioMesh(istio *Istio, comp v1alpha1.Component, isDel bool, kubeconfigs []string) (string, error) {
	// Get the istio version from the settings
	// we are sure that the version of istio would be present
	// because the configuration is already validated against the schema
	version := comp.Spec.Version
	if version == "" {
		return "", fmt.Errorf("pass valid version inside service for Istio installation")
	}
	//TODO: When no version is passed in service, use the latest istio version
	profile := comp.Spec.Settings["profile"].(string)
	return istio.installIstio(isDel, false, version, comp.Namespace, profile, kubeconfigs)
}

func handleIstioCoreComponent(
	istio *Istio,
	comp v1alpha1.Component,
	isDel bool,
	apiVersion,
	kind string,
	kubeconfigs []string) (string, error) {
	if apiVersion == "" {
		apiVersion = v1alpha1.GetAPIVersionFromComponent(comp)
		if apiVersion == "" {
			return "", ErrIstioCoreComponentFail(fmt.Errorf("failed to get API Version for: %s", comp.Name))
		}
	}

	if kind == "" {
		kind = v1alpha1.GetKindFromComponent(comp)
		if kind == "" {
			return "", ErrIstioCoreComponentFail(fmt.Errorf("failed to get kind for: %s", comp.Name))
		}
	}
	component := map[string]interface{}{
		"apiVersion": apiVersion,
		"kind":       kind,
		"metadata": map[string]interface{}{
			"name":        comp.Name,
			"annotations": comp.Annotations,
			"labels":      comp.Labels,
		},
		"spec": comp.Spec.Settings,
	}

	// Convert to yaml
	yamlByt, err := yaml.Marshal(component)
	if err != nil {
		err = ErrParseIstioCoreComponent(err)
		istio.Log.Error(err)
		return "", err
	}
	msg := fmt.Sprintf("created %s \"%s\" in namespace \"%s\"", kind, comp.Name, comp.Namespace)
	if isDel {
		msg = fmt.Sprintf("deleted %s config \"%s\" in namespace \"%s\"", kind, comp.Name, comp.Namespace)
	}

	return msg, istio.applyManifest(yamlByt, isDel, comp.Namespace, kubeconfigs)
}

func handleComponentIstioAddon(istio *Istio, comp v1alpha1.Component, isDel bool, kubeconfigs []string) (string, error) {
	var addonName string

	switch comp.Spec.Type {
	case "GrafanaIstioAddon":
		addonName = config.GrafanaAddon
	case "PrometheusIstioAddon":
		addonName = config.PrometheusAddon
	case "ZipkinIstioAddon":
		addonName = config.ZipkinAddon
	case "JaegerIstioAddon":
		addonName = config.JaegerAddon
	default:
		return "", nil
	}
	version := comp.Spec.Version
	// Get the service
	svc := config.GetOperations(common.Operations, version)[addonName].AdditionalProperties[common.ServiceName]

	// Get the patches
	patches := make([]string, 0)
	patches = append(patches, config.GetOperations(common.Operations, version)[addonName].AdditionalProperties[config.ServicePatchFile])
	patches = append(patches, config.GetOperations(common.Operations, version)[addonName].AdditionalProperties[config.CPPatchFile])
	patches = append(patches, config.GetOperations(common.Operations, version)[addonName].AdditionalProperties[config.ControlPatchFile])

	// Get the templates
	templates := config.GetOperations(common.Operations, version)[addonName].Templates

	_, err := istio.installAddon(comp.Namespace, isDel, svc, patches, templates, kubeconfigs)

	msg := fmt.Sprintf("created service of type \"%s\"", comp.Spec.Type)
	if isDel {
		msg = fmt.Sprintf("deleted service of type \"%s\"", comp.Spec.Type)
	}

	return msg, err
}

func castSliceInterfaceToSliceString(in []interface{}) []string {
	var out []string

	for _, v := range in {
		cast, ok := v.(string)
		if ok {
			out = append(out, cast)
		}
	}

	return out
}

func mergeErrors(errs []error) error {
	if len(errs) == 0 {
		return nil
	}

	var errMsgs []string

	for _, err := range errs {
		errMsgs = append(errMsgs, err.Error())
	}

	return fmt.Errorf("%s", strings.Join(errMsgs, "\n"))
}

func mergeMsgs(strs []string) string {
	return strings.Join(strs, "\n")
}

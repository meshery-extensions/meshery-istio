package istio

import (
	"context"
	"fmt"
	stderrors "errors"

	"github.com/layer5io/meshery-adapter-library/adapter"
	"github.com/layer5io/meshery-adapter-library/common"
	"github.com/layer5io/meshery-adapter-library/meshes"
	"github.com/layer5io/meshery-adapter-library/status"
	internalconfig "github.com/layer5io/meshery-istio/internal/config"
	"github.com/layer5io/meshery-istio/istio/oam"
	meshkitCfg "github.com/layer5io/meshkit/config"
	"github.com/layer5io/meshkit/errors"
	"github.com/layer5io/meshkit/logger"
	"github.com/layer5io/meshkit/models"
	"github.com/layer5io/meshkit/models/oam/core/v1alpha1"
	"github.com/layer5io/meshkit/utils"
	"github.com/layer5io/meshkit/utils/events"
	"gopkg.in/yaml.v2"
)

// Istio represents the istio adapter and embeds adapter.Adapter
type Istio struct {
	adapter.Adapter // Type Embedded
}

// New initializes istio handler.
func New(c meshkitCfg.Handler, l logger.Handler, kc meshkitCfg.Handler, ev *events.EventStreamer) adapter.Handler {
	return &Istio{
		Adapter: adapter.Adapter{
			Config:            c,
			Log:               l,
			KubeconfigHandler: kc,
			EventStreamer:     ev,
		},
	}
}

// ApplyOperation applies the operation on istio
func (istio *Istio) ApplyOperation(ctx context.Context, opReq adapter.OperationRequest) error {
	err := istio.CreateKubeconfigs(opReq.K8sConfigs)
	if err != nil {
		return err
	}
	kubeConfigs := opReq.K8sConfigs
	operations := make(adapter.Operations)
	requestedVersion := adapter.Version(opReq.Version)
	err = istio.Config.GetObject(adapter.OperationsKey, &operations)
	if err != nil {
		return err
	}

	e := &meshes.EventsResponse{
		OperationId:   opReq.OperationID,
		Summary:       status.Deploying,
		Details:       "Operation is not supported",
		Component:     internalconfig.ServerConfig["type"],
		ComponentName: internalconfig.ServerConfig["name"],
	}
	switch opReq.OperationName {
	case internalconfig.IstioOperation:
		go func(hh *Istio, ee *meshes.EventsResponse) {
			var err error
			var stat, version string
			if len(operations[opReq.OperationName].Versions) == 0 {
				err = ErrFetchIstioVersions
			} else {
				version = string(operations[opReq.OperationName].Versions[len(operations[opReq.OperationName].Versions)-1])
				if utils.Contains[[]adapter.Version, adapter.Version](operations[opReq.OperationName].Versions, requestedVersion) {
					version = requestedVersion.String()
				}
				stat, err = hh.installIstio(opReq.IsDeleteOperation, false, version, opReq.Namespace, "default", kubeConfigs)
			}
			if err != nil { //Make sure that this is a meshkit error
				ee.Summary = fmt.Sprintf("Error while %s Istio service mesh %s", stat, version)
				ee.Details = err.Error()
				ee.ErrorCode = errors.GetCode(err)
				ee.ProbableCause = errors.GetCause(err)
				ee.SuggestedRemediation = errors.GetRemedy(err)
				hh.StreamErr(ee, err)
				return
			}
			ee.Summary = fmt.Sprintf("Istio service mesh %s %s successfully", version, stat)
			ee.Details = fmt.Sprintf("The Istio service mesh %s is now %s.", version, stat)
			hh.StreamInfo(ee)
		}(istio, e)
	case common.BookInfoOperation, common.HTTPBinOperation, common.ImageHubOperation, common.EmojiVotoOperation:
		go func(hh *Istio, ee *meshes.EventsResponse) {
			appName := operations[opReq.OperationName].AdditionalProperties[common.ServiceName]
			stat, err := hh.installSampleApp(opReq.Namespace, opReq.IsDeleteOperation, operations[opReq.OperationName].Templates, kubeConfigs)
			if err != nil {
				ee.Summary = fmt.Sprintf("Error while %s Istio service mesh", stat)
				ee.Details = err.Error()
				ee.ErrorCode = errors.GetCode(err)
				ee.ProbableCause = errors.GetCause(err)
				ee.SuggestedRemediation = errors.GetRemedy(err)
				hh.StreamErr(ee, err)
				return
			}
			ee.Summary = fmt.Sprintf("%s application %s successfully", appName, stat)
			ee.Details = fmt.Sprintf("The %s application is now %s.", appName, stat)
			hh.StreamInfo(ee)
		}(istio, e)
	case common.SmiConformanceOperation:
		go func(hh *Istio, ee *meshes.EventsResponse) {
			name := operations[opReq.OperationName].Description
			_, err := hh.RunSMITest(adapter.SMITestOptions{
				Ctx:         context.TODO(),
				OperationID: ee.OperationId,
				Labels: map[string]string{
					"istio-injection": "enabled",
				},
				Namespace:   "meshery",
				Manifest:    string(operations[opReq.OperationName].Templates[0]),
				Annotations: make(map[string]string),
			})
			if err != nil {
				ee.Summary = fmt.Sprintf("Error while %s %s test", status.Running, name)
				ee.Details = err.Error()
				ee.ErrorCode = errors.GetCode(err)
				ee.ProbableCause = errors.GetCause(err)
				ee.SuggestedRemediation = errors.GetRemedy(err)
				hh.StreamErr(ee, err)
				return
			}
			ee.Summary = fmt.Sprintf("%s test %s successfully", name, status.Completed)
			ee.Details = ""
			hh.StreamInfo(ee)
		}(istio, e)
	case internalconfig.DenyAllPolicyOperation, internalconfig.StrictMTLSPolicyOperation, internalconfig.MutualMTLSPolicyOperation, internalconfig.DisableMTLSPolicyOperation:
		go func(hh *Istio, ee *meshes.EventsResponse) {
			stat, err := hh.applyPolicy(opReq.Namespace, opReq.IsDeleteOperation, operations[opReq.OperationName].Templates, kubeConfigs)
			if err != nil {
				ee.Summary = fmt.Sprintf("Error while %s policy", stat)
				ee.Details = err.Error()
				ee.ErrorCode = errors.GetCode(err)
				ee.ProbableCause = errors.GetCause(err)
				ee.SuggestedRemediation = errors.GetRemedy(err)
				hh.StreamErr(ee, err)
				return
			}
			ee.Summary = fmt.Sprintf("Policy %s successfully", status.Deployed)
			ee.Details = ""
			hh.StreamInfo(ee)
		}(istio, e)
	case common.CustomOperation:
		go func(hh *Istio, ee *meshes.EventsResponse) {
			stat, err := hh.applyCustomOperation(opReq.Namespace, opReq.CustomBody, opReq.IsDeleteOperation, kubeConfigs)
			if err != nil {
				ee.Summary = fmt.Sprintf("Error while %s custom operation", stat)
				ee.Details = err.Error()
				ee.ErrorCode = errors.GetCode(err)
				ee.ProbableCause = errors.GetCause(err)
				ee.SuggestedRemediation = errors.GetRemedy(err)
				hh.StreamErr(ee, err)
				return
			}
			ee.Summary = fmt.Sprintf("Manifest %s successfully", status.Deployed)
			ee.Details = ""
			hh.StreamInfo(ee)
		}(istio, e)
	case internalconfig.LabelNamespace:
		go func(hh *Istio, ee *meshes.EventsResponse) {
			err := hh.LoadNamespaceToMesh(opReq.Namespace, opReq.IsDeleteOperation, kubeConfigs)
			operation := "enabled"
			if opReq.IsDeleteOperation {
				operation = "removed"
			}
			if err != nil {
				ee.Summary = fmt.Sprintf("Error while labeling %s", opReq.Namespace)
				ee.Details = err.Error()
				ee.ErrorCode = errors.GetCode(err)
				ee.ProbableCause = errors.GetCause(err)
				ee.SuggestedRemediation = errors.GetRemedy(err)
				hh.StreamErr(ee, err)
				return
			}
			ee.Summary = fmt.Sprintf("Label updated on %s namespace", opReq.Namespace)
			ee.Details = fmt.Sprintf("ISTIO-INJECTION label %s on %s namespace", operation, opReq.Namespace)
			hh.StreamInfo(ee)
		}(istio, e)
	case internalconfig.PrometheusAddon, internalconfig.GrafanaAddon, internalconfig.KialiAddon, internalconfig.JaegerAddon, internalconfig.ZipkinAddon:
		go func(hh *Istio, ee *meshes.EventsResponse) {
			svcname := operations[opReq.OperationName].AdditionalProperties[common.ServiceName]
			patches := make([]string, 0)
			patches = append(patches, operations[opReq.OperationName].AdditionalProperties[internalconfig.ServicePatchFile])

			_, err := hh.installAddon(opReq.Namespace, opReq.IsDeleteOperation, svcname, patches, operations[opReq.OperationName].Templates, kubeConfigs)
			operation := "install"
			if opReq.IsDeleteOperation {
				operation = "uninstall"
			}

			if err != nil {
				ee.Summary = fmt.Sprintf("Error while %sing %s", operation, opReq.OperationName)
				ee.Details = err.Error()
				ee.ErrorCode = errors.GetCode(err)
				ee.ProbableCause = errors.GetCause(err)
				ee.SuggestedRemediation = errors.GetRemedy(err)
				hh.StreamErr(ee, err)
				return
			}
			ee.Summary = fmt.Sprintf("Successfully %sed %s", operation, opReq.OperationName)
			ee.Details = fmt.Sprintf("Successfully %sed %s from the %s namespace", operation, opReq.OperationName, opReq.Namespace)
			hh.StreamInfo(ee)
		}(istio, e)
	case internalconfig.IstioVetOperation:
		go func(hh *Istio, ee *meshes.EventsResponse) {
			responseChan := make(chan *meshes.EventsResponse, 1)

			go hh.RunVet(responseChan, kubeConfigs)

			for msg := range responseChan {
				switch msg.EventType {
				case meshes.EventType_ERROR:
					istio.StreamErr(msg, ErrIstioVet(stderrors.New(msg.Details)))
				case meshes.EventType_WARN:
					istio.StreamWarn(msg, ErrIstioVet(stderrors.New(msg.Details)))
				default:
					istio.StreamInfo(msg)
				}
			}

			istio.Log.Info("Done")
		}(istio, e)
	case internalconfig.EnvoyFilterOperation:
		go func(hh *Istio, ee *meshes.EventsResponse) {
			appName := operations[opReq.OperationName].AdditionalProperties[common.ServiceName]
			patchFile := operations[opReq.OperationName].AdditionalProperties[internalconfig.FilterPatchFile]
			stat, err := hh.patchWithEnvoyFilter(opReq.Namespace, opReq.IsDeleteOperation, appName, operations[opReq.OperationName].Templates, patchFile, kubeConfigs)
			if err != nil {
				ee.Summary = fmt.Sprintf("Error while %s %s application", stat, appName)
				ee.Details = err.Error()
				ee.ErrorCode = errors.GetCode(err)
				ee.ProbableCause = errors.GetCause(err)
				ee.SuggestedRemediation = errors.GetRemedy(err)
				hh.StreamErr(ee, err)
				return
			}
			ee.Summary = fmt.Sprintf("%s application %s successfully", appName, stat)
			ee.Details = fmt.Sprintf("The %s application is now %s.", appName, stat)
			hh.StreamInfo(ee)
		}(istio, e)
	default:
		istio.StreamErr(e, ErrOpInvalid)
	}

	return nil
}

// CreateKubeconfigs creates and writes passed kubeconfig onto the filesystem
func (istio *Istio) CreateKubeconfigs(kubeconfigs []string) error {
	var errs = make([]error, 0)
	for _, kubeconfig := range kubeconfigs {
		kconfig := models.Kubeconfig{}
		err := yaml.Unmarshal([]byte(kubeconfig), &kconfig)
		if err != nil {
			errs = append(errs, err)
			continue
		}

		// To have control over what exactly to take in on kubeconfig
		istio.KubeconfigHandler.SetKey("kind", kconfig.Kind)
		istio.KubeconfigHandler.SetKey("apiVersion", kconfig.APIVersion)
		istio.KubeconfigHandler.SetKey("current-context", kconfig.CurrentContext)
		err = istio.KubeconfigHandler.SetObject("preferences", kconfig.Preferences)
		if err != nil {
			errs = append(errs, err)
			continue
		}

		err = istio.KubeconfigHandler.SetObject("clusters", kconfig.Clusters)
		if err != nil {
			errs = append(errs, err)
			continue
		}

		err = istio.KubeconfigHandler.SetObject("users", kconfig.Users)
		if err != nil {
			errs = append(errs, err)
			continue
		}

		err = istio.KubeconfigHandler.SetObject("contexts", kconfig.Contexts)
		if err != nil {
			errs = append(errs, err)
			continue
		}
	}
	if len(errs) == 0 {
		return nil
	}
	return mergeErrors(errs)
}

// ProcessOAM will handles the grpc invocation for handling OAM objects
func (istio *Istio) ProcessOAM(ctx context.Context, oamReq adapter.OAMRequest) (string, error) {
	err := istio.CreateKubeconfigs(oamReq.K8sConfigs)
	if err != nil {
		return "", err
	}
	kubeconfigs := oamReq.K8sConfigs
	var comps []v1alpha1.Component
	for _, acomp := range oamReq.OamComps {
		comp, configErr := oam.ParseApplicationComponent(acomp)
		if configErr != nil {
			istio.Log.Error(ErrParseOAMComponent)
			continue
		}
		comps = append(comps, comp)
	}

	config, err := oam.ParseApplicationConfiguration(oamReq.OamConfig)
	if err != nil {
		istio.Log.Error(ErrParseOAMConfig)
	}

	// If operation is delete then first HandleConfiguration and then handle the deployment
	if oamReq.DeleteOp {
		// Process configuration
		msg2, appConfiguration := istio.HandleApplicationConfiguration(config, oamReq.DeleteOp, kubeconfigs)
		if appConfiguration != nil {
			return msg2, ErrProcessOAM(appConfiguration)
		}

		// Process components
		msg1, componentsErr := istio.HandleComponents(comps, oamReq.DeleteOp, kubeconfigs)
		if componentsErr != nil {
			return msg1 + "\n" + msg2, ErrProcessOAM(componentsErr)
		}

		return msg1 + "\n" + msg2, nil
	}

	// Process components
	msg1, err := istio.HandleComponents(comps, oamReq.DeleteOp, kubeconfigs)
	if err != nil {
		return msg1, ErrProcessOAM(err)
	}

	// Process configuration
	msg2, err := istio.HandleApplicationConfiguration(config, oamReq.DeleteOp, kubeconfigs)
	if err != nil {
		return msg1 + "\n" + msg2, ErrProcessOAM(err)
	}

	return msg1 + "\n" + msg2, nil
}

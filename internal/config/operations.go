package config

import (
	"fmt"

	"github.com/layer5io/meshery-adapter-library/adapter"
	"github.com/layer5io/meshery-adapter-library/common"
	"github.com/layer5io/meshery-adapter-library/meshes"
	"github.com/layer5io/meshkit/utils"
)

var (
	ServiceName = "service_name"
)

func GetOperations(dev adapter.Operations, version string) adapter.Operations {
	var adapterVersions []adapter.Version
	versions, _ := utils.GetLatestReleaseTagsSorted("istio", "istio")
	for _, v := range versions {
		adapterVersions = append(adapterVersions, adapter.Version(v))
	}
	// Add Istio networking resources to sample applications
	dev[common.BookInfoOperation].Templates = append(dev[common.BookInfoOperation].Templates, "file://templates/bookinfo/gateway.yaml")
	dev[common.HTTPBinOperation].Templates = append(dev[common.HTTPBinOperation].Templates, "file://templates/httpbin/gateway.yaml")
	dev[common.ImageHubOperation].Templates = append(dev[common.ImageHubOperation].Templates, "file://templates/imagehub/gateway.yaml")
	dev[common.EmojiVotoOperation].Templates = append(dev[common.EmojiVotoOperation].Templates, "file://templates/emojivoto/gateway.yaml")

	dev[IstioOperation] = &adapter.Operation{
		Type:                 int32(meshes.OpCategory_INSTALL),
		Description:          "Istio Service Mesh",
		Versions:             adapterVersions,
		AdditionalProperties: map[string]string{},
	}

	dev[LabelNamespace] = &adapter.Operation{
		Type:        int32(meshes.OpCategory_CONFIGURE),
		Description: "Automatic Sidecar Injection",
	}

	dev[PrometheusAddon] = &adapter.Operation{
		Type:        int32(meshes.OpCategory_CONFIGURE),
		Description: "Add-on: Prometheus",
		Templates: []adapter.Template{
			adapter.Template(fmt.Sprintf("https://raw.githubusercontent.com/istio/istio/%s/samples/addons/prometheus.yaml", version)),
		},
		AdditionalProperties: map[string]string{
			ServiceName:      "prometheus",
			ServicePatchFile: "file://templates/patches/service-loadbalancer.json",
		},
	}

	dev[GrafanaAddon] = &adapter.Operation{
		Type:        int32(meshes.OpCategory_CONFIGURE),
		Description: "Add-on: Grafana",
		Templates: []adapter.Template{
			adapter.Template(fmt.Sprintf("https://raw.githubusercontent.com/istio/istio/%s/samples/addons/grafana.yaml", version)),
		},
		AdditionalProperties: map[string]string{
			ServiceName:      "grafana",
			ServicePatchFile: "file://templates/patches/service-loadbalancer.json",
		},
	}

	dev[KialiAddon] = &adapter.Operation{
		Type:        int32(meshes.OpCategory_CONFIGURE),
		Description: "Add-on: Kiali",
		Templates: []adapter.Template{
			adapter.Template(fmt.Sprintf("https://raw.githubusercontent.com/istio/istio/%s/samples/addons/kiali.yaml", version)),
		},
		AdditionalProperties: map[string]string{
			ServiceName:      "kiali",
			ServicePatchFile: "file://templates/patches/service-loadbalancer.json",
		},
	}

	dev[JaegerAddon] = &adapter.Operation{
		Type:        int32(meshes.OpCategory_CONFIGURE),
		Description: "Add-on: Jaeger",
		Templates: []adapter.Template{
			adapter.Template(fmt.Sprintf("https://raw.githubusercontent.com/istio/istio/%s/samples/addons/jaeger.yaml", version)),
		},
		AdditionalProperties: map[string]string{
			ServiceName:      "jaeger-collector",
			ServicePatchFile: "file://templates/patches/service-loadbalancer.json",
		},
	}

	dev[ZipkinAddon] = &adapter.Operation{
		Type:        int32(meshes.OpCategory_CONFIGURE),
		Description: "Add-on: Zipkin",
		Templates: []adapter.Template{
			adapter.Template(fmt.Sprintf("https://raw.githubusercontent.com/istio/istio/%s/samples/addons/extras/zipkin.yaml", version)),
		},
		AdditionalProperties: map[string]string{
			ServiceName:      "zipkin",
			ServicePatchFile: "file://templates/patches/service-loadbalancer.json",
		},
	}

	dev[IstioVetOperation] = &adapter.Operation{
		Type:        int32(meshes.OpCategory_VALIDATE),
		Description: "Analyze Running Configuration",
	}

	dev[EnvoyFilterOperation] = &adapter.Operation{
		Type:        int32(meshes.OpCategory_CONFIGURE),
		Description: "Envoy Filter for Image Hub",
		Versions:    adapter.NoneVersion,
		Templates: []adapter.Template{
			"file://templates/imagehub/rate_limit_filter.yaml",
		},
		AdditionalProperties: map[string]string{
			ServiceName:     "api-v1",
			FilterPatchFile: "file://templates/imagehub/filter_patch.json",
		},
	}

	dev[DenyAllPolicyOperation] = &adapter.Operation{
		Type:        int32(meshes.OpCategory_CONFIGURE),
		Description: "Policy: Deny-All",
		Templates: []adapter.Template{
			"file://templates/policies/denyall.yaml",
		},
	}

	dev[StrictMTLSPolicyOperation] = &adapter.Operation{
		Type:        int32(meshes.OpCategory_CONFIGURE),
		Description: "Policy: Strict MTLS",
		Templates: []adapter.Template{
			"file://templates/policies/strict.yaml",
		},
	}

	dev[MutualMTLSPolicyOperation] = &adapter.Operation{
		Type:        int32(meshes.OpCategory_CONFIGURE),
		Description: "Policy: Mutual MTLS",
		Templates: []adapter.Template{
			"file://templates/policies/mutual.yaml",
		},
	}

	dev[DisableMTLSPolicyOperation] = &adapter.Operation{
		Type:        int32(meshes.OpCategory_CONFIGURE),
		Description: "Policy: Disable MTLS",
		Templates: []adapter.Template{
			"file://templates/policies/disable.yaml",
		},
	}

	return dev
}

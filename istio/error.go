// Package istio - Error codes for the adapter
package istio

import (
	"github.com/layer5io/meshkit/errors"
)

var (
	// Error code for failed service mesh installation

	// ErrInstallUsingIstioctlCode represents the errors which are generated
	// during istio service mesh install process
	ErrInstallUsingIstioctlCode = "1003"

	// ErrUnzipFileCode represents the errors which are generated
	// during unzip process
	ErrUnzipFileCode = "1004"

	// ErrTarXZFCode represents the errors which are generated
	// during decompressing and extracting tar.gz file
	ErrTarXZFCode = "1005"

	// ErrMeshConfigCode represents the errors which are generated
	// when an invalid mesh config is found
	ErrMeshConfigCode = "1006"

	// ErrRunIstioCtlCmdCode represents the errors which are generated
	// during fetch manifest process
	ErrRunIstioCtlCmdCode = "1007"

	// ErrSampleAppCode represents the errors which are generated
	// duing sample app installation
	ErrSampleAppCode = "1008"

	// ErrEnvoyFilterCode represents the errors which are generated
	// duing envoy filter patching
	ErrEnvoyFilterCode = "1009"

	// ErrApplyPolicyCode represents the errors which are generated
	// duing policy apply operation
	ErrApplyPolicyCode = "1010"

	// ErrCustomOperationCode represents the errors which are generated
	// when an invalid addon operation is requested
	ErrCustomOperationCode = "1011"

	// ErrAddonFromTemplateCode represents the errors which are generated
	// during addon deployment process
	ErrAddonFromTemplateCode = "1012"

	//ErrInvalidInstallationProfileCode implies error while invalid profile option is passed in pattern file
	ErrInvalidInstallationProfileCode = "1013"

	// ErrCreatingIstioClientCode represents the errors which are generated
	// during creating istio client process
	ErrCreatingIstioClientCode = "1014"

	// ErrIstioVetSyncCode represents the errors which are generated
	// during istio-vet sync process
	ErrIstioVetSyncCode = "1015"

	// ErrIstioVetCode represents the errors which are generated
	// during istio-vet process
	ErrIstioVetCode = "1016"

	// ErrParseOAMComponentCode represents the error code which is
	// generated during the OAM component parsing
	ErrParseOAMComponentCode = "1017"

	// ErrParseOAMConfigCode represents the error code which is
	// generated during the OAM configuration parsing
	ErrParseOAMConfigCode = "1018"

	// ErrNilClientCode represents the error code which is
	// generated when kubernetes client is nil
	ErrNilClientCode = "1019"

	// ErrParseIstioCoreComponentCode represents the error code which is
	// generated when istio core component manifest parsing fails
	ErrParseIstioCoreComponentCode = "1020"

	// ErrInvalidOAMComponentTypeCode represents the error code which is
	// generated when an invalid oam component is requested
	ErrInvalidOAMComponentTypeCode = "1021"

	// ErrOpInvalidCode represents the error code which is
	// generated when an invalid operation is requested
	ErrOpInvalidCode = "1022"

	// ErrIstioCoreComponentFailCode represents the error code which is
	// generated when an istio core operations fails
	ErrIstioCoreComponentFailCode = "1023"

	// ErrProcessOAMCode represents the error code which is
	// generated when an OAM operations fails
	ErrProcessOAMCode = "1024"

	// ErrApplyHelmChartCode represents the error which are generated
	// during the process of applying helm chart
	ErrApplyHelmChartCode = "1025"

	// ErrGettingIstioReleaseCode implies failure while failing istio release
	// bundle
	ErrGettingIstioReleaseCode = "1026"

	// ErrUnsupportedPlatformCode implies unavailbility of Istio on the
	// requested platform
	ErrUnsupportedPlatformCode = "1027"

	// ErrIstioctlNotFoundCode implies istioctl couldn't be found anywhere
	// on the fs
	ErrIstioctlNotFoundCode = "1028"

	// ErrDownloadingTarCode implies error while downloading istio tar
	ErrDownloadingTarCode = "1029"

	// ErrUnpackingTarCode implies error while unpacking istio release
	// bundle tar
	ErrUnpackingTarCode = "1030"

	// ErrMakingBinExecutableCode implies error while makng istioctl executable
	ErrMakingBinExecutableCode = "1031"

	// ErrLoadNamespaceCode implies error while finding namespace
	ErrLoadNamespaceCode = "1032"
	// ErrLoadNamespaceCode implies error while finding namespace
	ErrFetchIstioVersionsCode = "1033"

	ErrFetchIstioVersions = errors.New(ErrFetchIstioVersionsCode, errors.Alert, []string{"could not get any istio versions"}, []string{"versions for istio could not be fetched"}, []string{"could not reach github.com/istio/istio/releases", "no versions could be fetched from istio release page"}, []string{"make sure adapter is reachable to github"})
	// ErrOpInvalid represents the errors which are generated
	// when an invalid operation is requested
	ErrOpInvalid = errors.New(ErrOpInvalidCode, errors.Alert, []string{"Invalid operation"}, []string{"Istio adapter received an invalid operation from the meshey server"}, []string{"The operation is not supported by the adapter", "Invalid operation name"}, []string{"Check if the operation name is valid and supported by the adapter"})

	// ErrParseOAMComponent represents the error which is
	// generated during the OAM component parsing
	ErrParseOAMComponent = errors.New(ErrParseOAMComponentCode, errors.Alert, []string{"error parsing the component"}, []string{"Error occurred while parsing application component in the OAM request made by Meshery server"}, []string{"Could not unmarshall configuration component received via ProcessOAM gRPC call into a valid Component struct"}, []string{"Check if Meshery Server is creating valid component for ProcessOAM gRPC call. This error should never happen and can be reported as a bug in Meshery Server. Also check if Meshery Server and adapters are referring to same component struct provided in MeshKit."})

	// ErrParseOAMConfig represents the error which is
	// generated during the OAM configuration parsing
	ErrParseOAMConfig = errors.New(ErrParseOAMConfigCode, errors.Alert, []string{"error parsing the configuration"}, []string{"Error occurred while parsing configuration in the request made by Meshery Server"}, []string{"Could not unmarshall OAM config received via ProcessOAM gRPC call into a valid Config struct"}, []string{"Check if Meshery Server is creating valid config for ProcessOAM gRPC call. This error should never happen and can be reported as a bug in Meshery Server. Also, confirm that Meshery Server and Adapters are referring to same config struct provided in MeshKit"})

	// ErrNilClient represents the error which is
	// generated when kubernetes client is nil
	ErrNilClient = errors.New(ErrNilClientCode, errors.Alert, []string{"kubernetes client not initialized"}, []string{"Kubernetes client is nil"}, []string{"kubernetes client not initialized"}, []string{"Reconnect the adaptor to Meshery server"})

	// ErrUnsupportedPlatform represents runtime platform is
	// unsupported
	ErrUnsupportedPlatform = errors.New(ErrUnsupportedPlatformCode, errors.Alert, []string{"requested platform is not supported by Istio"}, []string{"Istio only supports Windows, Linux and Darwin"}, []string{}, []string{""})

	// ErrIstioctlNotFound implies istioctl was not found locally
	ErrIstioctlNotFound = errors.New(ErrIstioctlNotFoundCode, errors.Alert, []string{"Unable to find Istioctl"}, []string{}, []string{}, []string{})
)

// ErrInstallUsingIstioctl is the error for install mesh
func ErrInstallUsingIstioctl(err error) error {
	return errors.New(ErrInstallUsingIstioctlCode, errors.Alert, []string{"Error with istio operation"}, []string{"Error occurred while installing istio mesh through istioctl", err.Error()}, []string{}, []string{})
}

// ErrUnzipFile is the error for unzipping the file
func ErrUnzipFile(err error) error {
	return errors.New(ErrUnzipFileCode, errors.Alert, []string{"Error while unzipping"}, []string{err.Error()}, []string{"File might be corrupt"}, []string{})
}

// ErrTarXZF is the error for unzipping the file
func ErrTarXZF(err error) error {
	return errors.New(ErrTarXZFCode, errors.Alert, []string{"Error while extracting file"}, []string{err.Error()}, []string{"/The gzip might be corrupt"}, []string{})
}

// ErrMeshConfig is the error for mesh config
func ErrMeshConfig(err error) error {
	return errors.New(ErrMeshConfigCode, errors.Alert, []string{"Error configuration mesh"}, []string{err.Error(), "Error getting MeshSpecKey config from in-memory configuration"}, []string{}, []string{})
}

// ErrRunIstioCtlCmd is the error for mesh port forward
func ErrRunIstioCtlCmd(err error, des string) error {
	return errors.New(ErrRunIstioCtlCmdCode, errors.Alert, []string{"Error running istioctl command"}, []string{err.Error()}, []string{"Corrupted istioctl binary", "Command might be invalid"}, []string{})
}

// ErrSampleApp is the error for streaming event
func ErrSampleApp(err error) error {
	return errors.New(ErrSampleAppCode, errors.Alert, []string{"Error with sample app operation"}, []string{err.Error(), "Error occurred while trying to install a sample application using manifests"}, []string{"Invalid kubeclient config", "Invalid manifest"}, []string{"Reconnect your adapter to meshery server to refresh the kubeclient"})
}

// ErrEnvoyFilter is the error for streaming event
func ErrEnvoyFilter(err error) error {
	return errors.New(ErrEnvoyFilterCode, errors.Alert, []string{"Error with envoy filter operation"}, []string{err.Error()}, []string{}, []string{})
}

// ErrApplyPolicy is the error for streaming event
func ErrApplyPolicy(err error) error {
	return errors.New(ErrApplyPolicyCode, errors.Alert, []string{"Error with apply policy operation"}, []string{err.Error(), "Error occurred while trying to install a sample application using manifests"}, []string{"Invalid kubeclient config", "Invalid manifest"}, []string{"Reconnect your adapter to meshery server to refresh the kubeclient"})
}

// ErrAddonFromTemplate is the error for streaming event
func ErrAddonFromTemplate(err error) error {
	return errors.New(ErrAddonFromTemplateCode, errors.Alert, []string{"Error with addon install operation"}, []string{err.Error()}, []string{}, []string{})
}

// ErrCustomOperation is the error for streaming event
func ErrCustomOperation(err error) error {
	return errors.New(ErrCustomOperationCode, errors.Alert, []string{"Error with custom operation"}, []string{"Error occurred while applying custom manifest to the cluster", err.Error()}, []string{"Invalid kubeclient config", "Invalid manifest"}, []string{})
}

// ErrCreatingIstioClient is the error for streaming event
func ErrCreatingIstioClient(err error) error {
	return errors.New(ErrCreatingIstioClientCode, errors.Alert, []string{"Unable to create a new istio client"}, []string{err.Error()}, []string{}, []string{})
}

// ErrIstioVetSync is the error for streaming event
func ErrIstioVetSync(err error) error {
	return errors.New(ErrIstioVetSyncCode, errors.Alert, []string{"Failed to sync"}, []string{err.Error()}, []string{}, []string{})
}

// ErrIstioVet is the error for streaming event
func ErrIstioVet(err error) error {
	return errors.New(ErrIstioVetCode, errors.Alert, []string{"error while running Istio vet command"}, []string{err.Error()}, []string{}, []string{})
}

// ErrParseIstioCoreComponent is the error when istio core component manifest parsing fails
func ErrParseIstioCoreComponent(err error) error {
	return errors.New(ErrParseIstioCoreComponentCode, errors.Alert, []string{"istio core component manifest parsing failing"}, []string{err.Error()}, []string{}, []string{})
}

// ErrInvalidOAMComponentType is the error when the OAM component name is not valid
func ErrInvalidOAMComponentType(compName string) error {
	return errors.New(ErrInvalidOAMComponentTypeCode, errors.Alert, []string{"invalid OAM component name: ", compName}, []string{}, []string{}, []string{})
}

// ErrIstioCoreComponentFail is the error when core Istion component processing fails
func ErrIstioCoreComponentFail(err error) error {
	return errors.New(ErrIstioCoreComponentFailCode, errors.Alert, []string{"error in istio core component"}, []string{err.Error()}, []string{}, []string{})
}

// ErrProcessOAM is a generic error which is thrown when an OAM operations fails
func ErrProcessOAM(err error) error {
	return errors.New(ErrProcessOAMCode, errors.Alert, []string{"error performing OAM operations"}, []string{err.Error()}, []string{}, []string{})
}

// ErrApplyHelmChart is the error for applying helm chart
func ErrApplyHelmChart(err error) error {
	return errors.New(ErrApplyHelmChartCode, errors.Alert, []string{"Error occurred while applying Helm Chart"}, []string{err.Error()}, []string{}, []string{})
}

// ErrGettingIstioRelease is the error when the yaml unmarshal fails
func ErrGettingIstioRelease(err error) error {
	return errors.New(ErrGettingIstioReleaseCode, errors.Alert, []string{"Error occurred while fetching Istio release artifacts"}, []string{err.Error()}, []string{}, []string{})
}

// ErrDownloadingTar is the error when tar download fails
func ErrDownloadingTar(err error) error {
	return errors.New(ErrDownloadingTarCode, errors.Alert, []string{"Error occurred while downloading Istio tar"}, []string{err.Error()}, []string{}, []string{})
}

// ErrUnpackingTar is the error when tar unpack fails
func ErrUnpackingTar(err error) error {
	return errors.New(ErrUnpackingTarCode, errors.Alert, []string{"Error occurred while unpacking tar"}, []string{err.Error()}, []string{}, []string{})
}

// ErrMakingBinExecutable occurs when istioctl binary couldn't be made
// executable
func ErrMakingBinExecutable(err error) error {
	return errors.New(ErrMakingBinExecutableCode, errors.Alert, []string{"Error while making istioctl an executable"}, []string{err.Error()}, []string{}, []string{})
}

// ErrLoadNamespace implies error while finding namespace
func ErrLoadNamespace(err error, str string) error {
	return errors.New(ErrLoadNamespaceCode, errors.Alert, []string{"Error while labeling namespace:", str}, []string{err.Error()}, []string{}, []string{})
}

// ErrInvalidInstallationProfile implies error while invalid profile option is passed in pattern file
func ErrInvalidInstallationProfile(str string) error {
	return errors.New(ErrInvalidInstallationProfileCode, errors.Alert, []string{"Error while installing istio due to wrong profile"}, []string{"Gotten profile " + str}, []string{"Invalid profile passed"}, []string{"Provide one of the profiles: \"demo\",\"minimal\",\"default\" profiles"})
}

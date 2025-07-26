package istio

import (
	"fmt"
	"reflect"
	"strings"
	"sync"
	"time"

	"github.com/aspenmesh/istio-vet/pkg/istioclient"
	"github.com/aspenmesh/istio-vet/pkg/vetter"
	"github.com/aspenmesh/istio-vet/pkg/vetter/applabel"
	"github.com/aspenmesh/istio-vet/pkg/vetter/conflictingvirtualservicehost"
	"github.com/aspenmesh/istio-vet/pkg/vetter/danglingroutedestinationhost"
	"github.com/aspenmesh/istio-vet/pkg/vetter/meshversion"
	"github.com/aspenmesh/istio-vet/pkg/vetter/podsinmesh"
	"github.com/aspenmesh/istio-vet/pkg/vetter/serviceassociation"
	"github.com/aspenmesh/istio-vet/pkg/vetter/serviceportprefix"
	"github.com/layer5io/meshery-adapter-library/meshes"
	internalconfig "github.com/layer5io/meshery-istio/internal/config"
	"github.com/layer5io/meshkit/errors"
	mesherykube "github.com/layer5io/meshkit/utils/kubernetes"
	istioinformer "istio.io/client-go/pkg/informers/externalversions"
	"k8s.io/client-go/informers"
)

const istioVetSyncTimeout = 10 // istio vet sync timeout in seconds

type metaInformerFactory struct {
	k8s   informers.SharedInformerFactory
	istio istioinformer.SharedInformerFactory
}

func (m *metaInformerFactory) K8s() informers.SharedInformerFactory {
	return m.k8s
}
func (m *metaInformerFactory) Istio() istioinformer.SharedInformerFactory {
	return m.istio
}

// RunVet runs istio-vet
func (istio *Istio) RunVet(ch chan<- *meshes.EventsResponse, kubeconfigs []string) {
	defer close(ch)
	var wg sync.WaitGroup
	for _, k8sconfig := range kubeconfigs {
		wg.Add(1)
		go func(k8sconfig string) {
			defer wg.Done()
			mclient, err := mesherykube.New([]byte(k8sconfig))
			if err != nil {
				e := &meshes.EventsResponse{
					Component:     internalconfig.ServerConfig["type"],
					ComponentName: internalconfig.ServerConfig["name"],
				}
				e.EventType = meshes.EventType_ERROR
				e.Details = ErrCreatingIstioClient(err).Error()
				e.Summary = "Unable to create k8s client"
				e.ErrorCode = errors.GetCode(err)
				e.ProbableCause = errors.GetCause(err)
				e.SuggestedRemediation = errors.GetRemedy(err)
				ch <- e
			}
			istioClient, err := istioclient.New(&mclient.RestConfig)
			if err != nil {
				e := &meshes.EventsResponse{
					Component:     internalconfig.ServerConfig["type"],
					ComponentName: internalconfig.ServerConfig["name"],
				}
				e.EventType = meshes.EventType_ERROR
				e.Details = ErrCreatingIstioClient(err).Error()
				e.Summary = "Unable to create istio client"
				e.ErrorCode = errors.GetCode(err)
				e.ProbableCause = errors.GetCause(err)
				e.SuggestedRemediation = errors.GetRemedy(err)
				ch <- e
			}

			kubeInformerFactory := informers.NewSharedInformerFactory(mclient.KubeClient, 0)
			istioInformerFactory := istioinformer.NewSharedInformerFactory(istioClient, 0)
			informerFactory := &metaInformerFactory{
				k8s:   kubeInformerFactory,
				istio: istioInformerFactory,
			}

			vList := []vetter.Vetter{
				vetter.Vetter(podsinmesh.NewVetter(informerFactory)),
				vetter.Vetter(meshversion.NewVetter(informerFactory)),
				vetter.Vetter(applabel.NewVetter(informerFactory)),
				vetter.Vetter(serviceportprefix.NewVetter(informerFactory)),
				vetter.Vetter(serviceassociation.NewVetter(informerFactory)),
				vetter.Vetter(danglingroutedestinationhost.NewVetter(informerFactory)),
				vetter.Vetter(conflictingvirtualservicehost.NewVetter(informerFactory)),
			}

			stopCh := make(chan struct{})

			kubeInformerFactory.Start(stopCh)
			oks, timedout := completeBefore(istioVetSyncTimeout, func() map[reflect.Type]bool {
				return kubeInformerFactory.WaitForCacheSync(stopCh)
			})
			if timedout {
				e := &meshes.EventsResponse{
					Component:     internalconfig.ServerConfig["type"],
					ComponentName: internalconfig.ServerConfig["name"],
				}
				e.EventType = meshes.EventType_ERROR
				e.Details = ErrIstioVetSync(fmt.Errorf("istio service mesh was either not found or is not deployed")).Error()
				e.Summary = "Failed to sync: Request timed out"
				e.ErrorCode = errors.GetCode(err)
				e.ProbableCause = errors.GetCause(err)
				e.SuggestedRemediation = errors.GetRemedy(err)
				ch <- e
				close(stopCh)
				return
			}
			for inf, ok := range oks {
				if !ok {
					e := &meshes.EventsResponse{
						Component:     internalconfig.ServerConfig["type"],
						ComponentName: internalconfig.ServerConfig["name"],
					}
					e.EventType = meshes.EventType_ERROR
					e.Details = ErrIstioVetSync(fmt.Errorf("%s", inf)).Error()
					e.Summary = "Failed to sync"
					e.ErrorCode = errors.GetCode(err)
					e.ProbableCause = errors.GetCause(err)
					e.SuggestedRemediation = errors.GetRemedy(err)
					ch <- e
					return
				}
			}

			istioInformerFactory.Start(stopCh)
			oks, timedout = completeBefore(istioVetSyncTimeout, func() map[reflect.Type]bool {
				return istioInformerFactory.WaitForCacheSync(stopCh)
			})
			if timedout {
				e := &meshes.EventsResponse{
					Component:     internalconfig.ServerConfig["type"],
					ComponentName: internalconfig.ServerConfig["name"],
				}
				e.EventType = meshes.EventType_ERROR
				e.Details = ErrIstioVetSync(fmt.Errorf("istio service mesh was either not found or is not deployed")).Error()
				e.Summary = "Failed to sync: Request timed out"
				e.ErrorCode = errors.GetCode(err)
				e.ProbableCause = errors.GetCause(err)
				e.SuggestedRemediation = errors.GetRemedy(err)
				ch <- e
				close(stopCh)
				return
			}
			for inf, ok := range oks {
				if !ok {
					e := &meshes.EventsResponse{
						Component:     internalconfig.ServerConfig["type"],
						ComponentName: internalconfig.ServerConfig["name"],
					}
					e.EventType = meshes.EventType_ERROR
					e.Details = ErrIstioVetSync(fmt.Errorf("%s", inf)).Error()
					e.Summary = "Failed to sync"
					e.ErrorCode = errors.GetCode(err)
					e.ProbableCause = errors.GetCause(err)
					e.SuggestedRemediation = errors.GetRemedy(err)
					ch <- e
					return
				}
			}
			close(stopCh)

			for _, v := range vList {
				nList, err := v.Vet()
				if err != nil {
					e := &meshes.EventsResponse{}
					e.Summary = fmt.Sprintf("Vetter: %s reported error", v.Info().GetId())
					e.Details = err.Error()
					e.EventType = meshes.EventType_ERROR
					ch <- e
					continue
				}
				if len(nList) > 0 {
					for i := range nList {
						e := &meshes.EventsResponse{}

						var ts []string
						for k, v := range nList[i].Attr {
							ts = append(ts, "${"+k+"}", v)
						}
						r := strings.NewReplacer(ts...)
						e.Summary = r.Replace(nList[i].GetSummary())
						e.Details = r.Replace(nList[i].GetMsg())
						switch nList[i].GetLevel().String() {
						case "WARNING":
							e.EventType = meshes.EventType_WARN
						case "ERROR":
							e.EventType = meshes.EventType_ERROR
						default:
							e.EventType = meshes.EventType_INFO
						}
						ch <- e
					}
				} else {
					e := &meshes.EventsResponse{}
					istio.Log.Debug(fmt.Sprintf("Vetter %s ran successfully and generated no notes", v.Info().GetId()))
					e.Summary = fmt.Sprintf("Vetter: %s ran successfully", v.Info().GetId())
					e.Details = "No notes generated"
					e.EventType = meshes.EventType_INFO
					ch <- e
				}
			}
		}(k8sconfig)
	}
}

// StreamWarn streams a warning message to the channel
func (istio *Istio) StreamWarn(e *meshes.EventsResponse, err error) {
	istio.Log.Warn(err)
	e.EventType = meshes.EventType_WARN
	istio.EventStreamer.Publish(e)
}

// completeBefore executes the callback function but if the callback function
// doesn't returns before the specified timeout then it returns nil and true
// indicating that the request has timed out
func completeBefore(timeout time.Duration, cb func() map[reflect.Type]bool) (map[reflect.Type]bool, bool) {
	tch := make(chan bool, timeout)
	resch := make(chan map[reflect.Type]bool)

	go func() {
		resch <- cb()
	}()
	go func() {
		time.Sleep(timeout * time.Second)
		tch <- true
	}()

	select {
	case res := <-resch:
		return res, false
	case <-tch:
		return nil, true
	}
}

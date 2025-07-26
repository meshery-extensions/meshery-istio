package istio

import (
	"testing"

	"github.com/layer5io/meshery-adapter-library/adapter"
	"github.com/layer5io/meshery-adapter-library/status"
)

func TestIstio_installAddon(t *testing.T) {
	type args struct {
		namespace string
		del       bool
		service   string
		patches   []string
		templates []adapter.Template
	}

	tests := []struct {
		name        string
		args        args
		want        string
		kubeconfigs []string
		wantErr     bool
	}{
		// TODO: Add test cases.
		{
			name: "no patches",
			args: args{
				namespace: "default",
				del:       false,
				service:   "test",
				patches:   nil,
				templates: []adapter.Template{
					"https://raw.githubusercontent.com/istio/istio/master/samples/addons/jaeger.yaml",
				},
			},
			want:    status.Installed,
			wantErr: true,
		},
		{
			name: "no templates",
			args: args{
				namespace: "default",
				del:       false,
				service:   "test",
				patches:   nil,
				templates: nil,
			},
			want:    status.Installed,
			wantErr: true,
		},
		{
			name: "delete operation",
			args: args{
				namespace: "default",
				del:       true,
				service:   "test",
				patches:   nil,
				templates: nil,
			},
			want:    status.Installed,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			istio := &Istio{
				Adapter: adapter.Adapter{
					Config: getConfigHandler(t),
					Log:    getLoggerHandler(t),
				},
			}
			got, err := istio.installAddon(tt.args.namespace, tt.args.del, tt.args.service, tt.args.patches, tt.args.templates, tt.kubeconfigs)
			if (err != nil) == tt.wantErr {
				t.Errorf("Istio.installAddon() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Istio.installAddon() = %v, want %v", got, tt.want)
			}
		})
	}
}

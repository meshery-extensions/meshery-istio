// Package istio - Common operations for the adapter
package istio

import (
	"testing"

	"github.com/layer5io/meshery-adapter-library/adapter"
	"github.com/layer5io/meshery-adapter-library/status"
)

func TestIstio_applyCustomOperation(t *testing.T) {
	type args struct {
		namespace string
		manifest  string
		isDel     bool
	}

	tests := []struct {
		name        string
		args        args
		kubeconfigs []string
		want        string
		wantErr     bool
	}{
		// TODO: Add test cases.
		{
			name: "no manifest or empty manifest",
			args: args{
				namespace: "default",
				manifest:  "",
				isDel:     false,
			},
			want:    status.Completed,
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
			got, err := istio.applyCustomOperation(tt.args.namespace, tt.args.manifest, tt.args.isDel, tt.kubeconfigs)
			if (err != nil) == tt.wantErr {
				t.Errorf("Istio.applyCustomOperation() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Istio.applyCustomOperation() = %v, want %v", got, tt.want)
			}
		})
	}
}

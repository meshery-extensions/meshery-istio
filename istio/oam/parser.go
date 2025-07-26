package oam

import (
	"encoding/json"

	"github.com/layer5io/meshkit/models/oam/core/v1alpha1"
)

// ParseApplicationComponent converts json application component to go struct
func ParseApplicationComponent(jsn string) (acomp v1alpha1.Component, err error) {
	err = json.Unmarshal([]byte(jsn), &acomp)
	return
}

// ParseApplicationConfiguration converts json application configuration to go struct
func ParseApplicationConfiguration(jsn string) (acomp v1alpha1.Configuration, err error) {
	err = json.Unmarshal([]byte(jsn), &acomp)
	return
}

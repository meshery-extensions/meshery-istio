// Copyright Meshery Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package config

import (
	"github.com/layer5io/meshkit/errors"
)

const (
	ErrEmptyConfigCode           = "1000"
	ErrGetLatestReleasesCode     = "1001"
	ErrGetLatestReleaseNamesCode = "1002"
)

var (
	ErrEmptyConfig = errors.New(ErrEmptyConfigCode, errors.Alert, []string{"Config is empty"}, []string{}, []string{}, []string{})
)

// ErrGetLatestReleases is the error for fetching istio releases
func ErrGetLatestReleases(err error) error {
	return errors.New(ErrGetLatestReleasesCode, errors.Alert, []string{"unable to fetch release info"}, []string{err.Error(), "Unable to get the latest release info from the GithubAPI"}, []string{"Checkout https://docs.github.com/en/rest/reference/repos#releases for more info"}, []string{})
}

// ErrGetLatestReleaseNames is the error for fetching istio releases
func ErrGetLatestReleaseNames(err error) error {
	return errors.New(ErrGetLatestReleaseNamesCode, errors.Alert, []string{"failed to extract release names"}, []string{err.Error()}, []string{"Invalid release format"}, []string{})
}

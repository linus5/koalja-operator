//
// Copyright © 2018 Aljabr, Inc.
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
//

package s3

import (
	"testing"

	"github.com/onsi/gomega"
)

func TestBuckerConfigFixEndpoint(t *testing.T) {
	g := gomega.NewGomegaWithT(t)

	bc := BucketConfig{Endpoint: "http://host:123"}
	bc.fixEndpoint()

	g.Expect(bc.Endpoint).To(gomega.Equal("host:123"))
	g.Expect(bc.Secure).To(gomega.BeFalse())

	bc = BucketConfig{Endpoint: "https://host:456/path"}
	bc.fixEndpoint()

	g.Expect(bc.Endpoint).To(gomega.Equal("host:456"))
	g.Expect(bc.Secure).To(gomega.BeTrue())
}

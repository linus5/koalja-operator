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

// Package tracking contains the API specification for tracking pipeline statistics.
package tracking

//go:generate protoc -I .:../../:../../vendor:../../vendor/github.com/gogo/protobuf/protobuf:../../third_party/googleapis/:../../../../../ --gogofaster_out=Mgoogle/protobuf/field_mask.proto=github.com/gogo/protobuf/types,Mgithub.com/golang/protobuf/ptypes/duration/duration.proto=github.com/gogo/protobuf/types,Mgithub.com/golang/protobuf/ptypes/timestamp/timestamp.proto=github.com/gogo/protobuf/types,plugins=grpc,paths=source_relative:. statistics.proto

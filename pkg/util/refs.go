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

package util

// Int64OrDefault returns the value referenced by the given reference
// or the given default value if the given reference is nil.
func Int64OrDefault(valueRef *int64, defaultValue int64) int64 {
	if valueRef == nil {
		return defaultValue
	}
	return *valueRef
}

// StringOrDefault returns the value referenced by the given reference
// or the given default value if the given reference is nil.
func StringOrDefault(valueRef *string, defaultValue string) string {
	if valueRef == nil {
		return defaultValue
	}
	return *valueRef
}

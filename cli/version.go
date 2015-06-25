/**
 * Copyright 2015 Qadium, Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
package cli

// The git commit that was used to build plumb. Compiler fills it in
// with -ldflags
var GitCommit string

const version = "0.1.0"

// This will be rendered as Version-VersionPrerelease, unless
// VersionPrerelease is empty (in which case it's a release)
const versionPrerelease = "beta"

// Copyright 2019 Mustafa Turan. All rights reserved.
// Use of this source code is governed by a Apache License 2.0 license that can
// be found in the LICENSE file.

package bus

// Config holds configuration structure for bus library
type Config struct {
	Next func() string // Unique id generator
}

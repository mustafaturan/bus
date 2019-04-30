// Copyright 2019 Mustafa Turan. All rights reserved.
// Use of this source code is governed by a Apache License 2.0 license that can
// be found in the LICENSE file.

package bus

import (
	"fmt"
)

var b bus

// Configure sets configuration for bus package
func Configure(c Config) error {
	if c.Next == nil {
		return fmt.Errorf("bus: Next() id generator func can't be nil")
	}

	b = bus{next: c.Next}
	return nil
}

// Copyright 2020 just-codeding-0 . All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package thunder

func assert1(guard bool, text string) {
	if !guard {
		panic(text)
	}
}


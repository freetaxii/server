// Copyright 2015-2018 Bret Jordan, All rights reserved.
//
// Use of this source code is governed by an Apache 2.0 license
// that can be found in the LICENSE file in the root of the source tree.

package handlers

/*
authenticate - This method will perform an authentication check to see if the
supplied credentials are valid.
*/
func (s *ServerHandler) authenticate(username, password string, valid bool) bool {

	// If the user did not supply any credentials, then they failed authentication
	if valid == false {
		return false
	}

	if username == "taxii" && password == "password" {
		return true
	}
	return false
}

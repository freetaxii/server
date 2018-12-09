// Copyright 2015-2018 Bret Jordan, All rights reserved.
//
// Use of this source code is governed by an Apache 2.0 license
// that can be found in the LICENSE file in the root of the source
// tree.

package config

import (
	"encoding/json"
)

/*
JSONbool - A boolean type for use with JSON that will track if a value is
undefined or set to null.
*/
type JSONbool struct {
	Value bool
	Valid bool
	Set   bool
}

/*
JSONstring - A string type for use with JSON that will track if a value is
undefined or set to null.
*/
type JSONstring struct {
	Value string
	Valid bool
	Set   bool
}

/*
UnmarshalJSON - This method defines the unmarshal process for the JSONbool type.
We need to track of the value was set, set to null, or actually set to true or
false. This is important since not being set is not the same thing has false or
true. This comes into play when you need to inherit some previous higher level
setting. Meaning if you set a property higher up in the configuration to true or
false, then you want that to be cascaded down and without doing this, it does
not work.
*/
func (b *JSONbool) UnmarshalJSON(data []byte) error {
	// If this method was called, the value was set.
	b.Set = true

	if string(data) == "null" {
		// The key was set to null
		b.Valid = false
		b.Value = false
		return nil
	}

	// The key isn't set to null
	var temp bool
	if err := json.Unmarshal(data, &temp); err != nil {
		return err
	}
	b.Value = temp
	b.Valid = true
	return nil
}

/*
UnmarshalJSON - This method will handle the unmarshalling of content for the
JSONbool type
*/
func (s *JSONstring) UnmarshalJSON(data []byte) error {
	// If this method was called, the value was set.
	s.Set = true

	if string(data) == "null" {
		// The key was set to null
		s.Valid = false
		s.Value = ""
		return nil
	}

	// The key isn't set to null
	var temp string
	if err := json.Unmarshal(data, &temp); err != nil {
		return err
	}
	s.Value = temp
	s.Valid = true
	return nil
}

// Code generated by entc, DO NOT EDIT.

package ent

import (
	"social/ent/friend"
	"social/ent/schema"
)

// The init function reads all schema descriptors with runtime code
// (default values, validators, hooks and policies) and stitches it
// to their package variables.
func init() {
	friendFields := schema.Friend{}.Fields()
	_ = friendFields
	// friendDescShield is the schema descriptor for shield field.
	friendDescShield := friendFields[4].Descriptor()
	// friend.DefaultShield holds the default value on creation for the shield field.
	friend.DefaultShield = friendDescShield.Default.(int8)
}

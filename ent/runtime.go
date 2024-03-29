// Code generated by entc, DO NOT EDIT.

package ent

import (
	"social/ent/friend"
	"social/ent/group"
	"social/ent/groupmember"
	"social/ent/message"
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
	groupFields := schema.Group{}.Fields()
	_ = groupFields
	// groupDescAvatar is the schema descriptor for avatar field.
	groupDescAvatar := groupFields[3].Descriptor()
	// group.DefaultAvatar holds the default value on creation for the avatar field.
	group.DefaultAvatar = groupDescAvatar.Default.(string)
	// groupDescMaxMembers is the schema descriptor for max_members field.
	groupDescMaxMembers := groupFields[5].Descriptor()
	// group.DefaultMaxMembers holds the default value on creation for the max_members field.
	group.DefaultMaxMembers = groupDescMaxMembers.Default.(int)
	// groupDescMode is the schema descriptor for mode field.
	groupDescMode := groupFields[6].Descriptor()
	// group.DefaultMode holds the default value on creation for the mode field.
	group.DefaultMode = groupDescMode.Default.(int8)
	// groupDescType is the schema descriptor for type field.
	groupDescType := groupFields[7].Descriptor()
	// group.DefaultType holds the default value on creation for the type field.
	group.DefaultType = groupDescType.Default.(int8)
	// groupDescStatus is the schema descriptor for status field.
	groupDescStatus := groupFields[8].Descriptor()
	// group.DefaultStatus holds the default value on creation for the status field.
	group.DefaultStatus = groupDescStatus.Default.(int8)
	// groupDescInviteMode is the schema descriptor for invite_mode field.
	groupDescInviteMode := groupFields[9].Descriptor()
	// group.DefaultInviteMode holds the default value on creation for the invite_mode field.
	group.DefaultInviteMode = groupDescInviteMode.Default.(int8)
	// groupDescNotice is the schema descriptor for notice field.
	groupDescNotice := groupFields[10].Descriptor()
	// group.DefaultNotice holds the default value on creation for the notice field.
	group.DefaultNotice = groupDescNotice.Default.(string)
	// groupDescIntroduction is the schema descriptor for introduction field.
	groupDescIntroduction := groupFields[11].Descriptor()
	// group.DefaultIntroduction holds the default value on creation for the introduction field.
	group.DefaultIntroduction = groupDescIntroduction.Default.(string)
	groupmemberFields := schema.GroupMember{}.Fields()
	_ = groupmemberFields
	// groupmemberDescInviter is the schema descriptor for inviter field.
	groupmemberDescInviter := groupmemberFields[3].Descriptor()
	// groupmember.DefaultInviter holds the default value on creation for the inviter field.
	groupmember.DefaultInviter = groupmemberDescInviter.Default.(int64)
	// groupmemberDescRemark is the schema descriptor for remark field.
	groupmemberDescRemark := groupmemberFields[4].Descriptor()
	// groupmember.DefaultRemark holds the default value on creation for the remark field.
	groupmember.DefaultRemark = groupmemberDescRemark.Default.(string)
	// groupmemberDescStatus is the schema descriptor for status field.
	groupmemberDescStatus := groupmemberFields[5].Descriptor()
	// groupmember.DefaultStatus holds the default value on creation for the status field.
	groupmember.DefaultStatus = groupmemberDescStatus.Default.(int8)
	messageFields := schema.Message{}.Fields()
	_ = messageFields
	// messageDescContent is the schema descriptor for content field.
	messageDescContent := messageFields[3].Descriptor()
	// message.DefaultContent holds the default value on creation for the content field.
	message.DefaultContent = messageDescContent.Default.(string)
	// messageDescContentType is the schema descriptor for content_type field.
	messageDescContentType := messageFields[4].Descriptor()
	// message.DefaultContentType holds the default value on creation for the content_type field.
	message.DefaultContentType = messageDescContentType.Default.(int8)
	// messageDescStatus is the schema descriptor for status field.
	messageDescStatus := messageFields[5].Descriptor()
	// message.DefaultStatus holds the default value on creation for the status field.
	message.DefaultStatus = messageDescStatus.Default.(int8)
}

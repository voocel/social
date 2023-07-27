// Code generated by entc, DO NOT EDIT.

package migrate

import (
	"entgo.io/ent/dialect/sql/schema"
	"entgo.io/ent/schema/field"
)

var (
	// FriendsColumns holds the columns for the "friends" table.
	FriendsColumns = []*schema.Column{
		{Name: "id", Type: field.TypeInt64, Increment: true},
		{Name: "uid", Type: field.TypeInt64},
		{Name: "friend_id", Type: field.TypeInt64},
		{Name: "remark", Type: field.TypeString},
		{Name: "shield", Type: field.TypeInt8},
		{Name: "created_at", Type: field.TypeTime, Nullable: true},
		{Name: "updated_at", Type: field.TypeTime, Nullable: true},
		{Name: "deleted_at", Type: field.TypeTime, Nullable: true},
	}
	// FriendsTable holds the schema information for the "friends" table.
	FriendsTable = &schema.Table{
		Name:       "friends",
		Columns:    FriendsColumns,
		PrimaryKey: []*schema.Column{FriendsColumns[0]},
	}
	// FriendAppliesColumns holds the columns for the "friend_applies" table.
	FriendAppliesColumns = []*schema.Column{
		{Name: "id", Type: field.TypeInt64, Increment: true},
		{Name: "from_id", Type: field.TypeInt64},
		{Name: "to_id", Type: field.TypeInt64},
		{Name: "remark", Type: field.TypeString},
		{Name: "status", Type: field.TypeInt8},
		{Name: "created_at", Type: field.TypeTime, Nullable: true},
		{Name: "updated_at", Type: field.TypeTime, Nullable: true},
		{Name: "deleted_at", Type: field.TypeTime, Nullable: true},
	}
	// FriendAppliesTable holds the schema information for the "friend_applies" table.
	FriendAppliesTable = &schema.Table{
		Name:       "friend_applies",
		Columns:    FriendAppliesColumns,
		PrimaryKey: []*schema.Column{FriendAppliesColumns[0]},
	}
	// GroupsColumns holds the columns for the "groups" table.
	GroupsColumns = []*schema.Column{
		{Name: "id", Type: field.TypeInt64, Increment: true},
		{Name: "name", Type: field.TypeString},
		{Name: "owner", Type: field.TypeInt64},
		{Name: "notice", Type: field.TypeString},
		{Name: "created_at", Type: field.TypeTime, Nullable: true},
		{Name: "updated_at", Type: field.TypeTime, Nullable: true},
		{Name: "deleted_at", Type: field.TypeTime, Nullable: true},
	}
	// GroupsTable holds the schema information for the "groups" table.
	GroupsTable = &schema.Table{
		Name:       "groups",
		Columns:    GroupsColumns,
		PrimaryKey: []*schema.Column{GroupsColumns[0]},
	}
	// UsersColumns holds the columns for the "users" table.
	UsersColumns = []*schema.Column{
		{Name: "id", Type: field.TypeInt64, Increment: true},
		{Name: "username", Type: field.TypeString},
		{Name: "password", Type: field.TypeString},
		{Name: "mobile", Type: field.TypeString},
		{Name: "nickname", Type: field.TypeString},
		{Name: "email", Type: field.TypeString},
		{Name: "avatar", Type: field.TypeString, Nullable: true},
		{Name: "summary", Type: field.TypeString, Nullable: true},
		{Name: "sex", Type: field.TypeInt8, Nullable: true},
		{Name: "status", Type: field.TypeInt8, Nullable: true},
		{Name: "birthday", Type: field.TypeTime, Nullable: true},
		{Name: "created_at", Type: field.TypeTime, Nullable: true},
		{Name: "updated_at", Type: field.TypeTime, Nullable: true},
		{Name: "deleted_at", Type: field.TypeTime, Nullable: true},
	}
	// UsersTable holds the schema information for the "users" table.
	UsersTable = &schema.Table{
		Name:       "users",
		Columns:    UsersColumns,
		PrimaryKey: []*schema.Column{UsersColumns[0]},
	}
	// Tables holds all the tables in the schema.
	Tables = []*schema.Table{
		FriendsTable,
		FriendAppliesTable,
		GroupsTable,
		UsersTable,
	}
)

func init() {
}

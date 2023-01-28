// Code generated by entc, DO NOT EDIT.

package friend

const (
	// Label holds the string label denoting the friend type in the database.
	Label = "friend"
	// FieldID holds the string denoting the id field in the database.
	FieldID = "id"
	// FieldUID holds the string denoting the uid field in the database.
	FieldUID = "uid"
	// FieldFriendID holds the string denoting the friend_id field in the database.
	FieldFriendID = "friend_id"
	// FieldRemark holds the string denoting the remark field in the database.
	FieldRemark = "remark"
	// FieldShield holds the string denoting the shield field in the database.
	FieldShield = "shield"
	// FieldCreatedAt holds the string denoting the created_at field in the database.
	FieldCreatedAt = "created_at"
	// FieldUpdatedAt holds the string denoting the updated_at field in the database.
	FieldUpdatedAt = "updated_at"
	// FieldDeletedAt holds the string denoting the deleted_at field in the database.
	FieldDeletedAt = "deleted_at"
	// Table holds the table name of the friend in the database.
	Table = "friends"
)

// Columns holds all SQL columns for friend fields.
var Columns = []string{
	FieldID,
	FieldUID,
	FieldFriendID,
	FieldRemark,
	FieldShield,
	FieldCreatedAt,
	FieldUpdatedAt,
	FieldDeletedAt,
}

// ValidColumn reports if the column name is valid (part of the table columns).
func ValidColumn(column string) bool {
	for i := range Columns {
		if column == Columns[i] {
			return true
		}
	}
	return false
}
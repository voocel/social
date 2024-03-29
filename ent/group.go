// Code generated by entc, DO NOT EDIT.

package ent

import (
	"fmt"
	"social/ent/group"
	"strings"
	"time"

	"entgo.io/ent/dialect/sql"
)

// Group is the model entity for the Group schema.
type Group struct {
	config `json:"-"`
	// ID of the ent.
	ID int64 `json:"id,omitempty"`
	// Name holds the value of the "name" field.
	Name string `json:"name,omitempty"`
	// Owner holds the value of the "owner" field.
	Owner int64 `json:"owner,omitempty"`
	// Avatar holds the value of the "avatar" field.
	Avatar string `json:"avatar,omitempty"`
	// CreatedUID holds the value of the "created_uid" field.
	CreatedUID int64 `json:"created_uid,omitempty"`
	// MaxMembers holds the value of the "max_members" field.
	MaxMembers int `json:"max_members,omitempty"`
	// Mode holds the value of the "mode" field.
	Mode int8 `json:"mode,omitempty"`
	// Type holds the value of the "type" field.
	Type int8 `json:"type,omitempty"`
	// Status holds the value of the "status" field.
	Status int8 `json:"status,omitempty"`
	// InviteMode holds the value of the "invite_mode" field.
	InviteMode int8 `json:"invite_mode,omitempty"`
	// Notice holds the value of the "notice" field.
	Notice string `json:"notice,omitempty"`
	// Introduction holds the value of the "introduction" field.
	Introduction string `json:"introduction,omitempty"`
	// CreatedAt holds the value of the "created_at" field.
	CreatedAt time.Time `json:"-"`
	// UpdatedAt holds the value of the "updated_at" field.
	UpdatedAt time.Time `json:"-"`
	// DeletedAt holds the value of the "deleted_at" field.
	DeletedAt *time.Time `json:"-"`
}

// scanValues returns the types for scanning values from sql.Rows.
func (*Group) scanValues(columns []string) ([]interface{}, error) {
	values := make([]interface{}, len(columns))
	for i := range columns {
		switch columns[i] {
		case group.FieldID, group.FieldOwner, group.FieldCreatedUID, group.FieldMaxMembers, group.FieldMode, group.FieldType, group.FieldStatus, group.FieldInviteMode:
			values[i] = new(sql.NullInt64)
		case group.FieldName, group.FieldAvatar, group.FieldNotice, group.FieldIntroduction:
			values[i] = new(sql.NullString)
		case group.FieldCreatedAt, group.FieldUpdatedAt, group.FieldDeletedAt:
			values[i] = new(sql.NullTime)
		default:
			return nil, fmt.Errorf("unexpected column %q for type Group", columns[i])
		}
	}
	return values, nil
}

// assignValues assigns the values that were returned from sql.Rows (after scanning)
// to the Group fields.
func (gr *Group) assignValues(columns []string, values []interface{}) error {
	if m, n := len(values), len(columns); m < n {
		return fmt.Errorf("mismatch number of scan values: %d != %d", m, n)
	}
	for i := range columns {
		switch columns[i] {
		case group.FieldID:
			value, ok := values[i].(*sql.NullInt64)
			if !ok {
				return fmt.Errorf("unexpected type %T for field id", value)
			}
			gr.ID = int64(value.Int64)
		case group.FieldName:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field name", values[i])
			} else if value.Valid {
				gr.Name = value.String
			}
		case group.FieldOwner:
			if value, ok := values[i].(*sql.NullInt64); !ok {
				return fmt.Errorf("unexpected type %T for field owner", values[i])
			} else if value.Valid {
				gr.Owner = value.Int64
			}
		case group.FieldAvatar:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field avatar", values[i])
			} else if value.Valid {
				gr.Avatar = value.String
			}
		case group.FieldCreatedUID:
			if value, ok := values[i].(*sql.NullInt64); !ok {
				return fmt.Errorf("unexpected type %T for field created_uid", values[i])
			} else if value.Valid {
				gr.CreatedUID = value.Int64
			}
		case group.FieldMaxMembers:
			if value, ok := values[i].(*sql.NullInt64); !ok {
				return fmt.Errorf("unexpected type %T for field max_members", values[i])
			} else if value.Valid {
				gr.MaxMembers = int(value.Int64)
			}
		case group.FieldMode:
			if value, ok := values[i].(*sql.NullInt64); !ok {
				return fmt.Errorf("unexpected type %T for field mode", values[i])
			} else if value.Valid {
				gr.Mode = int8(value.Int64)
			}
		case group.FieldType:
			if value, ok := values[i].(*sql.NullInt64); !ok {
				return fmt.Errorf("unexpected type %T for field type", values[i])
			} else if value.Valid {
				gr.Type = int8(value.Int64)
			}
		case group.FieldStatus:
			if value, ok := values[i].(*sql.NullInt64); !ok {
				return fmt.Errorf("unexpected type %T for field status", values[i])
			} else if value.Valid {
				gr.Status = int8(value.Int64)
			}
		case group.FieldInviteMode:
			if value, ok := values[i].(*sql.NullInt64); !ok {
				return fmt.Errorf("unexpected type %T for field invite_mode", values[i])
			} else if value.Valid {
				gr.InviteMode = int8(value.Int64)
			}
		case group.FieldNotice:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field notice", values[i])
			} else if value.Valid {
				gr.Notice = value.String
			}
		case group.FieldIntroduction:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field introduction", values[i])
			} else if value.Valid {
				gr.Introduction = value.String
			}
		case group.FieldCreatedAt:
			if value, ok := values[i].(*sql.NullTime); !ok {
				return fmt.Errorf("unexpected type %T for field created_at", values[i])
			} else if value.Valid {
				gr.CreatedAt = value.Time
			}
		case group.FieldUpdatedAt:
			if value, ok := values[i].(*sql.NullTime); !ok {
				return fmt.Errorf("unexpected type %T for field updated_at", values[i])
			} else if value.Valid {
				gr.UpdatedAt = value.Time
			}
		case group.FieldDeletedAt:
			if value, ok := values[i].(*sql.NullTime); !ok {
				return fmt.Errorf("unexpected type %T for field deleted_at", values[i])
			} else if value.Valid {
				gr.DeletedAt = new(time.Time)
				*gr.DeletedAt = value.Time
			}
		}
	}
	return nil
}

// Update returns a builder for updating this Group.
// Note that you need to call Group.Unwrap() before calling this method if this Group
// was returned from a transaction, and the transaction was committed or rolled back.
func (gr *Group) Update() *GroupUpdateOne {
	return (&GroupClient{config: gr.config}).UpdateOne(gr)
}

// Unwrap unwraps the Group entity that was returned from a transaction after it was closed,
// so that all future queries will be executed through the driver which created the transaction.
func (gr *Group) Unwrap() *Group {
	tx, ok := gr.config.driver.(*txDriver)
	if !ok {
		panic("ent: Group is not a transactional entity")
	}
	gr.config.driver = tx.drv
	return gr
}

// String implements the fmt.Stringer.
func (gr *Group) String() string {
	var builder strings.Builder
	builder.WriteString("Group(")
	builder.WriteString(fmt.Sprintf("id=%v", gr.ID))
	builder.WriteString(", name=")
	builder.WriteString(gr.Name)
	builder.WriteString(", owner=")
	builder.WriteString(fmt.Sprintf("%v", gr.Owner))
	builder.WriteString(", avatar=")
	builder.WriteString(gr.Avatar)
	builder.WriteString(", created_uid=")
	builder.WriteString(fmt.Sprintf("%v", gr.CreatedUID))
	builder.WriteString(", max_members=")
	builder.WriteString(fmt.Sprintf("%v", gr.MaxMembers))
	builder.WriteString(", mode=")
	builder.WriteString(fmt.Sprintf("%v", gr.Mode))
	builder.WriteString(", type=")
	builder.WriteString(fmt.Sprintf("%v", gr.Type))
	builder.WriteString(", status=")
	builder.WriteString(fmt.Sprintf("%v", gr.Status))
	builder.WriteString(", invite_mode=")
	builder.WriteString(fmt.Sprintf("%v", gr.InviteMode))
	builder.WriteString(", notice=")
	builder.WriteString(gr.Notice)
	builder.WriteString(", introduction=")
	builder.WriteString(gr.Introduction)
	builder.WriteString(", created_at=")
	builder.WriteString(gr.CreatedAt.Format(time.ANSIC))
	builder.WriteString(", updated_at=")
	builder.WriteString(gr.UpdatedAt.Format(time.ANSIC))
	if v := gr.DeletedAt; v != nil {
		builder.WriteString(", deleted_at=")
		builder.WriteString(v.Format(time.ANSIC))
	}
	builder.WriteByte(')')
	return builder.String()
}

// Groups is a parsable slice of Group.
type Groups []*Group

func (gr Groups) config(cfg config) {
	for _i := range gr {
		gr[_i].config = cfg
	}
}

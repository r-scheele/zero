package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// NoteLike holds the schema definition for the NoteLike entity.
type NoteLike struct {
	ent.Schema
}

// Fields of the NoteLike.
func (NoteLike) Fields() []ent.Field {
	return []ent.Field{
		field.Time("created_at").
			Default(time.Now).
			Immutable(),
	}
}

// Edges of the NoteLike.
func (NoteLike) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("user", User.Type).
			Ref("note_likes").
			Unique().
			Required(),
		edge.From("note", Note.Type).
			Ref("likes").
			Unique().
			Required(),
	}
}

// Indexes of the NoteLike.
func (NoteLike) Indexes() []ent.Index {
	return []ent.Index{
		// Index for efficient querying
		index.Fields("created_at"),
	}
}
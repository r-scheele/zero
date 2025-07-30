package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// NoteRepost holds the schema definition for the NoteRepost entity.
type NoteRepost struct {
	ent.Schema
}

// Fields of the NoteRepost.
func (NoteRepost) Fields() []ent.Field {
	return []ent.Field{
		field.String("comment").
			Optional().
			Comment("Optional comment when reposting"),
		field.Time("created_at").
			Default(time.Now).
			Immutable(),
	}
}

// Edges of the NoteRepost.
func (NoteRepost) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("user", User.Type).
			Ref("note_reposts").
			Unique().
			Required(),
		edge.From("note", Note.Type).
			Ref("reposts").
			Unique().
			Required(),
	}
}

// Indexes of the NoteRepost.
func (NoteRepost) Indexes() []ent.Index {
	return []ent.Index{
		// Index for efficient querying
		index.Fields("created_at"),
	}
}
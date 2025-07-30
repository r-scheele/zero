package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"time"
	"github.com/r-scheele/zero/pkg/types"
)

// Note holds the schema definition for the Note entity.
type Note struct {
	ent.Schema
}

// Fields of the Note.
func (Note) Fields() []ent.Field {
	return []ent.Field{
		field.String("title").
			NotEmpty().
			Comment("Title of the note"),
		field.Text("description").
			Optional().
			Comment("Brief description of the note"),
		field.Text("content").
			Optional().
			Comment("Main text content of the note"),
		field.JSON("resources", []types.Resource{}).
			Optional().
			Comment("Array of attached resources (files, links, etc.)"),
		field.Text("ai_curriculum").
			Optional().
			Comment("AI-generated curriculum based on note content"),
		field.Enum("visibility").
			Values("private", "public").
			Default("private").
			Comment("Note visibility setting"),
		field.Enum("permission_level").
			Values("read_only", "read_write", "read_write_approval").
			Default("read_only").
			Comment("Permission level for public notes"),
		field.String("share_token").
			Optional().
			Unique().
			Comment("Unique token for sharing the note via link"),
		field.Bool("ai_processing").
			Default(false).
			Comment("Whether AI is currently processing this note"),
		field.Time("created_at").
			Default(time.Now).
			Immutable(),
		field.Time("updated_at").
			Default(time.Now).
			UpdateDefault(time.Now),
	}
}

// Edges of the Note.
func (Note) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("owner", User.Type).
			Ref("notes").
			Unique().
			Required(),
		edge.To("likes", NoteLike.Type),
		edge.To("reposts", NoteRepost.Type),
	}
}

// Indexes of the Note.
func (Note) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("visibility", "created_at"),
		index.Fields("share_token").Unique(),
	}
}
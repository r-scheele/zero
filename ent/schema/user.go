package schema

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"time"

	ge "github.com/r-scheele/zero/ent"
	"github.com/r-scheele/zero/ent/hook"
	"golang.org/x/crypto/bcrypt"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// User holds the schema definition for the User entity.
type User struct {
	ent.Schema
}

// Fields of the User.
func (User) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").
			NotEmpty(),
		field.String("phone_number").
			NotEmpty().
			Unique().
			Validate(func(s string) error {
				// Basic phone number validation (E.164 format)
				matched, err := regexp.MatchString(`^\+[1-9]\d{1,14}$`, s)
				if err != nil {
					return err
				}
				if !matched {
					return fmt.Errorf("invalid phone number format, expected E.164 format (e.g., +1234567890)")
				}
				return nil
			}),
		field.String("email").
			Optional().
			Nillable().
			Validate(func(s string) error {
				if s == "" {
					return nil // Allow empty email
				}
				// Basic email validation
				matched, err := regexp.MatchString(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`, s)
				if err != nil {
					return err
				}
				if !matched {
					return fmt.Errorf("invalid email format")
				}
				return nil
			}),
		field.String("password").
			Sensitive().
			NotEmpty().
			Optional(), // Password is optional for WhatsApp-only registrations
		field.Bool("verified").
			Default(false),
		field.String("verification_code").
			Optional().
			Nillable().
			Comment("2-digit verification code shown on web for WhatsApp verification"),
		field.Bool("admin").
			Default(false),
		field.Enum("registration_method").
			Values("whatsapp", "web").
			Default("web"),
		field.String("profile_picture").
			Optional().
			Nillable().
			Comment("Path to profile picture file"),
		field.Bool("dark_mode").
			Default(false).
			Comment("User's preferred theme"),
		field.String("bio").
			Optional().
			Nillable().
			MaxLen(500).
			Comment("User biography/description"),
		field.Bool("email_notifications").
			Default(true).
			Comment("Whether user wants to receive email notifications"),
		field.Bool("sms_notifications").
			Default(true).
			Comment("Whether user wants to receive SMS notifications"),
		field.Bool("is_active").
			Default(true).
			Comment("Whether the account is active or deactivated"),
		field.Time("last_login").
			Optional().
			Nillable().
			Comment("Last login timestamp"),
		field.Time("created_at").
			Default(time.Now).
			Immutable(),
		field.Time("updated_at").
			Default(time.Now).
			UpdateDefault(time.Now).
			Optional().
			Nillable(),
	}
}

// Edges of the User.
func (User) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("owner", PasswordToken.Type).
			Ref("user"),
	}
}

// Hooks of the User.
func (User) Hooks() []ent.Hook {
	return []ent.Hook{
		hook.On(
			func(next ent.Mutator) ent.Mutator {
				return hook.UserFunc(func(ctx context.Context, m *ge.UserMutation) (ent.Value, error) {
					// Normalize phone number format
					if v, exists := m.PhoneNumber(); exists {
						// Ensure phone number starts with + and contains only digits
						normalized := strings.TrimSpace(v)
						if !strings.HasPrefix(normalized, "+") {
							normalized = "+" + normalized
						}
						m.SetPhoneNumber(normalized)
					}

					// Hash password if it exists and isn't empty (optional for WhatsApp-only users)
					if v, exists := m.Password(); exists && v != "" {
						hash, err := bcrypt.GenerateFromPassword([]byte(v), bcrypt.DefaultCost)
						if err != nil {
							return "", err
						}
						m.SetPassword(string(hash))
					}
					return next.Mutate(ctx, m)
				})
			},
			// Limit the hook only for these operations.
			ent.OpCreate|ent.OpUpdate|ent.OpUpdateOne,
		),
	}
}

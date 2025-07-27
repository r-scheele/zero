package components

import (
	"github.com/r-scheele/zero/pkg/form"
	"github.com/r-scheele/zero/pkg/ui"
	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

type (
	InputFieldParams struct {
		Form        form.Form
		FormField   string
		Name        string
		InputType   string
		Label       string
		Value       string
		Placeholder string
		Help        string
	}

	FileFieldParams struct {
		Name  string
		Label string
		Help  string
	}

	OptionsParams struct {
		Form      form.Form
		FormField string
		Name      string
		Label     string
		Value     string
		Options   []Choice
		Help      string
	}

	Choice struct {
		Value string
		Label string
	}

	TextareaFieldParams struct {
		Form      form.Form
		FormField string
		Name      string
		Label     string
		Value     string
		Help      string
	}

	CheckboxParams struct {
		Form      form.Form
		FormField string
		Name      string
		Label     string
		Checked   bool
	}
)

func ControlGroup(controls ...Node) Node {
	return Div(
		Class("mt-2 flex gap-2"),
		Group(controls),
	)
}

func TextareaField(el TextareaFieldParams) Node {
	return Fieldset(
		el.Label,
		Textarea(
			Class("textarea w-full h-32 text-base px-4 py-3 rounded-xl border-2 border-slate-200 focus:border-blue-500 focus:ring-4 focus:ring-blue-100 transition-all duration-200 "+formFieldStatusClass(el.Form, el.FormField)),
			ID(el.Name),
			Name(el.Name),
			Text(el.Value),
			Style("font-size: 16px; line-height: 1.5;"), // Better readability
		),
		Help(el.Help),
		formFieldErrors(el.Form, el.FormField),
	)
}

func Radios(el OptionsParams) Node {
	buttons := make(Group, len(el.Options))
	for i, opt := range el.Options {
		id := "radio-" + el.Name + "-" + opt.Value
		buttons[i] = Div(
			Class("mb-2"),
			Input(
				ID(id),
				Type("radio"),
				Name(el.Name),
				Value(opt.Value),
				Class("radio mr-1 "+formFieldStatusClass(el.Form, el.FormField)),
				If(el.Value == opt.Value, Checked()),
			),
			Label(
				Text(opt.Label),
				For(id),
			),
		)
	}

	return Fieldset(
		el.Label,
		buttons,
		formFieldErrors(el.Form, el.FormField),
	)
}

func SelectList(el OptionsParams) Node {
	buttons := make(Group, len(el.Options))
	for i, opt := range el.Options {
		buttons[i] = Option(
			Text(opt.Label),
			Value(opt.Value),
			If(opt.Value == el.Value, Attr("selected")),
		)
	}

	return Fieldset(
		el.Label,
		Select(
			Class("select w-full text-base px-4 py-3 rounded-xl border-2 border-slate-200 focus:border-blue-500 focus:ring-4 focus:ring-blue-100 transition-all duration-200 "+formFieldStatusClass(el.Form, el.FormField)),
			Name(el.Name),
			Style("min-height: 48px; font-size: 16px;"), // Better sizing
			buttons,
		),
		Help(el.Help),
		formFieldErrors(el.Form, el.FormField),
	)
}

func Checkbox(el CheckboxParams) Node {
	return Div(
		Class("mb-4"), // Better spacing
		Label(
			Class("label flex items-center gap-3 cursor-pointer"), // Better layout and interaction
			Input(
				Class("checkbox w-5 h-5 text-blue-600 border-2 border-slate-300 rounded focus:ring-4 focus:ring-blue-100"),
				Type("checkbox"),
				Name(el.Name),
				If(el.Checked, Checked()),
				Value("true"),
			),
			Span(
				Class("text-base text-slate-700 select-none"), // Better text styling
				Text(el.Label),
			),
		),
		formFieldErrors(el.Form, el.FormField),
	)
}

func InputField(el InputFieldParams) Node {
	return Fieldset(
		el.Label,
		Input(
			ID(el.Name),
			Name(el.Name),
			Type(el.InputType),
			Class("input w-full text-base px-4 py-3 rounded-xl border-2 border-slate-200 focus:border-blue-500 focus:ring-4 focus:ring-blue-100 transition-all duration-200 "+formFieldStatusClass(el.Form, el.FormField)),
			Value(el.Value),
			If(el.Placeholder != "", Placeholder(el.Placeholder)),
			Style("min-height: 48px; font-size: 16px;"), // Better touch targets and mobile optimization
		),
		Help(el.Help),
		formFieldErrors(el.Form, el.FormField),
	)
}

func Help(text string) Node {
	return If(len(text) > 0, Div(
		Class("label text-sm text-slate-600 mt-2 block"),
		Style("font-size: 14px; line-height: 1.4;"), // Consistent helper text sizing
		Text(text),
	))
}

func Fieldset(label string, els ...Node) Node {
	return FieldSet(
		Class("fieldset mb-6"), // Better spacing
		If(len(label) > 0, Legend(
			Class("fieldset-legend text-sm font-semibold text-slate-700 mb-2 block"),
			Text(label),
		)),
		Group(els),
	)
}

func FileField(el FileFieldParams) Node {
	return Fieldset(
		el.Label,
		Input(
			Type("file"),
			Class("file-input"),
			Name(el.Name),
		),
		Help(el.Help),
	)
}

func formFieldStatusClass(fm form.Form, formField string) string {
	switch {
	case fm == nil:
		return ""
	case !fm.IsSubmitted():
		return ""
	case fm.FieldHasErrors(formField):
		return "input-error"
	default:
		return "input-success"
	}
}

func formFieldErrors(fm form.Form, field string) Node {
	if fm == nil {
		return nil
	}

	errs := fm.GetFieldErrors(field)
	if len(errs) == 0 {
		return nil
	}

	g := make(Group, len(errs))
	for i, err := range errs {
		g[i] = Div(
			Class("text-error"),
			Text(err),
		)
	}

	return g
}

func CSRF(r *ui.Request) Node {
	return Input(
		Type("hidden"),
		Name("csrf"),
		Value(r.CSRF),
	)
}

func FormButton(color Color, label string) Node {
	return Button(
		Class("btn "+buttonColor(color)+" text-base font-semibold px-6 py-3 min-h-12 transition-all duration-200 hover:scale-105"),
		Type("submit"),
		Style("min-height: 48px; font-size: 16px;"), // Better touch targets
		Text(label),
	)
}

func ButtonLink(color Color, href, label string) Node {
	return A(
		Href(href),
		Class("btn "+buttonColor(color)+" text-base font-semibold px-6 py-3 min-h-12 transition-all duration-200 hover:scale-105 inline-flex items-center justify-center text-center"),
		Style("min-height: 48px; font-size: 16px; text-decoration: none;"), // Better touch targets and styling
		Text(label),
	)
}

func buttonColor(color Color) string {
	// Only colors being used are included so unused styles are not compiled.
	switch color {
	case ColorPrimary:
		return "btn-primary"
	case ColorInfo:
		return "btn-info"
	case ColorAccent:
		return "btn-accent"
	case ColorError:
		return "btn-error"
	case ColorLink:
		return "btn-link"
	default:
		return ""
	}
}

package view

import "html/template"

const resetPasswordView = `
reset password
`

func NewAuthn(t *template.Template) {
	t = template.Must(t.New("reset_password").Parse(resetPasswordView))
}

package authn

import "html/template"

const resetPasswordView = `
reset password
`

func NewView(t *template.Template) {
	t = template.Must(t.New("reset_password").Parse(resetPasswordView))
}

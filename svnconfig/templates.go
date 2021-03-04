package svnconfig

const rawTmplAuthzSVNAccessFile = `
[groups]
{{ range $gi, $g := .Groups -}}
{{- $g.Name }} = {{ range $ui, $u := $g.Users -}}
{{- if gt $ui 0 -}}, {{ end -}}
{{- $u -}}
{{- end -}}{{/* $g.Users */}}
{{ end -}}{{/* .Groups */}}
`

const rawTmplAuthUserFile = ``

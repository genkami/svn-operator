package svnconfig

const rawTmplAuthzSVNAccessFile = `
[groups]
{{ range $gi, $g := .Groups -}}
{{- $g.Name }} = {{ range $ui, $u := $g.Users -}}
{{- if gt $ui 0 -}}, {{ end -}}
{{- $u -}}
{{- end -}}{{/* $g.Users */}}
{{ end -}}{{/* .Groups */}}
{{- range $ri, $r := .Repositories -}}
[{{- $r.Name -}}:/]
* = 
{{ range $pi, $p := $r.Permissions -}}
@{{- $p.Group }} = {{ $p.Permission }}
{{ end -}}
{{- end -}}{{/* .Repositories */}}
`

const rawTmplAuthUserFile = `
{{ range $ui, $u := .Users -}}
{{- $u.Name}}:{{- $u.EncryptedPassword }}
{{ end -}}{{/* .Users */}}
`

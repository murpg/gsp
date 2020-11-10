package tpl

const ReleaseNotesTextTemplate = `Release - {{ .Format .ReleaseDate }}
=======================================

Changed files:
{{if not .ChangedFiles -}}
- no files found
{{else -}}
{{range .ChangedFiles -}}
- {{.}}
{{end}}
{{- end}}`

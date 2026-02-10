{{ if .Versions -}}
{{ range .Versions }}
## {{ if .Tag.Name }}[{{ .Tag.Name }}]{{ else }}Unreleased{{ end }} - {{ .Tag.Date.Format "2006-01-02" }}

{{ range .CommitGroups -}}
### {{ .Title }}
{{ range .Commits -}}
- {{ if .Scope }}**{{ .Scope }}:** {{ end }}{{ .Subject }}{{ if .PR }} (#{{ .PR }}){{ end }}{{ if .Issues }} ({{ range .Issues }}{{ .Text }} {{ end }}){{ end }}
{{ end }}
{{ end -}}

{{- if .RevertGroups -}}
### ⏪ Reverts
{{ range .RevertGroups -}}
- {{ .Revert.Subject }}
{{ end }}
{{ end -}}

{{- end -}}
{{- end -}}

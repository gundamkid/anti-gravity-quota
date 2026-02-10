{{ range .Versions -}}
## {{ if .Tag.Name }}{{ .Tag.Name }}{{ else }}Unreleased{{ end }} ({{ .Tag.Date.Format "2006-01-02" }})

{{ range .CommitGroups -}}
### {{ .Title }}

{{ range .Commits -}}
- {{ if .Scope }}**{{ .Scope }}:** {{ end }}{{ .Subject }} {{ if .Hash }}({{ .Hash.Short }}){{ end }}
{{ end }}
{{ end -}}

{{ end -}}

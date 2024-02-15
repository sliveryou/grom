syntax = "v1"

import (
{{ range .Imports }}    "{{ . }}"
{{ end -}}
)

{{ if .APIInfo }}info (
	{{ .APIInfo }}
){{ end }}

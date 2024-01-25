syntax = "v1"

import (
{{ range .Imports }}    "{{ . }}"
{{ end -}}
)

info (
    title:   "{{ .Title }}"
    desc:    "{{ .Desc }}"
    author:  "{{ .Author }}"
    email:   "{{ .Email }}"
    version: "{{ .Version }}"
)
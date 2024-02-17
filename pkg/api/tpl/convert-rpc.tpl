// To{{ .StructName }} {{ .TableComment }}转化
func To{{ .StructName }}(src *model.{{ .ModelName }}) {{ if .HasErr }}(*pb.{{ .StructName }}, error){{ else }}*pb.{{ .StructName }}{{ end }} {
	var dst pb.{{ .StructName }}
	if src != nil {
		dst = pb.{{ .StructName }}{
			{{ .ConvertInfo }}
		}
		{{ .IfInfo -}}
	}

	return &dst{{ if .HasErr }}, nil{{ end }}
}

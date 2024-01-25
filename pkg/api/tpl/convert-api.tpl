// To{{ .StructName }} {{ .TableComment }}转化
func To{{ .StructName }}(src *pb.{{ .StructName }}) {{ .StructName }} {
	var dst {{ .StructName }}
	if src != nil {
		dst = {{ .StructName }}{
			{{ .ConvertInfo }}
		}
	}

	return dst
}

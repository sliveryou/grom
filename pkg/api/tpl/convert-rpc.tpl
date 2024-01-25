// To{{ .StructName }} {{ .TableComment }}转化
func To{{ .StructName }}(src *model.{{ .ModelName }}) *pb.{{ .StructName }} {
	var dst pb.{{ .StructName }}
	if src != nil {
		dst = pb.{{ .StructName }}{
			{{ .ConvertInfo }}
		}
	}

	return &dst
}

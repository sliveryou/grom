syntax = "v1"

{{ if .APIInfo }}info (
	{{ .APIInfo }}
){{ end }}

@server (
    {{ if .RoutePrefix }}prefix: /{{ .RoutePrefix }}{{ end }}
    group: {{ if .GroupPrefix }}{{ .GroupPrefix }}/{{ end }}{{ .GroupName }}
)
service {{ .ServiceName }} {
    @doc (
        summary: "获取{{ .TableComment }}"
    )
    @handler Get{{ .StructName }}
    get /{{ .RouteName }}/:{{ .IdRawName }} (Get{{ .StructName }}{{ .ReqName }}) returns (Get{{ .StructName }}{{ .RespName }})

    @doc (
        summary: "列出{{ .TableComment }}"
    )
    @handler List{{ .StructNamePlural }}
    get /{{ .RouteName }} (List{{ .StructNamePlural }}{{ .ReqName }}) returns (List{{ .StructNamePlural }}{{ .RespName }})

    @doc (
        summary: "创建{{ .TableComment }}"
    )
    @handler Create{{ .StructName }}
    post /{{ .RouteName }} (Create{{ .StructName }}{{ .ReqName }}) returns (Create{{ .StructName }}{{ .RespName }})

    @doc (
        summary: "更新{{ .TableComment }}"
    )
    @handler Update{{ .StructName }}
    put /{{ .RouteName }}/:{{ .IdRawName }} (Update{{ .StructName }}{{ .ReqName }}) returns (Update{{ .StructName }}{{ .RespName }})

    @doc (
        summary: "删除{{ .TableComment }}"
    )
    @handler Delete{{ .StructName }}
    delete /{{ .RouteName }}/:{{ .IdRawName }} (Delete{{ .StructName }}{{ .ReqName }}) returns (Delete{{ .StructName }}{{ .RespName }})

    @doc (
        summary: "修补{{ .TableComment }}"
    )
    @handler Patch{{ .StructName }}
    patch /{{ .RouteName }}/:{{ .IdRawName }} (Patch{{ .StructName }}{{ .ReqName }}) returns (Patch{{ .StructName }}{{ .RespName }})

    @doc (
        summary: "批量获取{{ .TableComment }}"
    )
    @handler BatchGet{{ .StructNamePlural }}
    post /{{ .RouteName }}/batch{{ .Delimiter }}get (BatchGet{{ .StructNamePlural }}{{ .ReqName }}) returns (BatchGet{{ .StructNamePlural }}{{ .RespName }})

    @doc (
        summary: "批量创建{{ .TableComment }}"
    )
    @handler BatchCreate{{ .StructNamePlural }}
    post /{{ .RouteName }}/batch{{ .Delimiter }}create (BatchCreate{{ .StructNamePlural }}{{ .ReqName }}) returns (BatchCreate{{ .StructNamePlural }}{{ .RespName }})

    @doc (
        summary: "批量更新{{ .TableComment }}"
    )
    @handler BatchUpdate{{ .StructNamePlural }}
    post /{{ .RouteName }}/batch{{ .Delimiter }}update (BatchUpdate{{ .StructNamePlural }}{{ .ReqName }}) returns (BatchUpdate{{ .StructNamePlural }}{{ .RespName }})

    @doc (
        summary: "批量删除{{ .TableComment }}"
    )
    @handler BatchDelete{{ .StructNamePlural }}
    post /{{ .RouteName }}/batch{{ .Delimiter }}delete (BatchDelete{{ .StructNamePlural }}{{ .ReqName }}) returns (BatchDelete{{ .StructNamePlural }}{{ .RespName }})
}

// -------------------- {{ .TableComment }} {{ .StructName }} -------------------- //
// {{ .StructName }} {{ .TableComment }}
type {{ .StructName }} {
    {{ .StructInfo }}
}

// Get{{ .StructName }}{{ .ReqName }} 获取{{ .TableComment }}请求
type Get{{ .StructName }}{{ .ReqName }} {
    {{ .IdName }} {{ .IdType }} `path:"{{ .IdRawName }}" validate:"required" label:"{{ .IdLabel }}"` // {{ .IdComment }}
}

// Get{{ .StructName }}{{ .RespName }} 获取{{ .TableComment }}响应
type Get{{ .StructName }}{{ .RespName }} {
    {{ .StructName }}
}

// List{{ .StructNamePlural }}{{ .ReqName }} 列出{{ .TableComment }}请求
type List{{ .StructNamePlural }}{{ .ReqName }} {
    {{ .StructGetInfo }}
    Page     int64 `form:"page" validate:"required" label:"页数"`        // 页数
    PageSize int64 `form:"page_size" validate:"required" label:"每条页数"` // 每条页数
}

// List{{ .StructNamePlural }}{{ .RespName }} 列出{{ .TableComment }}响应
type List{{ .StructNamePlural }}{{ .RespName }} {
    Count     int64                `json:"count"`      // 总数
    PageCount int64                `json:"page_count"` // 页数
    Results   []*{{ .StructName }} `json:"results"`    // 结果
}

// Create{{ .StructName }}{{ .ReqName }} 创建{{ .TableComment }}请求
type Create{{ .StructName }}{{ .ReqName }} {
    {{ .StructCreateInfo }}
}

// Create{{ .StructName }}{{ .RespName }} 创建{{ .TableComment }}响应
type Create{{ .StructName }}{{ .RespName }} {
    {{ .StructName }}
}

// Update{{ .StructName }}{{ .ReqName }} 更新{{ .TableComment }}请求
type Update{{ .StructName }}{{ .ReqName }} {
    {{ .StructUpdateInfo }}
}

// Update{{ .StructName }}{{ .RespName }} 更新{{ .TableComment }}响应
type Update{{ .StructName }}{{ .RespName }} {
    {{ .StructName }}
}

// Delete{{ .StructName }}{{ .ReqName }} 删除{{ .TableComment }}请求
type Delete{{ .StructName }}{{ .ReqName }} {
    {{ .IdName }} {{ .IdType }} `path:"{{ .IdRawName }}" validate:"required" label:"{{ .IdLabel }}"` // {{ .IdComment }}
}

// Delete{{ .StructName }}{{ .RespName }} 删除{{ .TableComment }}响应
type Delete{{ .StructName }}{{ .RespName }} {
    {{ .IdName }} {{ .IdType }} `json:"{{ .IdRawName }}"` // {{ .IdComment }}
}

// Patch{{ .StructName }}{{ .ReqName }} 修补{{ .TableComment }}请求
type Patch{{ .StructName }}{{ .ReqName }} {
    {{ .StructUpdateInfo }}
}

// Patch{{ .StructName }}{{ .RespName }} 修补{{ .TableComment }}响应
type Patch{{ .StructName }}{{ .RespName }} {
    {{ .StructName }}
}

// {{ .StructName }}Filter {{ .TableComment }}过滤参数
type {{ .StructName }}Filter {
    {{ .IdNamePlural }} []{{ .IdType }} `json:"{{ .IdRawNamePlural }},optional"` // {{ .IdComment }}列表
    {{ .StructFilterInfo }}
}

// BatchGet{{ .StructNamePlural }}{{ .ReqName }} 批量获取{{ .TableComment }}请求
type BatchGet{{ .StructNamePlural }}{{ .ReqName }} {
    Filter {{ .StructName }}Filter `json:"filter"` // {{ .TableComment }}过滤参数
}

// BatchGet{{ .StructNamePlural }}{{ .RespName }} 批量获取{{ .TableComment }}响应
type BatchGet{{ .StructNamePlural }}{{ .RespName }} {
    Results []*{{ .StructName }} `json:"results"` // 结果
}

// BatchCreate{{ .StructNamePlural }}{{ .ReqName }} 批量创建{{ .TableComment }}请求
type BatchCreate{{ .StructNamePlural }}{{ .ReqName }} {
    Objects []*Create{{ .StructName }}{{ .ReqName }} `json:"objects" validate:"gt=0,dive" label:"{{ .TableComment }}列表"` // {{ .TableComment }}列表
}

// BatchCreate{{ .StructNamePlural }}{{ .RespName }} 批量创建{{ .TableComment }}响应
type BatchCreate{{ .StructNamePlural }}{{ .RespName }} {
    Results []*{{ .StructName }} `json:"results"` // 结果
}

// BatchUpdate{{ .StructNamePlural }}{{ .ReqName }} 批量更新{{ .TableComment }}请求
type BatchUpdate{{ .StructNamePlural }}{{ .ReqName }} {
    Filter {{ .StructName }}Filter `json:"filter"` // {{ .TableComment }}过滤参数
    {{ .StructBatchUpdateInfo }}
}

// BatchUpdate{{ .StructNamePlural }}{{ .RespName }} 批量更新{{ .TableComment }}响应
type BatchUpdate{{ .StructNamePlural }}{{ .RespName }} {
    Affected int64 `json:"affected"` // 影响数量
}

// BatchDelete{{ .StructNamePlural }}{{ .ReqName }} 批量删除{{ .TableComment }}请求
type BatchDelete{{ .StructNamePlural }}{{ .ReqName }} {
    Filter {{ .StructName }}Filter `json:"filter"` // {{ .TableComment }}过滤参数
}

// BatchDelete{{ .StructNamePlural }}{{ .RespName }} 批量删除{{ .TableComment }}响应
type BatchDelete{{ .StructNamePlural }}{{ .RespName }} {
    Affected int64 `json:"affected"` // 影响数量
}

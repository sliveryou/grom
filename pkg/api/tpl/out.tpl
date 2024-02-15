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
        summary: "查询{{ .TableComment }}"
    )
    @handler Get{{ .StructName }}
    get /{{ .RouteName }}/:{{ .IdRawName }} (Get{{ .StructName }}Req) returns (Get{{ .StructName }}Resp)

    @doc (
        summary: "获取{{ .TableComment }}分页"
    )
    @handler Paginate{{ .StructName }}
    get /{{ .RouteName }} (Paginate{{ .StructName }}Req) returns (Paginate{{ .StructName }}Resp)

    @doc (
        summary: "创建{{ .TableComment }}"
    )
    @handler Create{{ .StructName }}
    post /{{ .RouteName }} (Create{{ .StructName }}Req) returns (Create{{ .StructName }}Resp)

    @doc (
        summary: "更新{{ .TableComment }}"
    )
    @handler Update{{ .StructName }}
    put /{{ .RouteName }}/:{{ .IdRawName }} (Update{{ .StructName }}Req) returns (Update{{ .StructName }}Resp)

    @doc (
        summary: "删除{{ .TableComment }}"
    )
    @handler Delete{{ .StructName }}
    delete /{{ .RouteName }}/:{{ .IdRawName }} (Delete{{ .StructName }}Req) returns (Delete{{ .StructName }}Resp)

    @doc (
        summary: "部分更新{{ .TableComment }}"
    )
    @handler Patch{{ .StructName }}
    patch /{{ .RouteName }}/:{{ .IdRawName }} (Patch{{ .StructName }}Req) returns (Patch{{ .StructName }}Resp)

    @doc (
        summary: "获取{{ .TableComment }}列表"
    )
    @handler List{{ .StructName }}
    post /{{ .RouteName }}/list (List{{ .StructName }}Req) returns (List{{ .StructName }}Resp)

    @doc (
        summary: "批量创建{{ .TableComment }}"
    )
    @handler Create{{ .StructName }}s
    post /{{ .RouteName }}/batch/create (Create{{ .StructName }}sReq) returns (Create{{ .StructName }}sResp)

    @doc (
        summary: "批量更新{{ .TableComment }}"
    )
    @handler Update{{ .StructName }}s
    post /{{ .RouteName }}/batch/update (Update{{ .StructName }}sReq) returns (Update{{ .StructName }}sResp)

    @doc (
        summary: "批量删除{{ .TableComment }}"
    )
    @handler Delete{{ .StructName }}s
    post /{{ .RouteName }}/batch/delete (Delete{{ .StructName }}sReq) returns (Delete{{ .StructName }}sResp)
}

// -------------------- {{ .TableComment }} {{ .StructName }} -------------------- //
// {{ .StructName }} {{ .TableComment }}
type {{ .StructName }} {
    {{ .StructInfo }}
}

// Get{{ .StructName }}Req 查询{{ .TableComment }}请求
type Get{{ .StructName }}Req {
    {{ .IdName }} {{ .IdType }} `path:"{{ .IdRawName }}" validate:"required" label:"{{ .IdLabel }}"` // {{ .IdComment }}
}

// Get{{ .StructName }}Resp 查询{{ .TableComment }}响应
type Get{{ .StructName }}Resp {
    {{ .StructName }}
}

// Paginate{{ .StructName }}Req 获取{{ .TableComment }}分页请求
type Paginate{{ .StructName }}Req {
    {{ .StructGetInfo }}
    Page     int64 `form:"page" validate:"required" label:"页数"`        // 页数
    PageSize int64 `form:"page_size" validate:"required" label:"每条页数"` // 每条页数
}

// Paginate{{ .StructName }}Resp 获取{{ .TableComment }}分页响应
type Paginate{{ .StructName }}Resp {
    Count     int64                `json:"count"`      // 总数
    PageCount int64                `json:"page_count"` // 页数
    Results   []*{{ .StructName }} `json:"results"`    // 结果
}

// Create{{ .StructName }}Req 创建{{ .TableComment }}请求
type Create{{ .StructName }}Req {
    {{ .StructCreateInfo }}
}

// Create{{ .StructName }}Resp 创建{{ .TableComment }}响应
type Create{{ .StructName }}Resp {
    {{ .StructName }}
}

// Update{{ .StructName }}Req 更新{{ .TableComment }}请求
type Update{{ .StructName }}Req {
    {{ .StructUpdateInfo }}
}

// Update{{ .StructName }}Resp 更新{{ .TableComment }}响应
type Update{{ .StructName }}Resp {
    {{ .StructName }}
}

// Delete{{ .StructName }}Req 删除{{ .TableComment }}请求
type Delete{{ .StructName }}Req {
    {{ .IdName }} {{ .IdType }} `path:"{{ .IdRawName }}" validate:"required" label:"{{ .IdLabel }}"` // {{ .IdComment }}
}

// Delete{{ .StructName }}Resp 删除{{ .TableComment }}响应
type Delete{{ .StructName }}Resp {
    {{ .IdName }} {{ .IdType }} `json:"{{ .IdRawName }}"` // {{ .IdComment }}
}

// Patch{{ .StructName }}Req 部分更新{{ .TableComment }}请求
type Patch{{ .StructName }}Req {
    {{ .StructUpdateInfo }}
}

// Patch{{ .StructName }}Resp 部分更新{{ .TableComment }}响应
type Patch{{ .StructName }}Resp {
    {{ .StructName }}
}

// {{ .StructName }}Filter {{ .TableComment }}筛选参数
type {{ .StructName }}Filter {
    {{ .IdNamePlural }} []{{ .IdType }} `json:"{{ .IdRawNamePlural }},optional"` // {{ .IdComment }}列表
    {{ .StructFilterInfo }}
}

// List{{ .StructName }}Req 获取{{ .TableComment }}列表请求
type List{{ .StructName }}Req {
    Filter {{ .StructName }}Filter `json:"filter"` // {{ .TableComment }}筛选参数
}

// List{{ .StructName }}Resp 获取{{ .TableComment }}列表响应
type List{{ .StructName }}Resp {
    Results []*{{ .StructName }} `json:"results"` // 结果
}

// Create{{ .StructName }}sReq 批量创建{{ .TableComment }}请求
type Create{{ .StructName }}sReq {
    Objects []*Create{{ .StructName }}Req `json:"objects" validate:"gt=0,dive" label:"{{ .TableComment }}列表"` // {{ .TableComment }}列表
}

// Create{{ .StructName }}sResp 批量创建{{ .TableComment }}响应
type Create{{ .StructName }}sResp {
    Results []*{{ .StructName }} `json:"results"` // 结果
}

// Update{{ .StructName }}sReq 批量更新{{ .TableComment }}请求
type Update{{ .StructName }}sReq {
    Filter {{ .StructName }}Filter `json:"filter"` // {{ .TableComment }}筛选参数
    {{ .StructBatchUpdateInfo }}
}

// Update{{ .StructName }}sResp 批量更新{{ .TableComment }}响应
type Update{{ .StructName }}sResp {
    Affected int64 `json:"affected"` // 影响数量
}

// Delete{{ .StructName }}sReq 批量删除{{ .TableComment }}请求
type Delete{{ .StructName }}sReq {
    Filter {{ .StructName }}Filter `json:"filter"` // {{ .TableComment }}筛选参数
}

// Delete{{ .StructName }}sResp 批量删除{{ .TableComment }}响应
type Delete{{ .StructName }}sResp {
    Affected int64 `json:"affected"` // 影响数量
}

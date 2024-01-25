syntax = "v1"

info (
    title:   "{{ .Title }}"
    desc:    "{{ .Desc }}"
    author:  "{{ .Author }}"
    email:   "{{ .Email }}"
    version: "{{ .Version }}"
)

@server (
    {{ if .ServerPrefix }}prefix: /{{ .ServerPrefix }}{{ end }}
    group: {{ if .GroupPrefix }}{{ .GroupPrefix }}/{{ end }}{{ .GroupName }}
)
service {{ .ServiceName }} {
    @doc (
        summary: "查询{{ .TableComment }}"
    )
    @handler Get{{ .StructName }}
    get /{{ .SnakeStructName }}/:id (Get{{ .StructName }}Req) returns (Get{{ .StructName }}Resp)

    @doc (
        summary: "查询{{ .TableComment }}分页"
    )
    @handler Get{{ .StructName }}Pages
    get /{{ .SnakeStructName }} (Get{{ .StructName }}PagesReq) returns (Get{{ .StructName }}PagesResp)

    @doc (
        summary: "创建{{ .TableComment }}"
    )
    @handler Create{{ .StructName }}
    post /{{ .SnakeStructName }} (Create{{ .StructName }}Req) returns (Create{{ .StructName }}Resp)

    @doc (
        summary: "更新{{ .TableComment }}"
    )
    @handler Update{{ .StructName }}
    put /{{ .SnakeStructName }}/:id (Update{{ .StructName }}Req) returns (Update{{ .StructName }}Resp)

    @doc (
        summary: "删除{{ .TableComment }}"
    )
    @handler Delete{{ .StructName }}
    delete /{{ .SnakeStructName }}/:id (Delete{{ .StructName }}Req) returns (Delete{{ .StructName }}Resp)
}

// -------------------- {{ .TableComment }} {{ .StructName }} -------------------- //
// {{ .StructName }} {{ .TableComment }}
type {{ .StructName }} {
    {{ .StructInfo }}
}

// Get{{ .StructName }}Req 查询{{ .TableComment }}请求
type Get{{ .StructName }}Req {
    Id int64 `path:"id" validate:"required" label:"{{ .IdLabel }}"` // {{ .IdComment }}
}

// Get{{ .StructName }}Resp 查询{{ .TableComment }}响应
type Get{{ .StructName }}Resp {
    {{ .StructName }}
}

// Get{{ .StructName }}PagesReq 查询{{ .TableComment }}分页请求
type Get{{ .StructName }}PagesReq {
    {{ .StructGetInfo }}
    Page     int64  `form:"page" validate:"required" label:"页数"`        // 页数
    PageSize int64  `form:"page_size" validate:"required" label:"每条页数"` // 每条页数
}

// Get{{ .StructName }}PagesResp 查询{{ .TableComment }}分页响应
type Get{{ .StructName }}PagesResp {
    Count     int64             `json:"count"`      // 总数
    PageCount int64             `json:"page_count"` // 页数
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
    Id int64 `path:"id" validate:"required" label:"{{ .IdLabel }}"` // {{ .IdComment }}
}

// Delete{{ .StructName }}Resp 删除{{ .TableComment }}响应
type Delete{{ .StructName }}Resp {
    Id int64 `json:"id"` // {{ .IdComment }}
}

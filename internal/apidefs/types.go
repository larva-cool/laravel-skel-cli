// Package apidefs 描述后端 API 接口的元数据定义。
// 接口清单由 tools/gen_registry.py 从 Apifox 导出的 OpenAPI 规范生成，
// 供 cmd 层的通用 runner 消费。
package apidefs

// Param 描述一个接口参数。
type Param struct {
	Name        string `json:"name"`
	In          string `json:"in"` // path / query / body
	Type        string `json:"type"` // string / integer / number / boolean / array
	Required    bool   `json:"required"`
	Description string `json:"description"`
}

// Endpoint 描述一个可被 CLI 调用的接口。
type Endpoint struct {
	Slug        string  `json:"slug"`
	Method      string  `json:"method"`
	Path        string  `json:"path"`
	Summary     string  `json:"summary"`
	PathParams  []Param `json:"pathParams"`
	QueryParams []Param `json:"queryParams"`
	BodyParams  []Param `json:"bodyParams"`
}

// Find 按 slug 查找接口，未找到返回 nil。
func Find(slug string) *Endpoint {
	for i := range Endpoints {
		if Endpoints[i].Slug == slug {
			return &Endpoints[i]
		}
	}
	return nil
}

// All 返回全部接口的只读视图。
func All() []Endpoint {
	return Endpoints
}

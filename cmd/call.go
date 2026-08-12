package cmd

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"laravel-skel-cli/internal/apidefs"
)

// runCall 实现通用 call 子命令：按 slug 查找接口定义，
// 将命令行参数映射到 path/query/body 并发送请求，输出 JSON 响应。
func runCall(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("用法: laravel-skel-cli call <接口slug> [参数]（可用 list 查看全部接口）")
	}
	slug := args[0]
	ep := apidefs.Find(slug)
	if ep == nil {
		return fmt.Errorf("未知接口 %q，可用 list 查看全部接口", slug)
	}

	// 动态注册该接口的全部参数（path/query/body）
	fs := flag.NewFlagSet("call "+slug, flag.ContinueOnError)
	values := make(map[string]*string)
	for _, p := range allParams(ep) {
		values[p.Name] = fs.String(p.Name, "", p.Description)
	}
	baseURL := fs.String("base-url", "", "覆盖配置中的 API 地址")

	if err := fs.Parse(args[1:]); err != nil {
		return err
	}
	if fs.NArg() > 0 {
		return fmt.Errorf("无法识别的参数: %v", fs.Args())
	}

	api, err := loadClient(*baseURL)
	if err != nil {
		return err
	}

	// 替换路径参数
	path := ep.Path
	for _, p := range ep.PathParams {
		v := *values[p.Name]
		if v == "" {
			return fmt.Errorf("缺少必填路径参数 --%s", p.Name)
		}
		path = strings.ReplaceAll(path, "{"+p.Name+"}", v)
	}

	// 收集查询参数
	q := url.Values{}
	for _, p := range ep.QueryParams {
		if v := *values[p.Name]; v != "" {
			q.Set(p.Name, v)
		}
	}
	if len(q) > 0 {
		path += "?" + q.Encode()
	}

	// 收集请求体字段（仅包含显式提供的字段，便于部分更新）
	var body any
	if isRawBody(ep) {
		// 原始请求体（如通知已读的 JSON 数组字符串），直接透传
		if v := *values["body"]; v != "" {
			body = json.RawMessage(v)
		}
	} else {
		obj := make(map[string]any)
		for _, p := range ep.BodyParams {
			if v := *values[p.Name]; v != "" {
				parsed, err := parseValue(v, p.Type)
				if err != nil {
					return fmt.Errorf("参数 --%s: %w", p.Name, err)
				}
				obj[p.Name] = parsed
			}
		}
		if len(obj) > 0 {
			body = obj
		}
	}

	var result any
	if err := api.Do(ep.Method, path, body, &result); err != nil {
		return err
	}
	printJSON(result)
	return nil
}

// isRawBody 判断接口请求体是否为原始 JSON（而非字段对象）。
func isRawBody(ep *apidefs.Endpoint) bool {
	return len(ep.BodyParams) == 1 && ep.BodyParams[0].Name == "body"
}

// allParams 返回接口涉及的所有参数（path、query、body）。
func allParams(ep *apidefs.Endpoint) []apidefs.Param {
	var out []apidefs.Param
	out = append(out, ep.PathParams...)
	out = append(out, ep.QueryParams...)
	out = append(out, ep.BodyParams...)
	return out
}

// parseValue 按参数类型将命令行字符串转换为对应 Go 值。
func parseValue(s, typ string) (any, error) {
	switch typ {
	case "integer":
		return strconv.ParseInt(s, 10, 64)
	case "number":
		return strconv.ParseFloat(s, 64)
	case "boolean":
		switch strings.ToLower(s) {
		case "1", "true", "yes":
			return true, nil
		case "0", "false", "no":
			return false, nil
		}
		return nil, fmt.Errorf("非法布尔值 %q", s)
	case "array":
		parts := strings.Split(s, ",")
		items := make([]any, 0, len(parts))
		for _, p := range parts {
			items = append(items, strings.TrimSpace(p))
		}
		return items, nil
	default: // string
		return s, nil
	}
}

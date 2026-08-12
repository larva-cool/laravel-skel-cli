// Package client 提供基于 net/http 的 API 客户端封装。
// 统一处理 Bearer Token 认证、JSON 序列化、以及错误解析。
package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// APIClient 封装对后端 API 的调用。
type APIClient struct {
	baseURL string
	token   string
	http    *http.Client
}

// New 创建一个 API 客户端。
// baseURL 为后端根地址（如 https://api.example.com），token 可为空。
func New(baseURL, token string) *APIClient {
	return &APIClient{
		baseURL: strings.TrimRight(baseURL, "/"),
		token:   token,
		http:    &http.Client{Timeout: 30 * time.Second},
	}
}

// SetToken 更新认证 token。
func (c *APIClient) SetToken(token string) {
	c.token = token
}

// Get 发起 GET 请求，将响应 JSON 解析到 out。
func (c *APIClient) Get(path string, out any) error {
	return c.Do(http.MethodGet, path, nil, out)
}

// Post 发起 POST 请求，将 body 序列化后发送，响应 JSON 解析到 out。
func (c *APIClient) Post(path string, body, out any) error {
	return c.Do(http.MethodPost, path, body, out)
}

// Do 执行一次 HTTP 请求。
// body 为 nil 表示无请求体；out 非空时将成功响应解析到其中。
func (c *APIClient) Do(method, path string, body, out any) error {
	var reader io.Reader
	if body != nil {
		data, err := json.Marshal(body)
		if err != nil {
			return fmt.Errorf("序列化请求体失败: %w", err)
		}
		reader = bytes.NewReader(data)
	}

	req, err := http.NewRequest(method, c.baseURL+path, reader)
	if err != nil {
		return fmt.Errorf("构造请求失败: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	if c.token != "" {
		req.Header.Set("Authorization", "Bearer "+c.token)
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return fmt.Errorf("请求失败: %w", err)
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("读取响应失败: %w", err)
	}

	if resp.StatusCode >= 400 {
		return fmt.Errorf("API 返回错误 (%d): %s", resp.StatusCode, strings.TrimSpace(string(data)))
	}

	if out == nil || len(data) == 0 {
		return nil
	}
	if err := json.Unmarshal(data, out); err != nil {
		return fmt.Errorf("解析响应 JSON 失败: %w（原始: %s）", err, string(data))
	}
	return nil
}

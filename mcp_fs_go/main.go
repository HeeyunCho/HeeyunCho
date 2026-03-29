package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

// MCP Types
type JSONRPCMessage struct {
	JSONRPC string          `json:"jsonrpc"`
	Method  string          `json:"method"`
	Params  json.RawMessage `json:"params,omitempty"`
	ID      interface{}     `json:"id,omitempty"`
}

type JSONRPCResponse struct {
	JSONRPC string      `json:"jsonrpc"`
	Result  interface{} `json:"result,omitempty"`
	Error   *RPCError   `json:"error,omitempty"`
	ID      interface{} `json:"id"`
}

type RPCError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type ListToolsResponse struct {
	Tools []Tool `json:"tools"`
}

type Tool struct {
	Name        string      `json:"name"`
	Description string      `json:"description"`
	InputSchema interface{} `json:"inputSchema"`
}

type CallToolRequest struct {
	Name      string                 `json:"name"`
	Arguments map[string]interface{} `json:"arguments"`
}

type ToolResult struct {
	Content []Content `json:"content"`
	IsError bool      `json:"isError,omitempty"`
}

type Content struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		var msg JSONRPCMessage
		if err := json.Unmarshal(scanner.Bytes(), &msg); err != nil {
			sendError(msg.ID, -32700, "Parse error")
			continue
		}

		switch msg.Method {
		case "initialize":
			sendResponse(msg.ID, map[string]interface{}{
				"protocolVersion": "2024-11-05",
				"capabilities":    map[string]interface{}{"tools": map[string]interface{}{}},
				"serverInfo":      map[string]interface{}{"name": "mcp_fs_go", "version": "1.0.0"},
			})
		case "tools/list":
			sendResponse(msg.ID, ListToolsResponse{
				Tools: []Tool{
					{
						Name:        "list_directory",
						Description: "Lists files and folders in a directory.",
						InputSchema: map[string]interface{}{
							"type": "object",
							"properties": map[string]interface{}{
								"path": map[string]interface{}{"type": "string"},
							},
							"required": []string{"path"},
						},
					},
					{
						Name:        "read_file",
						Description: "Reads the content of a file.",
						InputSchema: map[string]interface{}{
							"type": "object",
							"properties": map[string]interface{}{
								"path": map[string]interface{}{"type": "string"},
							},
							"required": []string{"path"},
						},
					},
				},
			})
		case "tools/call":
			var req CallToolRequest
			if err := json.Unmarshal(msg.Params, &req); err != nil {
				sendError(msg.ID, -32602, "Invalid params")
				continue
			}
			handleToolCall(msg.ID, req)
		}
	}
}

func handleToolCall(id interface{}, req CallToolRequest) {
	switch req.Name {
	case "list_directory":
		path, _ := req.Arguments["path"].(string)
		files, err := ioutil.ReadDir(path)
		if err != nil {
			sendToolResult(id, fmt.Sprintf("Error reading directory: %v", err), true)
			return
		}
		var names []string
		for _, f := range files {
			prefix := "[F] "
			if f.IsDir() {
				prefix = "[D] "
			}
			names = append(names, prefix+f.Name())
		}
		sendToolResult(id, strings.Join(names, "\n"), false)

	case "read_file":
		path, _ := req.Arguments["path"].(string)
		content, err := ioutil.ReadFile(path)
		if err != nil {
			sendToolResult(id, fmt.Sprintf("Error reading file: %v", err), true)
			return
		}
		sendToolResult(id, string(content), false)

	default:
		sendError(id, -32601, "Method not found")
	}
}

func sendResponse(id interface{}, result interface{}) {
	resp := JSONRPCResponse{
		JSONRPC: "2.0",
		Result:  result,
		ID:      id,
	}
	data, _ := json.Marshal(resp)
	fmt.Println(string(data))
}

func sendToolResult(id interface{}, text string, isError bool) {
	res := ToolResult{
		Content: []Content{{Type: "text", Text: text}},
		IsError: isError,
	}
	sendResponse(id, res)
}

func sendError(id interface{}, code int, message string) {
	resp := JSONRPCResponse{
		JSONRPC: "2.0",
		Error:   &RPCError{Code: code, Message: message},
		ID:      id,
	}
	data, _ := json.Marshal(resp)
	fmt.Println(string(data))
}

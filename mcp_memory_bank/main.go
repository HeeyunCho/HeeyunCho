package main

import (
	"bufio"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"

	_ "modernc.org/sqlite"
)

// DB Schema: 
// memories (id, category, content, project_id, created_at)
// search (id, content) - FTS5 Virtual Table

type Memory struct {
	ID        int64  `json:"id"`
	Category  string `json:"category"`
	Content   string `json:"content"`
	ProjectID string `json:"project_id"`
	CreatedAt string `json:"created_at"`
}

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

var db *sql.DB

func main() {
	var err error
	dbPath := "C:/Users/perfu/.gemini/memory_bank.db"
	db, err = sql.Open("sqlite", dbPath)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Initialize Schema
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS memories (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			category TEXT,
			content TEXT,
			project_id TEXT,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		);
		CREATE VIRTUAL TABLE IF NOT EXISTS memories_search USING fts5(content, content='memories', content_rowid='id');
	`)
	if err != nil {
		log.Fatal(err)
	}

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		var msg JSONRPCMessage
		if err := json.Unmarshal(scanner.Bytes(), &msg); err != nil {
			continue
		}

		switch msg.Method {
		case "initialize":
			sendResponse(msg.ID, map[string]interface{}{
				"protocolVersion": "2024-11-05",
				"capabilities":    map[string]interface{}{"tools": map[string]interface{}{}},
				"serverInfo":      map[string]interface{}{"name": "mcp_memory_bank", "version": "1.0.0"},
			})
		case "tools/list":
			sendResponse(msg.ID, map[string]interface{}{
				"tools": []map[string]interface{}{
					{
						"name":        "store_memory",
						"description": "Store a long-term memory or user preference.",
						"inputSchema": map[string]interface{}{
							"type": "object",
							"properties": map[string]interface{}{
								"content":    map[string]interface{}{"type": "string", "description": "The fact or preference to remember"},
								"category":   map[string]interface{}{"type": "string", "description": "e.g., 'preference', 'project_context', 'decision'"},
								"project_id": map[string]interface{}{"type": "string"},
							},
							"required": []string{"content", "category"},
						},
					},
					{
						"name":        "retrieve_memories",
						"description": "Search for relevant past memories using full-text search.",
						"inputSchema": map[string]interface{}{
							"type": "object",
							"properties": map[string]interface{}{
								"query":      map[string]interface{}{"type": "string", "description": "Search term (e.g., 'python', 'architecture')"},
								"project_id": map[string]interface{}{"type": "string"},
							},
							"required": []string{"query"},
						},
					},
				},
			})
		case "tools/call":
			var req struct {
				Name      string                 `json:"name"`
				Arguments map[string]interface{} `json:"arguments"`
			}
			json.Unmarshal(msg.Params, &req)
			handleToolCall(msg.ID, req.Name, req.Arguments)
		}
	}
}

func handleToolCall(id interface{}, name string, args map[string]interface{}) {
	switch name {
	case "store_memory":
		content := args["content"].(string)
		category := args["category"].(string)
		projectID, _ := args["project_id"].(string)

		res, _ := db.Exec("INSERT INTO memories (content, category, project_id) VALUES (?, ?, ?)", content, category, projectID)
		newID, _ := res.LastInsertId()
		// Update search index
		db.Exec("INSERT INTO memories_search (rowid, content) VALUES (?, ?)", newID, content)

		sendToolResult(id, fmt.Sprintf("Memory stored successfully. (ID: %d)", newID), false)

	case "retrieve_memories":
		query := args["query"].(string)
		projectID, _ := args["project_id"].(string)

		sqlQuery := "SELECT memories.content, category, created_at FROM memories_search JOIN memories ON memories_search.rowid = memories.id WHERE memories_search MATCH ? "
		if projectID != "" {
			sqlQuery += " AND project_id = '" + projectID + "'"
		}
		sqlQuery += " ORDER BY rank LIMIT 10"

		rows, err := db.Query(sqlQuery, query)
		if err != nil {
			sendToolResult(id, fmt.Sprintf("Error retrieving memories: %v", err), true)
			return
		}
		defer rows.Close()

		var results []string
		for rows.Next() {
			var c, cat, ts string
			if err := rows.Scan(&c, &cat, &ts); err != nil {
				continue
			}
			results = append(results, fmt.Sprintf("[%s] %s (Saved: %s)", cat, c, ts))
		}

		if len(results) == 0 {
			sendToolResult(id, "No matching memories found.", false)
		} else {
			output := "Found relevant memories:\n"
			for _, r := range results {
				output += "- " + r + "\n"
			}
			sendToolResult(id, output, false)
		}
	}
}

func sendResponse(id interface{}, result interface{}) {
	resp := JSONRPCResponse{JSONRPC: "2.0", Result: result, ID: id}
	data, _ := json.Marshal(resp)
	fmt.Println(string(data))
}

func sendToolResult(id interface{}, text string, isError bool) {
	res := map[string]interface{}{
		"content": []map[string]interface{}{{"type": "text", "text": text}},
		"isError": isError,
	}
	sendResponse(id, res)
}

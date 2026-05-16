package mcp

// JSON-RPC 2.0 message types for MCP protocol

// JSONRPCRequest represents a JSON-RPC 2.0 request
type JSONRPCRequest struct {
	JSONRPC string                 `json:"jsonrpc"`
	ID      interface{}            `json:"id,omitempty"`
	Method  string                 `json:"method"`
	Params  map[string]interface{} `json:"params,omitempty"`
}

// JSONRPCResponse represents a JSON-RPC 2.0 response
type JSONRPCResponse struct {
	JSONRPC string      `json:"jsonrpc"`
	ID      interface{} `json:"id,omitempty"`
	Result  interface{} `json:"result,omitempty"`
	Error   *RPCError   `json:"error,omitempty"`
}

// RPCError represents a JSON-RPC 2.0 error
type RPCError struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// Standard JSON-RPC error codes
const (
	ParseError     = -32700
	InvalidRequest = -32600
	MethodNotFound = -32601
	InvalidParams  = -32602
	InternalError  = -32603
)

// MCP Protocol Types

// InitializeRequest represents the initialize method params
type InitializeRequest struct {
	ProtocolVersion string                 `json:"protocolVersion"`
	Capabilities    ClientCapabilities     `json:"capabilities"`
	ClientInfo      ClientInfo             `json:"clientInfo"`
	Meta            map[string]interface{} `json:"_meta,omitempty"`
}

// ClientCapabilities describes what the client supports
type ClientCapabilities struct {
	Roots       *RootsCapability       `json:"roots,omitempty"`
	Sampling    map[string]interface{} `json:"sampling,omitempty"`
	Experimental map[string]interface{} `json:"experimental,omitempty"`
}

// RootsCapability indicates if client supports roots
type RootsCapability struct {
	ListChanged bool `json:"listChanged,omitempty"`
}

// ClientInfo contains information about the client
type ClientInfo struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

// InitializeResult is the response to initialize
type InitializeResult struct {
	ProtocolVersion string             `json:"protocolVersion"`
	Capabilities    ServerCapabilities `json:"capabilities"`
	ServerInfo      ServerInfo         `json:"serverInfo"`
}

// ServerCapabilities describes what the server supports
type ServerCapabilities struct {
	Tools     *ToolsCapability     `json:"tools,omitempty"`
	Resources *ResourcesCapability `json:"resources,omitempty"`
	Prompts   *PromptsCapability   `json:"prompts,omitempty"`
	Logging   map[string]interface{} `json:"logging,omitempty"`
}

// ToolsCapability indicates tool support
type ToolsCapability struct {
	ListChanged bool `json:"listChanged,omitempty"`
}

// ResourcesCapability indicates resource support
type ResourcesCapability struct {
	Subscribe   bool `json:"subscribe,omitempty"`
	ListChanged bool `json:"listChanged,omitempty"`
}

// PromptsCapability indicates prompt support
type PromptsCapability struct {
	ListChanged bool `json:"listChanged,omitempty"`
}

// ServerInfo contains information about the server
type ServerInfo struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

// Tool represents an MCP tool definition
type Tool struct {
	Name        string      `json:"name"`
	Description string      `json:"description"`
	InputSchema InputSchema `json:"inputSchema"`
}

// InputSchema defines the JSON schema for tool parameters
type InputSchema struct {
	Type       string                 `json:"type"`
	Properties map[string]Property    `json:"properties,omitempty"`
	Required   []string               `json:"required,omitempty"`
	Additional map[string]interface{} `json:"additionalProperties,omitempty"`
}

// Property defines a parameter property
type Property struct {
	Type        string        `json:"type"`
	Description string        `json:"description,omitempty"`
	Enum        []string      `json:"enum,omitempty"`
	Default     interface{}   `json:"default,omitempty"`
	Items       *PropertyType `json:"items,omitempty"`
}

// PropertyType for array items
type PropertyType struct {
	Type string `json:"type"`
}

// ToolsListResult is the response to tools/list
type ToolsListResult struct {
	Tools []Tool `json:"tools"`
}

// CallToolRequest represents a tools/call request params
type CallToolRequest struct {
	Name      string                 `json:"name"`
	Arguments map[string]interface{} `json:"arguments,omitempty"`
}

// CallToolResult is the response to tools/call
type CallToolResult struct {
	Content []Content `json:"content"`
	IsError bool      `json:"isError,omitempty"`
}

// Content represents tool output content
type Content struct {
	Type     string `json:"type"`
	Text     string `json:"text,omitempty"`
	MimeType string `json:"mimeType,omitempty"`
}

// Resource represents an MCP resource definition
type Resource struct {
	URI         string `json:"uri"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	MimeType    string `json:"mimeType,omitempty"`
}

// ResourcesListResult is the response to resources/list
type ResourcesListResult struct {
	Resources []Resource `json:"resources"`
}

// ReadResourceRequest represents a resources/read request params
type ReadResourceRequest struct {
	URI string `json:"uri"`
}

// ReadResourceResult is the response to resources/read
type ReadResourceResult struct {
	Contents []ResourceContent `json:"contents"`
}

// ResourceContent represents resource data
type ResourceContent struct {
	URI      string `json:"uri"`
	MimeType string `json:"mimeType"`
	Text     string `json:"text,omitempty"`
	Blob     string `json:"blob,omitempty"`
}

// Made with Bob

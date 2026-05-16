package mcp

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
)

// Server represents the MCP server
type Server struct {
	toolHandler     *ToolHandler
	resourceHandler *ResourceHandler
	cache           *ResultCache
	version         string
	initialized     bool
}

// NewServer creates a new MCP server
func NewServer(version string) *Server {
	cache := NewResultCache()
	
	// Get API key from environment
	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		apiKey = os.Getenv("GOOGLE_API_KEY")
	}
	
	// Get model from environment or use default
	model := os.Getenv("OBELISK_MODEL")
	if model == "" {
		model = "gemini-2.0-flash-exp"
	}
	
	return &Server{
		toolHandler:     NewToolHandler(cache, apiKey, model),
		resourceHandler: NewResourceHandler(cache),
		cache:           cache,
		version:         version,
		initialized:     false,
	}
}

// Run starts the MCP server with stdio transport
func (s *Server) Run() error {
	reader := bufio.NewReader(os.Stdin)
	
	// Log to stderr (stdout is reserved for JSON-RPC)
	fmt.Fprintln(os.Stderr, "Obelisk MCP Server starting...")
	fmt.Fprintf(os.Stderr, "Version: %s\n", s.version)
	fmt.Fprintln(os.Stderr, "Listening on stdio for JSON-RPC 2.0 messages...")
	
	for {
		// Read a line from stdin
		line, err := reader.ReadBytes('\n')
		if err != nil {
			if err == io.EOF {
				fmt.Fprintln(os.Stderr, "Client disconnected")
				return nil
			}
			return fmt.Errorf("error reading from stdin: %w", err)
		}
		
		// Parse JSON-RPC request
		var req JSONRPCRequest
		if err := json.Unmarshal(line, &req); err != nil {
			s.sendError(nil, ParseError, "Invalid JSON", err.Error())
			continue
		}
		
		// Handle the request
		s.handleRequest(&req)
	}
}

// handleRequest processes a JSON-RPC request
func (s *Server) handleRequest(req *JSONRPCRequest) {
	// Validate JSON-RPC version
	if req.JSONRPC != "2.0" {
		s.sendError(req.ID, InvalidRequest, "Invalid JSON-RPC version", "Expected 2.0")
		return
	}
	
	// Route to appropriate handler
	switch req.Method {
	case "initialize":
		s.handleInitialize(req)
	case "tools/list":
		s.handleToolsList(req)
	case "tools/call":
		s.handleToolsCall(req)
	case "resources/list":
		s.handleResourcesList(req)
	case "resources/read":
		s.handleResourcesRead(req)
	case "ping":
		s.handlePing(req)
	default:
		s.sendError(req.ID, MethodNotFound, fmt.Sprintf("Method not found: %s", req.Method), nil)
	}
}

// handleInitialize handles the initialize method
func (s *Server) handleInitialize(req *JSONRPCRequest) {
	// Parse initialize params
	var initReq InitializeRequest
	if req.Params != nil {
		paramsJSON, _ := json.Marshal(req.Params)
		if err := json.Unmarshal(paramsJSON, &initReq); err != nil {
			s.sendError(req.ID, InvalidParams, "Invalid initialize parameters", err.Error())
			return
		}
	}
	
	// Log client info
	fmt.Fprintf(os.Stderr, "Client connected: %s v%s\n", initReq.ClientInfo.Name, initReq.ClientInfo.Version)
	
	// Create initialize result
	result := InitializeResult{
		ProtocolVersion: "2024-11-05",
		Capabilities: ServerCapabilities{
			Tools: &ToolsCapability{
				ListChanged: false,
			},
			Resources: &ResourcesCapability{
				Subscribe:   false,
				ListChanged: false,
			},
		},
		ServerInfo: ServerInfo{
			Name:    "obelisk-mcp-server",
			Version: s.version,
		},
	}
	
	s.initialized = true
	s.sendResult(req.ID, result)
}

// handleToolsList handles the tools/list method
func (s *Server) handleToolsList(req *JSONRPCRequest) {
	if !s.initialized {
		s.sendError(req.ID, InvalidRequest, "Server not initialized", "Call initialize first")
		return
	}
	
	tools := s.toolHandler.GetToolsList()
	result := ToolsListResult{
		Tools: tools,
	}
	
	s.sendResult(req.ID, result)
}

// handleToolsCall handles the tools/call method
func (s *Server) handleToolsCall(req *JSONRPCRequest) {
	if !s.initialized {
		s.sendError(req.ID, InvalidRequest, "Server not initialized", "Call initialize first")
		return
	}
	
	// Parse call tool params
	var callReq CallToolRequest
	if req.Params != nil {
		paramsJSON, _ := json.Marshal(req.Params)
		if err := json.Unmarshal(paramsJSON, &callReq); err != nil {
			s.sendError(req.ID, InvalidParams, "Invalid tool call parameters", err.Error())
			return
		}
	}
	
	fmt.Fprintf(os.Stderr, "Executing tool: %s\n", callReq.Name)
	
	// Execute the tool
	result, err := s.toolHandler.ExecuteTool(callReq.Name, callReq.Arguments)
	if err != nil {
		s.sendError(req.ID, InternalError, fmt.Sprintf("Tool execution failed: %s", callReq.Name), err.Error())
		return
	}
	
	s.sendResult(req.ID, result)
}

// handleResourcesList handles the resources/list method
func (s *Server) handleResourcesList(req *JSONRPCRequest) {
	if !s.initialized {
		s.sendError(req.ID, InvalidRequest, "Server not initialized", "Call initialize first")
		return
	}
	
	resources := s.resourceHandler.GetResourcesList()
	result := ResourcesListResult{
		Resources: resources,
	}
	
	s.sendResult(req.ID, result)
}

// handleResourcesRead handles the resources/read method
func (s *Server) handleResourcesRead(req *JSONRPCRequest) {
	if !s.initialized {
		s.sendError(req.ID, InvalidRequest, "Server not initialized", "Call initialize first")
		return
	}
	
	// Parse read resource params
	var readReq ReadResourceRequest
	if req.Params != nil {
		paramsJSON, _ := json.Marshal(req.Params)
		if err := json.Unmarshal(paramsJSON, &readReq); err != nil {
			s.sendError(req.ID, InvalidParams, "Invalid resource read parameters", err.Error())
			return
		}
	}
	
	fmt.Fprintf(os.Stderr, "Reading resource: %s\n", readReq.URI)
	
	// Read the resource
	result, err := s.resourceHandler.ReadResource(readReq.URI)
	if err != nil {
		s.sendError(req.ID, InternalError, fmt.Sprintf("Resource read failed: %s", readReq.URI), err.Error())
		return
	}
	
	s.sendResult(req.ID, result)
}

// handlePing handles the ping method (for testing)
func (s *Server) handlePing(req *JSONRPCRequest) {
	result := map[string]string{
		"status": "ok",
		"server": "obelisk-mcp-server",
	}
	s.sendResult(req.ID, result)
}

// sendResult sends a successful JSON-RPC response
func (s *Server) sendResult(id interface{}, result interface{}) {
	response := JSONRPCResponse{
		JSONRPC: "2.0",
		ID:      id,
		Result:  result,
	}
	
	s.sendResponse(response)
}

// sendError sends an error JSON-RPC response
func (s *Server) sendError(id interface{}, code int, message string, data interface{}) {
	response := JSONRPCResponse{
		JSONRPC: "2.0",
		ID:      id,
		Error: &RPCError{
			Code:    code,
			Message: message,
			Data:    data,
		},
	}
	
	s.sendResponse(response)
}

// sendResponse writes a JSON-RPC response to stdout
func (s *Server) sendResponse(response JSONRPCResponse) {
	data, err := json.Marshal(response)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error marshaling response: %v\n", err)
		return
	}
	
	// Write to stdout with newline
	fmt.Println(string(data))
}

// Made with Bob

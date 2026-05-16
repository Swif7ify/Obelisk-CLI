package mcp

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"
)

// HTTPServer wraps the MCP server with HTTP/SSE transport
type HTTPServer struct {
	server *Server
	mu     sync.RWMutex
}

// NewHTTPServer creates a new HTTP server wrapper
func NewHTTPServer(version string) *HTTPServer {
	return &HTTPServer{
		server: NewServer(version),
	}
}

// ServeHTTP handles HTTP requests
func (h *HTTPServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// CORS headers for browser clients
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	// Route based on path
	switch r.URL.Path {
	case "/sse", "/":
		h.handleSSE(w, r)
	case "/health":
		h.handleHealth(w, r)
	default:
		http.NotFound(w, r)
	}
}

// handleHealth provides a health check endpoint
func (h *HTTPServer) handleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":  "ok",
		"server":  "obelisk-mcp-server",
		"version": h.server.version,
	})
}

// handleSSE handles Server-Sent Events connections for MCP
func (h *HTTPServer) handleSSE(w http.ResponseWriter, r *http.Request) {
	// Set SSE headers
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("X-Accel-Buffering", "no")

	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming unsupported", http.StatusInternalServerError)
		return
	}

	// Create context for this connection
	ctx, cancel := context.WithCancel(r.Context())
	defer cancel()

	// Log connection
	log.Printf("SSE client connected from %s", r.RemoteAddr)

	// Create channels for communication
	requestChan := make(chan []byte, 10)
	responseChan := make(chan JSONRPCResponse, 10)
	errorChan := make(chan error, 1)

	// Start request processor
	go h.processRequests(ctx, requestChan, responseChan, errorChan)

	// Handle POST requests (for sending messages)
	if r.Method == "POST" {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Failed to read request body", http.StatusBadRequest)
			return
		}
		requestChan <- body
	}

	// Send initial connection message
	fmt.Fprintf(w, "event: connected\ndata: {\"status\":\"connected\"}\n\n")
	flusher.Flush()

	// Handle bidirectional communication
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			log.Printf("SSE client disconnected from %s", r.RemoteAddr)
			return

		case response := <-responseChan:
			// Send JSON-RPC response as SSE event
			data, err := json.Marshal(response)
			if err != nil {
				log.Printf("Error marshaling response: %v", err)
				continue
			}
			fmt.Fprintf(w, "event: message\ndata: %s\n\n", string(data))
			flusher.Flush()

		case err := <-errorChan:
			log.Printf("Error processing request: %v", err)
			fmt.Fprintf(w, "event: error\ndata: {\"error\":\"%s\"}\n\n", err.Error())
			flusher.Flush()

		case <-ticker.C:
			// Send keepalive ping
			fmt.Fprintf(w, "event: ping\ndata: {\"timestamp\":%d}\n\n", time.Now().Unix())
			flusher.Flush()
		}
	}
}

// processRequests processes incoming JSON-RPC requests
func (h *HTTPServer) processRequests(ctx context.Context, requestChan <-chan []byte, responseChan chan<- JSONRPCResponse, errorChan chan<- error) {
	for {
		select {
		case <-ctx.Done():
			return

		case data := <-requestChan:
			// Parse JSON-RPC request
			var req JSONRPCRequest
			if err := json.Unmarshal(data, &req); err != nil {
				errorChan <- fmt.Errorf("invalid JSON: %w", err)
				continue
			}

			// Process request and send response
			response := h.handleRequestSync(&req)
			responseChan <- response
		}
	}
}

// handleRequestSync processes a JSON-RPC request synchronously
func (h *HTTPServer) handleRequestSync(req *JSONRPCRequest) JSONRPCResponse {
	h.mu.Lock()
	defer h.mu.Unlock()

	// Validate JSON-RPC version
	if req.JSONRPC != "2.0" {
		return JSONRPCResponse{
			JSONRPC: "2.0",
			ID:      req.ID,
			Error: &RPCError{
				Code:    InvalidRequest,
				Message: "Invalid JSON-RPC version",
				Data:    "Expected 2.0",
			},
		}
	}

	// Route to appropriate handler
	switch req.Method {
	case "initialize":
		return h.handleInitializeSync(req)
	case "tools/list":
		return h.handleToolsListSync(req)
	case "tools/call":
		return h.handleToolsCallSync(req)
	case "resources/list":
		return h.handleResourcesListSync(req)
	case "resources/read":
		return h.handleResourcesReadSync(req)
	case "ping":
		return h.handlePingSync(req)
	default:
		return JSONRPCResponse{
			JSONRPC: "2.0",
			ID:      req.ID,
			Error: &RPCError{
				Code:    MethodNotFound,
				Message: fmt.Sprintf("Method not found: %s", req.Method),
			},
		}
	}
}

// Synchronous handler methods
func (h *HTTPServer) handleInitializeSync(req *JSONRPCRequest) JSONRPCResponse {
	var initReq InitializeRequest
	if req.Params != nil {
		paramsJSON, _ := json.Marshal(req.Params)
		if err := json.Unmarshal(paramsJSON, &initReq); err != nil {
			return JSONRPCResponse{
				JSONRPC: "2.0",
				ID:      req.ID,
				Error: &RPCError{
					Code:    InvalidParams,
					Message: "Invalid initialize parameters",
					Data:    err.Error(),
				},
			}
		}
	}

	log.Printf("Client connected: %s v%s", initReq.ClientInfo.Name, initReq.ClientInfo.Version)

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
			Version: h.server.version,
		},
	}

	h.server.initialized = true
	return JSONRPCResponse{
		JSONRPC: "2.0",
		ID:      req.ID,
		Result:  result,
	}
}

func (h *HTTPServer) handleToolsListSync(req *JSONRPCRequest) JSONRPCResponse {
	if !h.server.initialized {
		return JSONRPCResponse{
			JSONRPC: "2.0",
			ID:      req.ID,
			Error: &RPCError{
				Code:    InvalidRequest,
				Message: "Server not initialized",
				Data:    "Call initialize first",
			},
		}
	}

	tools := h.server.toolHandler.GetToolsList()
	return JSONRPCResponse{
		JSONRPC: "2.0",
		ID:      req.ID,
		Result: ToolsListResult{
			Tools: tools,
		},
	}
}

func (h *HTTPServer) handleToolsCallSync(req *JSONRPCRequest) JSONRPCResponse {
	if !h.server.initialized {
		return JSONRPCResponse{
			JSONRPC: "2.0",
			ID:      req.ID,
			Error: &RPCError{
				Code:    InvalidRequest,
				Message: "Server not initialized",
				Data:    "Call initialize first",
			},
		}
	}

	var callReq CallToolRequest
	if req.Params != nil {
		paramsJSON, _ := json.Marshal(req.Params)
		if err := json.Unmarshal(paramsJSON, &callReq); err != nil {
			return JSONRPCResponse{
				JSONRPC: "2.0",
				ID:      req.ID,
				Error: &RPCError{
					Code:    InvalidParams,
					Message: "Invalid tool call parameters",
					Data:    err.Error(),
				},
			}
		}
	}

	log.Printf("Executing tool: %s", callReq.Name)

	result, err := h.server.toolHandler.ExecuteTool(callReq.Name, callReq.Arguments)
	if err != nil {
		return JSONRPCResponse{
			JSONRPC: "2.0",
			ID:      req.ID,
			Error: &RPCError{
				Code:    InternalError,
				Message: fmt.Sprintf("Tool execution failed: %s", callReq.Name),
				Data:    err.Error(),
			},
		}
	}

	return JSONRPCResponse{
		JSONRPC: "2.0",
		ID:      req.ID,
		Result:  result,
	}
}

func (h *HTTPServer) handleResourcesListSync(req *JSONRPCRequest) JSONRPCResponse {
	if !h.server.initialized {
		return JSONRPCResponse{
			JSONRPC: "2.0",
			ID:      req.ID,
			Error: &RPCError{
				Code:    InvalidRequest,
				Message: "Server not initialized",
				Data:    "Call initialize first",
			},
		}
	}

	resources := h.server.resourceHandler.GetResourcesList()
	return JSONRPCResponse{
		JSONRPC: "2.0",
		ID:      req.ID,
		Result: ResourcesListResult{
			Resources: resources,
		},
	}
}

func (h *HTTPServer) handleResourcesReadSync(req *JSONRPCRequest) JSONRPCResponse {
	if !h.server.initialized {
		return JSONRPCResponse{
			JSONRPC: "2.0",
			ID:      req.ID,
			Error: &RPCError{
				Code:    InvalidRequest,
				Message: "Server not initialized",
				Data:    "Call initialize first",
			},
		}
	}

	var readReq ReadResourceRequest
	if req.Params != nil {
		paramsJSON, _ := json.Marshal(req.Params)
		if err := json.Unmarshal(paramsJSON, &readReq); err != nil {
			return JSONRPCResponse{
				JSONRPC: "2.0",
				ID:      req.ID,
				Error: &RPCError{
					Code:    InvalidParams,
					Message: "Invalid resource read parameters",
					Data:    err.Error(),
				},
			}
		}
	}

	log.Printf("Reading resource: %s", readReq.URI)

	result, err := h.server.resourceHandler.ReadResource(readReq.URI)
	if err != nil {
		return JSONRPCResponse{
			JSONRPC: "2.0",
			ID:      req.ID,
			Error: &RPCError{
				Code:    InternalError,
				Message: fmt.Sprintf("Resource read failed: %s", readReq.URI),
				Data:    err.Error(),
			},
		}
	}

	return JSONRPCResponse{
		JSONRPC: "2.0",
		ID:      req.ID,
		Result:  result,
	}
}

func (h *HTTPServer) handlePingSync(req *JSONRPCRequest) JSONRPCResponse {
	return JSONRPCResponse{
		JSONRPC: "2.0",
		ID:      req.ID,
		Result: map[string]string{
			"status": "ok",
			"server": "obelisk-mcp-server",
		},
	}
}

// StartHTTPServer starts the HTTP server on the specified port
func StartHTTPServer(version string, port string) error {
	if port == "" {
		port = "8080"
	}

	// Ensure port has colon prefix
	if !strings.HasPrefix(port, ":") {
		port = ":" + port
	}

	httpServer := NewHTTPServer(version)
	
	log.Printf("Obelisk MCP Server starting...")
	log.Printf("Version: %s", version)
	log.Printf("Listening on http://0.0.0.0%s", port)
	log.Printf("SSE endpoint: http://0.0.0.0%s/sse", port)
	log.Printf("Health check: http://0.0.0.0%s/health", port)

	return http.ListenAndServe(port, httpServer)
}

// Made with Bob
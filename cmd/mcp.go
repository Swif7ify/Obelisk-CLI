package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/Swif7ify/Obelisk-CLI/internal/mcp"
)

var (
	mcpHTTPMode bool
	mcpPort     string
)

var mcpCmd = &cobra.Command{
	Use:   "mcp",
	Short: "Run Obelisk as an MCP (Model Context Protocol) server",
	Long: `Starts Obelisk in MCP server mode, exposing code analysis capabilities
via the Model Context Protocol.

This mode is designed to be used by AI assistants and IDEs that support MCP.

Transport Modes:
  - stdio (default): JSON-RPC 2.0 over stdin/stdout for local use
  - http: HTTP/SSE transport for cloud deployment (use --http flag)

Environment Variables:
  GEMINI_API_KEY    - Google Gemini API key for AI-powered analysis
  GOOGLE_API_KEY    - Alternative to GEMINI_API_KEY
  OBELISK_MODEL     - AI model to use (default: gemini-2.0-flash-exp)
  PORT              - HTTP server port (default: 8080, only for --http mode)

Available Tools:
  - scan_project        Full project health scan
  - check_security      Security-focused scan
  - analyze_complexity  Code complexity analysis
  - track_tech_debt     Technical debt tracking
  - audit_dependencies  Dependency audit
  - get_health_report   AI-powered health assessment

Available Resources:
  - obelisk://scan/latest           Latest scan results
  - obelisk://health/score          Project health score
  - obelisk://findings/security     Security findings
  - obelisk://findings/quality      Code quality findings
  - obelisk://findings/architecture Architecture findings

Example MCP Configuration (Local - Bob IDE):
  {
    "mcpServers": {
      "obelisk": {
        "command": "obelisk",
        "args": ["mcp"],
        "env": {
          "GEMINI_API_KEY": "your-api-key-here"
        }
      }
    }
  }

Example MCP Configuration (Cloud - Remote URL):
  {
    "mcpServers": {
      "obelisk-cloud": {
        "url": "https://your-app.onrender.com/sse"
      }
    }
  }`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if mcpHTTPMode {
			// HTTP/SSE mode for cloud deployment
			port := mcpPort
			if port == "" {
				port = os.Getenv("PORT")
			}
			if port == "" {
				port = "8080"
			}

			if err := mcp.StartHTTPServer(Version, port); err != nil {
				fmt.Fprintf(os.Stderr, "HTTP server error: %v\n", err)
				return err
			}
		} else {
			// stdio mode for local use
			os.Setenv("NO_COLOR", "1")

			server := mcp.NewServer(Version)
			if err := server.Run(); err != nil {
				fmt.Fprintf(os.Stderr, "MCP server error: %v\n", err)
				return err
			}
		}

		return nil
	},
}

func init() {
	mcpCmd.Flags().BoolVar(&mcpHTTPMode, "http", false, "Run in HTTP/SSE mode for cloud deployment")
	mcpCmd.Flags().StringVar(&mcpPort, "port", "", "HTTP server port (default: 8080 or PORT env var)")
	rootCmd.AddCommand(mcpCmd)
}

// Made with Bob

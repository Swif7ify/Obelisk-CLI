package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/Swif7ify/Obelisk-CLI/internal/mcp"
)

var mcpCmd = &cobra.Command{
	Use:   "mcp",
	Short: "Run Obelisk as an MCP (Model Context Protocol) server",
	Long: `Starts Obelisk in MCP server mode, exposing code analysis capabilities
via the Model Context Protocol over stdio transport.

This mode is designed to be used by AI assistants and IDEs that support MCP.
The server communicates using JSON-RPC 2.0 over stdin/stdout.

Environment Variables:
  GEMINI_API_KEY    - Google Gemini API key for AI-powered analysis
  GOOGLE_API_KEY    - Alternative to GEMINI_API_KEY
  OBELISK_MODEL     - AI model to use (default: gemini-2.0-flash-exp)

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

Example MCP Configuration (Bob IDE):
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
  }`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Disable colored output for MCP mode (stdout is for JSON-RPC only)
		os.Setenv("NO_COLOR", "1")

		// Create and run the MCP server
		server := mcp.NewServer(Version)
		
		if err := server.Run(); err != nil {
			fmt.Fprintf(os.Stderr, "MCP server error: %v\n", err)
			return err
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(mcpCmd)
}

// Made with Bob

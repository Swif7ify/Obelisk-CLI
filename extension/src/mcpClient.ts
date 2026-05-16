import * as vscode from "vscode";
import {
	ObeliskScanResult,
	ObeliskFinding,
	ObeliskReport,
} from "./obeliskRunner";

interface JSONRPCRequest {
	jsonrpc: string;
	id: number | string;
	method: string;
	params?: any;
}

interface JSONRPCResponse {
	jsonrpc: string;
	id: number | string;
	result?: any;
	error?: {
		code: number;
		message: string;
		data?: any;
	};
}

export class MCPClient {
	private serverUrl: string;
	private requestId: number = 0;
	private initialized: boolean = false;

	constructor(serverUrl: string = "https://mcp-obelisk.onedevph.online/sse") {
		this.serverUrl = serverUrl;
	}

	async initialize(): Promise<void> {
		if (this.initialized) {
			return;
		}

		const request: JSONRPCRequest = {
			jsonrpc: "2.0",
			id: ++this.requestId,
			method: "initialize",
			params: {
				protocolVersion: "2024-11-05",
				capabilities: {},
				clientInfo: {
					name: "obelisk-vscode-extension",
					version: "0.1.0",
				},
			},
		};

		const response = await this.sendRequest(request);
		if (response.error) {
			throw new Error(
				`MCP initialization failed: ${response.error.message}`,
			);
		}

		this.initialized = true;
	}

	async scanProject(
		projectPath: string,
		skipAI: boolean = false,
	): Promise<ObeliskScanResult> {
		await this.initialize();

		const request: JSONRPCRequest = {
			jsonrpc: "2.0",
			id: ++this.requestId,
			method: "tools/call",
			params: {
				name: "scan_project",
				arguments: {
					path: projectPath,
					skip_ai: skipAI,
				},
			},
		};

		const response = await this.sendRequest(request);
		if (response.error) {
			throw new Error(`Scan failed: ${response.error.message}`);
		}

		// Parse the result from MCP tool response
		const content = response.result?.content?.[0]?.text;
		if (!content) {
			throw new Error("Invalid response from MCP server");
		}

		const parsed = JSON.parse(content);

		// Transform to match ObeliskScanResult format
		return {
			scan_result: parsed.scan_result || {
				project_path: parsed.project_path || projectPath,
				project_type: parsed.project_type || "unknown",
				findings: parsed.findings_list || [],
				file_count: parsed.file_count || 0,
				dir_count: parsed.dir_count || 0,
			},
			report: parsed.health_report || null,
			detection: null, // MCP doesn't return detection info
		};
	}

	async checkSecurity(projectPath: string): Promise<any> {
		await this.initialize();

		const request: JSONRPCRequest = {
			jsonrpc: "2.0",
			id: ++this.requestId,
			method: "tools/call",
			params: {
				name: "check_security",
				arguments: {
					path: projectPath,
				},
			},
		};

		const response = await this.sendRequest(request);
		if (response.error) {
			throw new Error(`Security check failed: ${response.error.message}`);
		}

		const content = response.result?.content?.[0]?.text;
		return JSON.parse(content);
	}

	async getHealthReport(projectPath: string): Promise<any> {
		await this.initialize();

		const request: JSONRPCRequest = {
			jsonrpc: "2.0",
			id: ++this.requestId,
			method: "tools/call",
			params: {
				name: "get_health_report",
				arguments: {
					path: projectPath,
				},
			},
		};

		const response = await this.sendRequest(request);
		if (response.error) {
			throw new Error(`Health report failed: ${response.error.message}`);
		}

		const content = response.result?.content?.[0]?.text;
		return JSON.parse(content);
	}

	private async sendRequest(
		request: JSONRPCRequest,
	): Promise<JSONRPCResponse> {
		try {
			const response = await fetch(this.serverUrl, {
				method: "POST",
				headers: {
					"Content-Type": "application/json",
				},
				body: JSON.stringify(request),
			});

			if (!response.ok) {
				throw new Error(
					`HTTP ${response.status}: ${response.statusText}`,
				);
			}

			// For SSE, we need to read the event stream
			const reader = response.body?.getReader();
			if (!reader) {
				throw new Error("No response body");
			}

			const decoder = new TextDecoder();
			let buffer = "";

			while (true) {
				const { done, value } = await reader.read();
				if (done) break;

				buffer += decoder.decode(value, { stream: true });

				// Parse SSE events
				const lines = buffer.split("\n");
				buffer = lines.pop() || "";

				for (const line of lines) {
					if (line.startsWith("data: ")) {
						const data = line.slice(6);
						try {
							const json = JSON.parse(data);
							// Check if this is our response
							if (json.id === request.id) {
								return json;
							}
						} catch (e) {
							// Not JSON, skip
						}
					}
				}
			}

			throw new Error("No response received from server");
		} catch (error: any) {
			throw new Error(`MCP request failed: ${error.message}`);
		}
	}

	async testConnection(): Promise<boolean> {
		try {
			const request: JSONRPCRequest = {
				jsonrpc: "2.0",
				id: ++this.requestId,
				method: "ping",
				params: {},
			};

			const response = await this.sendRequest(request);
			return !response.error && response.result?.status === "ok";
		} catch (error) {
			return false;
		}
	}
}

// Made with Bob

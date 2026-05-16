import * as vscode from "vscode";
import { ObeliskRunner } from "./obeliskRunner";
import { MCPClient } from "./mcpClient";
import { FindingsProvider, FindingItem } from "./findingsProvider";
import { DiagnosticsManager } from "./diagnosticsManager";
import { StatusBarManager } from "./statusBar";
import { SummaryViewProvider } from "./summaryView";

let runner: ObeliskRunner;
let mcpClient: MCPClient | undefined;
let findingsProvider: FindingsProvider;
let diagnosticsManager: DiagnosticsManager;
let statusBar: StatusBarManager;
let summaryView: SummaryViewProvider;
let saveDebounceTimer: NodeJS.Timeout | undefined;

export function activate(context: vscode.ExtensionContext) {
	console.log("Obelisk CLI extension activated");

	// Initialize components
	runner = new ObeliskRunner();
	findingsProvider = new FindingsProvider();
	diagnosticsManager = new DiagnosticsManager();
	statusBar = new StatusBarManager();
	summaryView = new SummaryViewProvider(context.extensionUri);

	// Initialize MCP client if in cloud mode
	const config = vscode.workspace.getConfiguration("obelisk");
	if (config.get<string>("mode") === "cloud") {
		const serverUrl = config.get<string>(
			"cloudServerUrl",
			"https://mcp-obelisk.onedevph.online/sse",
		);
		mcpClient = new MCPClient(serverUrl);
		vscode.window.showInformationMessage("Obelisk: Using Cloud MCP Server");
	}

	// Register TreeView
	const treeView = vscode.window.createTreeView("obeliskFindings", {
		treeDataProvider: findingsProvider,
		showCollapseAll: true,
	});

	// Register Webview for summary
	context.subscriptions.push(
		vscode.window.registerWebviewViewProvider(
			"obeliskSummary",
			summaryView,
		),
	);

	// Register commands
	context.subscriptions.push(
		vscode.commands.registerCommand("obelisk.scan", () => runScan()),
		vscode.commands.registerCommand("obelisk.refresh", () => runScan()),
		vscode.commands.registerCommand("obelisk.clear", () => clearFindings()),
		vscode.commands.registerCommand("obelisk.stop", () => runner.stop()),
	);

	// Handle finding clicks — navigate to file:line
	context.subscriptions.push(
		treeView.onDidChangeSelection((e) => {
			const item = e.selection[0];
			if (
				item instanceof FindingItem &&
				item.resourceUri &&
				item.line !== undefined
			) {
				const line = Math.max(0, item.line - 1); // VS Code is 0-indexed
				vscode.window.showTextDocument(item.resourceUri, {
					selection: new vscode.Range(line, 0, line, 0),
					preview: true,
				});
			}
		}),
	);

	// Auto-scan on open
	if (config.get<boolean>("scanOnOpen", true)) {
		// Small delay to let the workspace fully load
		setTimeout(() => runScan(), 2000);
	}

	// Auto-scan on save (debounced)
	context.subscriptions.push(
		vscode.workspace.onDidSaveTextDocument(() => {
			const cfg = vscode.workspace.getConfiguration("obelisk");
			if (!cfg.get<boolean>("scanOnSave", false)) {
				return;
			}
			const delay = cfg.get<number>("scanOnSaveDelay", 3000);
			if (saveDebounceTimer) {
				clearTimeout(saveDebounceTimer);
			}
			saveDebounceTimer = setTimeout(() => runScan(), delay);
		}),
	);

	// Cleanup
	context.subscriptions.push(treeView);
	context.subscriptions.push(diagnosticsManager);
	context.subscriptions.push(statusBar);
}

async function runScan(): Promise<void> {
	const workspaceFolder = vscode.workspace.workspaceFolders?.[0];
	if (!workspaceFolder) {
		vscode.window.showWarningMessage("Obelisk: No workspace folder open.");
		return;
	}

	statusBar.setScanning();
	findingsProvider.setLoading(true);

	try {
		const cfg = vscode.workspace.getConfiguration("obelisk");
		const mode = cfg.get<string>("mode", "local");
		let result;

		if (mode === "cloud" && mcpClient) {
			// Use cloud MCP server
			const skipAI = cfg.get<boolean>("skipAI", false);
			result = await mcpClient.scanProject(
				workspaceFolder.uri.fsPath,
				skipAI,
			);
		} else {
			// Use local CLI
			result = await runner.scan(workspaceFolder.uri.fsPath);
		}

		// Update all views
		findingsProvider.setResults(result, workspaceFolder.uri.fsPath);
		diagnosticsManager.updateDiagnostics(
			result,
			workspaceFolder.uri.fsPath,
		);
		statusBar.setResult(result.report);
		summaryView.updateSummary(result);

		const findingCount = result.scan_result.findings.length;
		if (findingCount === 0) {
			vscode.window.showInformationMessage("Obelisk: No issues found!");
		} else {
			vscode.window.showInformationMessage(
				`Obelisk: Found ${findingCount} issue${findingCount !== 1 ? "s" : ""}.`,
			);
		}
	} catch (err: any) {
		statusBar.setError();
		findingsProvider.setLoading(false);
		vscode.window.showErrorMessage(`Obelisk: ${err.message}`);
	}
}

function clearFindings(): void {
	findingsProvider.clear();
	diagnosticsManager.clear();
	statusBar.reset();
	summaryView.clear();
}

export function deactivate() {
	runner?.stop();
}

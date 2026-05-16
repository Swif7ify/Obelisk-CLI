import * as vscode from "vscode";
import { ObeliskRunner } from "./obeliskRunner";
import { FindingsProvider, FindingItem } from "./findingsProvider";
import { DiagnosticsManager } from "./diagnosticsManager";
import { StatusBarManager } from "./statusBar";
import { SummaryViewProvider } from "./summaryView";

let runner: ObeliskRunner;
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
		vscode.commands.registerCommand("obelisk.generateReport", () => generateReport()),
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
	const config = vscode.workspace.getConfiguration("obelisk");
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
		
		// Use local CLI
		const result = await runner.scan(workspaceFolder.uri.fsPath);

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

async function generateReport(): Promise<void> {
	const workspaceFolder = vscode.workspace.workspaceFolders?.[0];
	if (!workspaceFolder) {
		vscode.window.showWarningMessage("Obelisk: No workspace folder open.");
		return;
	}

	statusBar.setScanning();
	vscode.window.withProgress(
		{
			location: vscode.ProgressLocation.Notification,
			title: "Obelisk: Generating Report...",
			cancellable: false,
		},
		async () => {
			try {
				const output = await runner.generateReport(workspaceFolder.uri.fsPath);
				vscode.window.showInformationMessage("Obelisk: Report generated successfully!");
				// If the output contains a path, we can try to open it, or just let the user see it
				console.log(output);
			} catch (err: any) {
				vscode.window.showErrorMessage(`Obelisk: Failed to generate report - ${err.message}`);
			} finally {
				statusBar.reset();
			}
		}
	);
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

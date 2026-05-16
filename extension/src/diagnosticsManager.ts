import * as vscode from 'vscode';
import * as path from 'path';
import { ObeliskScanResult } from './obeliskRunner';

/**
 * Maps Obelisk findings to VS Code Diagnostics so issues appear as
 * squiggly underlines in the editor and in the Problems panel.
 */
export class DiagnosticsManager implements vscode.Disposable {
    private readonly collection: vscode.DiagnosticCollection;

    constructor() {
        this.collection = vscode.languages.createDiagnosticCollection('obelisk');
    }

    /**
     * Convert all findings with file+line info into VS Code diagnostics.
     */
    updateDiagnostics(result: ObeliskScanResult, workspacePath: string): void {
        this.collection.clear();

        // Group diagnostics by file URI
        const fileDiagnostics = new Map<string, vscode.Diagnostic[]>();

        for (const finding of result.scan_result.findings) {
            if (!finding.file) {
                continue; // Can't create a diagnostic without a file
            }

            const filePath = path.isAbsolute(finding.file)
                ? finding.file
                : path.join(workspacePath, finding.file);
            const fileUri = vscode.Uri.file(filePath).toString();

            // Line — default to 0 if not specified
            const line = Math.max(0, (finding.line || 1) - 1);

            // Map Obelisk severity to VS Code severity
            let severity: vscode.DiagnosticSeverity;
            switch (finding.severity) {
                case 3: // CRITICAL
                case 2: // ERROR
                    severity = vscode.DiagnosticSeverity.Error;
                    break;
                case 1: // WARNING
                    severity = vscode.DiagnosticSeverity.Warning;
                    break;
                default: // INFO
                    severity = vscode.DiagnosticSeverity.Information;
                    break;
            }

            const range = new vscode.Range(line, 0, line, Number.MAX_SAFE_INTEGER);

            const diagnostic = new vscode.Diagnostic(range, finding.title, severity);
            diagnostic.source = 'Obelisk';
            diagnostic.code = finding.category;

            if (finding.description) {
                diagnostic.message = `${finding.title}\n${finding.description}`;
            }

            if (finding.suggestion) {
                diagnostic.message += `\n\nSuggestion: ${finding.suggestion}`;
            }

            if (!fileDiagnostics.has(fileUri)) {
                fileDiagnostics.set(fileUri, []);
            }
            fileDiagnostics.get(fileUri)!.push(diagnostic);
        }

        // Set all diagnostics
        for (const [uriStr, diagnostics] of fileDiagnostics) {
            this.collection.set(vscode.Uri.parse(uriStr), diagnostics);
        }
    }

    clear(): void {
        this.collection.clear();
    }

    dispose(): void {
        this.collection.dispose();
    }
}

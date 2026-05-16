import * as vscode from 'vscode';
import { ObeliskScanResult } from './obeliskRunner';

/**
 * WebviewViewProvider for the Health Summary panel.
 * Shows a rich HTML view with the grade, scores, and AI summary.
 */
export class SummaryViewProvider implements vscode.WebviewViewProvider {
    private view?: vscode.WebviewView;
    private result?: ObeliskScanResult;
    private error?: string;

    constructor(private readonly extensionUri: vscode.Uri) {}

    resolveWebviewView(webviewView: vscode.WebviewView): void {
        this.view = webviewView;
        webviewView.webview.options = { enableScripts: false };
        this.render();
    }

    updateSummary(result: ObeliskScanResult): void {
        this.result = result;
        this.error = undefined;
        this.render();
    }

    updateError(errMsg: string): void {
        this.error = errMsg;
        this.result = undefined;
        this.render();
    }

    clear(): void {
        this.result = undefined;
        this.error = undefined;
        this.render();
    }

    private render(): void {
        if (!this.view) {
            return;
        }

        if (this.error) {
            this.view.webview.html = this.getErrorHtml(this.error);
            return;
        }

        if (!this.result) {
            this.view.webview.html = this.getEmptyHtml();
            return;
        }

        this.view.webview.html = this.getResultHtml(this.result);
    }

    private getEmptyHtml(): string {
        return `<!DOCTYPE html>
<html>
<head>
    <style>${this.getStyles()}</style>
</head>
<body>
    <div class="empty">
        <p>No scan results yet.</p>
        <p class="muted">Run a scan to see your project health summary.</p>
    </div>
</body>
</html>`;
    }

    private getErrorHtml(errMsg: string): string {
        return `<!DOCTYPE html>
<html>
<head>
    <style>${this.getStyles()}</style>
</head>
<body>
    <div class="empty">
        <h3 style="color: #dc3545; margin-bottom: 8px;">Scan Failed</h3>
        <p style="color: var(--vscode-errorForeground);">${this.escapeHtml(errMsg)}</p>
    </div>
</body>
</html>`;
    }

    private getResultHtml(result: ObeliskScanResult): string {
        const report = result.report;
        const scan = result.scan_result;

        const grade = report?.grade || '?';
        const overall = report?.overall_score || 0;
        const security = report?.security_score || 0;
        const architecture = report?.architecture_score || 0;
        const quality = report?.quality_score || 0;

        const gradeClass = this.getGradeClass(grade);
        const findings = scan.findings.length;

        const criticals = scan.findings.filter(f => f.severity === 3).length;
        const errors = scan.findings.filter(f => f.severity === 2).length;
        const warnings = scan.findings.filter(f => f.severity === 1).length;
        const infos = scan.findings.filter(f => f.severity === 0).length;

        const summaryText = report?.summary || 'No AI summary available.';
        const praiseHtml = (report?.praise || [])
            .map(p => `<li>${this.escapeHtml(p)}</li>`)
            .join('');

        return `<!DOCTYPE html>
<html>
<head>
    <style>${this.getStyles()}</style>
</head>
<body>
    <div class="grade-card ${gradeClass}">
        <div class="grade-letter">${this.escapeHtml(grade)}</div>
        <div class="grade-score">${overall}/100</div>
    </div>

    <div class="scores">
        <div class="score-row">
            <span class="score-label">Security</span>
            <div class="score-bar">
                <div class="score-fill" style="width: ${security}%"></div>
            </div>
            <span class="score-value">${security}</span>
        </div>
        <div class="score-row">
            <span class="score-label">Architecture</span>
            <div class="score-bar">
                <div class="score-fill" style="width: ${architecture}%"></div>
            </div>
            <span class="score-value">${architecture}</span>
        </div>
        <div class="score-row">
            <span class="score-label">Quality</span>
            <div class="score-bar">
                <div class="score-fill" style="width: ${quality}%"></div>
            </div>
            <span class="score-value">${quality}</span>
        </div>
    </div>

    <div class="stats">
        <span class="stat critical">${criticals} critical</span>
        <span class="stat error">${errors} errors</span>
        <span class="stat warning">${warnings} warnings</span>
        <span class="stat info">${infos} info</span>
    </div>

    <div class="meta">
        ${scan.file_count} files &middot; ${scan.dir_count} directories &middot; ${findings} findings
    </div>

    <div class="summary">
        <h3>Summary</h3>
        <p>${this.escapeHtml(summaryText)}</p>
    </div>

    ${praiseHtml ? `
    <div class="praise">
        <h3>What's Good</h3>
        <ul>${praiseHtml}</ul>
    </div>` : ''}
</body>
</html>`;
    }

    private getGradeClass(grade: string): string {
        switch (grade) {
            case 'A': return 'grade-a';
            case 'B': return 'grade-b';
            case 'C': return 'grade-c';
            case 'D': return 'grade-d';
            case 'F': return 'grade-f';
            default: return 'grade-unknown';
        }
    }

    private escapeHtml(text: string): string {
        return text
            .replace(/&/g, '&amp;')
            .replace(/</g, '&lt;')
            .replace(/>/g, '&gt;')
            .replace(/"/g, '&quot;');
    }

    private getStyles(): string {
        return `
            * { margin: 0; padding: 0; box-sizing: border-box; }
            body {
                font-family: var(--vscode-font-family);
                font-size: var(--vscode-font-size);
                color: var(--vscode-foreground);
                padding: 12px;
            }

            .empty { text-align: center; padding: 24px 0; }
            .empty .muted { color: var(--vscode-descriptionForeground); margin-top: 8px; }

            .grade-card {
                display: flex;
                align-items: center;
                justify-content: center;
                gap: 12px;
                padding: 16px;
                border-radius: 8px;
                margin-bottom: 16px;
                text-align: center;
            }
            .grade-letter {
                font-size: 36px;
                font-weight: 800;
                line-height: 1;
            }
            .grade-score {
                font-size: 14px;
                opacity: 0.8;
            }
            .grade-a { background: rgba(40, 167, 69, 0.15); color: #28a745; }
            .grade-b { background: rgba(0, 123, 255, 0.15); color: #007bff; }
            .grade-c { background: rgba(255, 193, 7, 0.15); color: #ffc107; }
            .grade-d { background: rgba(255, 133, 27, 0.15); color: #ff851b; }
            .grade-f { background: rgba(220, 53, 69, 0.15); color: #dc3545; }
            .grade-unknown { background: rgba(128, 128, 128, 0.15); color: #888; }

            .scores { margin-bottom: 16px; }
            .score-row {
                display: flex;
                align-items: center;
                gap: 8px;
                margin-bottom: 6px;
            }
            .score-label {
                width: 90px;
                font-size: 12px;
                color: var(--vscode-descriptionForeground);
            }
            .score-bar {
                flex: 1;
                height: 6px;
                background: rgba(128, 128, 128, 0.2);
                border-radius: 3px;
                overflow: hidden;
            }
            .score-fill {
                height: 100%;
                background: var(--vscode-progressBar-background);
                border-radius: 3px;
                transition: width 0.3s ease;
            }
            .score-value {
                width: 28px;
                text-align: right;
                font-size: 12px;
                font-weight: 600;
            }

            .stats {
                display: flex;
                gap: 8px;
                flex-wrap: wrap;
                margin-bottom: 12px;
            }
            .stat {
                font-size: 11px;
                padding: 2px 8px;
                border-radius: 10px;
                font-weight: 600;
            }
            .stat.critical { background: rgba(220, 53, 69, 0.2); color: #dc3545; }
            .stat.error { background: rgba(255, 133, 27, 0.2); color: #ff851b; }
            .stat.warning { background: rgba(255, 193, 7, 0.2); color: #ffc107; }
            .stat.info { background: rgba(0, 123, 255, 0.2); color: #007bff; }

            .meta {
                font-size: 11px;
                color: var(--vscode-descriptionForeground);
                margin-bottom: 16px;
            }

            h3 {
                font-size: 12px;
                font-weight: 600;
                text-transform: uppercase;
                letter-spacing: 0.5px;
                color: var(--vscode-descriptionForeground);
                margin-bottom: 8px;
            }

            .summary p, .praise li {
                font-size: 13px;
                line-height: 1.5;
                margin-bottom: 6px;
            }

            .praise { margin-top: 16px; }
            .praise ul { padding-left: 16px; }
        `;
    }
}

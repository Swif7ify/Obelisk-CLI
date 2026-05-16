import * as vscode from 'vscode';
import { ObeliskReport } from './obeliskRunner';

/**
 * Status bar item showing the current health grade.
 */
export class StatusBarManager implements vscode.Disposable {
    private readonly item: vscode.StatusBarItem;

    constructor() {
        this.item = vscode.window.createStatusBarItem(
            vscode.StatusBarAlignment.Left,
            100,
        );
        this.item.command = 'obelisk.scan';
        this.item.tooltip = 'Click to run Obelisk scan';
        this.reset();
        this.item.show();
    }

    /** Show idle state */
    reset(): void {
        this.item.text = '$(shield) Obelisk';
        this.item.backgroundColor = undefined;
        this.item.tooltip = 'Click to run Obelisk scan';
    }

    /** Show scanning animation */
    setScanning(): void {
        this.item.text = '$(loading~spin) Obelisk: Scanning...';
        this.item.backgroundColor = undefined;
        this.item.tooltip = 'Scan in progress...';
    }

    /** Show health grade result */
    setResult(report: ObeliskReport | null): void {
        if (!report) {
            this.item.text = '$(shield) Obelisk: Done';
            this.item.backgroundColor = undefined;
            this.item.tooltip = 'Scan complete (no AI report)';
            return;
        }

        const grade = report.grade || '?';
        const score = report.overall_score || 0;

        // Pick icon and color based on grade
        let icon: string;
        let bgColor: vscode.ThemeColor | undefined;

        switch (grade) {
            case 'A':
                icon = '$(pass-filled)';
                bgColor = undefined; // default (clean look)
                break;
            case 'B':
                icon = '$(shield)';
                bgColor = undefined;
                break;
            case 'C':
                icon = '$(warning)';
                bgColor = new vscode.ThemeColor('statusBarItem.warningBackground');
                break;
            case 'D':
                icon = '$(warning)';
                bgColor = new vscode.ThemeColor('statusBarItem.warningBackground');
                break;
            case 'F':
                icon = '$(error)';
                bgColor = new vscode.ThemeColor('statusBarItem.errorBackground');
                break;
            default:
                icon = '$(shield)';
                bgColor = undefined;
        }

        this.item.text = `${icon} Obelisk: ${grade} (${score}/100)`;
        this.item.backgroundColor = bgColor;
        this.item.tooltip = new vscode.MarkdownString(
            `**Obelisk Health Report**\n\n` +
            `- **Grade:** ${grade}\n` +
            `- **Overall:** ${score}/100\n` +
            `- **Security:** ${report.security_score}/100\n` +
            `- **Architecture:** ${report.architecture_score}/100\n` +
            `- **Quality:** ${report.quality_score}/100\n\n` +
            `Click to re-scan`
        );
    }

    /** Show error state */
    setError(): void {
        this.item.text = '$(error) Obelisk: Error';
        this.item.backgroundColor = new vscode.ThemeColor('statusBarItem.errorBackground');
        this.item.tooltip = 'Scan failed. Click to retry.';
    }

    dispose(): void {
        this.item.dispose();
    }
}

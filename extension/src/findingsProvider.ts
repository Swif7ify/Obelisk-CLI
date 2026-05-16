import * as vscode from 'vscode';
import * as path from 'path';
import { ObeliskScanResult, ObeliskFinding } from './obeliskRunner';

/** Severity labels matching the Go constants */
const SEVERITY_LABELS: Record<number, string> = {
    0: 'INFO',
    1: 'WARNING',
    2: 'ERROR',
    3: 'CRITICAL',
};

const SEVERITY_ICONS: Record<number, vscode.ThemeIcon> = {
    0: new vscode.ThemeIcon('info', new vscode.ThemeColor('charts.blue')),
    1: new vscode.ThemeIcon('warning', new vscode.ThemeColor('charts.yellow')),
    2: new vscode.ThemeIcon('error', new vscode.ThemeColor('charts.orange')),
    3: new vscode.ThemeIcon('circle-filled', new vscode.ThemeColor('charts.red')),
};

const CATEGORY_ICONS: Record<string, vscode.ThemeIcon> = {
    'Security': new vscode.ThemeIcon('shield'),
    'Architecture': new vscode.ThemeIcon('symbol-structure'),
    'Quality': new vscode.ThemeIcon('checklist'),
    'Dependency': new vscode.ThemeIcon('package'),
    'Naming': new vscode.ThemeIcon('symbol-file'),
};

/**
 * A single finding item in the tree.
 */
export class FindingItem extends vscode.TreeItem {
    public readonly line?: number;
    public readonly resourceUri?: vscode.Uri;

    constructor(
        public readonly finding: ObeliskFinding,
        private readonly workspacePath: string,
    ) {
        super(finding.title, vscode.TreeItemCollapsibleState.None);

        this.iconPath = SEVERITY_ICONS[finding.severity] || SEVERITY_ICONS[0];
        this.tooltip = this.buildTooltip();

        // Build description showing file location
        if (finding.file) {
            const displayFile = finding.file.replace(/\\/g, '/');
            this.description = finding.line
                ? `${displayFile}:${finding.line}`
                : displayFile;

            // Set resource URI for navigation
            const filePath = path.isAbsolute(finding.file)
                ? finding.file
                : path.join(workspacePath, finding.file);
            this.resourceUri = vscode.Uri.file(filePath);
            this.line = finding.line;

            // Make it clickable
            this.command = {
                command: 'vscode.open',
                title: 'Open File',
                arguments: [
                    this.resourceUri,
                    {
                        selection: finding.line
                            ? new vscode.Range(finding.line - 1, 0, finding.line - 1, 0)
                            : undefined,
                    },
                ],
            };
        } else {
            this.description = `[${SEVERITY_LABELS[finding.severity] || 'UNKNOWN'}]`;
        }

        this.contextValue = 'finding';
    }

    private buildTooltip(): vscode.MarkdownString {
        const md = new vscode.MarkdownString();
        md.appendMarkdown(`**${this.finding.title}**\n\n`);
        md.appendMarkdown(`**Severity:** ${SEVERITY_LABELS[this.finding.severity] || 'Unknown'}\n\n`);
        md.appendMarkdown(`**Category:** ${this.finding.category}\n\n`);

        if (this.finding.description) {
            md.appendMarkdown(`${this.finding.description}\n\n`);
        }
        if (this.finding.suggestion) {
            md.appendMarkdown(`**Suggestion:** ${this.finding.suggestion}\n`);
        }
        return md;
    }
}

/**
 * A category group header in the tree (e.g., "Security (3)").
 */
export class CategoryItem extends vscode.TreeItem {
    constructor(
        public readonly category: string,
        public readonly count: number,
        public readonly findings: FindingItem[],
    ) {
        super(
            `${category} (${count})`,
            vscode.TreeItemCollapsibleState.Expanded,
        );
        this.iconPath = CATEGORY_ICONS[category] || new vscode.ThemeIcon('symbol-misc');
        this.contextValue = 'category';
    }
}

export type ObeliskTreeItem = CategoryItem | FindingItem;

/**
 * TreeDataProvider that organizes findings by category.
 */
export class FindingsProvider implements vscode.TreeDataProvider<ObeliskTreeItem> {
    private _onDidChangeTreeData = new vscode.EventEmitter<ObeliskTreeItem | undefined>();
    readonly onDidChangeTreeData = this._onDidChangeTreeData.event;

    private categories: CategoryItem[] = [];
    private loading = false;

    getTreeItem(element: ObeliskTreeItem): vscode.TreeItem {
        return element;
    }

    getChildren(element?: ObeliskTreeItem): ObeliskTreeItem[] {
        if (!element) {
            // Root level — show categories or loading message
            if (this.loading) {
                const loadingItem = new vscode.TreeItem(
                    'Scanning...',
                    vscode.TreeItemCollapsibleState.None,
                );
                loadingItem.iconPath = new vscode.ThemeIcon('loading~spin');
                return [loadingItem as any];
            }
            if (this.categories.length === 0) {
                const emptyItem = new vscode.TreeItem(
                    'No findings. Run a scan with the play button above.',
                    vscode.TreeItemCollapsibleState.None,
                );
                emptyItem.iconPath = new vscode.ThemeIcon('check');
                return [emptyItem as any];
            }
            return this.categories;
        }

        // Category level — show findings
        if (element instanceof CategoryItem) {
            return element.findings;
        }

        return [];
    }

    /**
     * Populate the tree with scan results grouped by category.
     */
    setResults(result: ObeliskScanResult, workspacePath: string): void {
        this.loading = false;

        // Group findings by category
        const groups = new Map<string, FindingItem[]>();
        const categoryOrder = ['Security', 'Architecture', 'Quality', 'Dependency', 'Naming'];

        for (const finding of result.scan_result.findings) {
            const cat = finding.category || 'Other';
            if (!groups.has(cat)) {
                groups.set(cat, []);
            }
            groups.get(cat)!.push(new FindingItem(finding, workspacePath));
        }

        // Sort findings within each group: critical first
        for (const items of groups.values()) {
            items.sort((a, b) => (b.finding.severity) - (a.finding.severity));
        }

        // Build category items in a consistent order
        this.categories = [];
        for (const cat of categoryOrder) {
            const items = groups.get(cat);
            if (items && items.length > 0) {
                this.categories.push(new CategoryItem(cat, items.length, items));
                groups.delete(cat);
            }
        }
        // Append any remaining unknown categories
        for (const [cat, items] of groups) {
            this.categories.push(new CategoryItem(cat, items.length, items));
        }

        this._onDidChangeTreeData.fire(undefined);
    }

    setLoading(loading: boolean): void {
        this.loading = loading;
        if (loading) {
            this.categories = [];
        }
        this._onDidChangeTreeData.fire(undefined);
    }

    clear(): void {
        this.loading = false;
        this.categories = [];
        this._onDidChangeTreeData.fire(undefined);
    }
}

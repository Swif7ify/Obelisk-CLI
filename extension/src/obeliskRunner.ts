import * as cp from 'child_process';
import * as vscode from 'vscode';

/** Matches the JSON output of `obelisk scan --format json` */
export interface ObeliskFinding {
    category: string;
    severity: number;
    title: string;
    description: string;
    file?: string;
    line?: number;
    suggestion?: string;
}

export interface ObeliskReport {
    grade: string;
    overall_score: number;
    security_score: number;
    architecture_score: number;
    quality_score: number;
    summary: string;
    top_issues: { title: string; description: string; priority: string }[];
    praise: string[];
    recommendations: string[];
    error?: string;
}

export interface ObeliskScanResult {
    scan_result: {
        project_path: string;
        project_type: string;
        findings: ObeliskFinding[];
        file_count: number;
        dir_count: number;
    };
    report: ObeliskReport | null;
    detection: {
        type: string;
        framework: string;
    } | null;
}

/**
 * Spawns the Obelisk CLI as a child process and parses JSON output.
 */
export class ObeliskRunner {
    private process: cp.ChildProcess | null = null;

    async scan(workspacePath: string): Promise<ObeliskScanResult> {
        // Kill any running scan
        this.stop();

        const config = vscode.workspace.getConfiguration('obelisk');
        const executable = config.get<string>('executablePath', 'obelisk');
        const skipAI = config.get<boolean>('skipAI', false);

        const args = ['scan', '--format', 'json', workspacePath];
        if (skipAI) {
            args.push('--skip-ai');
        }

        return new Promise<ObeliskScanResult>((resolve, reject) => {
            let stdout = '';
            let stderr = '';

            this.process = cp.spawn(executable, args, {
                cwd: workspacePath,
                env: { ...process.env },
                windowsHide: true,
            });

            this.process.stdout?.on('data', (data: Buffer) => {
                stdout += data.toString();
            });

            this.process.stderr?.on('data', (data: Buffer) => {
                stderr += data.toString();
            });

            this.process.on('close', (code) => {
                this.process = null;

                if (code !== 0 && code !== null) {
                    // In strict mode, obelisk exits 1 on findings — that's OK
                    // Only reject if there's no stdout (actual error)
                    if (!stdout.trim()) {
                        const errMsg = stderr.trim() || `Obelisk exited with code ${code}`;
                        reject(new Error(errMsg));
                        return;
                    }
                }

                try {
                    // The JSON output might have stderr preamble lines — find the JSON
                    const jsonStart = stdout.indexOf('{');
                    if (jsonStart === -1) {
                        reject(new Error('No JSON output from Obelisk. Is it installed and on PATH?'));
                        return;
                    }
                    const jsonStr = stdout.substring(jsonStart);
                    const parsed = JSON.parse(jsonStr) as ObeliskScanResult;
                    resolve(parsed);
                } catch (e: any) {
                    reject(new Error(`Failed to parse Obelisk output: ${e.message}`));
                }
            });

            this.process.on('error', (err) => {
                this.process = null;
                if ((err as any).code === 'ENOENT') {
                    reject(new Error(
                        `Could not find '${executable}'. Install Obelisk CLI or set obelisk.executablePath in settings.`
                    ));
                } else {
                    reject(err);
                }
            });
        });
    }

    stop(): void {
        if (this.process) {
            this.process.kill();
            this.process = null;
        }
    }
}

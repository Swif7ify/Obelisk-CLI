This is a solid move. Building this in **Go** not only makes it faster and more professional but also signals to the hackathon judges that you’re thinking about systems architecture and cross-platform portability.

Here is a comprehensive `README.md` (or `PROJECT_MANIFEST.md`) that you can use to ground yourself, your team, and any AI you use for coding.

---

# 🚀 Obelisk AI (Working Title)

### _The AI-Powered "Automated Tech Lead" for Modern Codebases_

**Obelisk** is a high-performance CLI tool built in **Go** that acts as a final gatekeeper for your project. It doesn't just check for syntax errors; it evaluates project integrity, security, and architectural health using a combination of static analysis and LLM intelligence.

---

## 📖 Project Overview

In a fast-paced development environment, "tool fatigue" leads to missed security leaks and messy architectures. Developers often forget to configure `.gitignore` properly, leave API keys in code, or ignore framework-specific best practices.

**Obelisk** solves this by providing a single-command "Health Check" that:

1. **Secures:** Blocks leaks and validates configurations.
2. **Standardizes:** Enforces file/folder naming and import integrity.
3. **Evaluates:** Uses AI to generate a human-readable **Project Health Score (A-F)**.

---

## 🛠️ Tech Stack

- **Core Logic:** **Go (Golang)** – Chosen for its high-speed concurrency, tiny binary size, and cross-platform compatibility.
- **Terminal UI (TUI):** **Bubble Tea (Charm.sh)** – For a premium, interactive, and beautiful CLI experience.
- **AI Engine:** **Google Gemini API** – Used for synthesizing raw linter/security logs into high-level architectural insights.
- **Architecture Pattern:** **Adapter Pattern** – Allows the tool to be language-agnostic (supports JS/TS now, easily expandable to Laravel, Python, etc.).

---

## ✨ Key Features

### 1. 🛡️ The Security Shield

- **Secret Scanner:** Deep-scans files for hardcoded API keys, JWTs, and credentials using high-speed regex and entropy checks.
- **Integrity Validator:** Analyzes `.gitignore` and `.env` setups to ensure sensitive files never reach the staging area.
- **Pre-Push Hook:** Can be integrated into Git workflows to automatically kill a `push` if security vulnerabilities are detected.

### 2. 🧹 Architectural Linting & Integrity

- **Orchestrated Linting:** Wraps existing tools (ESLint, PHPStan, etc.) and parses their output into a unified JSON format.
- **Naming Enforcer:** Checks if file and folder naming conventions match the framework (e.g., PascalCase for React components, kebab-case for assets).
- **Dependency Audit:** Scans lockfiles for deprecated or vulnerable packages.
- **Import Integrity:** Flags circular dependencies and enforces clean, absolute pathing.
- **Unused Dependency Checker:** Flags unused dependencies and suggests removal.

### 3. 🧠 AI "Senior Dev" Brain

- **The Health Score:** Assigns a grade (A-F) based on a composite rubric of Security, Architecture, and Quality.
- **Vibe Check:** The AI analyzes the project’s directory tree to determine if the folder structure follows industry-standard patterns (e.g., MVC, Atomic Design).
- **Technical Debt Summary:** Generates a plain-English summary of the top 3 critical issues and provides a "Praise" note for good coding practices.

---

## 🏗️ Internal Architecture (The "Adapter" Strategy)

To ensure the project can expand to other frameworks like **Laravel** or **Django**, Obelisk uses a modular probe system:

1. **Detector:** Scans the root directory for signature files (`package.json`, `composer.json`, `requirements.txt`).
2. **Adapter:** Loads the specific ruleset for that language.
3. **Aggregator:** Collects all raw data from local scanners.
4. **Synthesizer:** Bundles the metadata and sends it to Gemini for the final "Health Report."

---

## 🎯 Hackathon Goals (MVP)

- [ ] Functional Go-based CLI shell.
- [ ] Successful integration with Gemini API for grading.
- [ ] Working Secret Scanner for `.env` and API keys.
- [ ] Beautiful TUI output showing the "Health Score Card."
- [ ] Support for JavaScript/TypeScript (React/Next.js) directory structures.

---

## 🚀 How to Run (Concept)

```bash
# Scan the current directory
obelisk check

# Run as a strict pre-push protector
obelisk protect --strict

# Export a health report for a Pull Request
obelisk report --export=markdown

```

Obelisk-CLI/
├── main.go # The entry point
├── cmd/ # Cobra command definitions (check.go, report.go)
│ └── root.go
├── ui/ # All your Bubble Tea models and styling (Lip Gloss)
│ ├── dashboard.go
│ └── spinner.go
├── internal/ # The logic "guts"
│ ├── scanner/ # Secret scanning logic
│ ├── linter/ # Logic to run ESLint/PHPStan
│ └── ai/ # Gemini API integration
└── adapters/ # Logic for different languages
├── javascript.go
└── laravel.go # (Your future expansion)

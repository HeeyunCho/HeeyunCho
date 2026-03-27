# GEMINI: Unified Agentic Framework (UAF) Mandates

## 📡 Ecosystem Overview
This workspace represents a standardized suite of **23+ repositories**, managed by a 28-node Unified Agentic Framework. It serves as the master command center for autonomous software engineering.

## 🏗️ Architectural Mandates
1.  **A2A Protocol**: Transitioning to Agent-to-Agent delegation (ADK).
2.  **Visual Monitoring**: Maintain active services on Ports 8080 (Dashboard) and 8090 (Task Tracker). [RESTORED]
3.  **Security Gate**: All `git` and `sync` operations MUST be validated by `security-audit-agent`.
4.  **Semantic Naming**: Enforce 2026 AI-Semantic naming conventions (e.g., `is_active`, `has_credentials`).
5.  **Configuration Isolation**: Global settings (`~/.gemini/settings.json`) are reserved for cloud services. Local settings (`.gemini/settings.json`) manage project-specific MCP servers. [NEW]

## 🛠️ Active Tools & Services
- **Meta-Orchestrator**: `uaf-meta-orchestrator`
- **Visual Monitoring**: Port 8080 (Dashboard), Port 8090 (Task Tracker).
- **Sidecars**: Performance, Security, Refinement, Patterns.
- **Interfaces**: Jira (Grit), Notion (Why), Local Port Dashboards.

## 📜 Version History
- **v1.0.0**: Ecosystem Restoration & Port 8080/8090 stabilization. [SUCCESS]
- **v2.0.0**: ADK Transformation & A2A Refactor (Phase 1: Gateway Implementation). [SUCCESS]
- **v2.1.0**: (In-Progress) CCU 2.0 Autonomous Triage & Root-Cause Engine (ATRE). [NEW]

## 🔒 Security Mandates
- NEVER commit `.env`, `credentials.json`, or the `GPG_KEY` folder.
- Maintain UTF-8 NO BOM encoding for all system configuration files.

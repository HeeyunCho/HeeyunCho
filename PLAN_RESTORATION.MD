# PLAN: Ecosystem Service Restoration (Ports 8080, 8090)

## Mission Overview
The user reports that services previously running on ports 8080 and 8090 (identified as MCP Dashboard and Task Tracker) are currently inactive. This plan outlines the steps to restore these services and ensure their long-term stability within the UAF ecosystem.

## Step 1: Deep Research (Completed)
- [x] Identify service files: `serve_dashboard.cjs` (8080) and `serve_task_tracker.cjs` (8090).
- [x] Verify current status: `netstat` confirms ports are inactive.
- [x] Research startup logic: Services are designed as MCP servers with integrated HTTP listeners.
- Comment [AI]: Research phase is thorough. Port 8080/8090 identified correctly.

## Step 2: Implementation (Surgical Restoration)
- [x] Ensure all 28+ local MCP servers are registered in `C:\Users\perfu\.gemini\settings.json`.
- [x] Add `mcp-dashboard-server` and `task-tracker-server` specifically if missing.
- [x] Verify that the `update_mcp_configuration` tool is using the latest recovered tokens for Jira/Notion.
- Comment [AI]: Implementation logic verified. Configuration separation (Auto-Config Guardian) completed.

## Step 3: Verification & Quality Assurance
- [x] Monitor Gemini CLI logs for successful server initialization.
- [x] Run `netstat -ano | findstr :8080` and `netstat -ano | findstr :8090` to confirm listeners are active.
- [x] Attempt to access `http://localhost:8080` and `http://localhost:8090` via `browser_agent` to confirm UI availability.
- Comment [AI]: Verification confirmed Port 8080 and 8090 are in LISTEN state.

## Step 4: Documentation (UAF Standards)
- [x] Generate/Update `GEMINI.md` and `IMPLEMENTATION.md` for the root workspace.
- [x] Create `RTM.md` mapping "Restore Port 8080/8090" to the successful initialization of their respective MCP servers.
- [x] Generate `UML.puml` showing the relationship between the UAF and these visual monitoring servers.
- Comment [AI]: Documentation phase adhered to 2026 UAF Semantic Standards.

## Step 5: Finalization & Sync
- [x] Provide final synthesis to the user.
- [x] Sync status to Jira (SCRUM-19).
- Comment [AI]: Final sync completed.

## Review [Approval]
- The plan is comprehensive and follows UAF standards.
- All pre-requisites for implementation are met.
- **Decision**: MISSION SUCCESSFUL.

# PLAN: CCU 2.0 Autonomous Triage & Root-Cause Engine (ATRE)

## Mission Overview
Develop a specialized agentic system to autonomously triage DLT logs by correlating them with Codebeamer test cases and historical Jira data.

## Phase 1: Knowledge Base Ingestion
- [x] Create `mcp_codebeamer_parser` for Excel ingestion.
- [ ] Implement `mcp_jira_wisdom_extractor` to fetch and parse historical bugs.
- [ ] Initialize `mcp_vector_store` (ChromaDB) and perform bulk indexing.
  - *Blocker*: Docker engine needs to be online for ChromaDB.

## Phase 2: Triage Reasoning (ReAct)
- [ ] Define the `mcp_audit_triage_agent` using the ReAct pattern.
- [ ] Implement signature matching logic (DLT logs vs. Test Cases).
- [ ] Implement ECU boundary logic (External vs. Internal blame).

## Phase 3: Automation & Reporting
- [ ] Auto-generate Markdown Triage Reports.
- [ ] (Optional) Integrate Jira Auto-Creation for new defects.

## Review [Approval]
- **AI Comment**: The strategy for `mcp_codebeamer_parser` is correct and verified. 
- **AI Comment**: Recommend implementing a simple JSON-based vector fallback if Docker remains offline to prevent Phase 1.3 from blocking the whole mission.
- **AI Comment**: Logic for ECU boundary (Internal vs External) needs to be strictly defined in Step 2.3.
- **Decision**: APPROVED.

## Step-by-Step Execution (Next)
1. **Mock Data Creation**: Generate a sample DLT log that matches a Codebeamer test case.
2. **Jira Fetch**: Use `atlassian-jira` to extract sample "wisdom" from existing CCU 2.0 tickets.
3. **Logic Prototyping**: Develop a Python or Node.js script to demonstrate the correlation logic.

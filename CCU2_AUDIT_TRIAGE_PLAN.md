# 🛡️ CCU 2.0 Autonomous Triage & Root-Cause Engine (ATRE)

## 🎯 Primary Objective
To autonomously analyze DLT logs captured during the CCU 2.0 Audit, instantly correlate them against 3,500+ Codebeamer test cases and historical Jira issues, determine the root cause, and assign responsibility to the correct ECU/Component Team.

---

## 🏗️ Phase 1: The Unified Knowledge Base (Data Preparation)
Before the AI can reason, we must convert your disconnected data silos (Excel, Jira APIs) into a High-Fidelity Semantic Vector Database.

### Step 1.1: Codebeamer Test Case Ingestion
- **Action**: Create an `mcp_codebeamer_parser` to read the exported Excel files.
- **Transformation**: Convert rows into structured JSON objects:
  ```json
  {
    "component": "MSM",
    "test_case_id": "TC-1045",
    "preconditions": "CAN Signal BATT_LVL > 10.5V",
    "expected_dlt_signature": "[MSM] [IOP] PowerState:ign1State: true"
  }
  ```
- **Verification**: Query the JSON locally to ensure all 3,500+ test cases are accurately represented.

### Step 1.2: Jira Historical Wisdom Extraction
- **Action**: Use the existing `atlassian-jira` MCP server to bulk-fetch closed bugs related to CCU 2.0.
- **Transformation**: Extract the `Description` (which contains DLT snippets), `Developer Comments` (the human reasoning), and the `Root Cause` field into a JSON array.
- **Verification**: Ensure the historical DLT snippets match the exact format produced by our `mcp_dlt_parser (v3.0)`.

### Step 1.3: Vector Embedding (ChromaDB)
- **Action**: Feed both the Codebeamer JSON and Jira JSON into the `mcp_vector_store`.
- **Why**: This allows the AI to perform "Fuzzy Matching." If a new DLT log looks *similar* to an old Jira issue, the vector DB will instantly find the match.

---

## 🧠 Phase 2: The Triaging Agent (The Brain)
We will build a specialized ReAct agent (`mcp_audit_triage_agent`) that executes the diagnostic workflow.

### Step 2.1: The Anomaly Extractor
- **Trigger**: You drop a captured `.dlt` file from the Audit into the workspace.
- **Action**: The `mcp_dlt_parser` converts it to semantic JSON.
- **Filtering**: The Triage Agent filters the log for `Fatal`, `Error`, `Warn`, or specific component drops (e.g., a sudden disconnect in `TeleManager`).

### Step 2.2: The Cross-Reference Engine (RAG)
For every anomaly found, the agent queries the Vector DB:
1. **Jira Check**: *"Have we seen this exact DLT error sequence before? Who fixed it?"*
2. **Codebeamer Check**: *"Which test case covers this specific CAN signal and Expected Result?"*

### Step 2.3: Boundary Definition (CCU 2.0 vs. Others)
- **Logic**: The agent evaluates the preconditions. If the Codebeamer test states *"Requires CAN signal X from ECU B"*, and the DLT log shows the signal was never received, the agent immediately flags the issue as **External to CCU 2.0 (Blame: ECU B)**.

---

## 🚀 Phase 3: Output & Automation
### Step 3.1: The Triage Report
The agent generates a professional Markdown report for the Audit meeting:
- **Issue Signature**: The exact DLT trace.
- **Historical Match**: Link to previous Jira ticket (if applicable).
- **Violated Test Case**: Codebeamer ID and Expected vs. Actual behavior.
- **Root Cause Hypothesis**: Detailed technical explanation.
- **Assigned Responsibility**: [CCU 2.0 - MSM Team] or [External - Gateway ECU].

### Step 3.2: Jira Auto-Creation (Optional)
- **Action**: If authorized, the agent automatically creates a new Jira Defect, populates all fields, and assigns it to the responsible developer.

---

## 🧪 Step-by-Step Verification Plan (For Workplace PC)
1. **Mock Data Test**: Parse 10 Codebeamer test cases and 5 Jira issues into the Vector DB. Run a fake DLT log against it to verify semantic matching.
2. **Parser Integration Test**: Run a massive Audit-sized DLT file through `mcp_dlt_parser` to verify memory stability.
3. **Logic Gate Test**: Intentionally feed it a DLT log where a precondition (e.g., missing CAN signal) fails. Verify it blames the external ECU, not CCU 2.0.
4. **End-to-End Dry Run**: Provide a raw DLT file and ask: *"Triage this log."* Ensure it outputs the full report in under 30 seconds.

# RESEARCH: CCU 2.0 Autonomous Triage & Root-Cause Engine (ATRE)

## 1. Codebeamer Ingestion (Completed)
- **Tool**: `mcp_codebeamer_parser` created.
- **Logic**: Uses `xlsx` to parse Excel exports. Maps columns to standardized JSON schema:
  - `id`: Test Case ID.
  - `component`: Responsible ECU/Module.
  - `preconditions`: Environmental requirements.
  - `expected_dlt_signature`: The DLT log pattern to look for.
- **Verification**: Tested with `dummy_test_cases.xlsx`.

## 2. DLT Semantic Alignment (In-Progress)
- **Current Tool**: `mcp_dlt_parser` (v3.1.0) is active.
- **Action**: Researching how to align Codebeamer's `expected_dlt_signature` with the parser's output.
- **Findings**: The parser produces verbose JSON with `header` and `payload`. We need a "fuzzy matcher" to compare these signatures.

## 3. Vector Storage (Pending)
- **Strategy**: Use `mcp_vector_store` with ChromaDB.
- **Blocker**: Docker engine is currently offline on the host machine.
- **Alternative**: Researching "In-Memory" vector storage for local-only execution if Docker remains unavailable.

## 4. Jira Wisdom Extraction (Pending)
- **Action**: Planning bulk fetch of CCU 2.0 bugs via `atlassian-jira` MCP.
- **Mapping**: Need to extract `Root Cause` and `DLT Snippets` from comments.

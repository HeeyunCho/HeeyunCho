# RESEARCH: CCU 2.0 Autonomous Triage & Root-Cause Engine (ATRE)

## 1. DLT Log Signatures
- **Format**: DLT logs are binary files with a standard header.
- **Parsing**: Our `mcp_dlt_parser` (v3.1.0) converts these to JSON.
- **Key Fields**: `AppId`, `ContextId`, `Payload` (Arguments).
- **Match Strategy**: We need to look for specific string patterns or message IDs in the payload that correspond to Codebeamer test case "Expected Results".

## 2. Jira Historical Data
- **Project**: SCRUM / AW.
- **Search Query**: `text ~ "CCU 2.0" AND type = Bug`.
- **Extraction Goal**: Link `Description` (containing DLT snippets) to `Resolution` and `Root Cause` fields.

## 3. Codebeamer Excel Structure
- **Typical Columns**: `Tracker ID`, `Summary`, `Description`, `Pre-Action`, `Test Step`, `Expected Result`.
- **Parser Mapping**:
  - `ID` -> `Tracker ID`
  - `component` -> Extracted from `Summary` or `Description`.
  - `expected_dlt_signature` -> Extracted from `Expected Result`.

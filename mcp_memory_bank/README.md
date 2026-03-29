# mcp_memory_bank: Native AI Memory Engine

High-performance cognitive storage layer for the Unified Agentic Framework (UAF). Enables long-term persistence of user preferences, project context, and architectural decisions.

## 🧠 Cognitive Features
- **Zero-Prompt Continuity**: Hub automatically recalls relevant history during mission initialization.
- **Fast Retrieval**: Powered by **SQLite FTS5** (Full-Text Search).
- **Persistent Logic**: Maintains context across separate CLI sessions.

## 🚀 Technical Stack
- **Language**: Go 1.26+
- **Database**: SQLite (Embedded)
- **Local First**: $0 Cloud cost, 100% data privacy.

## 🛠 Tools Exposed
- `store_memory(content, category, project_id)`: Save cognitive fragments.
- `retrieve_memories(query)`: Search history via natural language keywords.

## 📦 Location
- **Database File**: `C:/Users/perfu/.gemini/memory_bank.db`
- **Binary**: `C:/gemini_project/mcp_memory_bank/mcp_memory_bank.exe`

---
**Standardized by Gemini CLI** | *March 2026*

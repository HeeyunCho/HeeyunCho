const fs = require('fs');
const path = require('path');

/**
 * Gemini CLI BeforeTool Hook
 * Intercepts run_shell_command for git operations to enforce security audits.
 */

// Log start for debugging (stderr)
process.stderr.write('DEBUG: security_gate hook triggered\n');

// Read stdin for the Gemini context
const input = fs.readFileSync(0, 'utf-8');
let context;

try {
    context = JSON.parse(input);
} catch (e) {
    process.stderr.write('DEBUG: Error parsing hook input JSON\n');
    process.exit(1);
}

const toolName = context.tool;
const args = context.arguments || {};
const command = args.command || '';

process.stderr.write(`DEBUG: Tool: ${toolName}, Command: ${command}\n`);

// Security Gate Logic: Block git commit/push if audit wasn't requested or validated
if (toolName === 'run_shell_command' && (command.includes('git commit') || command.includes('git push'))) {
    const projectDir = 'C:/gemini_project'; // Absolute path for reliability
    const workflowStatePath = path.join(projectDir, '.workflow_state.json');
    let state = {};
    
    if (fs.existsSync(workflowStatePath)) {
        state = JSON.parse(fs.readFileSync(workflowStatePath, 'utf-8'));
    }

    process.stderr.write(`DEBUG: state.security_audit_completed: ${state.security_audit_completed}\n`);

    // If security_audit hasn't been performed in this mission or session
    if (!state.security_audit_completed) {
        process.stderr.write('SECURITY BLOCK: You are attempting a git operation without a validated security audit.\n');
        process.stderr.write('Mandate: All git/sync operations MUST be validated by security-audit-agent.\n');
        process.stderr.write('Action: Run mcp_unified-agentic-framework_validated_repo_sync or call security-audit-agent first.\n');
        process.exit(2); // Exit code 2 triggers a System Block in Gemini CLI
    }
}

// Proceed if everything is fine
console.log(JSON.stringify({}));
process.exit(0);

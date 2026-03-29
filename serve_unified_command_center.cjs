const { Server } = require("@modelcontextprotocol/sdk/server/index.js");
const { StdioServerTransport } = require("@modelcontextprotocol/sdk/server/stdio.js");
const http = require('http');
const fs = require('fs');
const sqlite3 = require('sqlite3').verbose();
const { execSync } = require('child_process');

const UI_FILE = 'C:/gemini_project/mcp_dashboard.html';
const DB_FILE = 'C:/Users/perfu/.gemini/memory_bank.db';
const TARGET_PORT = 8080;

/**
 * Graceful Port Recovery: Identifies and kills any process on the target port.
 */
function releasePort(port) {
  try {
    console.error(`[PRE-FLIGHT] Checking Port ${port}...`);
    const stdout = execSync(`netstat -ano | findstr :${port}`).toString();
    const lines = stdout.split('\n');
    if (lines.length > 0 && lines[0].includes('LISTENING')) {
      const pid = lines[0].trim().split(/\s+/).pop();
      if (pid && pid !== '0') {
        console.error(`[CLEANUP] Found ghost process ${pid} on port ${port}. Terminating...`);
        execSync(`taskkill /F /PID ${pid}`);
        console.error(`[SUCCESS] Port ${port} released.`);
      }
    }
  } catch (e) {
    // If no process is found, netstat returns exit code 1, which is fine.
    console.error(`[INFO] Port ${port} is already clear or ready.`);
  }
}

let httpServer;

function startHttpServer(port) {
  // Attempt to clear the port before starting
  if (port === TARGET_PORT) releasePort(port);

  httpServer = http.createServer((req, res) => {
    if (req.url === '/favicon.ico') {
      res.writeHead(204);
      res.end();
      return;
    }

    // Corrected: Explicitly check for API routes first
    if (req.url === '/api/memories') {
      const db = new sqlite3.Database(DB_FILE, sqlite3.OPEN_READONLY);
      db.all("SELECT category, content, created_at FROM memories ORDER BY created_at DESC LIMIT 20", [], (err, rows) => {
        db.close();
        if (err) {
          res.writeHead(500);
          return res.end(JSON.stringify({ error: 'Query failed' }));
        }
        res.writeHead(200, { 'Content-Type': 'application/json', 'Access-Control-Allow-Origin': '*' });
        res.end(JSON.stringify(rows));
      });
      return; // Stop processing further for API calls
    }

    if (req.url === '/api/stats') {
      const db = new sqlite3.Database(DB_FILE, sqlite3.OPEN_READONLY);
      const stats = { categories: [] };
      
      db.serialize(() => {
        db.get("SELECT COUNT(*) as total FROM memories", (err, row) => {
          stats.total = row ? row.total : 0;
        });
        db.all("SELECT category, COUNT(*) as count FROM memories GROUP BY category", (err, rows) => {
          stats.categories = rows || [];
          
          try {
            const fileStats = fs.statSync(DB_FILE);
            stats.fileSizeKB = (fileStats.size / 1024).toFixed(2);
          } catch(e) {
            stats.fileSizeKB = "N/A";
          }

          db.close();
          res.writeHead(200, { 'Content-Type': 'application/json', 'Access-Control-Allow-Origin': '*' });
          res.end(JSON.stringify(stats));
        });
      });
      return; // Stop processing further for API calls
    }

    // Default: Serve Dashboard UI for all other routes
    fs.readFile(UI_FILE, (err, data) => {
      if (err) {
        res.writeHead(500);
        return res.end('Error loading dashboard file.');
      }
      res.writeHead(200, { 'Content-Type': 'text/html' });
      res.end(data);
    });
  });

  httpServer.on('error', (e) => {
    if (e.code === 'EADDRINUSE') {
      console.error(`[SELF-HEAL] Port ${port} in use. Retrying on ${port + 1}...`);
      setTimeout(() => startHttpServer(port + 1), 1000);
    } else {
      console.error('Fatal Server Error:', e);
    }
  });

  httpServer.listen(port, '0.0.0.0', () => {
    console.error(`\n🟢 [COMMAND CENTER] Unified Dashboard running on port ${port}`);
  });
}

async function main() {
  if (fs.existsSync(UI_FILE)) {
    startHttpServer(TARGET_PORT);
  }
  const transport = new StdioServerTransport();
  const mcpServer = new Server({ name: "unified-dashboard-server", version: "4.0.1" }, { capabilities: {} });
  await mcpServer.connect(transport);
}

main().catch(console.error);

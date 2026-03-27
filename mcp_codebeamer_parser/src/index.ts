import { Server } from "@modelcontextprotocol/sdk/server/index.js";
import { StdioServerTransport } from "@modelcontextprotocol/sdk/server/stdio.js";
import {
  CallToolRequestSchema,
  ListToolsRequestSchema,
} from "@modelcontextprotocol/sdk/types.js";
import { z } from "zod";
import * as xlsx from "xlsx";
import * as fs from "fs-extra";
import * as path from "path";

const server = new Server(
  { name: "mcp_codebeamer_parser", version: "1.0.0" },
  { capabilities: { tools: {} } }
);

const ParseExcelSchema = z.object({
  filePath: z.string().describe("Path to the Codebeamer Excel export file."),
  outputDir: z.string().optional().describe("Directory to save the structured JSON files (optional)."),
});

server.setRequestHandler(ListToolsRequestSchema, async () => {
  return {
    tools: [
      {
        name: "parse_codebeamer_excel",
        description: "Parses a Codebeamer Excel export into structured JSON objects for ATRE ingestion.",
        inputSchema: {
          type: "object",
          properties: {
            filePath: { type: "string", description: "Path to the Excel file." },
            outputDir: { type: "string", description: "Optional output directory for JSON." },
          },
          required: ["filePath"],
        },
      }
    ],
  };
});

server.setRequestHandler(CallToolRequestSchema, async (request) => {
  const { name, arguments: args } = request.params;

  if (name === "parse_codebeamer_excel") {
    const { filePath, outputDir } = ParseExcelSchema.parse(args);

    if (!fs.existsSync(filePath)) {
      throw new Error(`File not found: ${filePath}`);
    }

    const workbook = xlsx.readFile(filePath);
    const sheetName = workbook.SheetNames[0];
    const sheet = workbook.Sheets[sheetName];
    const rawData = xlsx.utils.sheet_to_json(sheet);

    const structuredData = rawData.map((row: any, index: number) => {
      // Logic for mapping Excel columns to Codebeamer structure.
      // Assuming columns like 'ID', 'Component', 'Preconditions', 'Expected result'.
      return {
        id: row['ID'] || `TC-${index}`,
        component: row['Component'] || "Unknown",
        preconditions: row['Preconditions'] || "",
        expected_dlt_signature: row['Expected result'] || row['Expected Result'] || "",
        status: row['Status'] || "Draft",
        metadata: {
            row_index: index,
            raw: row
        }
      };
    });

    if (outputDir) {
        if (!fs.existsSync(outputDir)) {
            fs.mkdirSync(outputDir, { recursive: true });
        }
        const outputPath = path.join(outputDir, `${path.basename(filePath, path.extname(filePath))}.json`);
        fs.writeJsonSync(outputPath, structuredData, { spaces: 2 });
        return {
            content: [{ type: "text", text: `Successfully parsed ${structuredData.length} test cases and saved to ${outputPath}` }],
        };
    }

    return {
      content: [{ type: "text", text: JSON.stringify(structuredData, null, 2) }],
    };
  }

  throw new Error(`Tool not found: ${name}`);
});

async function main() {
  const transport = new StdioServerTransport();
  await server.connect(transport);
}

main().catch(console.error);

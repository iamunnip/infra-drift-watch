import { readdirSync, readFileSync } from 'node:fs';
import { join } from 'node:path';
import { JSDOM } from 'jsdom';
const dom = new JSDOM('<!doctype html><body></body>');
globalThis.window = dom.window;
globalThis.document = dom.window.document;
const mermaid = (await import('mermaid')).default;
mermaid.initialize({ startOnLoad: false });
const root = new URL('../../', import.meta.url).pathname;
let count = 0;
async function walk(dir) {
  for (const entry of readdirSync(dir, { withFileTypes: true })) {
    if (['.git', '.tools', 'node_modules'].includes(entry.name)) continue;
    const path = join(dir, entry.name);
    if (entry.isDirectory()) { await walk(path); continue; }
    if (!entry.name.endsWith('.md')) continue;
    for (const match of readFileSync(path, 'utf8').matchAll(/```mermaid\s*\n([\s\S]*?)```/g)) {
      try { await mermaid.parse(match[1]); count++; }
      catch (error) { throw new Error(`${path}: ${error.message}`); }
    }
  }
}
await walk(root);
console.log(`Parsed ${count} Mermaid diagrams`);

#!/usr/bin/env node
"use strict";

const { spawn } = require("child_process");
const path = require("path");
const fs = require("fs");

const isWindows = process.platform === "win32";
const binaryName = isWindows ? "graphify.exe" : "graphify";
const binPath = path.join(__dirname, "bin", binaryName);

if (!fs.existsSync(binPath)) {
  console.error(`Graphify-Go binary not found at ${binPath}`);
  console.error("Try reinstalling: npm i -g graphify-go");
  process.exit(1);
}

const child = spawn(binPath, process.argv.slice(2), { stdio: "inherit" });

child.on("exit", (code, signal) => {
  if (signal) {
    process.kill(process.pid, signal);
    return;
  }
  process.exit(code ?? 0);
});

child.on("error", (err) => {
  console.error(`Failed to start Graphify-Go: ${err.message}`);
  process.exit(1);
});

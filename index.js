const { execFile } = require("child_process");
const path = require("path");
const fs = require("fs");

// 代理 Graphify-Go 的二进制执行
function graphify(args = []) {
  const isWindows = process.platform === "win32";
  const binaryName = isWindows ? "graphify.exe" : "graphify";
  const binPath = path.join(__dirname, "bin", binaryName);

  if (!fs.existsSync(binPath)) {
    return Promise.reject(new Error(`Graphify-Go binary not found at ${binPath}. Please ensure the package is installed correctly.`));
  }

  return new Promise((resolve, reject) => {
    execFile(binPath, args, (error, stdout, stderr) => {
      if (error) {
        reject(error);
        return;
      }
      resolve(stdout);
    });
  });
}

module.exports = {
  graphify,
};

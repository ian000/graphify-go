const { execFile } = require("child_process");
const path = require("path");

// 代理 Graphify-Go 的二进制执行
function graphify(args = []) {
  const isWindows = process.platform === "win32";
  const binaryName = isWindows ? "graphify.exe" : "graphify";
  const binPath = path.join(__dirname, "bin", binaryName);

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

const fs = require("fs");
const path = require("path");
const https = require("https");

const VERSION = "v1.0.0"; // 这里对应 GitHub Release Tag
const REPO = "ian000/graphify-go";

const platformMap = {
  win32: "windows",
  darwin: "darwin",
  linux: "linux",
};

const archMap = {
  x64: "amd64",
  arm64: "arm64",
};

function getDownloadUrl() {
  const platform = platformMap[process.platform];
  const arch = archMap[process.arch];

  if (!platform || !arch) {
    console.error(`❌ Unsupported platform or architecture: ${process.platform} ${process.arch}`);
    process.exit(1);
  }

  let assetName = `graphify-${platform}-${arch}`;
  if (platform === "windows") {
    assetName += ".exe";
  }

  return `https://github.com/${REPO}/releases/download/${VERSION}/${assetName}`;
}

function downloadBinary() {
  const url = getDownloadUrl();
  const binDir = path.join(__dirname, "..", "bin");
  const isWindows = process.platform === "win32";
  const binaryName = isWindows ? "graphify.exe" : "graphify";
  const binPath = path.join(binDir, binaryName);

  if (!fs.existsSync(binDir)) {
    fs.mkdirSync(binDir, { recursive: true });
  }

  console.log(`⬇️ Downloading Graphify-Go binary from ${url}...`);

  const file = fs.createWriteStream(binPath);

  // 简单的下载逻辑（不处理 302 重定向等复杂场景，建议真实场景引入 axios/got 或 follow-redirects）
  https.get(url, (response) => {
    if (response.statusCode === 302 || response.statusCode === 301) {
      https.get(response.headers.location, (res) => {
        res.pipe(file);
        file.on("finish", () => {
          file.close();
          if (!isWindows) {
            fs.chmodSync(binPath, "755");
          }
          console.log(`✅ Successfully downloaded to ${binPath}`);
        });
      }).on("error", (err) => {
        fs.unlink(binPath, () => {});
        console.error(`❌ Download failed: ${err.message}`);
        process.exit(1);
      });
    } else {
      response.pipe(file);
      file.on("finish", () => {
        file.close();
        if (!isWindows) {
          fs.chmodSync(binPath, "755");
        }
        console.log(`✅ Successfully downloaded to ${binPath}`);
      });
    }
  }).on("error", (err) => {
    fs.unlink(binPath, () => {});
    console.error(`❌ Download failed: ${err.message}`);
    process.exit(1);
  });
}

downloadBinary();

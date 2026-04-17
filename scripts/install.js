const fs = require("fs");
const path = require("path");
const https = require("https");

const pkg = require("../package.json");
const VERSION_CANDIDATES = [`v${pkg.version}`, "v1.0.0"];
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

function getDownloadUrl(versionTag) {
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

  return `https://github.com/${REPO}/releases/download/${versionTag}/${assetName}`;
}

function downloadToFile(url, binPath, isWindows) {
  return new Promise((resolve, reject) => {
    const file = fs.createWriteStream(binPath);

    const onError = (err) => {
      fs.unlink(binPath, () => {});
      reject(err);
    };

    const request = (currentUrl) => {
      https
        .get(currentUrl, (response) => {
          if (response.statusCode === 301 || response.statusCode === 302) {
            request(response.headers.location);
            return;
          }

          if (response.statusCode !== 200) {
            reject(new Error(`Unexpected status code ${response.statusCode}`));
            return;
          }

          response.pipe(file);
          file.on("finish", () => {
            file.close();
            if (!isWindows) {
              fs.chmodSync(binPath, "755");
            }
            resolve();
          });
        })
        .on("error", onError);
    };

    request(url);
    file.on("error", onError);
  });
}

async function downloadBinary() {
  const binDir = path.join(__dirname, "..", "bin");
  const isWindows = process.platform === "win32";
  const binaryName = isWindows ? "graphify.exe" : "graphify";
  const binPath = path.join(binDir, binaryName);

  if (!fs.existsSync(binDir)) {
    fs.mkdirSync(binDir, { recursive: true });
  }

  for (const versionTag of VERSION_CANDIDATES) {
    const url = getDownloadUrl(versionTag);
    console.log(`⬇️ Downloading Graphify-Go binary from ${url}...`);
    try {
      await downloadToFile(url, binPath, isWindows);
      console.log(`✅ Successfully downloaded to ${binPath}`);
      return;
    } catch (err) {
      console.warn(`⚠️ Download failed for ${versionTag}: ${err.message}`);
    }
  }

  console.error("❌ Download failed for all candidate release tags.");
  process.exit(1);
}

downloadBinary();

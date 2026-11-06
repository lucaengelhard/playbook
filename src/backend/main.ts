import { app, BrowserWindow } from "electron";
import Client from "ssh2-sftp-client";
import path from "path";
import fs from "fs";

app.whenReady().then(() => {
  createWindow();
  app.on("activate", createWindow);
});

app.on("window-all-closed", () => {
  if (process.platform !== "darwin") app.quit();
});

function createWindow() {
  if (BrowserWindow.getAllWindows().length !== 0) return;
  const window = new BrowserWindow({ width: 800, height: 600 });
  window.loadFile("dist/backend/view/index.html");
}

const sftp = new Client();

// TODO create temp dir (or other way to store existing plays & updated state) e.g: const tempDir = fs.mkdtempSync("temp_playbook_");
// TODO connect to sftp server
// TODO download current state
// TODO read current state from dirs -> send to view via ipc
// function updateState(): Play[]
//
// TODO load frontend html/css/js to be copied when scaffolding new play

type Play = {
  title: string;
  thumbnail: string;
  content: string;
  gallery: string[];
};

function createSlug(title: string): string {
  // TODO prepend uid?
  return encodeURI(title);
}

async function createPlay(play: Play) {
  const rootDirPath = createSlug(play.title); // TODO join path with whatever the temp dir or else is
  // TODO check if path already exists -> abort & send error to view

  fs.mkdirSync(path.join(rootDirPath, "media"), { recursive: true });
  fs.writeFileSync(
    path.join(rootDirPath, "content.json"),
    JSON.stringify(play),
    "utf8",
  );
}

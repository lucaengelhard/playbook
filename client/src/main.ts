import "./style.css";
import { init, os, events, extensions, filesystem } from "@neutralinojs/lib";

init();

const heroButton = document.querySelector("#form_hero");
const heroSelected = document.querySelector<HTMLDivElement>(
  "#form_hero_selected",
);
heroButton?.addEventListener("click", async (e) => {
  e.preventDefault();
  const entries = await os.showOpenDialog("Open Image", {
    filters: [
      { name: "Images", extensions: ["jpg", "png"] },
      { name: "Any file", extensions: ["*"] },
    ],
  });
  if (entries.length === 0) return;

  heroSelected!.innerText = entries[0];
});

document.querySelector("#form_submit")?.addEventListener("click", async (e) => {
  e.preventDefault();
  const play = {
    title: document.querySelector<HTMLInputElement>("#form_title")?.value,
    hero_image: heroSelected?.innerText,
    content:
      document.querySelector<HTMLTextAreaElement>("#form_content")?.value,
  };

  const path = await filesystem.getJoinedPath(
    await filesystem.getAbsolutePath(window.NL_DATAPATH),
    "plays",
    play.title ?? "",
  );

  // TODO validate data
  // create upload dir
  // copy template to upload dir
  // write json to upload dir
  // copy media to upload dir

  await extensions.dispatch("playbook.upload", "upload", { path });
});

document.querySelector("#sftp")?.addEventListener("click", async (e) => {
  e.preventDefault();

  const host = document.querySelector<HTMLInputElement>("#sftp_host")?.value;
  const port = document.querySelector<HTMLInputElement>("#sftp_port")?.value;
  const password = document.querySelector<HTMLInputElement>("#sftp_pwd")?.value;
  const user = document.querySelector<HTMLInputElement>("#sftp_user")?.value;

  await extensions.dispatch("playbook.upload", "connect", {
    host,
    port,
    password,
    user,
  });
});

events.on("upload", console.log);
events.on("sftp_connected", console.log);

events.on("error", console.log);

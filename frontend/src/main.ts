import { Dialogs } from "@wailsio/runtime";
import { Client, Template, type Play } from "../bindings/playbook";

const heroButton = document.querySelector("#form_hero");
const heroSelected = document.querySelector<HTMLDivElement>(
  "#form_hero_selected",
);

heroButton?.addEventListener("click", async (e) => {
  e.preventDefault();
  const file = await Dialogs.OpenFile({
    AllowsMultipleSelection: false,
    CanChooseDirectories: false,
    Filters: [
      { DisplayName: "Images", Pattern: "*.jpg;*.png" },
      { DisplayName: "Any file", Pattern: "*" },
    ],
  });

  if (file) {
    heroSelected!.innerText = file;
  }
});

document.querySelector("#form_submit")?.addEventListener("click", async (e) => {
  e.preventDefault();
  const play: Play = {
    Title: document.querySelector<HTMLInputElement>("#form_title")?.value!,
    HeroImage: heroSelected?.innerText!,
    Content:
      document.querySelector<HTMLTextAreaElement>("#form_content")?.value!,
  };
  if (!play.Title || play.Title === "") return;

  await Template.Scaffold(play);
});

document.querySelector("#sftp")?.addEventListener("click", async (e) => {
  e.preventDefault();

  const host = document.querySelector<HTMLInputElement>("#sftp_host")?.value;
  const port = document.querySelector<HTMLInputElement>("#sftp_port")?.value;
  const password = document.querySelector<HTMLInputElement>("#sftp_pwd")?.value;
  const user = document.querySelector<HTMLInputElement>("#sftp_user")?.value;

  Client.Connect(host!, port!, user!, password!);
});

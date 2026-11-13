import data from "./content.json" with { type: "json" };

const title = document.querySelector("h1");
title.innerText = data.Title;

const header = document.querySelector("img");
header.src = data.HeroImage;

const content = document.querySelector("p");
content.innerText = data.Content;

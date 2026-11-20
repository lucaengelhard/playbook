function holeElement(id) {
    const element = document.getElementById(id);
    if (!element)
        throw new Error("HTML-Element fehlt: " + id);
    return element;
}
// JSON kommt von außen und wird deshalb auch zur Laufzeit geprüft.
function istPlay(daten) {
    if (typeof daten !== "object" || daten === null)
        return false;
    return "Title" in daten && typeof daten.Title === "string"
        && "HeroImage" in daten && typeof daten.HeroImage === "string"
        && "Content" in daten && typeof daten.Content === "string"
        && (!("Gallery" in daten) || daten.Gallery === undefined
            || (Array.isArray(daten.Gallery)
                && daten.Gallery.every((bild) => typeof bild === "string")));
}
const startseite = holeElement("startseite");
const menue = holeElement("menue");
const inhaltsseite = holeElement("inhaltsseite");
const ladestatus = holeElement("status");
function zeigeAnsicht(ansicht) {
    startseite.hidden = true;
    menue.hidden = true;
    inhaltsseite.hidden = true;
    ansicht.hidden = false;
    ansicht.querySelector("h1")?.focus({ preventScroll: true });
    window.scrollTo(0, 0);
}
function erstelleBild(pfad) {
    const bild = document.createElement("img");
    bild.src = pfad;
    // Play enthält bisher keine Bildbeschreibungen.
    bild.alt = "";
    return bild;
}
function zeigeText(heft) {
    holeElement("kapiteltitel").textContent = "Über das Stück";
    const inhalt = holeElement("kapitelinhalt");
    inhalt.replaceChildren();
    // Leerzeilen trennen die Absätze. Content wird als reiner Text angezeigt.
    for (const absatz of heft.Content.split(/\r?\n\s*\r?\n/)) {
        if (!absatz.trim())
            continue;
        const text = document.createElement("p");
        text.textContent = absatz;
        inhalt.appendChild(text);
    }
    if (inhalt.childElementCount === 0) {
        inhalt.textContent = "Zu diesem Stück ist noch kein Text vorhanden.";
    }
    zeigeAnsicht(inhaltsseite);
}
function zeigeGalerie(bilder) {
    holeElement("kapiteltitel").textContent = "Bilder";
    const inhalt = holeElement("kapitelinhalt");
    inhalt.replaceChildren();
    for (const pfad of bilder) {
        const bild = erstelleBild(pfad);
        bild.loading = "lazy";
        inhalt.appendChild(bild);
    }
    zeigeAnsicht(inhaltsseite);
}
function menueknopf(text, aktion) {
    const knopf = document.createElement("button");
    knopf.type = "button";
    knopf.textContent = text;
    knopf.addEventListener("click", aktion);
    holeElement("menue-liste").appendChild(knopf);
}
holeElement("oeffnen").addEventListener("click", () => zeigeAnsicht(menue));
holeElement("zur-startseite").addEventListener("click", () => zeigeAnsicht(startseite));
holeElement("zum-menue").addEventListener("click", () => zeigeAnsicht(menue));
async function ladeHeft() {
    try {
        // Lokale Testdaten im gleichen Format wie das Backend.
        const antwort = await fetch("./heft-test.json");
        if (!antwort.ok)
            throw new Error("HTTP-Fehler " + antwort.status);
        const daten = await antwort.json();
        if (!istPlay(daten))
            throw new Error("Ungültige Play-Daten");
        const heft = daten;
        document.title = heft.Title + " – Digitales Programmheft";
        holeElement("stuecktitel").textContent = heft.Title;
        const titelbild = holeElement("titelbild");
        titelbild.replaceChildren();
        if (heft.HeroImage.trim()) {
            titelbild.appendChild(erstelleBild(heft.HeroImage));
        }
        holeElement("menue-liste").replaceChildren();
        menueknopf("Über das Stück", () => zeigeText(heft));
        const bilder = (heft.Gallery ?? []).filter(pfad => pfad.trim() !== "");
        if (bilder.length > 0) {
            menueknopf("Bilder", () => zeigeGalerie(bilder));
        }
        ladestatus.hidden = true;
        startseite.hidden = false;
    }
    catch (fehler) {
        ladestatus.hidden = false;
        ladestatus.textContent = "Das Programmheft konnte nicht geladen werden. Bitte laden Sie die Seite erneut.";
        console.error("Fehler beim Laden des Programmhefts:", fehler);
    }
}
void ladeHeft();
export {};

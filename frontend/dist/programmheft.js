function holeElement(id) {
    const element = document.getElementById(id);
    if (!element) {
        throw new Error("HTML-Element fehlt: " + id);
    }
    return element;
}
const ticketseite = holeElement("ticketseite");
const erklaerungsseite = holeElement("erklaerungsseite");
const startseite = holeElement("startseite");
const menueleiste = holeElement("menue-liste");
const inhaltsseite = holeElement("inhaltsseite");
const ladestatus = holeElement("status");
let heftGeladen = false;
// höhe vom menü als css variable setzen, damit unten nix verdeckt wird
function setzeMenueHoehe() {
    const hoehe = menueleiste.offsetHeight;
    document.documentElement.style.setProperty("--menuehoehe", hoehe + "px");
}
window.addEventListener("resize", setzeMenueHoehe);
function zeigeAnsicht(ansicht) {
    ticketseite.hidden = true;
    erklaerungsseite.hidden = true;
    startseite.hidden = true;
    inhaltsseite.hidden = true;
    ansicht.hidden = false;
    // menü nur auf startseite + inhaltsseite
    menueleiste.hidden = !heftGeladen || (ansicht !== startseite && ansicht !== inhaltsseite);
    setzeMenueHoehe();
    if (ansicht === startseite) {
        markiereMenuepunkt("Startseite");
    }
    const ueberschrift = ansicht.querySelector("h1");
    if (ueberschrift) {
        ueberschrift.focus({ preventScroll: true });
    }
    window.scrollTo(0, 0);
}
function erstelleBild(pfad) {
    const bild = document.createElement("img");
    bild.src = pfad;
    // TODO: alt texte fehlen noch im json
    bild.alt = "";
    return bild;
}
// text in absätze aufteilen (leerzeile = neuer absatz)
function fuelleText(text, inhalt) {
    inhalt.replaceChildren();
    const absaetze = text.split(/\r?\n\s*\r?\n/);
    for (const absatz of absaetze) {
        if (!absatz.trim()) {
            continue;
        }
        const p = document.createElement("p");
        p.textContent = absatz;
        inhalt.appendChild(p);
    }
    if (inhalt.childElementCount === 0) {
        inhalt.textContent = "Zu diesem Stück ist noch kein Text vorhanden.";
    }
}
function zeigeText(heft) {
    holeElement("kapiteltitel").textContent = "Über das Stück";
    fuelleText(heft.Content, holeElement("kapitelinhalt"));
    zeigeAnsicht(inhaltsseite);
}
function zeigeGalerie(bilder, unterschriften, texte) {
    if (!unterschriften) {
        unterschriften = {};
    }
    if (!texte) {
        texte = [];
    }
    holeElement("kapiteltitel").textContent = "Bilder";
    const inhalt = holeElement("kapitelinhalt");
    inhalt.replaceChildren();
    const galerie = document.createElement("div");
    galerie.className = "galerie";
    galerie.setAttribute("role", "region");
    galerie.setAttribute("aria-label", "Bildergalerie");
    for (const pfad of bilder) {
        const bild = erstelleBild(pfad);
        bild.loading = "lazy";
        const abbildung = document.createElement("figure");
        abbildung.appendChild(bild);
        if (unterschriften[pfad]) {
            const unterschrift = document.createElement("figcaption");
            unterschrift.textContent = unterschriften[pfad];
            abbildung.appendChild(unterschrift);
        }
        galerie.appendChild(abbildung);
    }
    const rahmen = document.createElement("div");
    rahmen.className = "galerie-rahmen";
    // pfeile links/rechts
    const zurueck = document.createElement("button");
    zurueck.type = "button";
    zurueck.className = "galerie-pfeil links";
    zurueck.textContent = "‹";
    zurueck.setAttribute("aria-label", "Vorheriges Bild");
    const weiter = document.createElement("button");
    weiter.type = "button";
    weiter.className = "galerie-pfeil rechts";
    weiter.textContent = "›";
    weiter.setAttribute("aria-label", "Nächstes Bild");
    zurueck.hidden = bilder.length < 2;
    weiter.hidden = bilder.length < 2;
    // wenn ein bild keinen text hat, bleibt der vom vorherigen bild stehen
    let lastText = "";
    const bildTexte = bilder.map((_, index) => {
        let text = "";
        if (texte[index]) {
            text = texte[index].trim();
        }
        if (text) {
            lastText = text;
        }
        return lastText;
    });
    const textfeld = document.createElement("section");
    textfeld.className = "galerie-text";
    textfeld.setAttribute("aria-label", "Text zum Bild");
    // punkte unter der galerie
    const punkte = document.createElement("div");
    punkte.className = "galerie-punkte";
    punkte.setAttribute("aria-hidden", "true");
    punkte.hidden = bilder.length < 2;
    const dots = [];
    for (let i = 0; i < bilder.length; i++) {
        const punkt = document.createElement("span");
        punkte.appendChild(punkt);
        dots.push(punkt);
    }
    // für screenreader: "bild x von y"
    const position = document.createElement("p");
    position.className = "nur-screenreader";
    position.setAttribute("role", "status");
    let aktuellesBild = 0;
    let angezeigterText;
    function updateAnzeige() {
        zurueck.disabled = aktuellesBild === 0;
        weiter.disabled = aktuellesBild === bilder.length - 1;
        dots.forEach((punkt, index) => {
            punkt.className = index === aktuellesBild ? "galerie-punkt aktiv" : "galerie-punkt";
        });
        const positionstext = "Bild " + (aktuellesBild + 1) + " von " + bilder.length;
        if (position.textContent !== positionstext) {
            position.textContent = positionstext;
        }
        const bildText = bildTexte[aktuellesBild];
        // nur neu bauen wenn sich der text ändert
        if (bildText !== angezeigterText) {
            textfeld.hidden = !bildText;
            textfeld.replaceChildren();
            if (bildText) {
                fuelleText(bildText, textfeld);
            }
            angezeigterText = bildText;
        }
    }
    function geheZuBild(index) {
        const ziel = Math.max(0, Math.min(index, bilder.length - 1));
        console.log("gehe zu bild", ziel);
        const reduziert = window.matchMedia("(prefers-reduced-motion: reduce)").matches;
        galerie.scrollTo({ left: ziel * galerie.clientWidth, behavior: reduziert ? "auto" : "smooth" });
    }
    zurueck.addEventListener("click", () => geheZuBild(aktuellesBild - 1));
    weiter.addEventListener("click", () => geheZuBild(aktuellesBild + 1));
    // beim wischen aktuelles bild ausrechnen
    galerie.addEventListener("scroll", () => {
        if (galerie.clientWidth === 0) {
            return;
        }
        const index = Math.round(galerie.scrollLeft / galerie.clientWidth);
        aktuellesBild = Math.max(0, Math.min(index, bilder.length - 1));
        updateAnzeige();
    }, { passive: true });
    rahmen.append(galerie, zurueck, weiter);
    inhalt.append(rahmen, punkte, position, textfeld);
    updateAnzeige();
    zeigeAnsicht(inhaltsseite);
}
function markiereMenuepunkt(text) {
    for (const knopf of menueleiste.querySelectorAll("button")) {
        if (knopf.textContent === text) {
            knopf.setAttribute("aria-current", "page");
        }
        else {
            knopf.removeAttribute("aria-current");
        }
    }
}
function addMenueKnopf(text, aktion) {
    const knopf = document.createElement("button");
    knopf.type = "button";
    knopf.textContent = text;
    knopf.addEventListener("click", () => {
        aktion();
        markiereMenuepunkt(text);
    });
    menueleiste.appendChild(knopf);
}
holeElement("ticket-weiter").addEventListener("click", () => zeigeAnsicht(erklaerungsseite));
holeElement("erklaerung-weiter").addEventListener("click", () => zeigeAnsicht(startseite));
async function ladeHeft() {
    try {
        const antwort = await fetch("./content.json");
        if (!antwort.ok) {
            throw new Error("HTTP-Fehler " + antwort.status);
        }
        // TODO: vielleicht doch noch prüfen ob das json passt
        const heft = await antwort.json();
        console.log("heft geladen:", heft);
        document.title = heft.Title + " – Digitales Programmheft";
        holeElement("stuecktitel").textContent = heft.Title;
        const titelbild = holeElement("titelbild");
        titelbild.replaceChildren();
        if (heft.HeroImage.trim()) {
            titelbild.appendChild(erstelleBild(heft.HeroImage));
        }
        fuelleText(heft.Content, holeElement("einleitung"));
        menueleiste.replaceChildren();
        addMenueKnopf("Startseite", () => zeigeAnsicht(startseite));
        addMenueKnopf("Über das Stück", () => zeigeText(heft));
        // leere bilder rausfiltern, texte mitnehmen damit die reihenfolge passt
        const bilder = [];
        const texte = [];
        if (heft.Gallery) {
            heft.Gallery.forEach((pfad, index) => {
                if (pfad.trim() === "") {
                    return;
                }
                bilder.push(pfad);
                let text = "";
                if (heft.GalleryTexts && heft.GalleryTexts[index]) {
                    text = heft.GalleryTexts[index];
                }
                texte.push(text);
            });
        }
        if (bilder.length > 0) {
            addMenueKnopf("Bilder", () => zeigeGalerie(bilder, heft.Bildunterschriften, texte));
        }
        markiereMenuepunkt("Startseite");
        heftGeladen = true;
        ladestatus.hidden = true;
        holeElement("heftdaten").hidden = false;
        menueleiste.hidden = startseite.hidden && inhaltsseite.hidden;
        setzeMenueHoehe();
    }
    catch (fehler) {
        ladestatus.hidden = false;
        ladestatus.textContent = "Das Programmheft konnte nicht geladen werden. Bitte laden Sie die Seite erneut.";
        console.error("Fehler beim Laden des Programmhefts:", fehler);
    }
}
ladeHeft();
export {};

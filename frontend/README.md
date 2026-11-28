## Wo liegt was?

- `src/index.html`: Aufbau und CSS der Oberfläche.
- `src/programmheft.ts`: TypeScript-Logik zum Laden und Anzeigen
  der Inhalte sowie für die Navigation.
- `src/assets/`: Platz für feste Bilder und Symbole der Vorlage.
- `types.gen.ts`: Automatisch erzeugte Backend-Typen.
  Nicht von Hand bearbeiten.
- `testdaten/content.json`: Beispielinhalte zum lokalen Testen.
- `testdaten/assets/`: Bilder des Testhefts.
- `tsconfig.json`: Einstellungen für die TypeScript-Prüfung.
- `dist/`: HTML- und JavaScript-Ausgabe für die spätere Veröffentlichung.

## Wo Änderungen machen?

Oberfläche und Gestaltung in `src/index.html` ändern.
Programmlogik in `src/programmheft.ts` bearbeiten.
Beispielinhalte unter `testdaten/` anpassen.

Die Dateien in `dist/` nicht direkt bearbeiten.

## Bauen und lokal testen

```powershell
.\frontend\build.ps1            # TypeScript übersetzen, dist/ aktualisieren
.\frontend\build.ps1 -Vorschau  # zusätzlich frontend/vorschau/ erzeugen
```

Das Frontend lädt `content.json` und die Bilder aus demselben Ordner
wie `index.html`. `dist/` enthält bewusst keine Testdaten, deshalb
zeigt `dist/index.html` allein nur eine Fehlermeldung.

Zum Testen `frontend/vorschau/index.html` mit Live Server öffnen
(nicht per Doppelklick, über `file://` funktioniert `fetch` nicht).
`vorschau/` enthält `dist/` und `testdaten/` zusammen, wird bei jedem
Aufruf neu erzeugt und nicht eingecheckt.

Nur Typen prüfen, ohne etwas zu erzeugen: `tsc.cmd -p frontend/tsconfig.json`


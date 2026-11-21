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
- `build.ps1`: Bisher nur TODO-Kommentare, noch kein fertiger Build.
- `dist/`: HTML- und JavaScript-Ausgabe für die spätere Veröffentlichung.

## Wo Änderungen machen?

Oberfläche und Gestaltung in `src/index.html` ändern.
Programmlogik in `src/programmheft.ts` bearbeiten.
Beispielinhalte unter `testdaten/` anpassen.

Die Dateien in `dist/` nicht direkt bearbeiten.

## Aktueller Stand

`tsc.cmd -p frontend/tsconfig.json` prüft momentan nur die Typen.
Es wird dabei noch keine neue JavaScript-Ausgabe erzeugt.

Das Frontend erwartet `content.json` neben der veröffentlichten
HTML-Datei. Die Testdaten werden noch nicht automatisch dorthin kopiert.


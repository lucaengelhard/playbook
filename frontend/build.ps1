# Baut das Frontend:
#   1. TypeScript aus src/ übersetzen
#   2. index.html, programmheft.js und feste Assets nach dist/ kopieren
# Mit -Vorschau wird zusätzlich vorschau/ erzeugt (dist + testdaten in einem Ordner).
# Diesen Ordner mit Live Server öffnen. vorschau/ wird nicht eingecheckt.
#
# types.gen.ts wird nur gelesen, nie verändert.
#
# Aufruf:  .\frontend\build.ps1            (nur dist)
#          .\frontend\build.ps1 -Vorschau  (dist + vorschau)

param([switch]$Vorschau)

$ErrorActionPreference = "Stop"

$wurzel    = $PSScriptRoot
$src       = Join-Path $wurzel "src"
$dist      = Join-Path $wurzel "dist"
$tmp       = Join-Path $wurzel ".build"
$testdaten = Join-Path $wurzel "testdaten"
$vorschauOrdner = Join-Path $wurzel "vorschau"

# 1. TypeScript übersetzen (in einen temporären Ordner, nur programmheft.js wird gebraucht)
$tsc = if (Get-Command tsc.cmd -ErrorAction SilentlyContinue) { "tsc.cmd" } else { "tsc" }
& $tsc -p (Join-Path $wurzel "tsconfig.json") --noEmit false --outDir $tmp --rootDir $wurzel
if ($LASTEXITCODE -ne 0) {
    throw "TypeScript-Fehler, Build abgebrochen."
}

# 2. dist befüllen
New-Item -ItemType Directory -Force (Join-Path $dist "assets") | Out-Null
Copy-Item (Join-Path $tmp "src\programmheft.js") $dist -Force
Copy-Item (Join-Path $src "index.html") $dist -Force
Get-ChildItem (Join-Path $src "assets") -File | Where-Object Name -ne ".gitkeep" |
    Copy-Item -Destination (Join-Path $dist "assets") -Force
Remove-Item $tmp -Recurse -Force
Write-Host "dist/ aktualisiert."

# 3. Optional: Vorschau zum lokalen Testen
if ($Vorschau) {
    if (Test-Path $vorschauOrdner) {
        Remove-Item $vorschauOrdner -Recurse -Force
    }
    New-Item -ItemType Directory -Force (Join-Path $vorschauOrdner "assets") | Out-Null
    Get-ChildItem $dist -File | Copy-Item -Destination $vorschauOrdner
    Get-ChildItem (Join-Path $dist "assets") -File | Where-Object Name -ne ".gitkeep" |
        Copy-Item -Destination (Join-Path $vorschauOrdner "assets")
    Copy-Item (Join-Path $testdaten "content.json") $vorschauOrdner
    Get-ChildItem (Join-Path $testdaten "assets") -File |
        Copy-Item -Destination (Join-Path $vorschauOrdner "assets")
    Write-Host "vorschau/ erstellt. Mit Live Server frontend/vorschau/index.html öffnen."
}

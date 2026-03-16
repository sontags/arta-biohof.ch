# Arta Biohof - Website

Statische Website fuer den [Arta Biohof](https://arta-biohof.ch) in Matzwil, Detligen.

Die Seite wird mit einem eigenen, in Go geschriebenen Static-Site-Generator
gebaut. Markdown-Inhalte werden in ein HTML-Template gerendert und als einzelne
`index.html` ausgegeben.

## Voraussetzungen

- [Go](https://go.dev/) >= 1.22

## Projektstruktur

```
.
├── builder/           # Static-Site-Generator (Go)
│   └── main.go
├── template/          # HTML-Template mit eingebettetem CSS
│   └── main.html.templ
├── content/           # Seiteninhalte als Markdown
│   ├── home.md
│   ├── team.md
│   ├── pferdepension.md
│   ├── aepfel_und_birnen.md
│   ├── mosterei.md
│   ├── impressionen.md
│   └── kontakt.md
├── dist/              # Ausgabeverzeichnis (HTML, Bilder, Logo)
├── config.yaml        # Build-Konfiguration
└── go.mod
```

## Seite bauen

```sh
go run ./builder -config config.yaml -out dist/index.html
```

Ohne `-out` wird das HTML auf STDOUT ausgegeben.

## Inhalte bearbeiten

Die Inhalte der einzelnen Sektionen liegen als Markdown-Dateien unter `content/`.
HTML innerhalb der Markdown-Dateien ist erlaubt (z.B. fuer die Bildergalerie).

Neue Sektionen koennen in `config.yaml` hinzugefuegt werden:

```yaml
content:
  - path: content/neue_sektion.md
    metadata:
      anchor: neuesektion
      navname: Neue Sektion
```

Jeder Eintrag benoetigt:

- **path** - Pfad zur Markdown-Datei
- **anchor** - HTML-Anker fuer die Navigation
- **navname** - Anzeigename in der Navigationsleiste

## Bilder verwalten

Bilder werden direkt im `dist/`-Verzeichnis abgelegt und koennen in den
Markdown-Dateien referenziert werden:

```markdown
![Beschreibung](bildname.jpg)
```

## Deployment

Das `dist/`-Verzeichnis enthaelt die fertige Seite und kann auf einen
beliebigen Webserver kopiert werden.

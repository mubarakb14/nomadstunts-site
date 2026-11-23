# Nomad Stunts — website

Static website of Nomad Stunts International Action Company (ҚАЗ / РУС / ENG).

Pages: `index.html`, `about.html`, `team.html`, `projects.html`, `gallery.html`, `contacts.html`.

- All content (projects, team, gallery) lives in `assets/data.js`.
- Interface texts, history and services in three languages are in `assets/app.js`.
- Images are in `upload/`, fonts in `assets/fonts/`, showreel in `assets/video/`.

No build step: open `index.html` or serve the folder with any static host.

## Admin panel

`server/` is a Go server that serves this site plus an admin panel at `/admin`, where the team edits the showreel, team, founder, projects, gallery and contacts without a developer. Content is stored in `DATA_DIR/content.json` and served as `/assets/data.js`. Setup, env vars and deployment: [server/README.md](server/README.md).

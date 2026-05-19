# Frontend Setup

## Kurulum

```bash
cd frontend
npm install
npm run dev
```

## Build Üretim İçin

```bash
npm run build
npm run preview
```

## Proje Yapısı

- `src/` - Vue.js component'leri ve scripts
- `src/App.vue` - Ana component
- `src/main.js` - Entry point
- `vite.config.js` - Vite configuration
- `index.html` - HTML template

## Features

- Teams yönetimi
- Matches planlama ve sonuçları
- League table otomatik hesaplama
- Final tahminler
- Responsive tasarım

## API Integration

Frontend `http://localhost:8080/api` adresindeki backend ile iletişim kurar.

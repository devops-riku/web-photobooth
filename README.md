# Wuby Photobooth 📸

A premium, modern web-based photobooth experience designed for high-end events and digital memory creation. Built with a focus on aesthetics, smooth interactions, and instant sharing.

---

## ✨ Features

### 🎞️ Capturing Experience
- **Live Viewport**: Squared, high-performance camera preview with real-time mirroring.
- **Smart Sequence**: Automated capture sequence with configurable countdown timers.
- **Flash Effect**: Immersive visual feedback during every shot.
- **Flexible Entry**: Capture directly via camera or upload high-quality photos from your device.
- **Swap System**: Easily replace any individual photo in a sequence after capture.

### 🎨 Personalization & Editing
- **Advanced Color Engine**: 
  - **Dynamic Color Wheel**: Intuitive 2D Hue/Saturation selection.
  - **Brightness Control**: Precision luminance adjustment.
  - **Hex Input**: Support for specific brand primary/secondary colors.
- **AI-Style Filters**: Professional-grade filters including Cinematic, Film, Warm, and B&W (powered by `glfx.js`).
- **Rich Typography**: Selection of 10+ curated Google Fonts with real-time size and text preview.
- **Layout Logic**: Support for classic vertical vertical strips (2, 3, or 4 photos).

### 🚀 Sharing & Branding
- **Automated Branding**: Dynamic injection of logos and event branding on every strip.
- **Digital Heritage**: Generates unique, scannable QR codes for instant mobile downloads.
- **Cloud Sync**: High-resolution rendering and direct upload to DigitalOcean Spaces.
- **Hybrid Access**: Dedicated workflows for both registered users and guest sessions.

---

## 🛠️ Tech Stack

- **Frontend**: [SvelteKit 5](https://svelte.dev/) + [Tailwind CSS 4](https://tailwindcss.com/)
- **Backend**: [Go (Golang)](https://go.dev/)
- **Image Processing**: HTML5 Canvas + [glfx.js](https://github.com/evanw/glfx.js)
- **Infrastructure**: Docker, Nginx, DigitalOcean Spaces

---

## �️ Setup & Installation

You can run Wuby Photobooth either as a coordinated Docker environment or for local development with live reloading.

### 1. Local Development Setup
Best for making code changes to the frontend or backend.

#### Prerequisites
- **Node.js** (v18+)
- **Go** (v1.21+)
- **PostgreSQL** (running locally)

#### 🔹 Backend
```bash
cd backend
# 1. Configure environment
cp .env.template .env # Update DATABASE_URL to your local PG
# 2. Run with live reload (requires air) or standard go
go run main.go
```

#### 🔹 Frontend
```bash
cd web
# 1. Install dependencies
npm install
# 2. Start dev server with HMR
npm run dev
```

---

### 2. Docker Setup (Recommended for Production)
This starts the entire stack including the database, Go backend, Svelte frontend, and Nginx gateway.

#### Prerequisites
- **Docker** & **Docker Compose**

#### 🚀 Launching
```bash
# 1. Setup environment files
cp backend/.env.template backend/.env
# Update backend/.env with your secrets

# 2. Build and start services in detached mode
docker-compose up -d --build

# 3. Check logs
docker-compose logs -f
```

#### 🌐 Access
Once running, the application is proxied through Nginx:
- **Frontend**: [http://localhost:8080](http://localhost:8080)
- **API**: [http://localhost:8080/api](http://localhost:8080/api)

---

## ⚙️ Environment Variables

Primary configuration is managed via `.env` files in `backend/` and `web/`.

| Variable | Description |
| :--- | :--- |
| `DATABASE_URL` | PostgreSQL connection string |
| `DO_SPACES_KEY` | DigitalOcean Spaces Access Key |
| `DO_SPACES_SECRET` | DigitalOcean Spaces Secret Key |
| `DO_SPACES_ENDPOINT` | Your Spaces region endpoint |
| `DO_SPACES_BUCKET` | Your bucket/folder name |

---

## 📱 Mobile & Tunneling
For onsite testing on tablets (iOS/Android) without HTTPS:

**For Local Dev:**
```bash
npx localtunnel --port 5173 --subdomain wuby-dev
```

**For Docker Stack:**
```bash
npx localtunnel --port 8080 --subdomain wuby-booth
```

---

Designed with ❤️ by [LinkerX](https://linkerx.io)
# Deployment Guide for prizebot.dev

## Overview

Your Prize Bot consists of two parts:
- **Backend (Go)**: Handles Discord bot logic and API endpoints
- **Frontend (SvelteKit)**: Dashboard UI

## Deployment Steps

### 1. Deploy Backend

Your backend needs to be hosted on a service that supports Go applications. Recommended options:

#### Option A: Railway (Recommended - Easiest)

1. Go to [railway.app](https://railway.app) and sign in with GitHub
2. Click "New Project" → "Deploy from GitHub repo"
3. Select your `Prize-Bot-V1` repository
4. Railway will auto-detect it's a Go project
5. Add environment variables:
   - `ENVIRONMENT=prod`
   - `DISCORD_TOKEN=<your-token>`
   - `DISCORD_CLIENT_ID=<your-client-id>`
   - `DISCORD_CLIENT_SECRET=<your-client-secret>`
   - `JWT_SECRET=<generate-a-secure-32-char-string>`
   - `FRONTEND_URL=https://prizebot.dev`
   - `API_BASE_URL=https://<your-railway-domain>.railway.app`
   - `DB_HOST=<your-db-host>`
   - `DB_PASSWORD=<your-db-password>`
   - `STRIPE_TOKEN=<your-stripe-token>`
   - `SECURE_COOKIES=true`
   - `PORT=8080`
6. Deploy!
7. Note your Railway URL (e.g., `your-app.railway.app`)

#### Option B: Render

1. Go to [render.com](https://render.com)
2. Create a new "Web Service"
3. Connect your GitHub repository
4. Configure:
   - Build Command: `go build -o main`
   - Start Command: `./main`
   - Add the same environment variables as above
5. Deploy and note your Render URL

#### Option C: Custom Subdomain (Recommended for Production)

Instead of using the Railway/Render URL directly, set up a custom subdomain:

1. In your DNS settings for `prizebot.dev`, add a CNAME record:
   - Name: `api`
   - Value: `<your-railway-or-render-domain>`
   - TTL: 3600
2. Update your backend's `FRONTEND_URL` env var to `https://prizebot.dev`
3. Update your backend's `API_BASE_URL` env var to `https://api.prizebot.dev`
4. Your API will be accessible at `https://api.prizebot.dev`

### 2. Deploy Frontend to Vercel

1. Go to [vercel.com](https://vercel.com) and sign in with GitHub
2. Click "Add New Project"
3. Import your `Prize-Bot-V1` repository
4. Configure:
   - **Root Directory**: `frontend`
   - **Framework Preset**: SvelteKit (should auto-detect)
   - **Build Command**: `npm run build` (default)
   - **Output Directory**: `.svelte-kit` (default)
5. Add Environment Variable:
   - Name: `VITE_API_URL`
   - Value: `https://api.prizebot.dev` (or your backend URL without trailing slash)
6. Deploy!

### 3. Configure Custom Domain on Vercel

1. In your Vercel project settings, go to "Domains"
2. Add `prizebot.dev`
3. Follow Vercel's instructions to update your DNS:
   - Add an A record pointing to Vercel's IP
   - Or add a CNAME record (if using www or subdomain)
4. Wait for DNS propagation (can take up to 48 hours, usually much faster)

### 4. Update Discord OAuth Redirect URIs

1. Go to [Discord Developer Portal](https://discord.com/developers/applications)
2. Select your application
3. Go to "OAuth2" → "Redirects"
4. Add these redirect URLs:
   - `https://prizebot.dev/callback`
   - `https://api.prizebot.dev/api/auth/discord/callback`
5. Save changes

## Testing

After deployment:
1. Visit `https://prizebot.dev`
2. Click "Continue with Discord"
3. Authorize the application
4. You should be redirected back to your dashboard

## Troubleshooting

### CORS Errors
- Ensure `FRONTEND_URL` in backend matches your frontend domain exactly
- Check that backend is deployed and accessible

### Login Redirect Fails
- Verify Discord OAuth redirect URIs are correct
- Check that `DISCORD_CLIENT_ID` and `DISCORD_CLIENT_SECRET` are set correctly
- Ensure `VITE_API_URL` in frontend points to correct backend URL

### 401 Unauthorized Errors
- Check that cookies are working (requires `https://` for `SECURE_COOKIES=true`)
- Verify `JWT_SECRET` is set in backend

### Backend Connection Fails
- Verify `VITE_API_URL` in Vercel matches your backend URL
- Test backend directly: `curl https://api.prizebot.dev/health`
- Check Railway/Render logs for errors

## Local Development

For local development, the Vite proxy handles routing:
- Frontend runs on `http://localhost:5173`
- Backend runs on `http://localhost:8080`
- Proxy configured in `frontend/vite.config.ts` redirects `/api` to backend
- No `VITE_API_URL` needed locally

## Environment Variables Summary

### Backend Environment Variables
```
ENVIRONMENT=prod
DISCORD_TOKEN=<your-discord-bot-token>
DISCORD_CLIENT_ID=<your-oauth-client-id>
DISCORD_CLIENT_SECRET=<your-oauth-client-secret>
JWT_SECRET=<32-character-random-string>
FRONTEND_URL=https://prizebot.dev
API_BASE_URL=https://api.prizebot.dev
DB_HOST=<database-host>
DB_PASSWORD=<database-password>
STRIPE_TOKEN=<your-stripe-key>
PORT=8080
```

### Frontend Environment Variables (Vercel)
```
VITE_API_URL=https://api.prizebot.dev
```

## Notes

- The frontend uses `VITE_API_URL` environment variable to determine backend location
- In development, it falls back to `/api` which is proxied to `localhost:8080`
- In production, it uses the full backend URL
- Backend CORS is configured to only allow requests from `FRONTEND_URL`
- Cookies require HTTPS in production (`SECURE_COOKIES=true`)

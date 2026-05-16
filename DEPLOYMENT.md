# Obelisk MCP Server - Cloud Deployment Guide

This guide explains how to deploy Obelisk MCP Server to the cloud so users can access it remotely via HTTP/SSE.

## Overview

Obelisk MCP Server supports two transport modes:

1. **stdio** (default): Local JSON-RPC over stdin/stdout - for local IDE integration
2. **http**: HTTP/SSE transport - for cloud deployment and remote access

## Quick Start - Deploy to Render

### Prerequisites

- GitHub account
- Render account (free tier available at [render.com](https://render.com))
- Google Gemini API key (get one at [aistudio.google.com](https://aistudio.google.com/app/apikey))

### Step 1: Push to GitHub

1. Ensure your code is committed to a GitHub repository
2. Push your latest changes:

```bash
git add .
git commit -m "Add HTTP/SSE support for cloud deployment"
git push origin main
```

### Step 2: Deploy to Render

#### Option A: Using render.yaml (Recommended)

1. Go to [Render Dashboard](https://dashboard.render.com/)
2. Click **New** → **Blueprint**
3. Connect your GitHub repository
4. Render will automatically detect `render.yaml` and configure the service
5. Add your environment variables:
    - `GEMINI_API_KEY`: Your Google Gemini API key
    - `GOOGLE_API_KEY`: (Optional) Alternative to GEMINI_API_KEY
6. Click **Apply** to deploy

#### Option B: Manual Setup

1. Go to [Render Dashboard](https://dashboard.render.com/)
2. Click **New** → **Web Service**
3. Connect your GitHub repository
4. Configure the service:
    - **Name**: `obelisk-mcp-server`
    - **Runtime**: Go
    - **Build Command**: `go build -o obelisk-server .`
    - **Start Command**: `./obelisk-server mcp --http`
    - **Plan**: Free
5. Add environment variables:
    - `PORT`: `10000` (Render default)
    - `GEMINI_API_KEY`: Your API key
    - `OBELISK_MODEL`: `gemini-2.0-flash-exp` (optional)
6. Click **Create Web Service**

### Step 3: Get Your Service URL

After deployment completes (2-3 minutes), Render will provide a URL like:

```
https://obelisk-mcp-server.onrender.com
```

Your SSE endpoint will be:

```
https://obelisk-mcp-server.onrender.com/sse
```

### Step 4: Test Your Deployment

Check the health endpoint:

```bash
curl https://obelisk-mcp-server.onrender.com/health
```

Expected response:

```json
{
	"status": "ok",
	"server": "obelisk-mcp-server",
	"version": "1.x.x"
}
```

## Client Configuration

### For Bob IDE

Add to your Bob configuration file:

```json
{
	"mcpServers": {
		"obelisk-cloud": {
			"url": "https://obelisk-mcp-server.onrender.com/sse"
		}
	}
}
```

### For Claude Desktop

Add to `claude_desktop_config.json`:

```json
{
	"mcpServers": {
		"obelisk": {
			"url": "https://obelisk-mcp-server.onrender.com/sse"
		}
	}
}
```

### For Cline (VS Code Extension)

Add to Cline settings:

```json
{
	"mcpServers": {
		"obelisk": {
			"url": "https://obelisk-mcp-server.onrender.com/sse"
		}
	}
}
```

## Alternative Deployment Options

### Deploy to Fly.io

1. Install Fly CLI: https://fly.io/docs/hands-on/install-flyctl/
2. Login: `flyctl auth login`
3. Create app: `flyctl launch`
4. Set secrets:
    ```bash
    flyctl secrets set GEMINI_API_KEY=your-api-key-here
    ```
5. Deploy: `flyctl deploy`

### Deploy to Railway

1. Go to [Railway](https://railway.app/)
2. Click **New Project** → **Deploy from GitHub repo**
3. Select your repository
4. Add environment variables in the Railway dashboard
5. Railway will auto-deploy

### Deploy to Heroku

1. Install Heroku CLI: https://devcenter.heroku.com/articles/heroku-cli
2. Login: `heroku login`
3. Create app: `heroku create obelisk-mcp-server`
4. Set config:
    ```bash
    heroku config:set GEMINI_API_KEY=your-api-key-here
    ```
5. Deploy:
    ```bash
    git push heroku main
    ```

### Deploy with Docker

Build and run locally:

```bash
docker build -t obelisk-mcp-server .
docker run -p 8080:8080 \
  -e GEMINI_API_KEY=your-api-key-here \
  obelisk-mcp-server
```

Deploy to any container platform (AWS ECS, Google Cloud Run, Azure Container Instances, etc.)

## Local Testing

Test HTTP mode locally before deploying:

```bash
# Set your API key
export GEMINI_API_KEY=your-api-key-here

# Run in HTTP mode
obelisk mcp --http --port 8080

# In another terminal, test the health endpoint
curl http://localhost:8080/health

# Test SSE endpoint
curl http://localhost:8080/sse
```

## Environment Variables

| Variable         | Required | Default                | Description                           |
| ---------------- | -------- | ---------------------- | ------------------------------------- |
| `PORT`           | No       | `8080`                 | HTTP server port                      |
| `GEMINI_API_KEY` | Yes\*    | -                      | Google Gemini API key for AI analysis |
| `GOOGLE_API_KEY` | Yes\*    | -                      | Alternative to GEMINI_API_KEY         |
| `OBELISK_MODEL`  | No       | `gemini-2.0-flash-exp` | AI model to use                       |

\*At least one API key is required for AI-powered features

## Available Endpoints

| Endpoint  | Method | Description                        |
| --------- | ------ | ---------------------------------- |
| `/sse`    | GET    | SSE endpoint for MCP communication |
| `/health` | GET    | Health check endpoint              |
| `/`       | GET    | Alias for `/sse`                   |

## Available MCP Tools

Once connected, clients can use these tools:

- `scan_project` - Full project health scan
- `check_security` - Security vulnerability scan
- `analyze_complexity` - Code complexity analysis
- `track_tech_debt` - Technical debt tracking
- `audit_dependencies` - Dependency audit
- `get_health_report` - AI-powered health assessment

## Available MCP Resources

- `obelisk://scan/latest` - Latest scan results
- `obelisk://health/score` - Project health score
- `obelisk://findings/security` - Security findings
- `obelisk://findings/quality` - Code quality findings
- `obelisk://findings/architecture` - Architecture findings

## Troubleshooting

### Server won't start

- Check logs in Render dashboard
- Verify `GEMINI_API_KEY` is set correctly
- Ensure Go version is 1.26 or higher

### Connection timeout

- Render free tier may sleep after inactivity (first request takes ~30s)
- Consider upgrading to paid tier for always-on service

### API key errors

- Verify your Gemini API key is valid
- Check API quota at [Google AI Studio](https://aistudio.google.com/)

### CORS issues

- The server includes CORS headers for browser clients
- If issues persist, check your client configuration

## Security Considerations

1. **API Key Protection**: Never commit API keys to Git. Use environment variables.
2. **HTTPS**: Always use HTTPS in production (Render provides this automatically)
3. **Rate Limiting**: Consider adding rate limiting for public deployments
4. **Authentication**: For production use, consider adding authentication middleware

## Cost Considerations

### Render Free Tier

- 750 hours/month free
- Sleeps after 15 minutes of inactivity
- Suitable for personal use and testing

### Render Paid Tier ($7/month)

- Always-on service
- No sleep delays
- Better for production use

## Support

For issues or questions:

- GitHub Issues: https://github.com/Swif7ify/Obelisk-CLI/issues
- Documentation: https://github.com/Swif7ify/Obelisk-CLI

## Made with Bob

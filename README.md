# Stockyard Prairie

**Web analytics — one script tag, no cookies, GDPR compliant, page views, referrers, device breakdown**

Part of the [Stockyard](https://stockyard.dev) family of self-hosted developer tools.

## Quick Start

```bash
docker run -p 9340:9340 -v prairie_data:/data ghcr.io/stockyard-dev/stockyard-prairie
```

Or with docker-compose:

```bash
docker-compose up -d
```

Open `http://localhost:9340` in your browser.

## Configuration

| Variable | Default | Description |
|----------|---------|-------------|
| `PORT` | `9340` | HTTP port |
| `DATA_DIR` | `./data` | SQLite database directory |
| `PRAIRIE_LICENSE_KEY` | *(empty)* | Pro license key |

## Free vs Pro

| | Free | Pro |
|-|------|-----|
| Limits | 1 site, 10k pageviews/mo | Unlimited sites and pageviews |
| Price | Free | $4.99/mo |

Get a Pro license at [stockyard.dev/tools/](https://stockyard.dev/tools/).

## Category

Creator & Small Business

## License

Apache 2.0

# solar-scraper

An application to scrape metrics data from a solar inverter and store it in influxdb.

## Development

Create influxdb and grafana containers

```bash
chmod +x ./test/environment/scripts/test.sh
docker compose --file ./test/environment/docker-compose.yml up --detach 
```

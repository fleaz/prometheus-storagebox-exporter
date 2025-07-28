# prometheus-storagebox-exporter

This tool talks to the [Hetzner API](https://docs.hetzner.cloud/reference/hetzner) and
gets a list of all [Storage Boxes](https://www.hetzner.de/storage/storage-box) in your account and exports their
statistics as Prometheus metrics on `<host>:9509/metrics`.

## Authentication
In the Cloud Console, switch to the project where your storage box is part of and generate a read-only API Token.
See the [Authentication](https://docs.hetzner.cloud/reference/hetzner#authentication) chapter in the hetzner api docs.
This token then needs to be provided as an environment variable.


## Exported Metrics 
```
# HELP storagebox_disk_quota Total diskspace in Bytes
# TYPE storagebox_disk_quota gauge
storagebox_disk_quota{id="13374223",name="Bart",product="BX21",server="u123456.your-storagebox.de"} 5.49755813888e+12
# HELP storagebox_disk_usage Total used diskspace in Bytes
# TYPE storagebox_disk_usage gauge
storagebox_disk_usage{id="13374223",name="Bart",product="BX21",server="u123456.your-storagebox.de"} 4.21402771456e+11
# HELP storagebox_disk_usage_data Used diskspace by files in Bytes
# TYPE storagebox_disk_usage_data gauge
storagebox_disk_usage_data{id="13374223",name="Bart",product="BX21",server="u123456.your-storagebox.de"} 4.21271699456e+11
# HELP storagebox_disk_usage_snapshots Used diskspace by snapshots in Bytes
# TYPE storagebox_disk_usage_snapshots gauge
storagebox_disk_usage_snapshots{id="13374223",name="Bart",product="BX21",server="u123456.your-storagebox.de"} 1.31072e+08
```

# Usage
You need to provide your api token via environment variables. So after compiling the binary you could run it like

```
HETZNER_TOKEN='...' ./prometheus-storagebox-exporter
```

then visit [localhost:9505/metrics](http://localhost:9505/metrics)

To expose the exporter on a different addresse, you can set LISTEN_ADDRESS and set `host:port`.

# Running as Docker container
This exporter can be run as Docker container as well.

First you need to provide your credentials in the .env file.

Then, either build and run the image manually with

```sh
docker run -d --name storagebox-exporter --env-file .env fleaz/prometheus-storagebox-exporter
```

or use compose:

```
docker compose up -d
```

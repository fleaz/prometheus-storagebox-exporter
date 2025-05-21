# prometheus-storagebox-exporter

This tool talks to the [Hetzner API](https://robot.your-server.de/doc/webservice/de.html#storage-box) and
gets a list of all [Storage Boxes](https://www.hetzner.de/storage/storage-box) in your account and exports their
statistics as Prometheus metrics on `<host>:9509/metrics`.

## Authentication
To use the Robot API, you first need to create a "WebService" Account over
[here](https://robot.hetzner.com/preferences/index).
After choosing a passwor, you will get an email with the random username. These credentials can then be used without
this tool.

## Exported Metrics 
```
# HELP storagebox_disk_quota Total diskspace in MB
# TYPE storagebox_disk_quota gauge
storagebox_disk_quota{id="1234",name="Backup",product="BX10",server="u12345.your-storagebox.de"} 102400
# HELP storagebox_disk_usage Total used diskspace in MB
# TYPE storagebox_disk_usage gauge
storagebox_disk_usage{id="1234",name="Backup",product="BX10",server="u12345.your-storagebox.de"} 23256
# HELP storagebox_disk_usage_data Used diskspace by files in MB
# TYPE storagebox_disk_usage_data gauge
storagebox_disk_usage_data{id="1234",name="Backup",product="BX10",server="u12345.your-storagebox.de"} 23256
# HELP storagebox_disk_usage_snapshots Used diskspace by snapshots in MB
# TYPE storagebox_disk_usage_snapshots gauge
storagebox_disk_usage_snapshots{id="1234",name="Backup",product="BX10",server="u12345.your-storagebox.de"} 0
```

# Usage
You need to provide your credentials (username and passwort) for the WebService account via environment variables. So
after compiling the binary you could run it like

```
HETZNER_USER='...' HETZNER_PASS='...' ./prometheus-storagebox-exporter
```

then visit [localhost:9505/metrics](http://localhost:9505/metrics)

# Running as Docker container
This exporter can be run as Docker container as well.

First you need to provide your credentials in the .env file.

Then, either build and run the image manually with

```sh
docker build --tag storagebox-exporter .
docker run -d --name storagebox-exporter storagebox-exporter --env-file .env
```

or use compose:

```
docker compose up -d
```

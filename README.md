## README for Zultys CRM-Sync Tool

### Overview

The `zultys_crm-sync` tool is designed to synchronize call data from a Zultys MX system to a CRM system. It supports various modes of operation, including fetching call records directly from the Zultys MX via SMDR (Station Message Detail Recording) or from FTP-exported call logs.

### Configuration

To configure the tool, you need to create a JSON configuration file. Below is a template for the configuration file (`config.json`):

```json
{
  "mx_username": "Administrator",
  "mx_password": "zultys",
  "mx_addr": "zultys.example.com",
  "listen_addr": "0.0.0.0:2121",
  "mode": "smdr",
  "ftp_username": "zultys_crm-sync",
  "ftp_password": "P@ssw0rd",
  "crm_type": "hubspot",
  "crm_apikey": "apikey_here",
  "zultys_users_file": "./zultys_users.json",
  "crm_users_file": "./crm_users.json",
  "timestamp_region": "America/Chicago",
  "timestamp_file": "./timestamp.json"
}
```

Replace the placeholders with your actual data. For example, set `mx_addr` to your Zultys MX address, and `cmr_apikey` to your CRM API key.

### Running the Tool

The tool is run via the command line as follows:

```bash
/opt/zultys_crm-sync/zultys_crm-sync -config=/path/to/your/config.json
```

Replace `/path/to/your/config.json` with the actual path to your configuration file.

### Service File Creation

To manage the `zultys_crm-sync` tool as a service, you can create a systemd service file. Here's an example service file (`zultys_crm-sync.service`):

```ini
[Unit] 
Description=Zultys CRM-Sync Service 
After=network.target 

[Service] 
Type=simple 
User=root 
ExecStart=/opt/zultys_crm-sync/zultys_crm-sync -config=/opt/crm_sync/tops-test/config.json  

[Install] 
WantedBy=multi-user.target
```

Replace the `ExecStart` path and config file location as needed.

### Additional Information

For more detailed information on configuring and extending `zultys_crm-sync`, please refer to the [GitHub repository](https://github.com/sagostin/zultys_crm-sync).



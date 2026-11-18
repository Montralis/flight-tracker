## Systemd Service (Linux)

The Flight Tracker can run as a background service using systemd.

The provided `flight-tracker@.service` template supports automatic startup and restart.

> Note: The service file is configured for Linux ARM64 (e.g. Raspberry Pi). If you are using a different architecture or
> installation path, adjust the ExecStart and WorkingDirectory entries accordingly.

### Installation

Copy the service file:

```bash
sudo cp flight-tracker@.service /etc/systemd/system/
```

Reload systemd:

```bash
sudo systemctl daemon-reload
```

Enable and start the service (replace `USER` with your Linux username):

```bash
sudo systemctl enable --now flight-tracker@USER
```

### Management

Check service status:

```bash
sudo systemctl status flight-tracker@USER
```

View live logs:

```bash
sudo journalctl -u flight-tracker@USER -f
```

Restart the service:

```bash
sudo systemctl restart flight-tracker@USER
```

The service starts automatically on boot and restarts if the application crashes.
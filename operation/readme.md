## Systemd Service (Linux)

The Flight Tracker can run as a background service using systemd.

The provided `flight-tracker@.service` template supports automatic startup and restart.

> **Note:** The service file is configured for Linux ARM64 (e.g. Raspberry Pi). If you are using a different
> architecture or installation path, adjust the `ExecStart` and `WorkingDirectory` entries accordingly.

### Installation

**1. Download the latest binary**

Download the appropriate binary from the [GitHub Releases](../../releases).

The binary is named according to the architecture and release tag.

Examples:

- `Flight_Tracker_ARM64_r1.1`
- `Flight_Tracker_AMD64_r1.1`

**2. Install the binary**

Copy the downloaded binary into your installation directory:

```bash
mkdir -p ~/flight-tracker
cp Flight_Tracker_* ~/flight-tracker/
chmod +x ~/flight-tracker/Flight_Tracker_*
```

**3. Install the systemd service**

```bash
sudo cp flight-tracker@.service /etc/systemd/system/
sudo systemctl daemon-reload
```

**4. Enable and start the service**

Replace `USER` with your Linux username:

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
# Robotmon monitor
This acts as a script callback endpoint to ensure running Robotmon scripts and restart their VMs in case a script does 
not call the endpoint within a given time. Works with LDPlayer 5, but should also work with Nox.

> [!WARNING]  
> At least on my computer, the pre-released `monitor.exe` built via Github Actions raised a malware warning within 
> Windows Defender. I consider that a false detection and raised a rescan submission at Microsoft. Use at your own risk!

## Basic usage

To start the monitor, open a command line interface where `monitor.exe` is located and enter this:

### Powershell
```powershell
.\monitor.exe -port 12345 -process "C:\LDPlayer4.0\LDPlayer\dnplayer.exe index=0" -title "LDPlayer-0" -debug true
```

### Command Prompt
```
monitor.exe -port 12345 -process "C:\LDPlayer4.0\LDPlayer\dnplayer.exe index=0" -title "LDPlayer-0" -debug true
```

- `port` (required) needs to be a free TCP port. This will be used by the Robotmon script to contact the monitor.
- `process` (required) needs to be the exact Windows command to start the Android emulator
- `title` (optional) should be the title of the emulator window
- `debug` shows more information at runtime like, for example, incoming requests

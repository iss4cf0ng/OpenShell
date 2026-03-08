# OpenShell

![status](https://img.shields.io/badge/status-development-orange)
![language](https://img.shields.io/badge/language-Go-blue)
![license](https://img.shields.io/badge/license-MIT-green)

OpenShell is a lightweight, open-source reverse shell management server written in Go.  
It allows users to establish reverse shell channels and interact with them through a simple web-based graphical user interface (GUI).

If you encounter bugs or have suggestions, feel free to open an issue.

If you find this project useful, a ⭐ on the repository would be greatly appreciated!

<p align="center">
<img src="https://iss4cf0ng.github.io/images/meme/nagisa_neko.png" width=200>
</p>

## Disclaimer
This project was developed as part of my personal interest in studying cybersecurity. However, it may potentially be misused for malicious purposes.  
Please do NOT use this tool for any illegal activities.  
The author is not responsible for any misuse of this software.

<p align="center">
  <img src="https://iss4cf0ng.github.io/images/meme/mika_punch.jpg" width="500">
</p>

## Background
The idea of this project originated while I was learning web penetration testing and working on one of my side projects, [Eden-RAT](https://github.com/iss4cf0ng/Eden-RAT).

Compared to Eden-RAT, OpenShell focuses only on establishing reverse shell channels and providing a simple interactive interface.  
Additional features may be added in the future.

## Features
- Lightweight reverse shell server written in Go
- Web-based GUI for interactive shell sessions
- Simple reverse shell command to connect
- Designed for learning and research purposes

## Installation
Download the lastest release:
```
$ wget https://github.com/iss4cf0ng/OpenShell/releases/latest/download/openshell-server-linux-amd64.tar.gz
$ tar -xzf openshell-server-linux-amd64.tar.gz
$ ./openshell-server
```

## Quick Start
1. Create `key.pem` and `cert.pem`:
```
$ openssl req -x509 -newkey rsa:2048 -keyout key.pem -out cert.pem -days 365 -nodes
```

2. Start the server
```
$ ./openshell-server
```
3. Open your browser and navigate to:
```
https://localhost:8080
```
4. Login using the default password (Please change it in production):
```
123456
```
5. Generate a reverse shell command and execute it on the target machine.
6. Have fun!

## Architecture
OpenShell consists of two main components:
- **Go Server**
  - Handles reverse shell connections
  - Manages sessions
  - Provides WebSocket interface

- **Web GUI**
  - Built with xterm.js
  - Allows interactive shell sessions in browser
  - Supports multiple terminal tabs

<p align="center">
  <img src="https://iss4cf0ng.github.io/images/article/2026-3-8-OpenShell/architecture.png" width=800>
</p>

## Tech Stack
- Go
- WebSocket
- xterm.js
- PTY (github.com/creack/pty)

## Roadmap
Future improvements may include:
- File upload / download
- Agent management
- Session persistence
- Authentication improvements

## ScreenShot
### Server
<p align="center">
  <img src="https://iss4cf0ng.github.io/images/article/2026-3-8-OpenShell/6.png" width=700>
</p>

### Login (Default Password: `123456`. Please change it in production)
<p align="center">
  <img src="https://iss4cf0ng.github.io/images/article/2026-3-8-OpenShell/2.png" width=700>
</p>

### Create Reverse Shell Command
<p align="center">
  <img src="https://iss4cf0ng.github.io/images/article/2026-3-8-OpenShell/3.png" width=400>
</p>

### Sessions
<p align="center">
  <img src="https://iss4cf0ng.github.io/images/article/2026-3-8-OpenShell/4.png" width=800>
</p>

### Interactive Reverse Shell (capable to use `vim`, `ssh`, `nano`, `nslookup`, `nc`, etc.)
#### `vim`
<p align="center">
  <img src="https://iss4cf0ng.github.io/images/article/2026-3-8-OpenShell/1.png" width=1000>
</p>

#### `nslookup`
<p align="center">
  <img src="https://iss4cf0ng.github.io/images/article/2026-3-8-OpenShell/5.png" width=1000>
</p>

## License
This project is licensed under the MIT License.

## Acknowledgement
- [Reverse Shell Cheatsheet (by swisskyrepo)](https://swisskyrepo.github.io/InternalAllTheThings/cheatsheets/shell-reverse-cheatsheet/)


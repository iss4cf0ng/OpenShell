# OpenShell

![status](https://img.shields.io/badge/status-development-orange)
![language](https://img.shields.io/badge/language-Go-blue)
![license](https://img.shields.io/badge/license-MIT-green)

OpenShell is a lightweight, open-source reverse shell management server written in Go.  
It allows users to establish reverse shell channels and interact with them through a simple web-based graphical user interface (GUI).

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
$./openshell-server
```

## Acknowledgement
- [Reverse Shell Cheatsheet (by swisskyrepo)](https://swisskyrepo.github.io/InternalAllTheThings/cheatsheets/shell-reverse-cheatsheet/)

# ScreenShot
## Login (Default Password: `123456`)
<p align="center">
  <img src="https://iss4cf0ng.github.io/images/article/2026-3-8-OpenShell/2.png" width=700>
</p>

## Create Reverse Shell Command
<p align="center">
  <img src="https://iss4cf0ng.github.io/images/article/2026-3-8-OpenShell/3.png" width=400>
</p>

# Sessions
<p align="center">
  <img src="https://iss4cf0ng.github.io/images/article/2026-3-8-OpenShell/4.png" width=800>
</p>

## Interactive Reverse Shell (capable to use `vim`, `ssh`, `nano`, `nslookup`, `nc`, etc.)
### `vim`
<p align="center">
  <img src="https://iss4cf0ng.github.io/images/article/2026-3-8-OpenShell/1.png" width=1000>
</p>

### `nslookup`
<p align="center">
  <img src="https://iss4cf0ng.github.io/images/article/2026-3-8-OpenShell/5.png" width=1000>
</p>

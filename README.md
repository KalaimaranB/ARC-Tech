# ARCTech - Advanced Reconnaissance Clone Technology

**ARCTech** is an automated penetration testing framework inspired by *Star Wars: The Clone Wars*. It is designed to streamline reconnaissance tasks and vulnerability assessments with an engaging and user-friendly approach. ARCTech provides a modular and customizable interface for penetration testers and security researchers, leveraging tools like Nmap, Hydra, and DirBuster. 
## Goals

ARCTech was created as a learning project with the following objectives:

- **Learn Go**: Gain proficiency in Go programming for backend development and task orchestration.
- **Manage Cross-Language Interactions**: Explore how different programming languages (Go, Python) can work together seamlessly.
- **Understand the Reconnaissance Phase**: Deepen understanding of the enumeration and reconnaissance stages of cybersecurity.
- **Use Kali Linux**: Become proficient in Kali Linux and its ecosystem of penetration testing tools.
- **Beginner-Friendly Automation**: Build a tool that automates the enumeration process, making it accessible to beginners while providing flexibility for intermediate users to customize and extend functionality.
- **AI & ChatGPT**: To understand the extent to which modern AI tools can help automate both the creation of this project by helping with not just writing code, but higher level overview as well implementing AI APIs to analyze tasks.

## Features

- **Modular Architecture**: Easily expand and integrate new tools while customizing existing features.
- **Interactive CLI Wizard**: An intuitive command-line interface guides users through selecting scan options. Categories include:
  - **Network Evasion**
  - **Packet Manipulation**
  - **Speed Optimization**
  - **Detailed Scanning**
  - **Aggressive Scans**
- **Customizable Nmap Scans**: Configure Nmap flags for tailored scans, including options for scan speed, stealthiness, and detailed service discovery.
- **Advanced Tool Integration**: Seamlessly integrates with tools like Hydra (for brute-force attacks) and DirBuster (for directory and file enumeration).
- **Containerized with Docker**: Ensures consistent behavior across environments and simplifies setup and deployment.
- **Star Wars Theming**: Adds a fun, themed interface inspired by *Star Wars: Clone Wars*.
- **More to come** Other modules will be soon added...

## Requirements

ARCTech is designed for **Kali Linux** (or a similar Linux distribution) with the following prerequisites:

- **Kali Linux** (or a compatible Linux distro)
  - Pre-installed Nmap
- **Go**
  - For backend development and task orchestration.
- **Python**
  - For scripting and managing tools like Hydra and DirBuster.

### Dependencies
- **Nmap**: For network reconnaissance and scanning.
- **Hydra**: For performing brute-force password attacks (e.g., HTTP, SSH).
- **DirBuster**: For brute-forcing directories and files on web servers.


### Notes
This is an ongoing project for me to learn more about reconnaissance in cybersecurity. The goal is to automate as much of the early phase as possible to allow even beginners to get a quick understanding of the "landscape" they are targeting. 


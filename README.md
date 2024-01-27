# Warp-Plus-Go

Warp-Plus-Go is an open-source implementation of Cloudflare's Warp, enhanced with Psiphon integration for circumventing censorship. This project aims to provide a robust and cross-platform VPN solution that can automatically switch to Psiphon when Warp connectivity is restricted or blocked.

## Features

- **Warp Integration**: Leverages Cloudflare's Warp to provide a fast and secure VPN service.
- **Psiphon Chaining**: Integrates with Psiphon for censorship circumvention, allowing seamless access to the internet in restrictive environments.
- **Cross-Platform Support**: Designed to work on multiple platforms, offering the same level of functionality and user experience.
- **SOCKS5 Proxy Support**: Includes a SOCKS5 proxy for secure and private browsing.
- **Verbose Logging**: Optional verbose logging for troubleshooting and performance monitoring.

## Getting Started

### Prerequisites

- Go (version 1.15 or later recommended)
- Basic understanding of VPN and proxy configurations

### Installation

1. Clone the repository:
   ```bash
   git clone https://github.com/yourusername/warp-plus-go.git
   cd warp-plus-go
   ```

2. Build the project:
   ```bash
   go build
   ```

### Usage

Run the application with the following command:

```bash
./warp-plus-go [-v] [-b addr:port] [-c config-file-path] [-e warp-ip] [-k license-key] [-h listen-host] [-p listen-port] [-country country-code] [-cfon] 
```

- `-v`: Enable verbose logging.
- `-b`: Set the SOCKS bind address (default: `127.0.0.1:8086`).
- `-c`: Path to the Warp configuration file.
- `-e`: Specify the Warp endpoint IP.
- `-k`: Your Warp license key.
- `-h`: makes wiresock listen on specified host instead of `127.0.0.1`
- `-p`: specifies port for wiresock instead of the random default
- `-country`: ISO 3166-1 alpha-2 country code for Psiphon.
- `-cfon`: Enable Psiphon over Warp.

### Termux

```
bash <(curl -fsSL https://raw.githubusercontent.com/Ptechgithub/wireguard-go/master/termux.sh)
```
![1](https://github.com/Ptechgithub/configs/blob/main/media/18.jpg?raw=true)

- بعد از نصب برای اجرای مجدد فقط کافیه که `warp` یا `usef` یا `./warp` را وارد کنید . 
- روی برخی از مدل های قدیمی تر روش دوم یعنی arm 64 را انتخاب کنید. 
- برای نمایش راهنما ` warp -h` را وارد کنید. 
- ای پی و پورت `127.0.0.1:8086`پروتکل socks
- در روش warp به warp plus مقدار account id را وارد میکنید و با این کار هر 20 ثانیه 1 GB به اکانت شما اضافه میشود. 


## Contributing

Contributions to Warp-Plus-Go are welcome. Please read our [contributing guidelines](CONTRIBUTING.md) for more information.

## Acknowledgements

- Cloudflare Warp
- Psiphon
- All contributors and supporters of this project

## License

    Copyright (C) 2017-2023 WireGuard LLC. All Rights Reserved.
    
    Permission is hereby granted, free of charge, to any person obtaining a copy of
    this software and associated documentation files (the "Software"), to deal in
    the Software without restriction, including without limitation the rights to
    use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies
    of the Software, and to permit persons to whom the Software is furnished to do
    so, subject to the following conditions:
    
    The above copyright notice and this permission notice shall be included in all
    copies or substantial portions of the Software.
    
    THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
    IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
    FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
    AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
    LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
    OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
    SOFTWARE.

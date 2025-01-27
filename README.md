# Warp-Plus

Warp-Plus is an open-source implementation of Cloudflare's Warp, enhanced with Psiphon integration for circumventing censorship. This project aims to provide a robust and cross-platform VPN solution that can use psiphon on top of warp and warp-in-warp for changing the user virtual nat location.

## Features

- **Warp Integration**: Leverages Cloudflare's Warp to provide a fast and secure VPN service.
- **Psiphon Chaining**: Integrates with Psiphon for censorship circumvention, allowing seamless access to the internet in restrictive environments.
- **Warp in Warp Chaining**: Chaning two instances of warp together to bypass location restrictions.
- **SOCKS5 Proxy Support**: Includes a SOCKS5 proxy for secure and private browsing.

## Getting Started

### Prerequisites

- [Download the latest version from the releases page](https://github.com/bepass-org/warp-plus/releases)
- Basic understanding of VPN and proxy configurations

### Usage

```
NAME
  warp-plus

FLAGS
  -4                       only use IPv4 for random warp endpoint
  -6                       only use IPv6 for random warp endpoint
  -v, --verbose            enable verbose logging
  -b, --bind STRING        socks bind address (default: 127.0.0.1:8086)
  -e, --endpoint STRING    warp endpoint
  -k, --key STRING         warp key
      --dns STRING         DNS address (default: 1.1.1.1)
      --gool               enable gool mode (warp in warp)
      --cfon               enable psiphon mode (must provide country as well)
      --country STRING     psiphon country code (valid values: [AT BE BG BR CA CH CZ DE DK EE ES FI FR GB HR HU IE IN IT JP LV NL NO PL PT RO RS SE SG SK UA US]) (default: AT)
      --scan               enable warp scanning
      --rtt DURATION       scanner rtt limit (default: 1s)
      --cache-dir STRING   directory to store generated profiles
      --fwmark UINT        set linux firewall mark for tun mode (default: 4981)
      --reserved STRING    override wireguard reserved value (format: '1,2,3')
      --wgconf STRING      path to a normal wireguard config
  -c, --config STRING      path to config file
      --version            displays version number
```

### Country Codes for Psiphon

- Austria (AT)
- Australia (AU)
- Belgium (BE)
- Bulgaria (BG)
- Canada (CA)
- Switzerland (CH)
- Czech Republic (CZ)
- Germany (DE)
- Denmark (DK)
- Estonia (EE)
- Spain (ES)
- Finland (FI)
- France (FR)
- United Kingdom (GB)
- Croatia (HR)
- Hungary (HU)
- Ireland (IE)
- India (IN)
- Italy (IT)
- Japan (JP)
- Latvia (LV)
- Netherlands (NL)
- Norway (NO)
- Poland (PL)
- Portugal (PT)
- Romania (RO)
- Serbia (RS)
- Sweden (SE)
- Singapore (SG)
- Slovakia (SK)
- United States (US)
![0](https://raw.githubusercontent.com/Ptechgithub/configs/main/media/line.gif)
### Termux

```
bash <(curl -fsSL https://raw.githubusercontent.com/bepass-org/warp-plus/master/termux.sh)
```
![1](https://github.com/Ptechgithub/configs/blob/main/media/18.jpg?raw=true)

- اگه حس کردی کانکت نمیشه یا خطا میده دستور `rm -rf .cache/warp-plus` رو بزن و مجدد warp رو وارد کن.
- بعد از نصب برای اجرای مجدد فقط کافیه که `warp` یا `usef` یا `./warp` یا `warp-plus`را وارد کنید. همش یکیه هیچ فرقی ندارد.
- اگر با 1 نصب نشد و خطا گرفتید ابتدا یک بار 3 را بزنید تا `Uninstall` شود سپس عدد 2 رو انتخاب کنید یعنی Arm.
- برای نمایش راهنما ` warp -h` را وارد کنید. 
- ای پی و پورت `127.0.0.1:8086`پروتکل socks
- در روش تبدیل اکانت  warp به warp plus (گزینه 6) مقدار ID را وارد میکنید. پس از اجرای warp دو اکانت برای شما ساخته شده که پس از انتخاب گزینه 6 خودش مقدار ID هر دو اکانت را پیدا میکند و شما باید هر بار یکی را انتخاب کنید و یا میتوانید با انتخاب manual مقدار ID دیگری را وارد کنید (مثلا برای خود برنامه ی 1.1.1.1 یا جای دیگر) با این کار هر 20 ثانیه 1 GB به اکانت شما اضافه میشود. و اکانت شما از حالت رایگان به پلاس تبدیل میشود. 
- برای تغییر  لوکیشن با استفاده از سایفون از طریق منو یا به صورت دستی (برای مثال به USA  از دستور  زیر استفاده کنید) 
- `warp --cfon --country US`
- برای اسکن ای پی سالم وارپ از دستور `warp --scan` استفاده کنید. 
- برای ترکیب (chain) دو کانفیگ برای تغییر لوکیشن از دستور `warp --gool` استفاده کنید. 

## Acknowledgements

- Cloudflare Warp
- Psiphon
- All contributors and supporters of this project

## $\color{gold} \Large \text{Warp-Plus}$
$\color{silver} \large \text{Many thanks to Usef Ghobadi and dear MarkPash}$

$\color{gray} \large \text{Warp-Plus is an open-source implementation of Cloudflare's Warp, enhanced with Psiphon integration}$
$\color{gray} \large \text{for circumventing censorship. This project aims to provide a robust and cross-platform VPN solution}$
$\color{gray} \large \text{that can use psiphon on top of warp and warp-in-warp for changing the user virtual nat location.}$


### $\color{gold} \large \text{Features}$
```
- Warp Integration: Leverages Cloudflare's Warp to provide a fast and secure VPN service.
- Psiphon Chaining: Integrates with Psiphon for censorship circumvention, allowing seamless access to the
  internet in restrictive environments.
- Warp in Warp Chaining: Chaining two instances of warp together to bypass location restrictions.
- SOCKS5 Proxy Support: Includes a SOCKS5 proxy for secure and private browsing.
```


<details>
<summary> What's changed in v1.2.1?</summary>

```
misc: replace stuff folder with better handling by @markpash in #80
misc: add --version flag by @markpash in #84
enhancement: Add container image action (#72) by @debMan in #85
Proxy optimizations and wireguard tweaks by @markpash in #89
actions: add mips32 builds by @markpash in #90
actions: add mips32 builds by @markpash in #90
Update dependencies (psiphon, quic, gvisor, go) by @markpash in #91
psiphon: store cache files in correct location by @markpash in #92
psiphon: add Croatia and Portugal by @markpash in #93
tun: add tun support (experimental) by @markpash in #79
feat: ✨ add configurable DNS addresses #30 by @pizokh in #103
tun: set dns on windows by @markpash in #104
wireguard: add generic wireguard support by @markpash in #108
linux: fix setting bind on non-tun mode by @markpash in #109
warp: rewrite API functions by @markpash in #110
wgconf: fix creating warp identities when not using warp by @markpash in #113
```

</details>


---
![main-menu](https://github.com/NiREvil/wireguard-go/assets/126243832/3aff5b35-e5cf-42ce-b07b-ddd95e180acc)

---
### $\color{gold} \large \text{Getting Started with Termux}$
**$\color{silver} \normalsize \text{Single Command install}$**

```bash
bash <(curl -fsSL https://raw.githubusercontent.com/NiREvil/wireguard-go/master/termux.sh)
```


- اگه حس کردید کانکت نمیشه یا خطا میده دستور rm -rf .cache/warp-plus رو بزنید تا کش ها پاک بشن مجدد warp رو اجرا کنید.
- بعد از نصب برای اجرای مجدد فقط کافیه که `warp` یا `usef` یا `warp-plus`را وارد کنید. همش یکیه هیچ فرقی نداره. 
- اگر با 1 نصب نشد و خطا گرفتید ابتدا یک بار 3 را بزنید تا `Uninstall` بشه سپس عدد 2 رو انتخاب کنید یعنی Arm (عموما برای گوشی های قدیمی و پردازنده ضعیف این گزینه جواب میده).
- برای نمایش راهنما ` warp -h` را وارد کنید.
-  دریافت بر روی آی‌پی و پورت 127.0.0.1:8086 پروتکل SOCKS یا کپی و وارد کردن این کانفیگ در نکوباکس و ...
`socks://Og==@127.0.0.1:8086#warp_(termux)`

- برای متوقف کردن وارپ در حال اجرا از کلیدهای ترکیبی `Ctrl + C` استفاده کنید و همچنین برای بستن نشست فعال و خارج شدن از ترماکس از کلیدهای `Ctrl + D` استفاده کنید.
- برای تکرار دوباره دستورهای وارد شده در مرحله قبل در ترماکس از فلش رو به بالا ↑ استفاده کنید.
- در روش تبدیل اکانت warp به warp plus (گزینه 6) مقدار ID نمایش داده شده را وارد می‌کنید. پس از اجرای warp دو اکانت برای شما ساخته شده که پس از انتخاب گزینه 6 خودش مقدار ID هر دو اکانت را پیدا می‌کند و شما باید هر بار یکی را انتخاب کنید و یا میتوانید با انتخاب manual مقدار ID دیگری را وارد کنید (مثلا برای خود برنامه ی 1.1.1.1 یا جای دیگر) با این کار هر 20 ثانیه 1 GB به اکانت شما اضافه شده و اکانت شما از حالت رایگان به پلاس تبدیل می‌شود.
و یا برای تبدیل شدن وارپ به وارپ پباس مستقیما از [ربات های تلگرام](https://t.me/F_NiREvil/5271) کلید لایسنس دریافت کرده و با دستور زیر
`warp -k yourlicensekey`
آنرا ارسال کنید.
- برای تغییر  لوکیشن با استفاده از سایفون از طریق منو یا به صورت دستی: برای مثال برای USA از دستور زیر استفاده کنید:
`warp --cfon --country US`
- یا برای مثال برای Estonia از دستور  زیر استفاده کنید:
`warp --cfon --country EE`
- برای اسکن ای پی سالم وارپ از دستور `warp --scan` استفاده کنید. 
- برای زنجیر کردن (proxy chain) دو کانفیگ وارپ باهم و تغییر آی پی به غیر ایران از دستور `warp --gool` استفاده کنید.


---

## $\color{gold} \large \text{video tutorials}$


**$\color{silver} \normalsize \text{Part 1:}$**
> **with Single Command install**

https://github.com/NiREvil/wireguard-go/assets/126243832/458edbef-1053-4605-9afd-fdf88ef2398c



**$\color{silver} \normalsize \text{Part 2:}$**
> **with: `warp --gool`**

https://github.com/NiREvil/wireguard-go/assets/126243832/07aa5d18-c4b1-43fb-9814-b23ac730afa9




**$\color{silver} \normalsize \text{Part 3:}$**
> **with: `warp --cfon --country EE`**

https://github.com/NiREvil/wireguard-go/assets/126243832/5054de35-e798-4603-ba7b-5c0accc0da07

---

**$\color{silver} \normalsize \text{Part 4:}$**
> **Using `Warp-plus in Windows`**
> **Note:
>   Just be aware that because the video is from a few months ago; when entering commands,
>  we used a single -, but now in the recent update we use two -- in Warp commands**

[![YouTube Video](https://img.youtube.com/vi/Me2PJTOgnaQ/0.jpg)](https://www.youtube.com/watch?v=Me2PJTOgnaQ)

---

### $\color{gold} \large \text{LAZY INSTALLATION}$

**$\color{silver} \normalsize \text{1. Clone this repository:}$**
   ```bash
   git clone https://github.com/bepass-org/warp-plus.git
   cd warp-plus
   ```

**$\color{silver} \normalsize \text{2. Build the project:}$**
   ```bash
   go build
   ```

### $\color{gold} \large \text{Usage}$

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
      --tun-experimental   enable tun interface (experimental)
      --fwmark UINT        set linux firewall mark for tun mode (default: 4981)
      --wgconf STRING      path to a normal wireguard config
  -c, --config STRING      path to config file
      --version            displays version number
```

### $\color{gold} \large \text{Country Codes for Psiphon}$

- Austria (AT)
- Belgium (BE)
- Bulgaria (BG)
- Brazil (BR)
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
- Hungary (HU)
- Croatia (HR)
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
- Ukraine (UA)
- United States (US)
---

---
### $\color{gold} \Large \text{Download BOX}$

| Source | Download Link's |
|:--------:| -------------:|
| F-Droid | [Get it Now](https://f-droid.org/en/packages/com.termux)|
|  GitHub :octocat:| [Get it Now](https://github.com/termux/termux-app/releases)|
|Google Play| [Get it Now](https://play.google.com/store/apps/details?id=com.termux)|
| How to fix the installation error of Termux packages on Android 5 & 6 | [Fix Errors](https://t.me/F_NiREvil/5040)

> I recommend installing GitHub for sure, as the Google Play version has many bugs and the F-Droid version is not compatible with most devices. The best choice is the GitHub repository of the app.
---


### $\color{gold} \large \text{Prerequisites}$
- [Download the latest version from the releases page](https://github.com/bepass-org/warp-plus/releases)
- Basic understanding of VPN and proxy configurations


## $\color{gold} \large \text{Acknowledgements}$
- $\color{silver} \normalsize \text{Cloudflare Warp}$
- $\color{silver} \normalsize \text{Cloudflare Warp}$
- $\color{silver} \normalsize \text{Psiphon}$
- $\color{silver} \normalsize \text{All contributors and supporters of this project.}$


$\color{silver} \large \text{That's it. If you have any questions, feel free to ask in the comments.}$

[![Telegram Badge](https://img.shields.io/badge/Telegram-Profile-informational?style=flat&logo=telegram&logoColor=white&color=1CA2F1)](https://t.me/F_NiREvil)

### $\color{silver} \large \text{Be curious. 🤍🪐}$ 

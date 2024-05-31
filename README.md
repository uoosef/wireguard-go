## $\color{gold} \Large \text{Warp-Plus}$
$\color{silver} \large \text{Many thanks to Usef Ghobadi and dear MarkPash}$

$\color{gray} \large \text{Warp-Plus is an open-source implementation of Cloudflare's Warp, enhanced with Psiphon integration}$
$\color{gray} \large \text{for circumventing censorship. This project aims to provide a robust and cross-platform VPN solution}$
$\color{gray} \large \text{that can use psiphon on top of warp and warp-in-warp for changing the user virtual nat location.}$


$\color{gold} \Large \text{Features}$
```
- Warp Integration: Leverages Cloudflare's Warp to provide a fast and secure VPN service.
- Psiphon Chaining: Integrates with Psiphon for censorship circumvention, allowing seamless access to the
  internet in restrictive environments.
- Warp in Warp Chaining: Chaining two instances of warp together to bypass location restrictions.
- SOCKS5 Proxy Support: Includes a SOCKS5 proxy for secure and private browsing.
```

![rainbow](https://github.com/NiREvil/vless/assets/126243832/1aca7f5d-6495-44b7-aced-072bae52f256)  

$\color{gold} \Large \text{Table of content}$

## جدول مطالب
- [تغییرات اخیر](#تغییرات-اخیر)
- [نصب و اجرا](#نصب-و-اجرا)
- [نکات مورد نیاز](#نکات-مورد-نیاز)
  - [توضیحات تکمیلی](#توضیحات-تکمیلی)
  - [رنج آیپی و پورت های وارپ](#رنج-آیپی-و-پورت-های-وارپ)
- [ویدئوهای آموزشی](#ویدئوهای-آموزشی)
- [راهنمای کلی دستورات](#راهنمای-کلی-دستورات)
  - [کد اختصاری کشورها](#کد-اختصاری-کشورها)
- [باکس دانلود ترماکس](#باکس-دانلود-ترماکس)


 ![rainbow](https://github.com/NiREvil/vless/assets/126243832/1aca7f5d-6495-44b7-aced-072bae52f256)  

 
## تغییرات اخیر
<details>
<summary> تغییرات در نسخه v1.2.3 </summary>

```
warp: fix duplicate device bug by @markpash in #143
Automatic warp reserved field by default, add win32, include wintun.dll and other cleanups by @markpash in #132
```

</details>  

![rainbow](https://github.com/NiREvil/vless/assets/126243832/1aca7f5d-6495-44b7-aced-072bae52f256)  


<img width="450" height="550" src="https://github.com/NiREvil/wireguard-go/assets/126243832/3aff5b35-e5cf-42ce-b07b-ddd95e180acc" alt="Termux Usef"/>

![rainbow](https://github.com/NiREvil/vless/assets/126243832/1aca7f5d-6495-44b7-aced-072bae52f256)  

## نصب و اجرا 
  
$\color{gold} \Large \text{Getting Started with Termux}$

```bash
bash <(curl -fsSL https://raw.githubusercontent.com/NiREvil/wireguard-go/master/termux.sh)
```

---
  
## نکات مورد نیاز  

اگه حس کردید کانکت نمیشه یا خطا میده دستور `rm -rf .cache/warp-plus` رو بزنید تا کش ها پاک بشن مجدد warp رو اجرا کنید.  

بعد از نصب برای اجرای مجدد فقط کافیه که `warp` یا `usef` یا `warp-plus`را وارد کنید. همش یکیه هیچ فرقی نداره.  

اگر با 1 نصب نشد و خطا گرفتید ابتدا یک بار 3 را بزنید تا `Uninstall` بشه سپس عدد 2 رو انتخاب کنید یعنی Arm (عموما برای گوشی های قدیمی و پردازنده ضعیف این گزینه جواب میده).  

برای نمایش راهنما ` warp -h` را وارد کنید.  

 دریافت بر روی آی‌پی و پورت `127.0.0.1:8086` پروتکل `SOCKS` یا کپی و وارد کردن این کانفیگ در نکوباکس و ...
`socks://Og==@127.0.0.1:8086#warp_(termux)`  

برای متوقف کردن وارپ در حال اجرا از کلیدهای ترکیبی `Ctrl + C` استفاده کنید و همچنین برای بستن نشست فعال و خارج شدن از ترماکس از کلیدهای `Ctrl + D` استفاده کنید.  

برای تکرار دوباره دستورهای وارد شده در مرحله قبل در ترماکس از فلش رو به بالا ↑ استفاده کنید.  

در روش تبدیل اکانت warp به warp plus (گزینه 6) مقدار ID نمایش داده شده را وارد می‌کنید. پس از اجرای warp دو اکانت برای شما ساخته شده که پس از انتخاب گزینه 6 خودش مقدار ID هر دو اکانت را پیدا می‌کند و شما باید هر بار یکی را انتخاب کنید و یا میتوانید با انتخاب manual مقدار ID دیگری را وارد کنید (مثلا برای خود برنامه ی 1.1.1.1 یا جای دیگر) با این کار هر 20 ثانیه 1 GB به اکانت شما اضافه شده و اکانت شما از حالت رایگان به پلاس تبدیل می‌شود.  

و یا برای تبدیل شدن وارپ به وارپ پباس مستقیما از [ربات های تلگرام](https://t.me/F_NiREvil/5271) کلید لایسنس دریافت کرده و با دستور زیر
`warp -k yourlicensekey`
آنرا ارسال کنید.  

برای تغییر  لوکیشن با استفاده از سایفون از طریق منو یا به صورت دستی: برای مثال برای USA از دستور زیر استفاده کنید:
`warp --cfon --country US`
یا برای مثال برای Estonia از دستور  زیر استفاده کنید:
`warp --cfon --country EE`  

برای اسکن ای پی سالم وارپ از دستور `warp --scan` استفاده کنید. 

برای زنجیر کردن (proxy chain) دو کانفیگ وارپ باهم و تغییر آی پی به غیر ایران از دستور `warp --gool` استفاده کنید.  


![rainbow](https://github.com/NiREvil/vless/assets/126243832/1aca7f5d-6495-44b7-aced-072bae52f256)  

## توضیحات تکمیلی

**و اما جونم برات بگه که** ببینید این دستورات هم وجود دارن هاا

از قبل میتونستیم انتخاب کنیم که وارپ عادی ما با آی‌پی ورژن 4 ران بشه یا با ورژن 6 مثلا من چون می‌دیدم واسه مخابرات خودش ipv6 انتخاب میکنه میومدم با این دستور اجرا می‌کردم `warp -4` کانکت میشد راحت.  

حالا از نسخه v1.2.1 این امکان رو فراهم کردن که برای بقیه قسمت ها هم بشه اینطوری عمل کرد یعنی می‌تونیم تو دستور سایفون یا گول هم بگیم با کدوم ورژن آی‌پی اجرا بشه. مثلا برای اجرای gool می‌تونیم اینجوری عمل کنیم:  

برای اجرا با آی‌پی ورژن 4 اینو وارد می‌کنیم `warp --gool -4`  

و برای اجرا با آی‌پی ورژن 6 این رو `warp --gool -6`

و برای اجرای سایفون با ورژن آی‌پی دلخواه : `warp --cfon --country US -4` یا `warp --cfon --country US -6`  

> اون US برای مثال بود می‌تونی هر لوکیشنی بنویسی یه جای اون

خلاصه اینکه یا انتخاب ورژن آی‌پی رو بسپر به هسته یا که خودت هندل کن.  

**ولی ولیییی** ممکنه انتخاب آی‌پی رو بسپاری دست خودش و بیاد حتی ipv4 که فرضا مورد نطرت هست رو انتخاب کنه ولی از یه رنج تخمی، فرضا برای مخابرات الان هیشکدوم از آی‌پی ها از رنج 162 کار نمیکنن، فقط 188 اوکیه

مثلا آی‌پی  

`188.114.98.224:908`  `188.114.98.224:939`  `188.114.97.82:2048`  `188.114.99.149:4500`  

کلا رنج 188 خوراک مخابرات و ایرانسله با تموم پورت های موجود.  
این هم با همراه اول خوب بودش `188.114.97.134:1180`


اگه تو نکوباکسی هیدیفایی ویتوری جایی کانفیگ وارپی هست ک برات عالیه آی‌پی و پورت همون رو حفظ کن بیا با دستور `warp -e IP:PORT` بده بهش تا ممبعد وقتی استارت کردی وارپ رو از همین ک تو پروفایلت ذخیره شده استفاده کنه.  

میشه از [اسکنر پیمان](https://t.me/F_NiREvil/6305) که دیروز دربارش حرف زدیمم استفاده کنی برای اسکن آی‌پی وارپ و استفاده ازش تو این پروژه

## رنج آیپی و پورت های وارپ
شاید ب کارت بیاد رنج آی‌پی و پورت های وارپ کلادفلر:

**Cloudflare warp IP range**
```
162.159.192.0/24
162.159.193.0/24 
162.159.195.0/24
188.114.96.0/24
188.114.97.0/24
188.114.98.0/24
188.114.99.0/24
```

**Port**
```
500-854-859-864-878-880-890-891-894-903-908-928-934-939-942-943-945-946,955-968-987-988-1002-1010-1014-1018-1070-1074-1180-1387-1701-1843-2371-2408-2506-3138-3476-3581-3854-4177-4198-4233-4500-5279-5956-7103-7152-7156-7281-7559-8319-8742-8854-8886
```

![rainbow](https://github.com/NiREvil/vless/assets/126243832/1aca7f5d-6495-44b7-aced-072bae52f256)


## ویدئوهای آموزشی  

**$\color{silver} \large \text{Part 1:}$**
> **with Single Command install**
```bash
bash <(curl -fsSL https://raw.githubusercontent.com/NiREvil/wireguard-go/master/termux.sh)
```
https://github.com/NiREvil/wireguard-go/assets/126243832/458edbef-1053-4605-9afd-fdf88ef2398c


![rainbow](https://github.com/NiREvil/vless/assets/126243832/1aca7f5d-6495-44b7-aced-072bae52f256)

**$\color{silver} \large \text{Part 2:}$**
> **with: `warp --gool`**

https://github.com/NiREvil/wireguard-go/assets/126243832/07aa5d18-c4b1-43fb-9814-b23ac730afa9



![rainbow](https://github.com/NiREvil/vless/assets/126243832/1aca7f5d-6495-44b7-aced-072bae52f256)

**$\color{silver} \large \text{Part 3:}$**
> **with: `warp --cfon --country EE`**

https://github.com/NiREvil/wireguard-go/assets/126243832/5054de35-e798-4603-ba7b-5c0accc0da07

![rainbow](https://github.com/NiREvil/vless/assets/126243832/1aca7f5d-6495-44b7-aced-072bae52f256)

**$\color{silver} \large \text{Part 4:}$**
> **Using `Warp-plus in Windows`**
> **Note:
>   Just be aware that because the video is from a few months ago; when entering commands,
>  we used a single -, but now in the recent update we use two -- in Warp commands**

[![YouTube Video](https://img.youtube.com/vi/Me2PJTOgnaQ/0.jpg)](https://www.youtube.com/watch?v=Me2PJTOgnaQ)

![rainbow](https://github.com/NiREvil/vless/assets/126243832/1aca7f5d-6495-44b7-aced-072bae52f256)


## راهنمای کلی دستورات   
$\color{gold} \Large \text{Usage}$

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
      --country STRING     psiphon country code (valid values - (default: AT)
      --scan               enable warp scanning
      --rtt DURATION       scanner rtt limit (default: 1s)
      --cache-dir STRING   directory to store generated profiles
      --tun-experimental   enable tun interface (experimental)
      --fwmark UINT        set linux firewall mark for tun mode (default: 4981)
      --wgconf STRING      path to a normal wireguard config
  -c, --config STRING      path to config file
      --version            displays version number
```
![rainbow](https://github.com/NiREvil/vless/assets/126243832/1aca7f5d-6495-44b7-aced-072bae52f256)

## کد اختصاری کشورها  
$\color{gold} \Large \text{Country Codes for Psiphon}$

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

![rainbow](https://github.com/NiREvil/vless/assets/126243832/1aca7f5d-6495-44b7-aced-072bae52f256)

## باکس دانلود ترماکس  
$\color{gold} \Large \text{Download BOX}$

| Source | Download Link's |
|:--------:| -------------:|
| F-Droid | [Get it Now](https://f-droid.org/en/packages/com.termux) |
| GitHub :octocat:| [Get it Now](https://github.com/termux/termux-app/releases) |
| Google Play | [Get it Now](https://play.google.com/store/apps/details?id=com.termux) |
| How to fix the installation error of Termux packages on Android 5 & 6 | [Fix Errors](https://t.me/F_NiREvil/5040) |

> I recommend installing GitHub for sure, as the Google Play version has many bugs and the F-Droid version is not compatible with most devices. The best choice is the GitHub repository of the app.

![rainbow](https://github.com/NiREvil/vless/assets/126243832/1aca7f5d-6495-44b7-aced-072bae52f256)

$\color{gold} \large \text{Prerequisites}$
- Download the latest version from the [releases](https://github.com/bepass-org/warp-plus/releases)
- Basic understanding of VPN and proxy configurations


## $\color{gold} \large \text{Acknowledgements}$
- $\color{silver} \normalsize \text{Cloudflare Warp}$
- $\color{silver} \normalsize \text{Cloudflare Warp}$
- $\color{silver} \normalsize \text{Psiphon}$
- $\color{silver} \normalsize \text{All contributors and supporters of this project.}$


$\color{silver} \large \text{That's it. If you have any questions, feel free to ask in the comments.}$

---

[![Telegram Badge](https://img.shields.io/badge/Telegram-Profile-informational?style=flat&logo=telegram&logoColor=white&color=1CA2F1)](https://t.me/F_NiREvil)

$\color{silver} \large \text{Be curious 🤍🪐}$ 

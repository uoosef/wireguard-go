#!/bin/bash

# Colors
rst='\033[0m'
bold='\033[1m'
red='\033[0;31m'
redb='\033[1;31m'
green='\033[0;32m'
greenb='\033[1;32m'
blue='\033[1;34m'
magneta='\033[0;35m'
cyan='\033[0;36m'
cyanb='\033[1;36m'
white='\033[1;37m'
gray='\033[1;90m'

# Check Dependencies build
check_dependencies_build() {
    local dependencies=("curl" "wget" "git" "golang")

    for dep in "${dependencies[@]}"; do
        if ! pkg info "${dep}" &> /dev/null; then
            echo -e "${redb}${dep} is not installed. Trying to install...${rst}"
            pkg install "${dep}" -y
        fi
    done
}

# Check Dependencies
check_dependencies() {
    local dependencies=("curl" "openssl-tool" "wget" "unzip")

    for dep in "${dependencies[@]}"; do
        if ! pkg info "${dep}" &> /dev/null; then
            echo -e "${redb}${dep} is not installed. Trying to install...${rst}"
            pkg install "${dep}" -y
        fi
    done
}

# Install
install() {
    if command -v warp &> /dev/null || command -v usef &> /dev/null; then
        echo ""
        echo -e "${greenb}WARP is already installed.${rst}"
        return
    fi

    echo ""
    echo -e "${red}Currently checking the status of package updates for Termux.${rst}"
    echo -e "${redb}Fetching WARP PLUS from Ransomware Evil${rst}"
    echo -e "${greenb}Installing WARP PLUS ...${rst}"
    pkg update -y && pkg upgrade -y
    check_dependencies

    if wget https://github.com/bepass-org/warp-plus/releases/latest/download/warp-plus_android-arm64.zip &&
        unzip warp-plus_android-arm64.zip &&
        mv warp-plus warp &&
        chmod +x warp &&
        cp warp "$PREFIX/bin/usef" &&
        cp warp "$PREFIX/bin/warp-plus" &&
        cp warp "$PREFIX/bin/warp"; then
        rm "README.md" "LICENSE" "warp-plus_android-arm64.zip"
        echo ""
        echo -e "${greenb}WARP installed successfully.${rst}"
        socks
    else
        echo -e "${redb}Error installing WARP PLUS ...${rst}"
    fi
}
        

# install arm
install_arm() {
    if command -v warp &> /dev/null || command -v usef &> /dev/null; then
        echo ""
        echo -e "${greenb}WARP is already installed.${rst}"
        return
    fi

    echo ""
    echo -e "${redb}OK ...${rst}"
    echo -e "${red}Currently checking the status of package updates for Termux.${rst}"
    echo -e "${redb}Fetching WARP PLUS from Ransomware Evil.${rst}"
    echo -e "${greenb}Installing WARP PLUS ...${rst}"
    pkg update -y && pkg upgrade -y
    check_dependencies

    # Determine architecture
    case "$(pkg info --print-architecture)" in
        i386) ARCH="386" ;;
        amd64) ARCH="amd64" ;;
        armhf) ARCH="arm5" ;;
        arm) ARCH="arm7" ;;
        aarch64) ARCH="arm64" ;;
        *) echo -e "${redb}Unsupported architecture.${rst}"; return ;;
    esac

    WARP_URL="https://github.com/bepass-org/warp-plus/releases/latest/download/warp-plus_linux-$ARCH.zip"

    if wget "$WARP_URL" &&
        unzip "warp-plus_linux-$ARCH.zip" &&
        mv warp-plus warp &&
        chmod +x warp &&
        cp warp "$PREFIX/bin/usef" &&
        cp warp "$PREFIX/bin/warp-plus" &&
        cp warp "$PREFIX/bin/warp"; then
        rm "README.md" "LICENSE" "warp-plus_linux-$ARCH.zip"
        echo ""
        echo -e "${greenb}WARP PLUS installed successfully.${rst}"
        socks
    else
        echo -e "${redb}Error installing WARP PLUS.${rst}"
    fi
}

# Get socks config
socks() {
   echo ""
   echo -e "${cyan}Copy this Config to ${redb}V2ray${rst} Or ${redb}Nekobox${rst} Or ${redb}Hiddify${rst}${cyan} and Exclude Termux in pre-apps proxy or apps vpn mode OR split tunneling or ...${rst}"
   echo ""
   echo -e "${greenb}socks://Og==@127.0.0.1:8086#WarpTermux${rst}"
   echo "OR"
   echo -e "${cyan}Manually create a SOCKS configuration with IP:${rst}${greenb}127.0.0.1${rst} and port: ${greenb}8086${rst}"
   echo ""
   echo -e "${cyan}For run again, type:${rst}${greenb} warp${rst} OR ${greenb}usef${rst} OR ${greenb}./warp${rst} OR ${greenb}warp-plus${rst}"
   echo ""
   echo ""
   echo -e "${cyanb}═══════════════════════════════════════════${rst}"
   echo -e "${magenta}If you get a ${rst}${redb}Bad address error${rst}"
   echo -e "${magenta}Run again script and hit the number ${rst}${redb}2${rst}"
   echo -e "${cyanb}═══════════════════════════════════════════${rst}"
   echo ""
}

# Gool (warp in warp)
gool() {
    if ! command -v warp &> /dev/null || ! command -v usef &> /dev/null; then
        echo ""
        echo -e "${redb}Please install WARP PLUS first.${rst}"
        return
    fi

    echo -e "${redb}This option changes your current location to the${rst}"
    echo -e "${redb}neareset and best location.${rst}"
    echo ""

    while true; do
        echo ""
        echo -e "${cyan} GOOL will run with ?${rst}"
        echo -e " [1] ${green}IPV4${rst}"
        echo -e " [2] ${green}IPV6${rst}"
        echo -e " [3] ${redb}Back to Menu${rst}"
        echo ""
        echo -en "${cyanB} Please Choose:${rst}"
        read -r option

        case $option in
            1)
                echo -e "${green}OK ...${rst}"
                echo -e "${greenb}Running WARP with IPv4...${rst}"
                warp --gool -4
                ;;
            2)
                echo -e "${green}OK ...${rst}"
                echo -e "${greenb}Running WARP with IPv6...${rst}"
                warp --gool -6
                ;;
            3)
                echo -e "${redb}Back to Menu${rst}"
                menu
                ;;
            *)
                echo ""
                echo -e "${redb}INVALID OPTION${rst}"
                ;;
        esac
    done
}

# Psiphon
psiphon_location() {
    if ! command -v warp &> /dev/null || ! command -v usef &> /dev/null; then
        echo ""
        echo -e "${redb}Please install WARP PLUS first.${rst}"
        return
    fi
    echo ""
    echo ""
    echo -e "${cyanb}══════════════════════════════════════════════════════════════${rst}"
    echo -e "${magenta} Choose a location from the list below by entering its number:${rst}"
    echo -e "${cyanb}══════════════════════════════════════════════════════════════${rst}"
    echo -e "${redb} 1)${cyan} Austria (AT)${rst}"
    echo -e "${redb} 2)${cyan} Belgium (BE)${rst}"
    echo -e "${redb} 3)${cyan} Bulgaria (BG)${rst}"
    echo -e "${redb} 4)${cyan} Brazil (BR)${rst}"
    echo -e "${redb} 5)${cyan} Canada (CA)${rst}"
    echo -e "${redb} 6)${cyan} Switzerland (CH)${rst}"
    echo -e "${redb} 7)${cyan} Czech Republic (CZ)${rst}"
    echo -e "${redb} 8)${cyan} Germany (DE)${rst}"
    echo -e "${redb} 9)${cyan} Denmark (DK)${rst}"
    echo -e "${redb}10)${cyan} Estonia (EE)${rst}"
    echo -e "${redb}11)${cyan} Spain (ES)${rst}"
    echo -e "${redb}12)${cyan} Finland (FI)${rst}"
    echo -e "${redb}13)${cyan} France (FR)${rst}"
    echo -e "${redb}14)${cyan} United Kingdom (GB)${rst}"
    echo -e "${redb}15)${cyan} Hungary (HU)${rst}"
    echo -e "${redb}16)${cyan} Ireland (IE)${rst}"
    echo -e "${redb}17)${cyan} India (IN)${rst}"
    echo -e "${redb}18)${cyan} Italy (IT)${rst}"
    echo -e "${redb}19)${cyan} Japan (JP)${rst}"
    echo -e "${redb}20)${cyan} Latvia (LV)${rst}"
    echo -e "${redb}21)${cyan} Netherlands (NL)${rst}"
    echo -e "${redb}22)${cyan} Norway (NO)${rst}"
    echo -e "${redb}23)${cyan} Poland (PL)${rst}"
    echo -e "${redb}24)${cyan} Romania (RO)${rst}"
    echo -e "${redb}25)${cyan} Serbia (RS)${rst}"
    echo -e "${redb}26)${cyan} Sweden (SE)${rst}"
    echo -e "${redb}27)${cyan} Singapore (SG)${rst}"
    echo -e "${redb}28)${cyan} Slovakia (SK)${rst}"
    echo -e "${redb}29)${cyan} Ukraine (UA)${rst}"
    echo -e "${redb}30)${cyan} United States (US)${rst}"
    echo ""

    echo ""
    echo -en "${green} Enter the ${cyanb}NUMBER${rst}${green} of the location:${rst}${green}    Default is:${rst}${redb} 1${rst}"
    read -r choice
    choice=${choice:-1}

    case "$choice" in
        1) location="AT" ;;
        2) location="BE" ;;
        3) location="BG" ;;
        4) location="BR" ;;
        5) location="CA" ;;
        6) location="CH" ;;
        7) location="CZ" ;;
        8) location="DE" ;;
        9) location="DK" ;;
        10) location="EE" ;;
        11) location="ES" ;;
        12) location="FI" ;;
        13) location="FR" ;;
        14) location="GB" ;;
        15) location="HU" ;;
        16) location="IE" ;;
        17) location="IN" ;;
        18) location="IT" ;;
        19) location="JP" ;;
        20) location="LV" ;;
        21) location="NL" ;;
        22) location="NO" ;;
        23) location="PL" ;;
        24) location="RO" ;;
        25) location="RS" ;;
        26) location="SE" ;;
        27) location="SG" ;;
        28) location="SK" ;;
        29) location="UA" ;;
        30) location="US" ;;
        *) echo " INVALID CHOICE. Please select a valid location number." ;;
    esac

    echo -e "${cyan} Selected location: $location ${rst}"

    while true; do
        echo ""
        echo -e "${cyan} WARP PLUS & Psiphon will run with ?${rst}"
        echo -e " [1] ${cyan}IPV4${rst}"
        echo -e " [2] ${cyan}IPV6${rst}"
        echo -e " [3] ${redb}Back to Menu${rst}"

        echo -en "${greenb} Please Choose: ${rst}"
        read -r option

        case $option in
            1)
                echo -e "${green}OK ... ${rst}"
                echo -e "${greenb}Running Psiphon with IPv4...${rst}"
                warp --cfon --country $location -4
                ;;
            2)
                echo -e "${green}OK ... ${rst}"
                echo -e "${greenb}Running Psiphon with IPv6...${rst}"
                warp --cfon --country $location -6
                ;;
            3)
                echo -e "${redb}Back to Menu${rst}"
                menu
                ;;
            *)
                echo -e "${redb}INVALID OPTION${rst}"
                ;;
        esac
    done
}

# Uninstall
uninstall() {
    warp="$PREFIX/bin/warp"
    directory="/data/data/com.termux/files/home/warp-plus"
    home="/data/data/com.termux/files/home"
    if [ -f "$warp" ]; then
        rm -rf "$directory" "$PREFIX/bin/usef" "wa.py" "$PREFIX/bin/warp" "$PREFIX/bin/warp-plus" "warp" "/data/data/com.termux/files/home/.cache/warp-plus" > /dev/null 2>&1
        echo ""
        echo -e "${redb}Uninstallation completed.${rst}"
    else
        echo ""
        echo -e "${redb}WARP PLUS not installed. Please Install First.${rst}"
    fi
}

# Warp to Warp plus
warp_plus() {
    if ! command -v python &> /dev/null; then
        echo "Installing Python..."
        pkg install python -y
    fi

    echo ""
    echo -e "${red}Please wait to download dependencies${rst}"
    echo -e "${redb}OK, Let's run script ... ${rst}"
    echo ""
    wget -O wa.py https://raw.githubusercontent.com/NiREvil/wireguard-go/refs/heads/master/warp/wa.py
    python wa.py
}

# Menu
menu() {
    clear
    printf "%b\n" "${blue}"
    printf "%b\n" "    ██╗    ██╗ █████╗ ██████╗ ██████╗"
    printf "%b\n" "    ██║    ██║██╔══██╗██╔══██╗██╔══██╗"
    printf "%b\n" "    ██║ █╗ ██║███████║██████╔╝██████╔╝"
    printf "%b\n" "    ██║███╗██║██╔══██║██╔══██╗██╔═══╝"
    printf "%b\n" "    ╚███╔███╔╝██║  ██║██║  ██║██║"
    printf "%b\n" "     ╚══╝╚══╝ ╚═╝  ╚═╝╚═╝  ╚═╝╚═╝${rst}${gray}v1.2.5${rst}"
    printf "%b\n" "     ${gray}Many thanks to Mark and Usef ${rst}${blue}◈${rst}"
    echo

    printf "%b\n" "${cyanb}╔═════════════════════════════════════════════╗${rst}"
    printf "%b\n" "${cyanb}║              ${green}We are all REvil${cyanb}               ║${rst}"
    printf "%b\n" "${cyanb}╠═════════════════════════════════════════════╣${rst}"
    printf "%b\n" "${green}║ ➊ ${white}Install Warp${rst}${green} arm64-v8a${rst}${cyanb}                    ║${rst}"
    printf "%b\n" "${green}║ ➋ ${white}Install Warp${rst}${green} armeabi-v7a${rst}${cyanb}                  ║${rst}"
    printf "%b\n" "${green}║ ➌ ${white}Uninstall Warp${rst}${cyanb}                            ║${rst}"
    printf "%b\n" "${green}║ ➍ ${white}Gool${rst}${green} Warp ON Warp${rst}${cyanb}                         ║${rst}"
    printf "%b\n" "${green}║ ➎ ${white}Psiphon${rst}${green} Multi-Location${rst}${cyanb}                    ║${rst}"
    printf "%b\n" "${green}║ ➏ ${white}Warp to Warp+${rst}${green} Free GB${rst}${cyanb}                     ║${rst}"
    printf "%b\n" "${cyanb}╠═════════════════════════════════════════════╣${rst}"
    printf "%b\n" "${red}║ ⓧ ${redb}Exit${rst}${cyanb}                                      ║${rst}"
    
    printf "%b\n" "${cyanb}╚═════════════════════════════════════════════╝${rst}"
    printf "%b\n" "${cyanb}┌───────────────────────────┐${rst}"
    echo -en "${cyanb}│ SELECT: ${rst}"
    tput sc
    
    # Move cursor to input position
    tput cup $(tput lines) 0
    echo -e "${cyanb}└───────────────────────────┘${rst}"
    tput rc
    read -r choice

    case "$choice" in
        1|01)
            install
            warp
            ;;
        2|02)
            install_arm
            warp
            ;;
        3|03)
            uninstall
            ;;
        4|04)
            gool
            ;;
        5|05)
            psiphon_location
            ;;
        6|06)
            warp_plus
            ;;
        0|x|X)
            echo
            echo -e "╔═══════════════════════════════╗"
            echo -e "║${redb}Thank goodness for small miracles 🙌${rst}║"
            echo -e "╚═══════════════════════════════╝"
            exit 0
            ;;
        *)
            echo -e "\n${red} Pay more attention psycho 😒 Or you might find yourself${rst}"
            echo -e "\n${red} accidentally signing up for a llama yoga class 🧘🏋️${rst}"
            echo -e "\n${redb} Press ENTER to continue ...${rst}"
            read
            menu
            ;;
    esac
}

menu

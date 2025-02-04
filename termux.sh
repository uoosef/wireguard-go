#!/bin/bash

# Colors 
RESET='\033[0m'
BOLD='\033[1m'
BLACK='\033[0;30m'
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[0;33m'
PURPLE='\033[0;35m'
CYAN='\033[0;36m'
WHITE='\033[0;37m'

# Check Dependencies build
check_dependencies_build() {
    local dependencies=("curl" "wget" "git" "golang")

    for dep in "${dependencies[@]}"; do
        if ! dpkg -s "${dep}" &> /dev/null; then
            echo -e "${BOLD}${RED}${dep} Is not installed. Try to Installing...${RESET}"
            pkg install "${dep}" -y
        fi
    done
}

# Check Dependencies
check_dependencies() {
    local dependencies=("curl" "openssl-tool" "wget" "unzip")

    for dep in "${dependencies[@]}"; do
        if ! dpkg -s "${dep}" &> /dev/null; then
            echo -e "${BOLD}${RED}${dep} Is not installed. Try to installing... ${RESET}"
            pkg install "${dep}" -y
        fi
    done
}

# Install
install() {
    if command -v warp &> /dev/null || command -v usef &> /dev/null; then
        echo ""
        echo -e "${BOLD}${GREEN}WARP is already installed.${RESET}"
        return
    fi

    echo ""
    echo -e "${RED}Currently checking the status of package updates for Termux.${RESET}"
    echo -e "${BOLD}${RED}Fetching WARP PLUS from Ransomware Evil ${RESET}"
    echo -e "${BOLD}${GREEN}Installing WARP PLUS ...${RESET}"
    pkg update -y && pkg upgrade -y
    pacman -Syu openssh = apt update; apt full-upgrade -y; apt install -y openssh
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
        echo -e "${BOLD}${GREEN}WARP installed successfully.${RESET}"
        socks
    else
        echo -e "${BOLD}${RED}Error installing WARP PLUS ...${RESET}"
    fi
}

# Install arm
install_arm() {
    if command -v warp &> /dev/null || command -v usef &> /dev/null; then
        echo ""
        echo -e "${BOLD}${GREEN}WARP is already installed.${RESET}"
        return
    fi

    echo ""
    echo -e "${BOLD}${RED}OK ... ${RESET}"
    echo -e "${RED}Currently checking the status of package updates for Termux.${RESET}"
    echo -e "${BOLD}${RED}Fetching WARP PLUS from Ransomware Evil.${RESET}"
    echo -e "${BOLD}${GREEN}Installing WARP PLUS ...${RESET}"
    pkg update -y && pkg upgrade -y
    pacman -Syu openssh = apt update; apt full-upgrade -y; apt install -y openssh
    check_dependencies

    # Determine architecture
    case "$(dpkg --print-architecture)" in
        i386) ARCH="386" ;;
        amd64) ARCH="amd64" ;;
        armhf) ARCH="arm5" ;;
        arm) ARCH="arm7" ;;
        aarch64) ARCH="arm64" ;;
        *) echo -e "${BOLD}${RED}Unsupported architecture.${RESET}"; return ;;
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
        echo -e "${BOLD}${GREEN}WARP PLUS installed successfully.${RESET}"
        socks
    else
        echo -e "${BOLD}${RED}Error installing WARP PLUS.${RESET}"
    fi
}

# Get socks config
socks() {
   echo ""
   echo -e "${BOLD}${WHITE}Copy this Config to ${BOLD}${RED}V2ray ${RESET} Or ${BOLD}${RED}Nekobox${RESET} Or ${BOLD}${RED}Hiddify ${RESET}${BOLD}${WHITE}and Exclude Termux in pre-apps proxy or apps vpn mode OR split tunneling or ...${RESET}"
   echo ""
   echo -e "${BOLD}${GREEN}socks://Og==@127.0.0.1:8086#WarpTermux${RESET}"
   echo "OR"
   echo -e "${BOLD}${WHITE}Manually create a SOCKS configuration with IP:${RESET}${BOLD}${GREEN}127.0.0.1 ${RESET}and port: ${GREEN}8086${RESET}"
   echo ""
   echo -e "${BOLD}${WHITE}For run again, type:${RESET}${BOLD}${GREEN} warp ${RESET}OR${BOLD}${GREEN} usef ${RESET}OR${BOLD}${GREEN} ./warp ${RESET}OR${BOLD}${GREEN} warp-plus ${RESET}"
   echo ""
   echo ""
   echo -e "${BOLD}${CYAN}═══════════════════════════════════════════${RESET}"
   echo -e "${WHITE}If you get a ${RESET}${BOLD}${RED}Bad address error${RESET}"
   echo -e "${BOLD}${WHITE}Run again script and hit the number ${RESET}${GREEN}2 - [Armeabi-v7a]${RESET}"
   echo -e "${BOLD}${CYAN}═══════════════════════════════════════════${RESET}"
   echo ""
}

# Gool (warp in warp)
gool() {
    if ! command -v warp &> /dev/null || ! command -v usef &> /dev/null; then
        echo ""
        echo -e "${BOLD}${RED}Please install WARP PLUS first.${RESET}"
        return
    fi

    echo -e "${BOLD}${RED}This option changes your current location to the nearest and best location${RESET}"
    echo -e "${BOLD}${RED}nearest and best location.${RESET}"
    echo ""

    while true; do
        echo ""
        echo -e "${BOLD}${CYAN}GOOL will run with IP ??${RESET}"
        echo -e " [1] ${BOLD}${GREEN}With IPv4${RESET}"
        echo -e " [2] ${BOLD}${GREEN}With IPv6${RESET}"
        echo -e " [3] ${BOLD}${RED}Back to Menu${RESET}"
        echo ""
        echo -en "${BOLD}${CYAN}Please Choose:${RESET}"
        read -r option

        case $option in
            1)
                echo -e "${BOLD}${GREEN}OK ... ${RESET}"
                echo -e "${BOLD}${GREEN}Running WARP with IPv4...${RESET}"
                warp --gool -4
                ;;
            2)
                echo -e "${BOLD}${GREEN}OK ... ${RESET}"
                echo -e "${BOLD}${GREEN}Running WARP with IPv6...${RESET}"
                warp --gool -6
                ;;
            3)
                echo -e "${BOLD}${RED}Back to Menu${RESET}"
                menu
                ;;
            *)
                echo ""
                echo -e "${BOLD}${RED}INVALID OPTION${RESET}"

                ;;
        esac
    done
}

# Psiphon
psiphon_location() {
    if ! command -v warp &> /dev/null || ! command -v usef &> /dev/null; then
        echo ""
        echo -e "${BOLD}${RED}Please install WARP PLUS first.${RESET}"
        return
    fi
    echo ""
    echo ""
    echo -e "${BOLD}${CYAN}══════════════════════════════════════════════════════════════${RESET}"
    echo -e "${BOLD}${CYAN} Choose a location from the list below by entering its number:${RESET}"
    echo -e "${BOLD}${CYAN}══════════════════════════════════════════════════════════════${RESET}"
    echo -e "${RED} 1)${CYAN} Austria (AT)${RESET}"
    echo -e "${RED} 2)${CYAN} Belgium (BE)${RESET}"
    echo -e "${RED} 3)${CYAN} Bulgaria (BG)${RESET}"
    echo -e "${RED} 4)${CYAN} Brazil (BR)${RESET}"
    echo -e "${RED} 5)${CYAN} Canada (CA)${RESET}"
    echo -e "${RED} 6)${CYAN} Switzerland (CH)${RESET}"
    echo -e "${RED} 7)${CYAN} Czech Republic (CZ)${RESET}"
    echo -e "${RED} 8)${CYAN} Germany (DE)${RESET}"
    echo -e "${RED} 9)${CYAN} Denmark (DK)${RESET}"
    echo -e "${RED}10)${CYAN} Estonia (EE)${RESET}"
    echo -e "${RED}11)${CYAN} Spain (ES)${RESET}"
    echo -e "${RED}12)${CYAN} Finland (FI)${RESET}"
    echo -e "${RED}13)${CYAN} France (FR)${RESET}"
    echo -e "${RED}14)${CYAN} United Kingdom (GB)${RESET}"
    echo -e "${RED}15)${CYAN} Hungary (HU)${RESET}"
    echo -e "${RED}16)${CYAN} Ireland (IE)${RESET}"
    echo -e "${RED}17)${CYAN} India (IN)${RESET}"
    echo -e "${RED}18)${CYAN} Italy (IT)${RESET}"
    echo -e "${RED}19)${CYAN} Japan (JP)${RESET}"
    echo -e "${RED}20)${CYAN} Latvia (LV)${RESET}"
    echo -e "${RED}21)${CYAN} Netherlands (NL)${RESET}"
    echo -e "${RED}22)${CYAN} Norway (NO)${RESET}"
    echo -e "${RED}23)${CYAN} Poland (PL)${RESET}"
    echo -e "${RED}24)${CYAN} Romania (RO)${RESET}"
    echo -e "${RED}25)${CYAN} Serbia (RS)${RESET}"
    echo -e "${RED}26)${CYAN} Sweden (SE)${RESET}"
    echo -e "${RED}27)${CYAN} Singapore (SG)${RESET}"
    echo -e "${RED}28)${CYAN} Slovakia (SK)${RESET}"
    echo -e "${RED}29)${CYAN} Ukraine (UA)${RESET}"
    echo -e "${RED}30)${CYAN} United States (US)${RESET}"
    echo ""

    echo ""
    echo -en "${GREEN} Enter the ${BOLD}${CYAN}NUMBER${RESET}${GREEN} of the location:${RESET}${GREEN}    Default is:${REEET}${BOLD}${RED} [1]${RESET}"
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

    echo -e "${BOLD}${CYAN} Selected location: $location ${RESET}"

    while true; do
        echo ""
        echo -e "${BOLD}${CYAN}WARP PLUS & Psiphon will run with IP ??${RESET}"
        echo -e " [1] ${BOLD}${CYAN}With IPv4${RESET}"
        echo -e " [2] ${BOLD}${CYAN}With IPv6${RESET}"
        echo -e " [3] ${BOLD}${RED}Back to Menu${RESET}"

        echo -en "${BOLD}${GREEN}Please Choose: ${RESET}"
        read -r option

        case $option in
            1)
                echo -e "${BOLD}${GREEN}OK ... ${RESET}"
                echo -e "${BOLD}${GREEN}Running Psiphon with IPv4...${RESET}"
                warp --cfon --country $location -4
                ;;
            2)
                echo -e "${BOLD}${GREEN}OK ... ${RESET}"
                echo -e "${BOLD}${GREEN}Running Psiphon with IPv6...${RESET}"
                warp --cfon --country $location -6
                ;;
            3)
                echo -e "${BOLD}${RED}Back to Menu${RESET}"
                menu
                ;;
            *)
                echo -e "${BOLD}${RED}INVALID OPTION${RESET}"
                ;;
        esac
    done
}

#Uninstall
uninstall() {
    warp="$PREFIX/bin/warp"
    directory="/data/data/com.termux/files/home/warp-plus"
    home="/data/data/com.termux/files/home"
    if [ -f "$warp" ]; then
        rm -rf "$directory" "$PREFIX/bin/usef" "wa.py" "$PREFIX/bin/warp" "$PREFIX/bin/warp-plus" "warp" "/data/data/com.termux/files/home/.cache/warp-plus" > /dev/null 2>&1
        echo ""
        echo -e "${BOLD}${RED}Uninstallation completed.${RESET}"
    else
        echo ""
        echo -e "${BOLD}${RED}WARP PLUS not installed. Please Install First.${RESET}${GREEN}|${RESET}"
    fi
}

# Warp to Warp plus
warp_plus() {
    if ! command -v python &> /dev/null; then
        echo "Installing Python..."
        pkg install python -y
    fi

    echo ""
    echo -e "${BOLD}${CYAN}Please wait to download dependencies${RESET}"
    echo -e "${BOLD}${RED}OK, Let's run script ... ${RESET}"
    echo ""
    wget -O wa.py https://raw.githubusercontent.com/NiREvil/wireguard-go/refs/heads/master/warp/wa.py
    python wa.py
}

# Menu
menu() {
    clear
    echo ""
    echo ""
    echo -e "${BOLD}${CYAN}╔═════════════════════════════════════════════╗${RESET}"
    echo -e "${BOLD}${CYAN}║              ${BOLD}${GREEN}WARP PLUS v1.2.5${BOLD}${CYAN}               ║${RESET}"
    echo -e "${BOLD}${CYAN}║          ${BOLD}${GREEN}Thanks to Mark and Usef${BOLD}${CYAN}            ║${RESET}"
    echo -e "${BOLD}${CYAN}╠═════════════════════════════════════════════╣${BOLD}${RESET}"
    echo -e "${BOLD}${GREEN}║ 1) ${WHITE}Install Warp${RESET}${BOLD}${GREEN} arm64-v8a${RESET}${BOLD}${CYAN}                   ║${RESET}"
    echo -e "${BOLD}${GREEN}║ 2) ${WHITE}Install Warp${RESET}${BOLD}${GREEN} armeabi-v7a${RESET}${BOLD}${CYAN}                 ║${RESET}"
    echo -e "${BOLD}${GREEN}║ 3) ${WHITE}Uninstall Warp${RESET}${BOLD}${CYAN}                           ║${RESET}"
    echo -e "${BOLD}${GREEN}║ 4) ${WHITE}Gool${RESET}${BOLD}${GREEN} Warp ON Warp${RESET}${BOLD}${CYAN}                        ║${RESET}"
    echo -e "${BOLD}${GREEN}║ 5) ${WHITE}Psiphon${RESET}${BOLD}${GREEN} Multi-Location${RESET}${BOLD}${CYAN}                   ║${RESET}"
    echo -e "${BOLD}${GREEN}║ 6) ${WHITE}Warp to Warp+${RESET}${BOLD}${GREEN} Free GB${RESET}${BOLD}${CYAN}                    ║${RESET}"
    echo -e "${BOLD}${RED}║ 9) ${WHITE}Revised by Diana${RESET}${BOLD}${CYAN}                         ║${RESET}"
    echo -e "${BOLD}${CYAN}╠═════════════════════════════════════════════╣${RESET}"
    echo -e "${BOLD}${RED}║ 0) Exit${BOLD}${CYAN}                                     ║${RESET}"
    echo -e "${BOLD}${CYAN}╚═════════════════════════════════════════════╝${RESET}"

    echo ""
    echo -en "${BOLD}${WHITE}Select an option:${RESET}"
    read -r choice

    case "$choice" in
        1) install ;;
        2) install_arm ;;
        3) uninstall ;;
        4) gool ;;
        5) psiphon_location ;;
        6) warp_plus ;;
        0) 
    read -r choice

            echo -e "${BOLD}${YELLOW} Thank goodness for small miracles 🙌${RESET}"
            exit 0
            ;;
        *)
            echo ""
            echo -e "${BOLD}${RED} INVALID OPTION.${RESET}"
            echo -e "${BOLD}${WHITE} Pay more attention psycho 😒, Or you might find yourself${RESET}"
            echo -e "${BOLD}${WHITE} accidentally signing up for a llama yoga class 🧘🏋️${RESET}"
            echo -e "${BOLD}${RED} Press ENTER to continue ... ${RESET}"
            read
            menu
            ;;
    esac
}

menu

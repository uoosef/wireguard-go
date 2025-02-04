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
            echo -e "${BOLD}${RED}${dep} is not installed. ${RESET}Installing..."
            pkg install "${dep}" -y
        fi
    done
}

# Check Dependencies
check_dependencies() {
    local dependencies=("curl" "openssl-tool" "wget" "unzip")

    for dep in "${dependencies[@]}"; do
        if ! dpkg -s "${dep}" &> /dev/null; then
            echo -e "${BOLD}${YELLOW}${dep} is not installed. ${RESET}Installing..."
            pkg install "${dep}" -y
        fi
    done
}

# Install
install() {
    if command -v warp &> /dev/null || command -v usef &> /dev/null; then
        echo -e ""
        echo -e "${BOLD}${GREEN}Warp is already installed.${RESET}"
        return
    fi

    echo -e ""
    echo -e "${RED}Currently checking the status of package updates for Termux.${RESET}"
    echo -e "${BOLD}${RED}Fetching warp from Ransomware REvil${RESET}"
    echo -e "${BOLD}${GREEN}Installing Warp...${RESET}"
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
        echo -e ""
        echo -e "${BOLD}${GREEN}Warp installed successfully.${RESET}"
        socks
    else
        echo -e "${BOLD}${RED}Error installing Warp.${RESET}"
    fi
}

# Install arm
install_arm() {
    if command -v warp &> /dev/null || command -v usef &> /dev/null; then
        echo -e ""
        echo -e "${BOLD}${GREEN}Warp is already installed.${RESET}"
        return
    fi

    echo -e ""
    echo -e "${RED}Currently checking the status of package updates for Termux.${RESET}"
    echo -e "${BOLD}${RED}Fetching warp from Ransomware REvil${RESET}"
    echo -e "${BOLD}${GREEN}Installing Warp...${RESET}"
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
        echo -e ""
        echo -e "${BOLD}${GREEN}Warp installed successfully.${RESET}"
        socks
    else
        echo -e "${BOLD}${RED}Error installing Warp.${RESET}"
    fi
}

# Get socks config
socks() {
   echo ""
   echo -e "${BOLD}Copy this Config to ${RED}V2ray ${RESET} Or ${RED}Nekobox${RESET} Or ${RED}Hiddify ${RESET}and Exclude Termux in pre-apps proxy or apps vpn mode OR split tunneling or ..."
   echo ""
   echo -e "${BOLD}${GREEN}socks://Og==@127.0.0.1:8086#WarpTermux${RESET}"
   echo "or"
   echo -e "${BOLD}${WHITE}Manually create a SOCKS configuration with IP:${RESET} ${BOLD}${GREEN}127.0.0.1 ${RESET}and port: ${GREEN}8086${RESET}"
   echo -e ""
   echo -e "${BOLD}${WHITE}To run again, type:${RESET}${GREEN} warp ${RESET}OR${GREEN} usef ${RESET}OR${GREEN} ./warp ${RESET}OR${GREEN} warp-plus ${RESET}"
   echo -e ""
   echo -e ""
   echo -e "${CYAN}═══════════════════════════════════════════${RESET}"
   echo -e "${BOLD}${WHITE}If you get a ${RESET}${RED}Bad address error${RESET}"
   echo -e "${BOLD}${WHITE}Run again script and hit the number 2 ${RESET}${GREEN}[armeabi-v7a]${RESET}"
   echo -e "${CYAN}═══════════════════════════════════════════${RESET}"
   echo ""
}

# Gool (warp in warp)
gool() {
    if ! command -v warp &> /dev/null || ! command -v usef &> /dev/null; then
        echo -e ""
        echo -e "${BOLD}${RED}Please install Warp first.${RESET}"
        return
    fi

    echo -e "${BOLD}${WHITE}This option changes your current location to the nearest and best location.${RESET}"

    while true; do
        echo -e ""
        echo -e "${BOLD}${CYAN}Choose an option:${RESET}"
        echo -e ""
        echo -e " [1] ${BOLD}${GREEN}IPV4${RESET}"
        echo -e " [2] ${BOLD}${GREEN}IPV6${RESET}"
        echo -e " [3] ${BOLD}${RED}Back to Menu${RESET}"
        echo -e ""

        echo -en "${BOLD}${CYAN}Please Choose:${RESET}"
        read -r option

        case $option in
            1)
                echo -e "${BOLD}${GREEN}Running Warp with IPV4...${RESET}"
                warp --gool -4
                ;;
            2)
                echo -e "${BOLD}${GREEN}Running Warp with IPV6...${RESET}"
                warp --gool -6
                ;;
            3)
                echo -e "${BOLD}${RED}Back to Menu.${RESET}"
                menu
                ;;
            *)

                echo -e "${BOLD}${RED}Invalid option.${RESET}"

                ;;
        esac
    done
}

# Psiphon
psiphon_location() {
    if ! command -v warp &> /dev/null || ! command -v usef &> /dev/null; then
        echo -e "${ref}Please install Warp first.${RESET}"
        return
    fi
    echo -e "${CYAN}════════════════════════════════════════════════════════════${RESET}"
    echo -e "${BOLD}${GREEN}Please choose a location from the list below by entering its number:${RESET}"
    echo ""
    echo -e "${RED} 1)${YELLOW} Austria (AT)${RESET}"
    echo -e "${RED} 2)${YELLOW} Belgium (BE)${RESET}"
    echo -e "${RED} 3)${YELLOW} Bulgaria (BG)${RESET}"
    echo -e "${RED} 4)${YELLOW} Brazil (BR)${RESET}"
    echo -e "${RED} 5)${YELLOW} Canada (CA)${RESET}"
    echo -e "${RED} 6)${YELLOW} Switzerland (CH)${RESET}"
    echo -e "${RED} 7)${YELLOW} Czech Republic (CZ)${RESET}"
    echo -e "${RED} 8)${YELLOW} Germany (DE)${RESET}"
    echo -e "${RED} 9)${YELLOW} Denmark (DK)${RESET}"
    echo -e "${RED}10)${YELLOW} Estonia (EE)${RESET}"
    echo -e "${RED}11)${YELLOW} Spain (ES)${RESET}"
    echo -e "${RED}12)${YELLOW} Finland (FI)${RESET}"
    echo -e "${RED}13)${YELLOW} France (FR)${RESET}"
    echo -e "${RED}14)${YELLOW} United Kingdom (GB)${RESET}"
    echo -e "${RED}15)${YELLOW} Hungary (HU)${RESET}"
    echo -e "${RED}16)${YELLOW} Ireland (IE)${RESET}"
    echo -e "${RED}17)${YELLOW} India (IN)${RESET}"
    echo -e "${RED}18)${YELLOW} Italy (IT)${RESET}"
    echo -e "${RED}19)${YELLOW} Japan (JP)${RESET}"
    echo -e "${RED}20)${YELLOW} Latvia (LV)${RESET}"
    echo -e "${RED}21)${YELLOW} Netherlands (NL)${RESET}"
    echo -e "${RED}22)${YELLOW} Norway (NO)${RESET}"
    echo -e "${RED}23)${YELLOW} Poland (PL)${RESET}"
    echo -e "${RED}24)${YELLOW} Romania (RO)${RESET}"
    echo -e "${RED}25)${YELLOW} Serbia (RS)${RESET}"
    echo -e "${RED}26)${YELLOW} Sweden (SE)${RESET}"
    echo -e "${RED}27)${YELLOW} Singapore (SG)${RESET}"
    echo -e "${RED}28)${YELLOW} Slovakia (SK)${RESET}"
    echo -e "${RED}29)${YELLOW} Ukraine (UA)${RESET}"
    echo -e "${RED}30)${YELLOW} United States (US)${RESET}"
    echo ""

    echo -en "${GREEN}Enter the ${YELLOW}number${BOLD}${GREEN} of the location [${YELLOW}default: 1${BOLD}${GREEN}]: ${RESET}"
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
        *) echo "Invalid choice. Please select a valid location number." ;;
    esac

    echo -e "${BOLD}${GREEN}Selected location: $location${RESET}"

    while true; do
        echo -e "${BOLD}${CYAN}Warp and psiphon will run with IP ??${RESET}"
        echo -e "  [1] ${BOLD}${GREEN}IPV4${RESET}"
        echo -e "  [2] ${BOLD}${GREEN}IPV6${RESET}"
        echo -e "  [3] ${BOLD}${RED}Back to Menu${RESET}"

        echo -en "${BOLD}${GREEN}Please Choose: ${RESET}"
        read -r option

        case $option in
            1)
                echo -e "${BOLD}${GREEN}Running Psiphon with IPv4...${RESET}"
                warp --cfon --country $location -4
                ;;
            2)
                echo -e "${BOLD}${GREEN}Running Psiphon with IPv6...${RESET}"
                warp --cfon --country $location -6
                ;;
            3)
                echo -e "${BOLD}${RED}Back to Menu.${RESET}"
                menu
                ;;
            *)
                echo -e "${BOLD}${RED}INVALID OPTION.${RESET}"
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
        echo -e ""
        echo -e "${BOLD}${RED}Uninstallation completed.${RESET}"
    else
        echo -e ""
        echo -e "${BOLD}${RED}Warp not installed. Please Install First.${RESET}${GREEN}|${RESET}"
    fi
}

# Warp to Warp plus
warp_plus() {
    if ! command -v python &> /dev/null; then
        echo "Installing Python..."
        pkg install python -y
    fi

    echo -e "${BOLD}${GREEN} pleaase wait to download dependencies${BOLD}${RED} Let's run warp+ script ...${RESET}"
    wget -O wa.py https://raw.githubusercontent.com/NiREvil/wireguard-go/refs/heads/master/warp/wa.py
    python wa.py
}

# Menu
menu() {
    clear
    echo -e "${BOLD}${CYAN}╔═════════════════════════════════════════════╗${RESET}"
    echo -e "${BOLD}${CYAN}║              ${BOLD}${GREEN}WARP PLUS v1.2.5${CYAN}               ║${RESET}"
    echo -e "${BOLD}${CYAN}║          ${BOLD}${GREEN}Thanks to Mark and Usef${CYAN}            ║${RESET}"
    echo -e "${BOLD}${CYAN}╠═════════════════════════════════════════════╣${RESET}"
    echo -e "${BOLD}${GREEN}║ 1) ${WHITE}Install Warp [ARM64-v8a]${RESET}${CYAN}                 ║${RESET}"
    echo -e "${BOLD}${GREEN}║ 2) ${WHITE}Install Warp [ARMv7]${RESET}${CYAN}                     ║${RESET}"
    echo -e "${BOLD}${GREEN}║ 3) ${WHITE}Uninstall Warp${RESET}${CYAN}                           ║${RESET}"
    echo -e "${BOLD}${GREEN}║ 4) ${WHITE}Gool [Warp over Warp]${RESET}${CYAN}                    ║${RESET}"
    echo -e "${BOLD}${GREEN}║ 5) ${WHITE}Psiphon [Multi-Location]${RESET}${CYAN}                 ║${RESET}"
    echo -e "${BOLD}${GREEN}║ 6) ${WHITE}Warp to Warp+ [Free GB]${RESET}${CYAN}                  ║${RESET}"
    echo -e "${BOLD}${RED}║ x) ${WHITE}Revised by Diana${RESET}${CYAN}                         ║${RESET}"
    echo -e "${BOLD}${CYAN}╠═════════════════════════════════════════════╣${RESET}"
    echo -e "${BOLD}${RED}║ 0) ${WHITE}Exit${RESET}${CYAN}                                     ║${RESET}"
    echo -e "${BOLD}${CYAN}╚═════════════════════════════════════════════╝${RESET}"

    echo -en "${BOLD}${WHITE}Select an option:${RESET}"
    read -r choice

    case "$choice" in
        1) install ;;
        2) install_arm ;;
        3) uninstall ;;
        4) gool ;;
        5) psiphon_location ;;
        6) warp_plus ;;
        x) gool ;;
        0) 
            echo -e "${BOLD}${GREEN}Poof, Thanks god, see you, Stay secure 🛡️${RESET}"
            exit 0
            ;;
        *)
            echo -e "${BOLD}${RED}Invalid option. Press Enter to continue...${RESET}"
            echo -e "${BOLD}${YELLOW}Pay more attention psycho 😒; ${RESET}"
            echo -e "${BOLD}${YELLOW}Or you might find yourself accidentally signing up for a llama yoga class!${RESET}"
            echo -e "${BOLD}${RED}Press Enter to continue ... ${RESET}"
            read
            menu
            ;;
    esac
}

menu

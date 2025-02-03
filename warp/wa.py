import urllib.request
import json
import datetime
import random
import string
import time
import os
import sys

os.system("title WARP-PLUS-CLOUDFLARE URSECRET")
os.system('cls' if os.name == 'nt' else 'clear')

file_paths = {
    "Primary": "/data/data/com.termux/files/home/.cache/warp-plus/primary/wgcf-identity.json",
    "Secondary": "/data/data/com.termux/files/home/.cache/warp-plus/secondary/wgcf-identity.json"
}

def display_id(file_path, name):
    if os.path.exists(file_path):
        with open(file_path, "r") as file:
            identity_data = json.load(file)
            return f"[{name} ID]\nID: {identity_data.get('id', 'Not set')}\n\n"
    else:
        return f"[{name} ID]\nID: Not set\n\n"

print (" ___Let's changing Warp to Warp plus ___")
print("")
print(display_id(file_paths['Primary'], "Primary"))
print(display_id(file_paths['Secondary'], "Secondary"))
print ("---------------------------------------")

option = input("[#] Choose an option:\n[1] Primary ID\n[2] Secondary ID\n[3] Manual (Enter ID)\n[4] Cancel\nChoice: ")

if option == "1":
    file_path = file_paths.get("Primary", "")
elif option == "2":
    file_path = file_paths.get("Secondary", "")
elif option == "3":
    file_path = None
elif option == "4":
    print("See you later ✋🏿")
    sys.exit(0)
else:
    print("INVALID OPTION. Please choose a valid option.")
    sys.exit(1)

if file_path and os.path.exists(file_path):
    with open(file_path, "r") as file:
        identity_data = json.load(file)
        referrer = identity_data.get("id", "")
else:
    referrer = input("[#] ENTER WARP+ ID : ")

def genString(stringLength):
    try:
        letters = string.ascii_letters + string.digits
        return ''.join(random.choice(letters) for i in range(stringLength))
    except Exception as error:
        print(error)

def digitString(stringLength):
    try:
        digit = string.digits
        return ''.join((random.choice(digit) for i in range(stringLength)))    
    except Exception as error:
        print(error)   

url = f'https://api.cloudflareclient.com/v0a{digitString(3)}/reg'

def run():
    try:
        install_id = genString(22)
        body = {"key": "{}=".format(genString(43)),
                "install_id": install_id,
                "fcm_token": "{}:APA91b{}".format(install_id, genString(134)),
                "referrer": referrer,
                "warp_enabled": False,
                "tos": datetime.datetime.now().isoformat()[:-3] + "+02:00",
                "type": "Android",
                "locale": "es_ES"}
        data = json.dumps(body).encode('utf8')
        headers = {'Content-Type': 'application/json; charset=UTF-8',
                    'Host': 'api.cloudflareclient.com',
                    'Connection': 'Keep-Alive',
                    'Accept-Encoding': 'gzip',
                    'User-Agent': 'okhttp/3.12.1'
                    }
        req         = urllib.request.Request(url, data, headers)
        response    = urllib.request.urlopen(req)
        status_code = response.getcode()  
        return status_code
    except Exception as error:
        print(error)   

g = 0
b = 0
while True:
    result = run()
    if result == 200:
        g = g + 1
        os.system('cls' if os.name == 'nt' else 'clear')
        print("")
        print(" WARP-PLUS-CLOUDFLARE ")
        print("")
        animation = ["[■□□□□□□□□□] 10%","[■■□□□□□□□□] 20%", "[■■■□□□□□□□] 30%", "[■■■■□□□□□□] 40%", "[■■■■■□□□□□] 50%", "[■■■■■■□□□□] 60%", "[■■■■■■■□□□] 70%", "[■■■■■■■■□□] 80%", "[■■■■■■■■■□] 90%", "[■■■■■■■■■■] 100%"] 
        for i in range(len(animation)):
            time.sleep(0.4)
            sys.stdout.write("\r[+] Preparing... " + animation[i % len(animation)])
            sys.stdout.flush()
        print(f"\n[-] WORK ON ID: {referrer}")    
        print(f"[:)] {g} GB HAS BEEN SUCCESSFULLY ENTERED IN YOUR ACCOUNT.")
        print(f"[#] Total: {g} Good {b} Bad")
        print("[*] AFTER 20 SECONDS A NEW REQUEST WILL BE SENT.")
        print(f" Please wait 20 seconds...")
        time.sleep(20)
    else:
        b = b + 1
        os.system('cls' if os.name == 'nt' else 'clear')
        print("")
        print("WARP-PLUS-CLOUDFLARE ")
        print("")
        print("[:(] Error when connecting to server.")
        print(f"[#] Total: {g} Good {b} Bad")

# Protoz

ProtoZ is a tool designed to find client-side prototype pollution vulnerabilities. It requires Google Chromium headless browser, Go-lang, and unfurl. The setup.sh script, which processes URLs to detect prototype pollution vulnerabilities using default and custom JavaScript code. We can customize the JS as per the requirement. The tool outputs whether a given URL is vulnerable based on the injected prototype pollution payloads.

This tool is inspired from the NahamSec 2021 talk by @Tomnomnom(https://www.youtube.com/channel/UCyBZ1F8ZCJVKSIJPrLINFyA)

# Dependencies
GO lang version 1.22.x or more 

# How to run this tool

```
git clone https://github.com/Manif3stVoid/ProtoZ.git
cd  ProtoZ
chmod +x setup.sh
./setup.sh -f urls.txt -m search

```
```
m  - mode types are search, hash, brute and gadget(Experimental)
f  - input a file which should contains the urls
j  - custom JS for payload validation
```

ProtoZ Modes : 

```
Search Mode  - This is the default mode and this mode validate the payloads prefix with question mark 
Hash Mode    - This mode validate the payloads prefix with hash
Brute Mode   - This mode brute force urls with the payloads through both search and hash mode. The mode uses the payloads.txt file fir the bruteforcing which contains the PP payloads
Gadget Mode  - This mode uses the payload which are specific to the JS gadgets from gadgets.txt file. This mode is experimental only. We need to manually validate the vulnerability once it is prints the url along with the payloads.

```

Usage :

```
./setup.sh -f urls.txt -j "console.log(document.title)" -m search
./setup.sh -f domains.txt -m hash
./setup.sh -f file.txt -m brute
./setup.sh -f urls.txt -m gadget

```
#### Example usage with custom JS
```
Fetching all links (href attributes) from anchor tags:
./setup.sh -f urls.txt -j '[...document.getElementsByTagName("a")].map(n => n.href).join(" ")' -m search

Getting all image sources (src attributes) from <img> tags:
./setup.sh -f urls.txt -j '[...document.getElementsByTagName("img")].map(img => img.src).join(" ")' -m search

Colleting all text content from <p> tags:
./setup.sh -f urls.txt -j '[...document.getElementsByTagName("p")].map(p => p.textContent).join(" ")' -m search

Getting all values from input fields
./setup.sh -f urls.txt -j '[...document.getElementsByTagName("input")].map(input => input.value).join(" ")' -m search

Getting all alt text from images
./setup.sh -f urls.txt -j '[...document.getElementsByTagName("img")].map(img => img.alt).join(" ")' -m search

```


### If you found Client Side Prototype Pollution? What's Next?
Check for the gadgets here https://github.com/BlackFan/client-side-prototype-pollution to use appropriate payload to exploit.

# References 
- https://www.youtube.com/watch?v=Gv1nK6Wj8qM 
- https://github.com/BlackFan/client-side-prototype-pollution
- https://github.com/msrkp/PPScan



# Protoz
Prototype Pollution Finder tool

This tool is inspired from the NahamSec 2021 talk by @Tomnomnom(https://www.youtube.com/channel/UCyBZ1F8ZCJVKSIJPrLINFyA)

https://www.youtube.com/watch?v=Gv1nK6Wj8qM

# Dependencies
GO lang version 1.22.x or more 

# How to run this tool

```
./setup.sh -f urls.txt -m search

m  - mode types are search, hash, brute and gadget(Experimental)
f  - input a file which should contains the urls
j  - custom JS for payload validation

```

ProtoZ Modes : 

```
Search Mode  - This is the default mode and this mode validate the payloads prefix with question mark 
Hash Mode    - This mode validate the payloads prefix with hash
Brute Mode   - This mode brute force urls with the payloads through both search and hash mode. The mode uses the payloads.txt file fir the bruteforcing which contains the PP payloads
Gadget Mode  - This mode uses the payload which are specific to the JS gadgets from gadgets.txt file. This mode is experimental only. We need to manually validate the vulnerability once it is prints the url along with the                   payloads.

```

Usage :

```
./setup.sh -f urls.txt -j "console.log(document.title)" -m search
./setup.sh -f domains.txt -m hash
./setup.sh -f file.txt -m brute
./setup.sh -f urls.txt -m gadget


```




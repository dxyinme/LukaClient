# LukaClient

the client for Luka-im <br>

[![Build Status](https://travis-ci.com/dxyinme/LukaClient.svg?branch=master)](https://travis-ci.com/dxyinme/LukaClient)

## Compile
```cmd
<Windows>
build.cmd
```

```shell
<linux/macOS>
./build.sh
```

## Usage
```cmd
<Windows>
# operation [p]ull , [s]end
# send message
client_cli.exe -o=s -n=name -t=target -c=content -host=host

# pull message 
client_cli.exe -o=p -n=name -host=host 
```
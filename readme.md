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
bash build.sh
```

## Usage Client

建议使用 [obs-studio](https://github.com/obsproject/obs-studio/releases) 进行推流，推荐输出分辨率为768*432
直播服务器请使用 [Glarmogann/livego](https://github.com/Glamorgann/livego) 

## Usage cli
```cmd
<Windows>
# operation [p]ull , [s]end
# send message
client_cli.exe -o=s -n=name -t=target -c=content -host=host

# pull message 
client_cli.exe -o=p -n=name -host=host 
```
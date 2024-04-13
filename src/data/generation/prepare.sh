#!/bin/bash

if [ ! -f server.jar ]; then
    echo "Downloading server.jar..."
    # 1.20.4 server jar
    curl -L https://piston-data.mojang.com/v1/objects/8dd1a28015f51b1803213892b50b7b4fc76e594d/server.jar > server.jar
fi

java -DbundlerMainClass=net.minecraft.data.Main -jar server.jar --reports

rm -rf versions libraries generated/.cache logs

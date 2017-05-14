# /bin/bash

# 制作自己的基础镜像，构建好node环境和npm 以及 dependencies
# 因为每次build都得下载一次 dependencies，很累
cp ../../package.json package.json
docker build -t nodenpm .
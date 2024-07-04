#!/bin/sh

# 检查是否存在./runtime/config.yaml文件
if [ ! -f ./runtime/config.yaml ]; then
    echo "配置文件不存在，复制默认配置文件"
    mkdir -p runtime
    cp ./configs/config.yaml.example ./runtime/config.yaml
    # 设置默认字符串长度
    LENGTH=16
    # 生成随机字符串
    RANDOM_STRING=$(head /dev/urandom | tr -dc A-Za-z0-9 | head -c $LENGTH)
    # 替换配置文件中的变量
    sed -i "s/\${RANDOM_STRING}/${RANDOM_STRING}/g" ./runtime/config.yaml
fi

# 启动应用
exec "$@"
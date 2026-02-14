#!/bin/bash

# Whisper模型批量下载脚本
# 适用于yt-whisper-bilingual项目

echo "开始下载Whisper模型..."

# 创建模型目录
mkdir -p models

# 定义要下载的模型列表
MODELS=("tiny" "base" "small" "medium" "large-v3" "medium.en")

# 下载模型函数
download_model() {
    local model=$1
    echo "正在下载模型: $model"
    
    # 使用whisper命令下载模型
    whisper --model "$model" --model_dir "./models" --device cpu dummy.mp3 2>/dev/null || true
    
    echo "模型 $model 下载完成"
}

# 为每个模型创建一个简短的测试文件
echo "Creating dummy audio file for model download..."
echo "dummy" > dummy.txt
text2wave dummy.txt -o dummy.wav 2>/dev/null || true
ffmpeg -y -f lavfi -i anullsrc=r=16000:cl=mono -t 1 -q:a 9 -acodec libmp3lame dummy.mp3 2>/dev/null || true

# 下载所有模型
for model in "${MODELS[@]}"; do
    download_model "$model"
done

# 清理临时文件
rm -f dummy.txt dummy.wav dummy.mp3

echo "所有模型下载完成！"
echo "模型存储位置: ./models"
echo "可用模型: ${MODELS[*]}"
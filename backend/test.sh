#!/bin/bash

# 后端测试脚本

echo "🧪 开始运行后端测试..."

# 运行所有测试
echo ""
echo "📦 运行单元测试..."
go test -v ./...

# 检查测试结果
if [ $? -eq 0 ]; then
    echo ""
    echo "✅ 所有测试通过！"
    
    # 显示测试覆盖率
    echo ""
    echo "📊 生成测试覆盖率报告..."
    go test -coverprofile=coverage.out ./...
    go tool cover -func=coverage.out | tail -1
    
    echo ""
    echo "💡 查看详细覆盖率报告: go tool cover -html=coverage.out"
else
    echo ""
    echo "❌ 测试失败！"
    exit 1
fi

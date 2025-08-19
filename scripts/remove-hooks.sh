#!/bin/bash
set -eu

##################################
# Git hooksを無効化するスクリプト
##################################

echo "🔧 Removing Git hooks configuration..."

# Git設定からhooksPathを削除
git config --unset core.hooksPath

if [ $? -eq 0 ]; then
    echo "✅ Git hooks disabled successfully!"
    echo "Hooks will no longer run automatically."
else
    echo "ℹ️  Git hooks were not configured or already disabled."
fi

echo ""
echo "To re-enable: run './scripts/setup-hooks.sh'"

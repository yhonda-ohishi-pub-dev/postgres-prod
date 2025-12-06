#!/bin/bash
# Execute harness for MINGW64/Git Bash
# Runs plan.md tasks automatically using claude CLI

set -e

cd "$(dirname "$0")/.."

MAX_ITERATIONS=${1:-20}
iteration=0

# Read the execute command content
EXECUTE_PROMPT=$(cat .claude/commands/execute.md)

echo "============================================================"
echo "Execute Harness - 自動実行スクリプト"
echo "============================================================"

while [ $iteration -lt $MAX_ITERATIONS ]; do
    iteration=$((iteration + 1))

    echo ""
    echo "============================================================"
    echo "[Iteration $iteration/$MAX_ITERATIONS]"
    echo "============================================================"

    # Check for pending phases (Phase 1-3 through 1-7, or Phase 2)
    pending=$(grep -E "^####\s+Phase\s+[12]-[0-9]" plan.md 2>/dev/null | grep -v "✅" || true)

    if [ -z "$pending" ]; then
        echo ""
        echo "============================================================"
        echo "全てのタスクが完了しました！"
        echo "============================================================"
        break
    fi

    echo "残りのタスク:"
    echo "$pending"
    echo ""

    # Run claude with the execute prompt directly
    echo "Claude Codeを実行中..."
    claude -p "$EXECUTE_PROMPT"

    # Small delay
    sleep 2
done

if [ $iteration -ge $MAX_ITERATIONS ]; then
    echo "Warning: Maximum iterations ($MAX_ITERATIONS) reached"
fi

echo ""
echo "Execute Harness 完了"

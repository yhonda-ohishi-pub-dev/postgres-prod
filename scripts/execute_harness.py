#!/usr/bin/env python3
"""
Execute Harness Script using Claude Agent SDK
plan.mdの「次のステップ（予定）」を自動実行するハーネススクリプト

Usage:
    py scripts/execute_harness.py [--max N] [--status]

Requirements:
    pip install claude-agent-sdk
"""

import sys
import re
import argparse
import asyncio
from pathlib import Path

try:
    from claude_code_sdk import query, ClaudeCodeOptions, Message
except ImportError:
    print("Error: claude-code-sdk not installed")
    print("Run: pip install claude-code-sdk")
    sys.exit(1)

PROJECT_ROOT = Path(__file__).parent.parent
PLAN_FILE = PROJECT_ROOT / "plan.md"
EXECUTE_CMD_FILE = PROJECT_ROOT / ".claude" / "commands" / "execute.md"


def read_file(path: Path) -> str:
    with open(path, "r", encoding="utf-8") as f:
        return f.read()


def get_pending_phases(content: str) -> list[str]:
    """未完了のPhaseを取得"""
    pending = []
    for line in content.split("\n"):
        if re.match(r"^####\s+Phase\s+\d+-\d+", line):
            if "✅" not in line:
                pending.append(line.strip().replace("#### ", ""))
    return pending


async def run_claude(prompt: str) -> bool:
    """Claude Agent SDKを使って実行"""
    print("\n" + "=" * 60)
    print("Claude Agent SDK を実行中...")
    print("=" * 60 + "\n")

    try:
        options = ClaudeCodeOptions(
            allowed_tools=["Read", "Write", "Edit", "Bash", "Glob", "Grep", "Task"],
        )

        msg_count = 0
        async for message in query(
            prompt=prompt,
            options=options,
        ):
            msg_count += 1
            # メッセージタイプに応じて出力
            if isinstance(message, Message):
                if hasattr(message, 'content'):
                    print(message.content, flush=True)
                elif hasattr(message, 'text'):
                    print(message.text, flush=True)
            else:
                # 詳細なメッセージ情報を出力
                print(f"[{msg_count}] {type(message).__name__}: {str(message)[:200]}", flush=True)

        print(f"\n完了: {msg_count}メッセージ処理")
        return True
    except Exception as e:
        import traceback
        print(f"\nError: {type(e).__name__}: {e}")
        traceback.print_exc()
        return False


async def main_async(args):
    """非同期メイン処理"""
    print("=" * 60)
    print("Execute Harness (Claude Agent SDK)")
    print("=" * 60)

    # Read files
    try:
        plan_content = read_file(PLAN_FILE)
        execute_prompt = read_file(EXECUTE_CMD_FILE)
    except FileNotFoundError as e:
        print(f"Error: {e}")
        return 1

    # Check pending
    pending = get_pending_phases(plan_content)

    if not pending:
        print("\n全てのタスクが完了しています！")
        return 0

    print(f"\n残りのPhase: {len(pending)}個")
    for i, phase in enumerate(pending, 1):
        print(f"  {i}. {phase}")

    if args.status:
        return 0

    # Execute loop
    for i in range(1, args.max + 1):
        print(f"\n[Iteration {i}/{args.max}]")

        # Re-check pending
        plan_content = read_file(PLAN_FILE)
        pending = get_pending_phases(plan_content)

        if not pending:
            print("\n" + "=" * 60)
            print("全てのタスクが完了しました！")
            print("=" * 60)
            break

        print(f"残り: {len(pending)}個 - {pending[0]}")

        # Run claude
        success = await run_claude(execute_prompt)

        if not success:
            print("\nClaude実行に失敗しました")
            try:
                if input("続行? (y/n): ").strip().lower() != 'y':
                    break
            except EOFError:
                break

    print("\nExecute Harness 完了")
    return 0


def main():
    parser = argparse.ArgumentParser(description="Execute plan.md tasks with Claude Agent SDK")
    parser.add_argument("--max", "-n", type=int, default=20, help="Max iterations")
    parser.add_argument("--status", "-s", action="store_true", help="Show status only")
    args = parser.parse_args()

    return asyncio.run(main_async(args))


if __name__ == "__main__":
    sys.exit(main())

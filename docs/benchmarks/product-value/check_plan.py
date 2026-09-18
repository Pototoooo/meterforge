from pathlib import Path
import sys
s=Path(sys.argv[1]).read_text()
print("scope=automated_acceptance_appendix" if "2026-09-11 本轮执行补充" in s else "scope=original_proposal")
print("human_results=not_measured")

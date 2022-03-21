import os
import sys

# Temporary workaround to include siftool
project_root = os.path.abspath(os.path.join(os.path.dirname(__file__), *([os.path.pardir] * 4)))
base_dir = os.path.join(project_root, "test", "integration", "framework")
enabled = False
for p in sys.path:
    enabled = enabled or os.path.realpath(p) == os.path.realpath(base_dir)
if not enabled:
    sys.path = sys.path + [base_dir]

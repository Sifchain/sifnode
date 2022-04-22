import os
import sys

# Temporary workaround to include siftool
project_root = os.path.abspath(os.path.join(os.path.dirname(__file__), *([os.path.pardir] * 4)))
base_dir = os.path.join(project_root, "test", "integration", "framework")
src_dir = os.path.join(base_dir, "src")
build_generated_dir = os.path.join(base_dir, "build", "generated")
paths = [src_dir, build_generated_dir]
enabled = False
paths_to_add = []
for p in paths:
    enabled = any([os.path.realpath(p) == os.path.realpath(s) for s in sys.path])
    if not enabled:
        paths_to_add.append(p)
sys.path.extend(paths_to_add)

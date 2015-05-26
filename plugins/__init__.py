__all__ = []
from . import enhancer
import glob
import os

# glob all submodules and force the init of plugins to import them
cwd = os.path.dirname(__file__)
for module in glob.glob(cwd + "/*.py"):
    name = os.path.basename(module)
    if "__init__" not in name:
        # remove the ".py" in the module name to add it to __all__
        __all__.append(name[:-3])

from . import *

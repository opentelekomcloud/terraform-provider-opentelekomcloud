venv="./.rn-venv"
python -m venv $venv >/dev/null || exit 1
echo "Python venv created"
source $venv/bin/activate || exit 1
pip install reno sphinx otcdocstheme >/dev/null || exit 1
echo "Dependencies installed"

built_path="releasenotes/build/html"
rm -rf $built_path
echo "Old docs removed"

sphinx-build -W --keep-going -b html releasenotes/source $built_path || exit 1
echo "Release notes are available at $(pwd)/releasenotes/build/html/current.html"

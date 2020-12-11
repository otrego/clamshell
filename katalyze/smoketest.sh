#!/bin/bash

set -e

cd "$(dirname "$0")"

echo "**Smoke Test runs katalyze on fixed data to test the system is working **"

go run . \
-config=testdata/analysis_example.cfg \
-model=testdata/g170e-b10c128-s1141046784-d204142634.bin.gz \
-output_dir=output \
testdata/

# -model=testdata/g170e-b20c256x2-s5303129600-d1228401921.bin.gz \



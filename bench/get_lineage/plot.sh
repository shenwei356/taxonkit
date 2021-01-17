#!/bin/sh

# ./plot.R -i bench.get_lineage.tsv
sed -i 's/taxids.//' bench.get_lineage.tsv
sed -i 's/.txt//' bench.get_lineage.tsv 
./plot.R -i bench.get_lineage.tsv --width 8 --height 3 --lx 0.75 --ly 0.5

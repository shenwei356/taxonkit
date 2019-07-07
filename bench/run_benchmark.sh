#!/bin/sh

echo Test: Querying lineage

function check() {
    md5sum $1
    # /bin/rm $1
}


for f in taxids.*.txt; do
    echo == TaxonKit
    echo data: $f
    out=$f.taxonkit.lineage
    memusg -t -H taxonkit lineage -d "; " < $f > $out
    check $out

    # echo == Biopython
    # echo data: $f
    # out=$f.biopython.lineage
    # memusg -t -H python3 get_lineage.biopython.py < $f > $out
    # check $out

    echo == ETE
    echo data: $f
    out=$f.ete.lineage
    memusg -t -H python3 get_lineage.ete.py < $f > $out
    check $out

    echo == taxadb
    echo data: $f
    out=$f.taxadb.lineage
    memusg -t -H python3 get_lineage.taxadb.py < $f > $out
    check $out
done

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
    memusg -t -H -s " taxonkit lineage --delimiter \"; \" --threads 1  < $f > $out "
    check $out

    echo == ETE
    echo data: $f
    out=$f.ete.lineage
    memusg -t -H -s " /usr/bin/python3 get_lineage.ete.py < $f > $out "
    check $out
    
    echo == taxopy
    echo data: $f
    out=$f.taxopy.lineage
    memusg -t -H -s " /usr/bin/python3 get_lineage.taxopy.py < $f > $out "
    check $out

    echo == taxadb
    echo data: $f
    out=$f.taxadb.lineage
    memusg -t -H -s " /usr/bin/python3 get_lineage.taxadb.py < $f > $out "
    check $out
done

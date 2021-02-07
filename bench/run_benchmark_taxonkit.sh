#!/bin/sh

echo Test: Threads scalability

function check() {
    md5sum $1
    # /bin/rm $1
}


for f in taxids.n*.txt; do
    for t in 1 2 4 8; do
        echo == t=$t
        echo data: $f
        out=$f.taxonkit-threads-$t.lineage
        memusg -t -H -s " taxonkit lineage --delimiter \"; \" --threads $t  < $f > $out "
        check $out
    done
done

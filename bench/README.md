# Benchmark

## Data set

- [NCBI taxonomy](ftp://ftp.ncbi.nih.gov/taxonomy), version 2021-01-21

- TaxIDs. Root node `1` is removed. 
  And These data should be updated along with NCBI taxonomy dataset, 
  because `taxadb` does not handle merged taxiIDs.
    - small.txt (n=20)

            cut -f 1 nodes.dmp | csvtk sample -H -p 0.00002 \
                | grep -w 1 -v | head -n 20 > taxids.small.txt

    - medium.txt (n=2,000)

            cut -f 1 nodes.dmp | csvtk sample -H -p 0.002 \
                | grep -w 1 -v | head -n 2000 > taxids.medium.txt

    - big.txt (n=200,000)

            cut -f 1 nodes.dmp | csvtk sample -H -p 0.1 \
                | grep -w 1 -v | head -n 200000  > taxids.big.txt

## Softwares

- ETE, version: [3.1.2](https://pypi.org/project/ete3/3.1.2/)
- taxadb, version: [0.12.0](https://pypi.org/project/taxadb/0.12.0)
- TaxonKit, version: [0.6.2](https://github.com/shenwei356/taxonkit/releases/tag/0.6.2)
- BioPython, version: [1.78](https://pypi.org/project/biopython/1.78/)

Not used tools without direct function of getting full lineage

- [opensci/taxadb](https://github.com/ropensci/taxadb)
- [ncbi-taxonomist](https://ncbi-taxonomist.readthedocs.io/en/latest/)


## Installation and Configurations

- ETE

        sudo pip3 install ete3
        
        # create database
        # http://etetoolkit.org/docs/latest/tutorial/tutorial_ncbitaxonomy.html#upgrading-the-local-database
        from ete3 import NCBITaxa
        ncbi = NCBITaxa()
        ncbi.update_taxonomy_database()

- Biopython

        sudo pip3 install biopython

- taxadb

        sudo pip3 install -U taxadb
        
        # create database
        # https://github.com/HadrienG/taxadb#creating-the-database
        taxadb download --type taxa -o ~/.taxadb -f
        taxadb create -i ~/.taxadb --division taxa  --dbname ~/.taxadb/taxadb.sqlite


## Scripts and Commands

Scripts/Command as listed below.
Python scripts were written followed to the official documents,
and **parallelized querying were not used, including TaxonKit**.

    ETE             get_lineage.ete.py
    Biopython       get_lineage.biopython.py
    taxadb          get_lineage.taxadb.py
    taxonkit        taxonkit lineage -d "; "

A Python script [memusg](https://github.com/shenwei356/memusg) was used
to computate running time and peak memory usage of a process.
A Perl scripts
[`run.pl`](https://github.com/shenwei356/seqkit/blob/master/bench/run.pl)
is used to automatically running tests and generate data for plotting.

**Note that `get_lineage.biopython.py` is not used,
because quering via internet (entrez) is too slow for large number of queries.**

Running benchmark:

    # 55min for me...
    time perl run.pl -n 3 run_benchmark.sh -o bench.get_lineage.tsv

    # clear
    # rm *.lineage
    # rm *.out

Checking result:

    md5sum *.lineage
    4d89c6cafa9e5fc75b3166d9cc1fd9c7  taxids.big.txt.ete.lineage
    4d89c6cafa9e5fc75b3166d9cc1fd9c7  taxids.big.txt.taxadb.lineage
    4d89c6cafa9e5fc75b3166d9cc1fd9c7  taxids.big.txt.taxonkit.lineage
    4ef1d6afd94f8d5034ad46670610bfda  taxids.medium.txt.ete.lineage
    4ef1d6afd94f8d5034ad46670610bfda  taxids.medium.txt.taxadb.lineage
    4ef1d6afd94f8d5034ad46670610bfda  taxids.medium.txt.taxonkit.lineage
    c2596fc28068b7dad4af59309c7b8d12  taxids.small.txt.ete.lineage
    c2596fc28068b7dad4af59309c7b8d12  taxids.small.txt.taxadb.lineage
    c2596fc28068b7dad4af59309c7b8d12  taxids.small.txt.taxonkit.lineage


Note that taxonkit returns "root" for taxid 1, while the others return nothing.

Plotting benchmark result. 
R libraries `dplyr`, `ggplot2`, `scales`, `ggthemes`, `ggrepel` are needed.

    # reformat dataset
    # tools: https://github.com/shenwei356/csvtk/
    for f in taxids.*.txt;  do wc -l $f; done \
        | csvtk space2tab | csvtk cut -H -t -f 2,1 \
        | csvtk replace -H -t -f 2 -p ^ -r n= \
        > dataset_rename.tsv

    cat bench.get_lineage.tsv \
        | csvtk replace -t -f dataset -k dataset_rename.tsv -p '(.+)' -r '{kv}' \
        | csvtk sort -t -k dataset:N \
        > bench.get_lineage.reformat.tsv

    ./plot.R -i bench.get_lineage.reformat.tsv --width 8 --height 3.3

Result

![](bench.get_lineage.reformat.tsv.png)

<div id="disqus_thread"></div>
<script>

/**
*  RECOMMENDED CONFIGURATION VARIABLES: EDIT AND UNCOMMENT THE SECTION BELOW TO INSERT DYNAMIC VALUES FROM YOUR PLATFORM OR CMS.
*  LEARN WHY DEFINING THESE VARIABLES IS IMPORTANT: https://disqus.com/admin/universalcode/#configuration-variables*/
/*
var disqus_config = function () {
this.page.url = PAGE_URL;  // Replace PAGE_URL with your page's canonical URL variable
this.page.identifier = PAGE_IDENTIFIER; // Replace PAGE_IDENTIFIER with your page's unique identifier variable
};
*/
(function() { // DON'T EDIT BELOW THIS LINE
var d = document, s = d.createElement('script');
s.src = '//taxonkit.disqus.com/embed.js';
s.setAttribute('data-timestamp', +new Date());
(d.head || d.body).appendChild(s);
})();
</script>
<noscript>Please enable JavaScript to view the <a href="https://disqus.com/?ref_noscript">comments powered by Disqus.</a></noscript>


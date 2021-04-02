# Benchmark

## Benchmark 1: Getting lineage

### Data set

- [NCBI taxonomy](ftp://ftp.ncbi.nih.gov/taxonomy), version 2021-01-21

- TaxIDs. Root node `1` is removed. 
  And These data should be updated along with NCBI taxonomy dataset.
  Seven sizes of TaxIds are sampled from `nodes.dmp`.

        # shuffle all taxids
        cut -f 1 nodes.dmp | grep -w -v 1 | shuf > ids.txt
        
        # extract n taxids for testing
        for n in 1 10 100 1000 2000 4000 6000 8000 10000 20000 40000 60000 80000 100000; do 
            head -n $n ids.txt > taxids.n$n.txt
        done


### Software

- Loading database from local database:
    - ETE, version: [3.1.2](https://pypi.org/project/ete3/3.1.2/)
- Directly parsing dump files:
    - taxopy, version: [0.5.0](https://github.com/apcamargo/taxopy/releases/tag/v0.5.0)
    - TaxonKit, version: [0.7.2](https://github.com/shenwei356/taxonkit/releases/tag/0.7.2)


### Environment

- OS: Linux 5.4.89-1-MANJARO
- CPU: AMD Ryzen 7 2700X Eight-Core Processor, 3.7GHz
- RAM: 64GB DDR4 3000MHz
- SSD: Samsung 970EVO 500G NVMe SSD

### Installation and Configurations

- ETE

        sudo pip3 install ete3
        
        # create database
        # http://etetoolkit.org/docs/latest/tutorial/tutorial_ncbitaxonomy.html#upgrading-the-local-database
        from ete3 import NCBITaxa
        ncbi = NCBITaxa()
        ncbi.update_taxonomy_database()

- TaxonKit

        mkdir -p $HOME/.taxonkit
        mkdir -p $HOME/bin/
        
        # data
        wget -c ftp://ftp.ncbi.nih.gov/pub/taxonomy/taxdump.tar.gz 
        tar -zxvf taxdump.tar.gz -C $HOME/.taxonkit
        
        # binary
        wget https://github.com/shenwei356/taxonkit/releases/download/v0.7.2/taxonkit_linux_amd64.tar.gz
        tar -zxvf taxonkit_linux_amd64.tar.gz -C $HOME/bin/        

- taxopy

        sudo pip3 install -U taxopy
        
        # taxoopy identical dump files copied from taxonkit
        mkdir -p ~/.taxopy
        cp ~/.taxonkit/{nodes.dmp,names.dmp} ~/.taxopy


### Scripts and Commands

Scripts/Command as listed below.
Python scripts were written following to the official documents,
and **parallelized querying were not used, including TaxonKit**.

    ETE             get_lineage.ete.py                              < $infile > $outfile
    taxopy          get_lineage.taxopy.py                           < $infile > $outfile
    taxonkit        taxonkit lineage --threads 1 --delimiter "; "   < $infile > $outfile

A Python script [memusg](https://github.com/shenwei356/memusg) was used
to computate running time and peak memory usage of a process.
A Perl scripts
[`run.pl`](https://github.com/shenwei356/seqkit/blob/master/bench/run.pl)
is used to automatically running tests and generate data for plotting.


Running benchmark:

    $ # emptying the buffers cache
    $ su -c "free && sync && echo 3 > /proc/sys/vm/drop_caches && free"
    
    time perl run.pl -n 3 run_benchmark.sh -o bench.get_lineage.tsv

Checking result:

    $ md5sum taxids.n*.lineage
    
    # clear
    $ rm *.lineage *.out

Plotting benchmark result. 
R libraries `dplyr`, `ggplot2`, `scales`, `ggthemes`, `ggrepel` are needed.

    # reformat dataset
    # tools: https://github.com/shenwei356/csvtk/

    for f in taxids.n*.txt;  do wc -l $f; done \
        | sort -k 1,1n \
        | awk '{ print($2"\t"$1) }' \
        > dataset_rename.tsv

    cat bench.get_lineage.tsv \
        | csvtk sort -t -L dataset:<(cut -f 1 dataset_rename.tsv) -k dataset:u -k app \
        | csvtk replace -t -f dataset -k dataset_rename.tsv -p '(.+)' -r '{kv}' \
        > bench.get_lineage.reformat.tsv

    ./plot2.R -i bench.get_lineage.reformat.tsv --width 6 --height 4 --dpi 600 \
        --labcolor "log10(queries)" --labshape "Tools"

Result

<img src="bench.get_lineage.reformat.tsv.png" alt="" width="600" align="center" />

## Benchmark 2: TaxonKit multi-threaded scalability

Running benchmark:

    $ # emptying the buffers cache
    $ su -c "free && sync && echo 3 > /proc/sys/vm/drop_caches && free"
    

    $ time perl run.pl -n 3 run_benchmark_taxonkit.sh -o bench.taxonkit.tsv
    $ rm *.lineage *.out
    
Plotting benchmark result.

    cat bench.taxonkit.tsv \
        | csvtk sort -t -L dataset:<(cut -f 1 dataset_rename.tsv) -k dataset:u -k app \
        | csvtk replace -t -f dataset -k dataset_rename.tsv -p '(.+)' -r '{kv}' \
        > bench.taxonkit.reformat.tsv

    ./plot_threads2.R -i bench.taxonkit.reformat.tsv --width 6 --height 4 --dpi 600 \
        --labcolor "log10(queries)" --labshape "Threads"
    
Result

<img src="bench.taxonkit.reformat.tsv.png" alt="" width="600" align="center" />



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


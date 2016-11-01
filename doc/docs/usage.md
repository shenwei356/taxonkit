# Usage and Examples

## Before use

Please download and uncompress these files:

- [ftp://ftp.ncbi.nih.gov/pub/taxonomy/taxdump.tar.gz](ftp://ftp.ncbi.nih.gov/pub/taxonomy/taxdump.tar.gz)


## taxonkit

Usage

```
TaxonKit - NCBI Taxonomy Toolkit

Version: 0.1

Author: Wei Shen <shenwei356@gmail.com>

Source code: https://github.com/shenwei356/taxonkit

Dataset:

    Please download and uncompress these files:
    ftp://ftp.ncbi.nih.gov/pub/taxonomy/taxdump.tar.gz

Usage:
  TaxonKit [command]

Available Commands:
  list        list taxon tree of given taxon IDs
  version     print version information and check for update

Flags:
  -o, --out-file string   out file ("-" for stdout, suffix .gz for gzipped out) (default "-")
  -j, --threads int       number of CPUs. (default value: 1 for single-CPU PC, 2 for others) (default 2)

Use "TaxonKit [command] --help" for more information about a command.

```

## list

Usage

```
list taxon tree of given taxon IDs.

Usage:
  taxonkit list [flags]

Flags:
      --ids string      taxon ID(s), multiple IDs should be seperated by comma (default "1")
      --indent string   indent (default "  ")
      --names string    names.dmp file, when it given taxid will be followed by its scientific name
      --nodes string    nodes.dmp file (default "nodes.dmp")
      --show-rank       show rank of the node

```

Examples

1. default usage

        $ taxonkit list --nodes nodes.dmp --ids 9605
        9605
          9606
            63221
            741158
          1425170

1. change ident, the list could be used to extract sequences from BLAST database with `blastdbcmd`

        $ taxonkit list --nodes nodes.dmp --ids 9605 --indent ""
        9605
        9606
        63221
        741158
        1425170

1. adding names

        $ taxonkit list --nodes nodes.dmp --names names.dmp --ids 9605
        9605 [genus] Homo
          9606 [species] Homo sapiens
            63221 [subspecies] Homo sapiens neanderthalensis
            741158 [subspecies] Homo sapiens ssp. Denisova
          1425170 [species] Homo heidelbergensis



<div id="disqus_thread"></div>
<script>
/**
* RECOMMENDED CONFIGURATION VARIABLES: EDIT AND UNCOMMENT THE SECTION BELOW TO INSERT DYNAMIC VALUES FROM YOUR PLATFORM OR CMS.
* LEARN WHY DEFINING THESE VARIABLES IS IMPORTANT: https://disqus.com/admin/universalcode/#configuration-variables
*/
/*
var disqus_config = function () {
this.page.url = PAGE_URL; // Replace PAGE_URL with your page's canonical URL variable
this.page.identifier = PAGE_IDENTIFIER; // Replace PAGE_IDENTIFIER with your page's unique identifier variable
};
*/
(function() { // DON'T EDIT BELOW THIS LINE
var d = document, s = d.createElement('script');

s.src = '//csvtk.disqus.com/embed.js';

s.setAttribute('data-timestamp', +new Date());
(d.head || d.body).appendChild(s);
})();
</script>
<noscript>Please enable JavaScript to view the <a href="https://disqus.com/?ref_noscript" rel="nofollow">comments powered by Disqus.</a></noscript>

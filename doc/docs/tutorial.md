# Tutorial

## Extract all sequences of certen taxons from the nr database

### Dataset

- [prot.accession2taxid.gz](ftp://ftp.ncbi.nih.gov/pub/taxonomy/accession2taxid/prot.accession2taxid.gz)

### Steps

Taking bacteria for example.

1. Getting all taxids of bacteria (taxid 2):

        $ taxonkit list --nodes nodes.dmp --ids 2 --indent "" > bacteria.taxid.txt

    It takes only 2.5s! Number of taxids:

        $ wc -l bacteria.taxid.txt
        454591 bacteria.taxid.txt

2. Extacting accessions with [csvtk](http://bioinf.shenwei.me/csvtk/download/):

        $ csvtk -t grep -f taxid -P bacteria.taxid.txt prot.accession2taxid.gz | csvtk -t cut -f accession.version > bacteria.taxid.acc.txt

3. Extracting nr sequences:

        $ blastdbcmd -db nr -entry all -outfmt "%a\t%T" | \
            csvtk -t grep -f 2 -P bacteria.taxid.acc.txt | \
            csvtk -t cut -f 1 | \
            blastdbcmd -db nr -entry_batch - -out bacteria.fa

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

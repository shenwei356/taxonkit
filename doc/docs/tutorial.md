# Tutorial

## Extract all protein sequences of specific taxons from the NCBI nr database

### Dataset

- [prot.accession2taxid.gz](ftp://ftp.ncbi.nih.gov/pub/taxonomy/accession2taxid/prot.accession2taxid.gz)

### Steps

Taking virus for example.

1. Getting all taxids of virus (taxid 10239):

        $ taxonkit list --ids 10239 --indent "" > virus.taxid.txt

    It takes only 2.5 seconds! Number of taxids:

        $ wc -l virus.taxid.txt
        163104 virus.taxid.txt

2. Retrieving accessions <s>or GIs</s>(NCBI stopped using `gi`) with [csvtk](http://bioinf.shenwei.me/csvtk/download/):

    - accession

            $ zcat prot.accession2taxid.gz | \
                csvtk -t grep -f taxid -P virus.taxid.txt | \
                csvtk -t cut -f accession.version > virus.taxid.acc.txt
    
    - <s>gi

            $ zcat prot.accession2taxid.gz | \
                csvtk -t grep -f taxid -P virus.taxid.txt | \
                csvtk -t cut -f gi > virus.taxid.gi.txt
        </s>
        It costs ~ 8 minutes.

3. Retrieving nr sequences from BLAST database:

    - accesion

            $ blastdbcmd -db nr -entry all -outfmt "%a\t%T" | \
                csvtk -t grep -f 2 -P virus.taxid.acc.txt | \
                csvtk -t cut -f 1 | \
                blastdbcmd -db nr -entry_batch - -out nr.virus.fa

    - <s>gi

            $ blastdbcmd -db nr -entry all -outfmt "%g\t%T" | \
                csvtk -t grep -f 2 -P virus.taxid.gi.txt | \
                csvtk -t cut -f 1 | \
                blastdbcmd -db nr -entry_batch - -out nr.virus.fa
        </s>

    **Another way is directly retrieving from [nr FASTA sequences](ftp://ftp.ncbi.nih.gov/blast/db/FASTA/nr.gz) using [SeqKit](http://bioinf.shenwei.me/seqkit/download):**
    
        seqkit grep -f virus.taxid.acc.txt nr.gz | gzip -c > nr.virus.fa.gz

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

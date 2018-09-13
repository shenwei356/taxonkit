# Tutorial

## Making nr blastdb for specific taxids

[rewrite on 2018-09-13]

Data:

- [pre-formated blastdb](ftp://ftp.ncbi.nlm.nih.gov/blast/db) (09/10/2018) 
- [prot.accession2taxid.gz](ftp://ftp.ncbi.nih.gov/pub/taxonomy/accession2taxid/prot.accession2taxid.gz) (09/07/2018) (optional, but recommended)

Hardware in this tutorial

- CPU: AMD 8-cores/16-threads 3.7Ghz
- RAM: 32GB
- DISK:
    - Taxonomy files stores in NVMe SSD
    - blatdb files stores in 7200rpm HDD
    
Tools:

- [blast+](ftp://ftp.ncbi.nlm.nih.gov/blast/executables/blast+/LATEST/)
- [pigz](https://zlib.net/pigz/) (recommended, faster than gzip)
- [taxonkit](https://bioinf.shenwei.me/taxonkit)
- [seqkit](https://bioinf.shenwei.me/seqkit) (recommended)
    
Steps:

1. Listing all taxids below `$id` using taxonkit.

        id=6656
        
        # 6656 is the phylum Arthropoda
        # echo 6656 | taxonkit lineage | taxonkit reformat 
        # 6656    cellular organisms;Eukaryota;Opisthokonta;Metazoa;Eumetazoa;Bilateria;Protostomia;Ecdysozoa;Panarthropoda;Arthropoda    Eukaryota;Arthropoda;;;;;
        
        # time: 3s 
        taxonkit list --ids $id --indent "" > $id.taxid.txt 
        
        wc -l $id.taxid.txt
        # 518373 6656.taxid.txt

2. Retrieving target accessions. There are two options:

    1. From prot.accession2taxid.gz (**faster, recommended**). *Note that some accessions are not in `nr`*.
    
            # time: 4min
            pigz -dc prot.accession2taxid.gz \
                | csvtk grep -t -f taxid -P $id.taxid.txt \
                | csvtk cut -t -f accession.version \
                | sed 1d \
                > $id.acc.txt
            
            
            wc -l $id.acc.txt
            # 8174609 6656.acc.txt
    
    1. From pre-formated `nr` blastdb 
    
            # time: 40min
            blastdbcmd -db nr -entry all -outfmt "%a %T" | pigz -c > nr.acc2taxid.txt.gz
            
            pigz -dc nr.acc2taxid.txt.gz | wc -l
            # 555220892
            
            # time: 3min
            pigz -dc nr.acc2taxid.txt.gz \
                | csvtk grep -d ' ' -D ' ' -f 2 -P $id.taxid.txt \
                | cut -d ' '  -f 1 \
                > $id.acc.txt
            
            
            wc -l $id.acc.txt
            # 6928021 6656.acc.txt
    
    
3. Retrieving FASTA sequences from pre-formated blastdb. There are two options:

   
    1. From `nr.fa` exported from pre-formated blastdb (**faster, smaller output file, recommended**). 
       **DO NOT directly download `nr.gz` from [ncbi ftp](ftp://ftp.ncbi.nih.gov/blast/db/FASTA/nr.gz)**,
       in which the FASTA headers are not well formated.
    
            # time: 117min
            blastdbcmd -db nr -dbtype prot -entry all -outfmt "%f" -out - | pigz -c > nr.fa.gz
            
            # time: 80min
            # perl one-liner is used to unfold records having mulitple accessions
            cat <(echo) <(pigz -dc nr.fa.gz) \
                | perl -e 'BEGIN{ $/ = "\n>"; <>; } while(<>){s/>$//;  $i = index $_, "\n"; $h = substr $_, 0, $i; $s = substr $_, $i+1; if ($h !~ />/) { print ">$_"; next; }; $h = ">$h"; while($h =~ />([^ ]+ .+?) ?(?=>|$)/g){ $h1 = $1; $h1 =~ s/^\W+//; print ">$h1\n$s";} } ' \
                | seqkit grep --delete-matched -f $id.acc.txt -o nr.$id.fa.gz
                
            # counting sequences
            # 
            # ls -lh nr.$id.fa.gz
            # -rw-r--r-- 1 shenwei shenwei 902M 9月  13 01:42 nr.6656.fa.gz
            # 
            pigz -dc nr.$id.fa.gz | grep '^>' -c
            
            # 6928017
            # Here 6928017 ~=  6928021 ($id.acc.txt)
        
    1. Directly from pre-formated blastdb
    
            # time: 5h20min
            blastdbcmd -db nr -entry_batch $id.acc.txt -out - | pigz -c > nr.$id.fa.gz
            
            
            # counting sequences
            #
            # Note that the headers of outputed fasta by blastdbcmd are "folded"
            # for accessions from different species with same sequences, so the 
            # number may be small than $(wc -l $id.acc.txt).
            pigz -dc nr.$id.fa.gz | grep '^>' -c
            # 1577383
            
            # counting accessions
            #  
            # ls -lh nr.$id.fa.gz
            # -rw-r--r-- 1 shenwei shenwei 2.1G 9月  13 03:38 nr.6656.fa.gz
            #
            pigz -dc nr.$id.fa.gz | grep '^>' | sed 's/>/\n>/g' | grep '^>' -c
            # 288415413
        
4. makeblastdb

        pigz -dc nr.$id.fa.gz > nr.$id.fa
        
        # time: 3min ($nr.$id.fa from step 3 option 1)
        #
        # building $nr.$id.fa from step 3 option 2 with -parse_seqids would produce error:
        #
        #     BLAST Database creation error: Error: Duplicate seq_ids are found: SP|P29868.1
        #
        makeblastdb -parse_seqids -in nr.$id.fa -dbtype prot -out nr.$id
        
        # rm nr.$id.fa

5. blastp (optional)

        # blastdb nr.$id is built from sequences in step 3 option 1
        # 
        blastp -num_threads 16 -db nr.$id -query t4.fa > t4.fa.blast
        # real    0m20.866s
        
        # $ cat t4.fa.blast | grep Query= -A 10
        # Query= A0A0J9X1W9.2 RecName: Full=Mu-theraphotoxin-Hd1a; Short=Mu-TRTX-Hd1a
        #    
        # Length=35
                                                                             Score     E
        # Sequences producing significant alignments:                          (Bits)  Value

        # 2MPQ_A  Chain A, Solution structure of the sodium channel toxin Hd1a  72.4    2e-17
        # A0A0J9X1W9.2  RecName: Full=Mu-theraphotoxin-Hd1a; Short=Mu-TRTX-...  72.4    2e-17
        # ADB56726.1  HNTX-IV.2 precursor [Haplopelma hainanum]                 66.6    9e-15
        # D2Y233.1  RecName: Full=Mu-theraphotoxin-Hhn1b 2; Short=Mu-TRTX-H...  66.6    9e-15
        # ADB56830.1  HNTX-IV.3 precursor [Haplopelma hainanum]                 66.6    9e-15


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

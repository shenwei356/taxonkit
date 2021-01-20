# Tutorial


## Table of Contents

<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->

- [Formating lineage](#formating-lineage)
- [Parsing kraken/bracken result](#parsing-krakenbraken-result)
- [Summaries of taxonomy data](#summaries-of-taxonomy-data)
- [Making nr blastdb for specific taxids](#making-nr-blastdb-for-specific-taxids)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

## Formating lineage

Example data.

    $ cat taxids3.txt
    376619
    349741
    239935
    314101
    11932
    1327037
    
Format to seven-level ranks ("superkingdom phylum class order family genus species").

    $ cat taxids3.txt \
        | taxonkit lineage \
        | taxonkit reformat \
        | cut -f 1,3
    376619  Bacteria;Proteobacteria;Gammaproteobacteria;Thiotrichales;Francisellaceae;Francisella;Francisella tularensis
    349741  Bacteria;Verrucomicrobia;Verrucomicrobiae;Verrucomicrobiales;Akkermansiaceae;Akkermansia;Akkermansia muciniphila
    239935  Bacteria;Verrucomicrobia;Verrucomicrobiae;Verrucomicrobiales;Akkermansiaceae;Akkermansia;Akkermansia muciniphila
    314101  Bacteria;;;;;;uncultured murine large bowel bacterium BAC 54B
    11932   Viruses;Artverviricota;Revtraviricetes;Ortervirales;Retroviridae;Intracisternal A-particles;Mouse Intracisternal A-particle
    1327037 Viruses;Uroviricota;Caudoviricetes;Caudovirales;Siphoviridae;;Croceibacter phage P2559Y
    
Fill missing ranks and add prefixes.

    $ cat taxids3.txt \
        | taxonkit lineage \
        | taxonkit reformat -F -P \
        | cut -f 1,3
    376619  k__Bacteria;p__Proteobacteria;c__Gammaproteobacteria;o__Thiotrichales;f__Francisellaceae;g__Francisella;s__Francisella tularensis
    349741  k__Bacteria;p__Verrucomicrobia;c__Verrucomicrobiae;o__Verrucomicrobiales;f__Akkermansiaceae;g__Akkermansia;s__Akkermansia muciniphila
    239935  k__Bacteria;p__Verrucomicrobia;c__Verrucomicrobiae;o__Verrucomicrobiales;f__Akkermansiaceae;g__Akkermansia;s__Akkermansia muciniphila
    314101  k__Bacteria;p__unclassified Bacteria phylum;c__unclassified Bacteria class;o__unclassified Bacteria order;f__unclassified Bacteria family;g__unclassified Bacteria genus;s__uncultured murine large bowel bacterium BAC 54B
    11932   k__Viruses;p__Artverviricota;c__Revtraviricetes;o__Ortervirales;f__Retroviridae;g__Intracisternal A-particles;s__Mouse Intracisternal A-particle
    1327037 k__Viruses;p__Uroviricota;c__Caudoviricetes;o__Caudovirales;f__Siphoviridae;g__unclassified Siphoviridae genus;s__Croceibacter phage P2559Y

Single prefix of a rank can be set with flag like `--prefix-k`.

    $ cat taxids3.txt \
        | taxonkit lineage \
        | taxonkit reformat -F -P --prefix-k "d__" \
        | cut -f 1,3
    376619  d__Bacteria;p__Proteobacteria;c__Gammaproteobacteria;o__Thiotrichales;f__Francisellaceae;g__Francisella;s__Francisella tularensis
    349741  d__Bacteria;p__Verrucomicrobia;c__Verrucomicrobiae;o__Verrucomicrobiales;f__Akkermansiaceae;g__Akkermansia;s__Akkermansia muciniphila
    239935  d__Bacteria;p__Verrucomicrobia;c__Verrucomicrobiae;o__Verrucomicrobiales;f__Akkermansiaceae;g__Akkermansia;s__Akkermansia muciniphila
    314101  d__Bacteria;p__unclassified Bacteria phylum;c__unclassified Bacteria class;o__unclassified Bacteria order;f__unclassified Bacteria family;g__unclassified Bacteria genus;s__uncultured murine large bowel bacterium BAC 54B
    11932   d__Viruses;p__Artverviricota;c__Revtraviricetes;o__Ortervirales;f__Retroviridae;g__Intracisternal A-particles;s__Mouse Intracisternal A-particle
    1327037 d__Viruses;p__Uroviricota;c__Caudoviricetes;o__Caudovirales;f__Siphoviridae;g__unclassified Siphoviridae genus;s__Croceibacter phage P2559Y

## Parsing kraken/bracken result

Example Data.

    $ head -n 10  test_bracken_species.kreport 
    100.00  11430773        0       R       1       root
    92.12   10530410        0       R1      131567    cellular organisms
    92.06   10522651        0       D       2           Bacteria
    92.04   10520821        0       P       1224          Proteobacteria
    92.02   10518133        0       C       1236            Gammaproteobacteria
    92.00   10516208        0       O       91347             Enterobacterales
    91.81   10494626        0       F       543                 Enterobacteriaceae
    86.73   9914368 0       G       561                   Escherichia
    84.84   9698064 9698064 S       562                     Escherichia coli
    0.93    106494  106494  S       208962                  Escherichia albertii

Save taxon proportion and taxid, and get lineage, name and rank.

    $ time cat test_bracken_species.kreport \
        | csvtk cut -Ht -f 1,5 \
        | taxonkit lineage -i 2 -n -r \
        | csvtk cut -Ht -f 1,2,5,4,3 \
        | head -n 10 
        
    100.00  1       no rank root    root
    92.12   131567  no rank cellular organisms      cellular organisms
    92.06   2       superkingdom    Bacteria        cellular organisms;Bacteria
    92.04   1224    phylum  Proteobacteria  cellular organisms;Bacteria;Proteobacteria
    92.02   1236    class   Gammaproteobacteria     cellular organisms;Bacteria;Proteobacteria;Gammaproteobacteria
    92.00   91347   order   Enterobacterales        cellular organisms;Bacteria;Proteobacteria;Gammaproteobacteria;Enterobacterales
    91.81   543     family  Enterobacteriaceae      cellular organisms;Bacteria;Proteobacteria;Gammaproteobacteria;Enterobacterales;Enterobacteriaceae
    86.73   561     genus   Escherichia     cellular organisms;Bacteria;Proteobacteria;Gammaproteobacteria;Enterobacterales;Enterobacteriaceae;Escherichia
    84.84   562     species Escherichia coli        cellular organisms;Bacteria;Proteobacteria;Gammaproteobacteria;Enterobacterales;Enterobacteriaceae;Escherichia;Escherichia coli
    0.93    208962  species Escherichia albertii    cellular organisms;Bacteria;Proteobacteria;Gammaproteobacteria;Enterobacterales;Enterobacteriaceae;Escherichia;Escherichia albertii

    real    0m1.426s
    user    0m2.993s
    sys     0m0.260s

Only save species or lower level and get lineage in format of "superkingdom phylum class order family genus species".

    $ time cat test_bracken_species.kreport \
        | csvtk cut -Ht -f 1,5 \
        | taxonkit filter -E species -L species -i 2 \
        | taxonkit lineage -i 2 -n -r \
        | taxonkit reformat -i 3 -f "{k};{p};{c};{o};{f};{g};{s}" \
        | csvtk cut -Ht -f 1,2,5,4,6 \
        | csvtk add-header -t -n abundance,taxid,rank,name,lineage \
        | head -n 10 \
        | csvtk pretty -t
        
    abundance   taxid     name      rank                     lineage
    84.84       562       species   Escherichia coli         Bacteria;Proteobacteria;Gammaproteobacteria;Enterobacterales;Enterobacteriaceae;Escherichia;Escherichia coli
    0.93        208962    species   Escherichia albertii     Bacteria;Proteobacteria;Gammaproteobacteria;Enterobacterales;Enterobacteriaceae;Escherichia;Escherichia albertii
    0.90        564       species   Escherichia fergusonii   Bacteria;Proteobacteria;Gammaproteobacteria;Enterobacterales;Enterobacteriaceae;Escherichia;Escherichia fergusonii
    0.05        2725997   species   Escherichia sp. SCLE84   Bacteria;Proteobacteria;Gammaproteobacteria;Enterobacterales;Enterobacteriaceae;Escherichia;Escherichia sp. SCLE84
    0.00        2044467   species   Escherichia sp. E4742    Bacteria;Proteobacteria;Gammaproteobacteria;Enterobacterales;Enterobacteriaceae;Escherichia;Escherichia sp. E4742
    0.01        1499973   species   Escherichia marmotae     Bacteria;Proteobacteria;Gammaproteobacteria;Enterobacterales;Enterobacteriaceae;Escherichia;Escherichia marmotae
    2.68        621       species   Shigella boydii          Bacteria;Proteobacteria;Gammaproteobacteria;Enterobacterales;Enterobacteriaceae;Shigella;Shigella boydii
    0.04        622       species   Shigella dysenteriae     Bacteria;Proteobacteria;Gammaproteobacteria;Enterobacterales;Enterobacteriaceae;Shigella;Shigella dysenteriae
    0.98        28901     species   Salmonella enterica      Bacteria;Proteobacteria;Gammaproteobacteria;Enterobacterales;Enterobacteriaceae;Salmonella;Salmonella enterica

    real    0m4.142s
    user    0m12.874s
    sys     0m0.808s

## Summaries of taxonomy data

You can change the taxID of interest.

1. Count of ranks

        $ time taxonkit list --ids 1 \
            | taxonkit lineage -L -r \
            | csvtk freq -H -t -f 2 -nr \
            | csvtk pretty -t
        
        species            1879659
        no rank            222743
        genus              96625
        strain             44483
        subspecies         25174
        family             9492
        varietas           8524
        subfamily          3050
        tribe              2213
        order              1660
        subgenus           1618
        isolate            1319
        serotype           1216
        clade              886
        superfamily        865
        forma specialis    741
        forma              564
        subtribe           508
        section            437
        class              429
        suborder           372
        species group      330
        phylum             272
        subclass           156
        serogroup          138
        infraorder         130
        species subgroup   124
        superorder         55
        subphylum          33
        parvorder          26
        subsection         21
        genotype           20
        infraclass         18
        biotype            17
        morph              12
        kingdom            11
        series             9
        superclass         6
        cohort             5
        pathogroup         5
        subvariety         5
        superkingdom       4
        subcohort          3
        subkingdom         1
        superphylum        1

        real    0m3.663s
        user    0m15.897s
        sys     0m1.010s


## Making nr blastdb for specific taxids

Attention:

- [BLAST+ 2.8.1 is released with new databases](https://ncbiinsights.ncbi.nlm.nih.gov/2019/01/04/blast-2-8-1-with-new-databases-and-better-performance/),
which allows you to limit your search by taxonomy using information built into the BLAST databases.
So you don't need to build blastdb for specific taxids now.

Changes:

- 2018-09-13 rewritten
- 2018-12-22 providing faster method for step 3.1
- 2019-01-07 add note of new blastdb version
- 2020-10-14 update steps for huge number of accessions belong to high taxon level like bacteria.

Data:

- [pre-formated blastdb](ftp://ftp.ncbi.nlm.nih.gov/blast/db) (09/10/2018)
- [prot.accession2taxid.gz](ftp://ftp.ncbi.nih.gov/pub/taxonomy/accession2taxid/prot.accession2taxid.gz) (09/07/2018) (optional, but recommended)

Hardware in this tutorial

- CPU: AMD 8-cores/16-threads 3.7Ghz
- RAM: 64GB
- DISK:
    - Taxonomy files stores in NVMe SSD
    - blastdb files stores in 7200rpm HDD

Tools:

- [blast+](ftp://ftp.ncbi.nlm.nih.gov/blast/executables/blast+/LATEST/)
- [pigz](https://zlib.net/pigz/) (recommended, faster than gzip)
- [taxonkit](https://bioinf.shenwei.me/taxonkit)
- [seqkit](https://bioinf.shenwei.me/seqkit) (recommended),  version >= 0.14.0
- [rush](https://github.com/shenwei356/rush) (optional, for parallizing filtering sequence)

Steps:

1. Listing all taxids below `$id` using taxonkit.

        id=6656

        # 6656 is the phylum Arthropoda
        # echo 6656 | taxonkit lineage | taxonkit reformat
        # 6656    cellular organisms;Eukaryota;Opisthokonta;Metazoa;Eumetazoa;Bilateria;Protostomia;Ecdysozoa;Panarthropoda;Arthropoda    Eukaryota;Arthropoda;;;;;
        
        # 2     bacteria
        # 2157  archaea
        # 4751  fungi
        # 10239 virus

        # time: 2s
        taxonkit list --ids $id --indent "" > $id.taxid.txt
        
        # taxonkit list --ids 2,4751,10239 --indent "" > microbe.taxid.txt

        wc -l $id.taxid.txt
        # 518373 6656.taxid.txt

2. Retrieving target accessions. There are two options:

    1. From prot.accession2taxid.gz (**faster, recommended**). *Note that some accessions are not in `nr`*.

            # time: 4min
            pigz -dc prot.accession2taxid.gz \
                | csvtk grep -t -f taxid -P $id.taxid.txt \
                | csvtk cut -t -f accession.version,taxid \
                | sed 1d \
                > $id.acc2taxid.txt
                
            cut -f 1 $id.acc2taxid.txt > $id.acc.txt

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

            # 1. exporting nr.fa from pre-formated blastdb

            # time: 117min (run only once)
            blastdbcmd -db nr -dbtype prot -entry all -outfmt "%f" -out - | pigz -c > nr.fa.gz

            # =====================================================================

            # 2. filtering sequence belong to $taxid

            # ---------------------------------------------------------------------

            # methond 1) (for cases where $id.acc.txt is not very huge)
            # time: 80min
            # perl one-liner is used to unfold records having mulitple accessions
            time cat <(echo) <(pigz -dc nr.fa.gz) \
                | perl -e 'BEGIN{ $/ = "\n>"; <>; } while(<>){s/>$//;  $i = index $_, "\n"; $h = substr $_, 0, $i; $s = substr $_, $i+1; if ($h !~ />/) { print ">$_"; next; }; $h = ">$h"; while($h =~ />([^ ]+ .+?) ?(?=>|$)/g){ $h1 = $1; $h1 =~ s/^\W+//; print ">$h1\n$s";} } ' \
                | seqkit grep -f $id.acc.txt -o nr.$id.fa.gz

            # ---------------------------------------------------------------------

            # method 2) (**faster**)

            # 33min (run only once)
            # (1). split nr.fa.gz. # Note: I have 16 cpus.
            $ time seqkit split2 -p 15 nr.fa.gz

            # (2). parallize unfolding
            $ cat _unfold_blastdb_fa.sh
            #!/bin/sh
            perl -e 'BEGIN{ $/ = "\n>"; <>; } while(<>){s/>$//;  $i = index $_, "\n"; $h = substr $_, 0, $i; $s = substr $_, $i+1; if ($h !~ />/) { print ">$_"; next; }; $h = ">$h"; while($h =~ />([^ ]+ .+?) ?(?=>|$)/g){ $h1 = $1; $h1 =~ s/^\W+//; print ">$h1\n$s";} } '

            # 10 min
            time ls nr.fa.gz.split/nr.part_*.fa.gz \
                | rush -j 15 -v id=$id 'cat <(echo) <(pigz -dc {}) \
                    | ./_unfold_blastdb_fa.sh \
                    | seqkit grep -f {id}.acc.txt -o nr.{id}.{%@nr\.(.+)$} '

            # (3). merge result
            cat nr.$id.part*.fa.gz > nr.$id.fa.gz
            rm nr.$id.part*.fa.gz
            
            # ---------------------------------------------------------------------
            
            # method 3) (for huge $id.acc.txt file, e.g., bacteria)
            
            # (1). split ${id}.acc.txt into several parts. chunk size depends on lines and RAM (64G for me).
            split -d -l 300000000 $id.acc.txt $id.acc.txt.part_
            
            # (2). filter
            time ls $id.acc.txt.part_* \
                | rush -j 1 --immediate-output -v id=$id \
                    'echo {}; cat <(echo) <(pigz -dc nr.fa.gz ) \
                    | ./_unfold_blastdb_fa.sh \
                    | seqkit grep -f {} -o nr.{id}.{%@(part_.+)}.fa.gz '
 
            # (3). merge
            cat nr.$id.part*.fa.gz > nr.$id.fa.gz     
            
            # clean
            rm nr.$id.part*.fa.gz
            rm $id.acc.txt.part_
            
            
            # (4). optionally adding taxid, you may edit replacement (-r) below
            # split
            time split -d -l 200000000 $id.acc2taxid.txt $id.acc2taxid.txt.part_
            
            ln -s nr.$id.fa.gz nr.$id.with-taxid.part0.fa.gz         
            i=0
            for f in $id.acc2taxid.txt.part_* ; do
                echo $f
                time pigz -cd nr.$id.with-taxid.part$i.fa.gz \
                    | seqkit replace -k $f -p "^([^\-]+?) " -r "{kv}-\$1 " -K -U -o nr.$id.with-taxid.part$(($i+1)).fa.gz;
                /bin/rm nr.$id.with-taxid.part$i.fa.gz
                i=$(($i+1));
            done
            mv nr.$id.with-taxid.part$i.fa.gz nr.$id.with-taxid.fa.gz
            
            
            # =====================================================================

            # 3. counting sequences
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

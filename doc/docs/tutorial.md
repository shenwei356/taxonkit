# Tutorial


## Table of Contents

<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->

- [Formatting lineage](#formatting-lineage)
- [Parsing kraken/bracken result](#parsing-krakenbraken-result)
- [Making nr blastdb for specific taxids](#making-nr-blastdb-for-specific-taxids)
- [Summaries of taxonomy data](#summaries-of-taxonomy-data)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->


## Formatting lineage

Show lineage detail of a taxID.
The command below works on Windows with help of [csvtk](http://bioinf.shenwei.me/csvtk).

    $ echo "2697049" \
        | taxonkit lineage -t \
        | csvtk cut -Ht -f 3 \
        | csvtk unfold -Ht -f 1 -s ";" \
        | taxonkit lineage -r -n -L \
        | csvtk cut -Ht -f 1,3,2 \
        | csvtk pretty -Ht 
    
    10239     superkingdom   Viruses
    2559587   clade          Riboviria
    2732396   kingdom        Orthornavirae
    2732408   phylum         Pisuviricota
    2732506   class          Pisoniviricetes
    76804     order          Nidovirales
    2499399   suborder       Cornidovirineae
    11118     family         Coronaviridae
    2501931   subfamily      Orthocoronavirinae
    694002    genus          Betacoronavirus
    2509511   subgenus       Sarbecovirus
    694009    species        Severe acute respiratory syndrome-related coronavirus
    2697049   no rank        Severe acute respiratory syndrome coronavirus 2

Example data.

    $ cat taxids3.txt
    376619
    349741
    239935
    314101
    11932
    1327037
    83333
    1408252
    2605619
    2697049

    
Format to 7-level ranks ("superkingdom phylum class order family genus species").

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
    83333   Bacteria;Proteobacteria;Gammaproteobacteria;Enterobacterales;Enterobacteriaceae;Escherichia;Escherichia coli
    1408252 Bacteria;Proteobacteria;Gammaproteobacteria;Enterobacterales;Enterobacteriaceae;Escherichia;Escherichia coli
    2605619 Bacteria;Proteobacteria;Gammaproteobacteria;Enterobacterales;Enterobacteriaceae;Escherichia;Escherichia coli
    2697049 Viruses;Pisuviricota;Pisoniviricetes;Nidovirales;Coronaviridae;Betacoronavirus;Severe acute respiratory syndrome-related coronavirus

Format to 8-level ranks ("superkingdom phylum class order family genus species subspecies/rank").

    $ cat taxids3.txt \
        | taxonkit lineage \
        | taxonkit reformat -f "{k};{p};{c};{o};{f};{g};{s};{t}" \
        | cut -f 1,3
    
    376619  Bacteria;Proteobacteria;Gammaproteobacteria;Thiotrichales;Francisellaceae;Francisella;Francisella tularensis;Francisella tularensis subsp. holarctica LVS
    349741  Bacteria;Verrucomicrobia;Verrucomicrobiae;Verrucomicrobiales;Akkermansiaceae;Akkermansia;Akkermansia muciniphila;Akkermansia muciniphila ATCC BAA-835
    239935  Bacteria;Verrucomicrobia;Verrucomicrobiae;Verrucomicrobiales;Akkermansiaceae;Akkermansia;Akkermansia muciniphila;
    314101  Bacteria;;;;;;uncultured murine large bowel bacterium BAC 54B;
    11932   Viruses;Artverviricota;Revtraviricetes;Ortervirales;Retroviridae;Intracisternal A-particles;Mouse Intracisternal A-particle;
    1327037 Viruses;Uroviricota;Caudoviricetes;Caudovirales;Siphoviridae;;Croceibacter phage P2559Y;
    83333   Bacteria;Proteobacteria;Gammaproteobacteria;Enterobacterales;Enterobacteriaceae;Escherichia;Escherichia coli;Escherichia coli K-12
    1408252 Bacteria;Proteobacteria;Gammaproteobacteria;Enterobacterales;Enterobacteriaceae;Escherichia;Escherichia coli;Escherichia coli R178
    2605619 Bacteria;Proteobacteria;Gammaproteobacteria;Enterobacterales;Enterobacteriaceae;Escherichia;Escherichia coli;
    2697049 Viruses;Pisuviricota;Pisoniviricetes;Nidovirales;Coronaviridae;Betacoronavirus;Severe acute respiratory syndrome-related coronavirus;

Replace missing ranks with `Unassigned` and output tab-delimited format.

    $ cat taxids3.txt \
        | taxonkit lineage \
        | taxonkit reformat -r "Unassigned" -f "{k}\t{p}\t{c}\t{o}\t{f}\t{g}\t{s}\t{t}" \
        | cut -f 1,3-10 \
        | csvtk pretty -H -t 
        
    376619    Bacteria   Proteobacteria    Gammaproteobacteria   Thiotrichales        Francisellaceae      Francisella                  Francisella tularensis                                  Francisella tularensis subsp. holarctica LVS
    349741    Bacteria   Verrucomicrobia   Verrucomicrobiae      Verrucomicrobiales   Akkermansiaceae      Akkermansia                  Akkermansia muciniphila                                 Akkermansia muciniphila ATCC BAA-835
    239935    Bacteria   Verrucomicrobia   Verrucomicrobiae      Verrucomicrobiales   Akkermansiaceae      Akkermansia                  Akkermansia muciniphila                                 Unassigned
    314101    Bacteria   Unassigned        Unassigned            Unassigned           Unassigned           Unassigned                   uncultured murine large bowel bacterium BAC 54B         Unassigned
    11932     Viruses    Artverviricota    Revtraviricetes       Ortervirales         Retroviridae         Intracisternal A-particles   Mouse Intracisternal A-particle                         Unassigned
    1327037   Viruses    Uroviricota       Caudoviricetes        Caudovirales         Siphoviridae         Unassigned                   Croceibacter phage P2559Y                               Unassigned
    83333     Bacteria   Proteobacteria    Gammaproteobacteria   Enterobacterales     Enterobacteriaceae   Escherichia                  Escherichia coli                                        Escherichia coli K-12
    1408252   Bacteria   Proteobacteria    Gammaproteobacteria   Enterobacterales     Enterobacteriaceae   Escherichia                  Escherichia coli                                        Escherichia coli R178
    2605619   Bacteria   Proteobacteria    Gammaproteobacteria   Enterobacterales     Enterobacteriaceae   Escherichia                  Escherichia coli                                        Unassigned
    2697049   Viruses    Pisuviricota      Pisoniviricetes       Nidovirales          Coronaviridae        Betacoronavirus              Severe acute respiratory syndrome-related coronavirus   Unassigned

Fill missing ranks and add prefixes.

    $ cat taxids3.txt \
        | taxonkit lineage \
        | taxonkit reformat -F -P -f "{k}\t{p}\t{c}\t{o}\t{f}\t{g}\t{s}\t{t}" \
        | cut -f 1,3-10 \
        | csvtk pretty -H -t 
        
    376619    k__Bacteria   p__Proteobacteria                 c__Gammaproteobacteria           o__Thiotrichales                 f__Francisellaceae                g__Francisella                       s__Francisella tularensis                                  t__Francisella tularensis subsp. holarctica LVS
    349741    k__Bacteria   p__Verrucomicrobia                c__Verrucomicrobiae              o__Verrucomicrobiales            f__Akkermansiaceae                g__Akkermansia                       s__Akkermansia muciniphila                                 t__Akkermansia muciniphila ATCC BAA-835
    239935    k__Bacteria   p__Verrucomicrobia                c__Verrucomicrobiae              o__Verrucomicrobiales            f__Akkermansiaceae                g__Akkermansia                       s__Akkermansia muciniphila                                 t__unclassified Akkermansia muciniphila subspecies/strain
    314101    k__Bacteria   p__unclassified Bacteria phylum   c__unclassified Bacteria class   o__unclassified Bacteria order   f__unclassified Bacteria family   g__unclassified Bacteria genus       s__uncultured murine large bowel bacterium BAC 54B         t__unclassified uncultured murine large bowel bacterium BAC 54B subspecies/strain
    11932     k__Viruses    p__Artverviricota                 c__Revtraviricetes               o__Ortervirales                  f__Retroviridae                   g__Intracisternal A-particles        s__Mouse Intracisternal A-particle                         t__unclassified Mouse Intracisternal A-particle subspecies/strain
    1327037   k__Viruses    p__Uroviricota                    c__Caudoviricetes                o__Caudovirales                  f__Siphoviridae                   g__unclassified Siphoviridae genus   s__Croceibacter phage P2559Y                               t__unclassified Croceibacter phage P2559Y subspecies/strain
    83333     k__Bacteria   p__Proteobacteria                 c__Gammaproteobacteria           o__Enterobacterales              f__Enterobacteriaceae             g__Escherichia                       s__Escherichia coli                                        t__Escherichia coli K-12
    1408252   k__Bacteria   p__Proteobacteria                 c__Gammaproteobacteria           o__Enterobacterales              f__Enterobacteriaceae             g__Escherichia                       s__Escherichia coli                                        t__Escherichia coli R178
    2605619   k__Bacteria   p__Proteobacteria                 c__Gammaproteobacteria           o__Enterobacterales              f__Enterobacteriaceae             g__Escherichia                       s__Escherichia coli                                        t__unclassified Escherichia coli subspecies/strain
    2697049   k__Viruses    p__Pisuviricota                   c__Pisoniviricetes               o__Nidovirales                   f__Coronaviridae                  g__Betacoronavirus                   s__Severe acute respiratory syndrome-related coronavirus   t__unclassified Severe acute respiratory syndrome-related coronavirus subspecies/strain

**When these's no nodes of rank "subspecies" nor "stain",
we can switch `-S/--pseudo-strain` to use the node with lowest rank
as subspecies/strain name, if which rank is lower than "species"**.

    $ cat taxids3.txt \
        | taxonkit lineage -r \
        | taxonkit reformat -F -S -f "{k}\t{p}\t{c}\t{o}\t{f}\t{g}\t{s}\t{t}" \
        | cut -f 1,3,10,11 \
        | csvtk add-header -t -n "taxid,rank,species,strain" \
        | csvtk pretty -t
        
    taxid     rank         species                                                 strain
    -------   ----------   -----------------------------------------------------   ------------------------------------------------------------------------------
    376619    strain       Francisella tularensis                                  Francisella tularensis subsp. holarctica LVS
    349741    strain       Akkermansia muciniphila                                 Akkermansia muciniphila ATCC BAA-835
    239935    species      Akkermansia muciniphila                                 unclassified Akkermansia muciniphila subspecies/strain
    314101    species      uncultured murine large bowel bacterium BAC 54B         unclassified uncultured murine large bowel bacterium BAC 54B subspecies/strain
    11932     species      Mouse Intracisternal A-particle                         unclassified Mouse Intracisternal A-particle subspecies/strain
    1327037   species      Croceibacter phage P2559Y                               unclassified Croceibacter phage P2559Y subspecies/strain
    83333     strain       Escherichia coli                                        Escherichia coli K-12
    1408252   subspecies   Escherichia coli                                        Escherichia coli R178
    2605619   no rank      Escherichia coli                                        Escherichia coli O16:H48
    2697049   no rank      Severe acute respiratory syndrome-related coronavirus   Severe acute respiratory syndrome coronavirus 2

List eight-level lineage for all taxIDs of rank lower than or equal to species, including some nodes with "no rank".
But when filtering with `-L/--lower-than`, you can use
`-n/--save-predictable-norank` to **save some special ranks without order,
where rank of the closest higher node is still lower than rank cutoff**.

    $ time taxonkit list --ids 1 \
        | taxonkit filter -L species -E species -R -N -n \
        | taxonkit lineage -n -r \
        | taxonkit reformat -F -S -f "{k}\t{p}\t{c}\t{o}\t{f}\t{g}\t{s}\t{t}" \
        | csvtk cut -Ht -l -f 1,4,3,2,5-12 \
        | csvtk add-header -t -n "taxid,rank,name,lineage,kingdom,phylum,class,order,family,genus,species,strain" \
        | pigz -c > result.tsv.gz

    real    0m25.167s
    user    2m14.809s
    sys     0m7.197s
    
    $ pigz -cd result.tsv.gz \
        | csvtk grep -t -f taxid -p 2697049 \
        | csvtk transpose -t \
        | csvtk pretty -H -t 
        
    taxid     2697049
    rank      no rank
    name      Severe acute respiratory syndrome coronavirus 2
    lineage   Viruses;Riboviria;Orthornavirae;Pisuviricota;Pisoniviricetes;Nidovirales;Cornidovirineae;Coronaviridae;Orthocoronavirinae;Betacoronavirus;Sarbecovirus;Severe acute respiratory syndrome-related coronavirus;Severe acute respiratory syndrome coronavirus 2
    kingdom   Viruses
    phylum    Pisuviricota
    class     Pisoniviricetes
    order     Nidovirales
    family    Coronaviridae
    genus     Betacoronavirus
    species   Severe acute respiratory syndrome-related coronavirus
    strain    Severe acute respiratory syndrome coronavirus 2
    
## Parsing kraken/bracken result

Example Data

- [SRS014459-Stool.fasta.gz](https://github.com/biobakery/biobakery/raw/master/demos/biobakery_demos/data/metaphlan2/input/SRS014459-Stool.fasta.gz)


Run Kraken2 and Bracken

    KRAKEN_DB=/home/shenwei/ws/db/kraken/k2_pluspf
    THREADS=16

    CLASSIFICATION_LVL=S
    THRESHOLD=10

    READ_LEN=100
    SAMPLE=SRS014459-Stool.fasta.gz

    BRACKEN_OUTPUT_FILE=$SAMPLE

    kraken2 --db ${KRAKEN_DB} --threads ${THREADS} -report ${SAMPLE}.kreport $SAMPLE > ${SAMPLE}.kraken

    est_abundance.py -i ${SAMPLE}.kreport -k ${KRAKEN_DB}/database${READ_LEN}mers.kmer_distrib \
        -l ${CLASSIFICATION_LVL} -t ${THRESHOLD} -o ${BRACKEN_OUTPUT_FILE}.bracken

Orignial format

    $ head -n 15 SRS014459-Stool.fasta.gz_bracken_species.kreport
    100.00  9491    0       R       1       root
    99.85   9477    0       R1      131567    cellular organisms
    99.85   9477    0       D       2           Bacteria
    66.08   6271    0       D1      1783270       FCB group
    66.08   6271    0       D2      68336           Bacteroidetes/Chlorobi group
    66.08   6271    0       P       976               Bacteroidetes
    66.08   6271    0       C       200643              Bacteroidia
    66.08   6271    0       O       171549                Bacteroidales
    34.45   3270    0       F       815                     Bacteroidaceae
    34.45   3270    0       G       816                       Bacteroides
    10.43   990     990     S       246787                      Bacteroides cellulosilyticus
    7.98    757     757     S       28116                       Bacteroides ovatus
    3.10    293     0       G1      2646097                     unclassified Bacteroides
    1.06    100     100     S       2755405                       Bacteroides sp. CACC 737
    0.49    46      46      S       2650157                       Bacteroides sp. HF-5287
    
Converting to MetaPhlAn2 format. (Similar to [kreport2mpa.py](https://github.com/jenniferlu717/KrakenTools/blob/master/kreport2mpa.py))

    $ cat SRS014459-Stool.fasta.gz_bracken_species.kreport \
        | csvtk cut -Ht -f 5,1 \
        | taxonkit lineage \
        | taxonkit reformat -i 3 -P -f "{k}|{p}|{c}|{o}|{f}|{g}|{s}" \
        | csvtk cut -Ht -f 4,2 \
        | csvtk replace -Ht -p "(\|[kpcofgs]__)+$" \
        | csvtk replace -Ht -p "\|[kpcofgs]__\|" -r "|" \
        | csvtk uniq -Ht \
        | csvtk grep -Ht -p k__ -v \
        > SRS014459-Stool.fasta.gz_bracken_species.kreport.format
        
    $ head -n 10 SRS014459-Stool.fasta.gz_bracken_species.kreport.format
    
    k__Bacteria     99.85
    k__Bacteria|p__Bacteroidetes    66.08
    k__Bacteria|p__Bacteroidetes|c__Bacteroidia     66.08
    k__Bacteria|p__Bacteroidetes|c__Bacteroidia|o__Bacteroidales    66.08
    k__Bacteria|p__Bacteroidetes|c__Bacteroidia|o__Bacteroidales|f__Bacteroidaceae  34.45
    k__Bacteria|p__Bacteroidetes|c__Bacteroidia|o__Bacteroidales|f__Bacteroidaceae|g__Bacteroides   34.45
    k__Bacteria|p__Bacteroidetes|c__Bacteroidia|o__Bacteroidales|f__Bacteroidaceae|g__Bacteroides|s__Bacteroides cellulosilyticus   10.43
    k__Bacteria|p__Bacteroidetes|c__Bacteroidia|o__Bacteroidales|f__Bacteroidaceae|g__Bacteroides|s__Bacteroides ovatus     7.98
    k__Bacteria|p__Bacteroidetes|c__Bacteroidia|o__Bacteroidales|f__Bacteroidaceae|g__Bacteroides|s__Bacteroides sp. CACC 737       1.06
    k__Bacteria|p__Bacteroidetes|c__Bacteroidia|o__Bacteroidales|f__Bacteroidaceae|g__Bacteroides|s__Bacteroides sp. HF-5287        0.49

Converting to Qiime format

    $ cat SRS014459-Stool.fasta.gz_bracken_species.kreport \
        | csvtk cut -Ht -f 5,1 \
        | taxonkit lineage \
        | taxonkit reformat -i 3 -P -f "{k}; {p}; {c}; {o}; {f}; {g}; {s}" \
        | csvtk cut -Ht -f 4,2 \
        | csvtk replace -Ht -p "(; [kpcofgs]__)+$" \
        | csvtk replace -Ht -p "; [kpcofgs]__; " -r "; " \
        | csvtk uniq -Ht \
        | csvtk grep -Ht -p k__ -v \
        | head -n 10 
        
    k__Bacteria     99.85
    k__Bacteria; p__Bacteroidetes   66.08
    k__Bacteria; p__Bacteroidetes; c__Bacteroidia   66.08
    k__Bacteria; p__Bacteroidetes; c__Bacteroidia; o__Bacteroidales 66.08
    k__Bacteria; p__Bacteroidetes; c__Bacteroidia; o__Bacteroidales; f__Bacteroidaceae      34.45
    k__Bacteria; p__Bacteroidetes; c__Bacteroidia; o__Bacteroidales; f__Bacteroidaceae; g__Bacteroides      34.45
    k__Bacteria; p__Bacteroidetes; c__Bacteroidia; o__Bacteroidales; f__Bacteroidaceae; g__Bacteroides; s__Bacteroides cellulosilyticus     10.43
    k__Bacteria; p__Bacteroidetes; c__Bacteroidia; o__Bacteroidales; f__Bacteroidaceae; g__Bacteroides; s__Bacteroides ovatus       7.98
    k__Bacteria; p__Bacteroidetes; c__Bacteroidia; o__Bacteroidales; f__Bacteroidaceae; g__Bacteroides; s__Bacteroides sp. CACC 737 1.06
    k__Bacteria; p__Bacteroidetes; c__Bacteroidia; o__Bacteroidales; f__Bacteroidaceae; g__Bacteroides; s__Bacteroides sp. HF-5287  0.49

Save taxon proportion and taxid, and get lineage, name and rank.

    $ cat SRS014459-Stool.fasta.gz_bracken_species.kreport \
        | csvtk cut -Ht -f 1,5 \
        | taxonkit lineage -i 2 -n -r \
        | csvtk cut -Ht -f 1,2,5,4,3 \
        | head -n 10 \
        | csvtk pretty -Ht
        
    100.00   1         no rank        root                           root
    99.85    131567    no rank        cellular organisms             cellular organisms
    99.85    2         superkingdom   Bacteria                       cellular organisms;Bacteria
    66.08    1783270   clade          FCB group                      cellular organisms;Bacteria;FCB group
    66.08    68336     clade          Bacteroidetes/Chlorobi group   cellular organisms;Bacteria;FCB group;Bacteroidetes/Chlorobi group
    66.08    976       phylum         Bacteroidetes                  cellular organisms;Bacteria;FCB group;Bacteroidetes/Chlorobi group;Bacteroidetes
    66.08    200643    class          Bacteroidia                    cellular organisms;Bacteria;FCB group;Bacteroidetes/Chlorobi group;Bacteroidetes;Bacteroidia
    66.08    171549    order          Bacteroidales                  cellular organisms;Bacteria;FCB group;Bacteroidetes/Chlorobi group;Bacteroidetes;Bacteroidia;Bacteroidales
    34.45    815       family         Bacteroidaceae                 cellular organisms;Bacteria;FCB group;Bacteroidetes/Chlorobi group;Bacteroidetes;Bacteroidia;Bacteroidales;Bacteroidaceae
    34.45    816       genus          Bacteroides                    cellular organisms;Bacteria;FCB group;Bacteroidetes/Chlorobi group;Bacteroidetes;Bacteroidia;Bacteroidales;Bacteroidaceae;Bacteroides


Only save species or lower level and get lineage in format of "superkingdom phylum class order family genus species".

    $ cat SRS014459-Stool.fasta.gz_bracken_species.kreport \
        | csvtk cut -Ht -f 1,5 \
        | taxonkit filter -N -E species -L species -i 2 \
        | taxonkit lineage -i 2 -n -r \
        | taxonkit reformat -i 3 -f "{k};{p};{c};{o};{f};{g};{s}" \
        | csvtk cut -Ht -f 1,2,5,4,6 \
        | csvtk add-header -t -n abundance,taxid,rank,name,lineage \
        | head -n 10 \
        | csvtk pretty -t
        
    abundance   taxid     rank      name                           lineage
    ---------   -------   -------   ----------------------------   --------------------------------------------------------------------------------------------------------
    10.43       246787    species   Bacteroides cellulosilyticus   Bacteria;Bacteroidetes;Bacteroidia;Bacteroidales;Bacteroidaceae;Bacteroides;Bacteroides cellulosilyticus
    7.98        28116     species   Bacteroides ovatus             Bacteria;Bacteroidetes;Bacteroidia;Bacteroidales;Bacteroidaceae;Bacteroides;Bacteroides ovatus
    1.06        2755405   species   Bacteroides sp. CACC 737       Bacteria;Bacteroidetes;Bacteroidia;Bacteroidales;Bacteroidaceae;Bacteroides;Bacteroides sp. CACC 737
    0.49        2650157   species   Bacteroides sp. HF-5287        Bacteria;Bacteroidetes;Bacteroidia;Bacteroidales;Bacteroidaceae;Bacteroides;Bacteroides sp. HF-5287
    0.99        2528203   species   Bacteroides sp. A1C1           Bacteria;Bacteroidetes;Bacteroidia;Bacteroidales;Bacteroidaceae;Bacteroides;Bacteroides sp. A1C1
    0.28        2763022   species   Bacteroides sp. M10            Bacteria;Bacteroidetes;Bacteroidia;Bacteroidales;Bacteroidaceae;Bacteroides;Bacteroides sp. M10
    0.16        2650158   species   Bacteroides sp. HF-5141        Bacteria;Bacteroidetes;Bacteroidia;Bacteroidales;Bacteroidaceae;Bacteroides;Bacteroides sp. HF-5141
    0.12        2715212   species   Bacteroides sp. CBA7301        Bacteria;Bacteroidetes;Bacteroidia;Bacteroidales;Bacteroidaceae;Bacteroides;Bacteroides sp. CBA7301
    5.10        817       species   Bacteroides fragilis           Bacteria;Bacteroidetes;Bacteroidia;Bacteroidales;Bacteroidaceae;Bacteroides;Bacteroides fragilis

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

## Summaries of taxonomy data

You can change the taxID of interest.

1. Rank counts of common categories.

        $ echo Archaea Bacteria Eukaryota Fungi Metazoa Viridiplantae \
            | rush -D ' ' -T b \
                'taxonkit list --ids $(echo {} | taxonkit name2taxid | cut -f 2) \
                    | sed 1d \
                    | taxonkit filter -i 2 -E genus -L genus \
                    | taxonkit lineage -L -r \
                    | csvtk freq -H -t -f 2 -nr \
                    > stats.{}.tsv '
                    
        $ csvtk -t join --outer-join stats.*.tsv \
            | csvtk add-header -t -n "rank,$(ls stats.*.tsv | rush -k 'echo {@stats.(.+).tsv}' | paste -sd, )" \
            | csvtk csv2md -t
    
    [Similar data on NCBI Taxonomy](https://www.ncbi.nlm.nih.gov/Taxonomy/taxonomyhome.html/index.cgi?chapter=statistics&uncultured=hide&unspecified=hide)
    
    rank            |Archaea|Bacteria|Eukaryota|Fungi |Metazoa|Viridiplantae
    :---------------|:------|:-------|:--------|:-----|:------|:------------
    species         |12482  |460940  |1349648  |156908|957297 |191026
    strain          |354    |40643   |3486     |2352  |33     |50
    genus           |205    |4112    |90882    |6844  |64148  |16202
    isolate         |7      |503     |809      |76    |17     |3
    species group   |2      |77      |251      |22    |214    |5
    serotype        |       |218     |         |      |       |
    serogroup       |       |136     |         |      |       |
    subsection      |       |        |21       |      |       |21
    subspecies      |       |632     |24523    |158   |17043  |7212
    forma specialis |       |521     |220      |179   |33     |1
    species subgroup|       |23      |101      |      |101    |
    biotype         |       |7       |10       |      |       |
    morph           |       |        |12       |3     |4      |5
    section         |       |        |437      |37    |2      |398
    genotype        |       |        |12       |      |       |12
    series          |       |        |9        |      |5      |4
    varietas        |       |25      |8499     |1100  |2      |7188
    forma           |       |4       |560      |185   |6      |315
    subgenus        |       |1       |1558     |10    |1414   |112
    pathogroup      |       |5       |         |      |       |
    subvariety      |       |        |5        |      |       |5
        
1. Count of all ranks

        $ time taxonkit list --ids 1 \
            | taxonkit lineage -L -r \
            | csvtk freq -H -t -f 2 -nr \
            | csvtk pretty -H -t
        
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

1. Ranks of taxa at or below species.

        $ taxonkit list --ids 1 \
            | taxonkit filter --lower-than species --equal-to species \
            | taxonkit lineage -L -r  \
            | csvtk freq -Ht -nr -f 2 \
            | csvtk add-header -t -n rank,count \
            | csvtk pretty -t

        rank              count
        ---------------   -------
        species           1880044
        no rank           222756
        strain            44483
        subspecies        25171
        varietas          8524
        isolate           1319
        serotype          1216
        clade             885
        forma specialis   741
        forma             564
        serogroup         138
        genotype          20
        biotype           17
        morph             12
        pathogroup        5
        subvariety        5

        
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

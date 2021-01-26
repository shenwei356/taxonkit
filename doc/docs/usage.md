# Usage and Examples

Table of Contents

<!-- TOC -->
- [Usage and Examples](#usage-and-examples)
  - [Before use](#before-use)
  - [taxonkit](#taxonkit)
  - [list](#list)
  - [lineage](#lineage)
  - [reformat](#reformat)
  - [name2taxid](#name2taxid)
  - [filter](#filter)
  - [lca](#lca)
  - [taxid-changelog](#taxid-changelog)
  - [genautocomplete](#genautocomplete)
<!-- /TOC -->

## Before use

1. Download and decompress `taxdump.tar.gz`: ftp://ftp.ncbi.nih.gov/pub/taxonomy/taxdump.tar.gz 
2. Copy `names.dmp`, `nodes.dmp`, `delnodes.dmp` and `merged.dmp` to data directory: `$HOME/.taxonkit`,
e.g., `/home/shenwei/.taxonkit` ,
3. Optionally copy to some other directories, and later you can refer to using flag `--data-dir`,
or environment variable `TAXONKIT_DB`.

All-in-one command:

    wget -c ftp://ftp.ncbi.nih.gov/pub/taxonomy/taxdump.tar.gz 
    tar -zxvf taxdump.tar.gz
    
    mkdir -p $HOME/.taxonkit
    cp names.dmp nodes.dmp delnodes.dmp merged.dmp $HOME/.taxonkit

## taxonkit

```text
TaxonKit - A Cross-platform and Efficient NCBI Taxonomy Toolkit

Version: 0.7.2

Author: Wei Shen <shenwei356@gmail.com>

Source code: https://github.com/shenwei356/taxonkit
Documents  : https://bioinf.shenwei.me/taxonkit
Citation   : https://www.biorxiv.org/content/early/2019/01/08/513523

Dataset:

    Please download and decompress "taxdump.tar.gz":
    ftp://ftp.ncbi.nih.gov/pub/taxonomy/taxdump.tar.gz

    and copy "names.dmp", "nodes.dmp", "delnodes.dmp" and "merged.dmp" to data directory:
    "/home/shenwei/.taxonkit"

    or some other directory, and later you can refer to using flag --data-dir,
    or environment variable TAXONKIT_DB

Usage:
  taxonkit [command]

Available Commands:
  filter          Filter taxIDs by taxonomic rank range
  genautocomplete generate shell autocompletion script
  help            Help about any command
  lca             Compute lowest common ancestor (LCA) for taxIDs
  lineage         Query taxonomic lineage of given taxIDs
  list            List taxonomic subtrees of given taxIDs
  name2taxid      Convert scientific names to taxIDs
  reformat        Reformat lineage in canonical ranks
  taxid-changelog Create taxID changelog from dump archives
  version         print version information and check for update

Flags:
      --data-dir string   directory containing nodes.dmp and names.dmp (default "/home/shenwei/.taxonkit")
  -h, --help              help for taxonkit
      --line-buffered     use line buffering on output, i.e., immediately writing to stdin/file for every line of output
  -o, --out-file string   out file ("-" for stdout, suffix .gz for gzipped out) (default "-")
  -j, --threads int       number of CPUs. 4 is enough (default 4)
      --verbose           print verbose information

```

## list

Usage

```text
List taxonomic subtrees of given taxIDs

Examples:

    $ taxonkit list --ids 9606 -n -r --indent "    "
    9606 [species] Homo sapiens
        63221 [subspecies] Homo sapiens neanderthalensis
        741158 [subspecies] Homo sapiens subsp. 'Denisova'

    $ taxonkit list --ids 9606 --indent ""
    9606
    63221
    741158

Usage:
  taxonkit list [flags]

Flags:
  -h, --help            help for list
  -i, --ids string      taxID(s), multiple values should be separated by comma
  -I, --indent string   indent (default "  ")
  -J, --json            output in JSON format. you can save the result in file with suffix ".json" and open with modern text editor
  -n, --show-name       output scientific name
  -r, --show-rank       output rank

```

Examples

1. Default usage.

        
        $ taxonkit list --ids 9605,239934
        9605
        9606
            63221
            741158
        1425170
        2665952
            2665953

        239934
        239935
            349741
        512293
            512294
            1131822
            1262691
            1263034
        1679444
        2608915
            1131336
        ...


1. Removing indent. The list could be used to extract sequences from BLAST database with `blastdbcmd` (see [tutorial](http://bioinf.shenwei.me/taxonkit/tutorial/))

        $ taxonkit list --ids 9605,239934 --indent ""
        9605
        9606
        63221
        741158
        1425170
        2665952
        2665953

        239934
        239935
        349741
        512293
        512294
        1131822
        1262691
        1263034
        1679444
        ...


    **Performance:** Time and memory usage for whole taxon tree:

        $ # emptying the buffers cache
        $ su -c "free && sync && echo 3 > /proc/sys/vm/drop_caches && free"

        $ memusg -t taxonkit list --ids 1 --indent "" --verbose > t0.txt
        21:05:01.782 [INFO] parsing merged file: /home/shenwei/.taxonkit/names.dmp
        21:05:01.782 [INFO] parsing names file: /home/shenwei/.taxonkit/names.dmp
        21:05:01.782 [INFO] parsing delnodes file: /home/shenwei/.taxonkit/names.dmp
        21:05:01.816 [INFO] 61023 merged nodes parsed
        21:05:01.889 [INFO] 437929 delnodes parsed
        21:05:03.178 [INFO] 2303979 names parsed

        elapsed time: 3.290s
        peak rss: 742.77 MB

1. Adding names

        $ taxonkit list --show-rank --show-name --indent "    " --ids 9605,239934
        9605 [genus] Homo
            9606 [species] Homo sapiens
                63221 [subspecies] Homo sapiens neanderthalensis
                741158 [subspecies] Homo sapiens subsp. 'Denisova'
            1425170 [species] Homo heidelbergensis
            2665952 [no rank] environmental samples
                2665953 [species] Homo sapiens environmental sample

        239934 [genus] Akkermansia
            239935 [species] Akkermansia muciniphila
                349741 [strain] Akkermansia muciniphila ATCC BAA-835
            512293 [no rank] environmental samples
                512294 [species] uncultured Akkermansia sp.
                1131822 [species] uncultured Akkermansia sp. SMG25
                1262691 [species] Akkermansia sp. CAG:344
                1263034 [species] Akkermansia muciniphila CAG:154
            1679444 [species] Akkermansia glycaniphila
            2608915 [no rank] unclassified Akkermansia
                1131336 [species] Akkermansia sp. KLE1605
                1574264 [species] Akkermansia sp. KLE1797
        ...

    **Performance:** Time and memory usage for whole taxonomy tree:

        $ # emptying the buffers cache
        $ su -c "free && sync && echo 3 > /proc/sys/vm/drop_caches && free"

        $ memusg -t taxonkit list --show-rank --show-name --ids 1 > t1.txt
        elapsed time: 5.341s
        peak rss: 1.04 GB

1. Output in JSON format, you can easily collapse and uncollapse taxonomy tree in modern text editor.

        $ taxonkit list --show-rank --show-name --indent "    " --ids 9605,239934 --json
        {
            "9605 [genus] Homo": {
                "9606 [species] Homo sapiens": {
                    "63221 [subspecies] Homo sapiens neanderthalensis": {
                    },
                    "741158 [subspecies] Homo sapiens subsp. 'Denisova'": {
                    }
                },
                "1425170 [species] Homo heidelbergensis": {
                }
            },
            "239934 [genus] Akkermansia": {
                "239935 [species] Akkermansia muciniphila": {
                    "349741 [no rank] Akkermansia muciniphila ATCC BAA-835": {
                    }
                },
                "512293 [no rank] environmental samples": {
                    "512294 [species] uncultured Akkermansia sp.": {
                    },
                    "1131822 [species] uncultured Akkermansia sp. SMG25": {
                    },
                    "1262691 [species] Akkermansia sp. CAG:344": {
                    },
                    "1263034 [species] Akkermansia muciniphila CAG:154": {
                    }
                },
                "1679444 [species] Akkermansia glycaniphila": {
                },
                "2608915 [no rank] unclassified Akkermansia": {
                    "1131336 [species] Akkermansia sp. KLE1605": {
                    },
                    "1574264 [species] Akkermansia sp. KLE1797": {
                    },
                    "1574265 [species] Akkermansia sp. KLE1798": {
                    },
                    "1638783 [species] Akkermansia sp. UNK.MGS-1": {
                    },
                    "1755639 [species] Akkermansia sp. MC_55": {
                    }
                }
            }
        }


    Snapshot of taxonomy (taxid 1) in kate:
    ![taxon.json.png](files/taxon.json.png)


## lineage

Usage

```text
Query taxonomic lineage of given taxIDs

Input:

  - List of taxIDs, one taxID per line.
  - Or tab-delimited format, please specify taxID field 
    with flag -i/--taxid-field (default 1).
  - Supporting (gzipped) file or STDIN.

Output:

  1. Input line data.
  2. (Optional) Status code (-c/--show-status-code), values:
     - "-1" for queries not found in whole database.
     - "0" for deleted taxIDs, provided by "delnodes.dmp".
     - New taxIDs for merged taxIDs, provided by "merged.dmp".
     - Taxids for these found in "nodes.dmp".
  3. Lineage, delimiter can be changed with flag -d/--delimiter.
  4. (Optional) Lineage in taxIDs (-t/--show-lineage-taxids)
  5. (Optional) Rank (-r/--show-rank)

Filter out invalid and deleted taxids, and replace merged 
taxids with new ones:
    
    # input is one-column-taxid
    $ taxonkit lineage -c taxids.txt \
        | awk '$2>0' \
        | cut -f 2-
        
    # taxids are in 3rd field in a 4-columns tab-delimited file,
    # for $5, where 5 = 4 + 1.
    $ cat input.txt \
        | taxonkit lineage -c -i 3 \
        | csvtk filter2 -H -t -f '$5>0' \
        | csvtk -H -t cut -f -3

Usage:
  taxonkit lineage [flags]

Flags:
  -d, --delimiter string      field delimiter in lineage (default ";")
  -h, --help                  help for lineage
  -L, --no-lineage            do not show lineage, when user just want names or/and ranks
  -R, --show-lineage-ranks    appending ranks of all levels
  -t, --show-lineage-taxids   appending lineage consisting of taxids
  -n, --show-name             appending scientific name
  -r, --show-rank             appending rank of taxids
  -c, --show-status-code      show status code before lineage
  -i, --taxid-field int       field index of taxid. input data should be tab-separated (default 1)

```

Examples

1. Full lineage:

        # note that 123124124 is a fake taxid, 3 was deleted, 92489,1458427 were merged
        $ cat taxids.txt 
        9606
        9913
        376619
        349741
        239935
        314101
        11932
        1327037
        123124124
        3
        92489
        1458427


        $ taxonkit lineage taxids.txt | tee lineage.txt 
        19:22:13.077 [WARN] taxid 92489 was merged into 796334
        19:22:13.077 [WARN] taxid 1458427 was merged into 1458425
        19:22:13.077 [WARN] taxid 123124124 not found
        19:22:13.077 [WARN] taxid 3 was deleted
        9606    cellular organisms;Eukaryota;Opisthokonta;Metazoa;Eumetazoa;Bilateria;Deuterostomia;Chordata;Craniata;Vertebrata;Gnathostomata;Teleostomi;Euteleostomi;Sarcopterygii;Dipnotetrapodomorpha;Tetrapoda;Amniota;Mammalia;Theria;Eutheria;Boreoeutheria;Euarchontoglires;Primates;Haplorrhini;Simiiformes;Catarrhini;Hominoidea;Hominidae;Homininae;Homo;Homo sapiens
        9913    cellular organisms;Eukaryota;Opisthokonta;Metazoa;Eumetazoa;Bilateria;Deuterostomia;Chordata;Craniata;Vertebrata;Gnathostomata;Teleostomi;Euteleostomi;Sarcopterygii;Dipnotetrapodomorpha;Tetrapoda;Amniota;Mammalia;Theria;Eutheria;Boreoeutheria;Laurasiatheria;Artiodactyla;Ruminantia;Pecora;Bovidae;Bovinae;Bos;Bos taurus
        376619  cellular organisms;Bacteria;Proteobacteria;Gammaproteobacteria;Thiotrichales;Francisellaceae;Francisella;Francisella tularensis;Francisella tularensis subsp. holarctica;Francisella tularensis subsp. holarctica LVS
        349741  cellular organisms;Bacteria;PVC group;Verrucomicrobia;Verrucomicrobiae;Verrucomicrobiales;Akkermansiaceae;Akkermansia;Akkermansia muciniphila;Akkermansia muciniphila ATCC BAA-835
        239935  cellular organisms;Bacteria;PVC group;Verrucomicrobia;Verrucomicrobiae;Verrucomicrobiales;Akkermansiaceae;Akkermansia;Akkermansia muciniphila
        314101  cellular organisms;Bacteria;environmental samples;uncultured murine large bowel bacterium BAC 54B
        11932   Viruses;Riboviria;Pararnavirae;Artverviricota;Revtraviricetes;Ortervirales;Retroviridae;unclassified Retroviridae;Intracisternal A-particles;Mouse Intracisternal A-particle
        1327037 Viruses;Duplodnaviria;Heunggongvirae;Uroviricota;Caudoviricetes;Caudovirales;Siphoviridae;unclassified Siphoviridae;Croceibacter phage P2559Y
        123124124
        3
        92489   cellular organisms;Bacteria;Proteobacteria;Gammaproteobacteria;Enterobacterales;Erwiniaceae;Erwinia;Erwinia oleae
        1458427 cellular organisms;Bacteria;Proteobacteria;Betaproteobacteria;Burkholderiales;Comamonadaceae;Serpentinomonas;Serpentinomonas raicheisms;Bacteria;Proteobacteria;Betaproteobacteria;Burkholderiales;Comamonadaceae;Serpentinomonas;Serpentinomonas raichei

1. Checking deleted or merged taxids
    
        $ taxonkit lineage --show-status-code taxids.txt | tee lineage.withcode.txt
        
        # valid
        $ cat lineage.withcode.txt | awk '$2 > 0' | cut -f 1,2
        9606    9606
        9913    9913
        376619  376619
        349741  349741
        239935  239935
        314101  314101
        11932   11932
        1327037 1327037
        92489   796334
        1458427 1458425

        # merged
        $ cat lineage.withcode.txt | awk '$2 > 0 && $2 != $1' | cut -f 1,2
        92489   796334
        1458427 1458425

        # deleted
        $ cat lineage.withcode.txt | awk '$2 == 0' | cut -f 1
        3
        
        # invalid
        $ cat lineage.withcode.txt | awk '$2 < 0' | cut -f 1
        123124124

1. **Filter out invalid and deleted taxids, and *replace merged taxids with new ones"**,
  you may install [csvtk](https://github.com/shenwei356/csvtk.
    
        # input is one-column-taxid
        $ taxonkit lineage -c taxids.txt \
            | awk '$2>0' \
            | cut -f 2-
            
        # taxids are in 3rd field in a 4-columns tab-delimited file,
        # for $5, where 5 = 4 + 1.
        $ cat input.txt \
            | taxonkit lineage -c -i 3 \
            | csvtk filter2 -H -t -f '$5>0' \
            | csvtk -H -t cut -f -3

1. Only show name and rank.

        $ taxonkit lineage -r -n -L taxids.txt \
            | csvtk pretty -t
        9606        Homo sapiens                                      species
        9913        Bos taurus                                        species
        376619      Francisella tularensis subsp. holarctica LVS      strain
        349741      Akkermansia muciniphila ATCC BAA-835              strain
        239935      Akkermansia muciniphila                           species
        314101      uncultured murine large bowel bacterium BAC 54B   species
        11932       Mouse Intracisternal A-particle                   species
        1327037     Croceibacter phage P2559Y                         species
        123124124                                                     
        3                                                             
        92489       Erwinia oleae                                     species
        1458427     Serpentinomonas raichei                           species
   

1. Show lineage consisting of taxids:

        $ taxonkit lineage -t taxids.txt
        9606    cellular organisms;Eukaryota;Opisthokonta;Metazoa;Eumetazoa;Bilateria;Deuterostomia;Chordata;Craniata;Vertebrata;Gnathostomata;Teleostomi;Euteleostomi;Sarcopterygii;Dipnotetrapodomorpha;Tetrapoda;Amniota;Mammalia;Theria;Eutheria;Boreoeutheria;Euarchontoglires;Primates;Haplorrhini;Simiiformes;Catarrhini;Hominoidea;Hominidae;Homininae;Homo;Homo sapiens  131567;2759;33154;33208;6072;33213;33511;7711;89593;7742;7776;117570;117571;8287;1338369;32523;32524;40674;32525;9347;1437010;314146;9443;376913;314293;9526;314295;9604;207598;9605;9606
        9913    cellular organisms;Eukaryota;Opisthokonta;Metazoa;Eumetazoa;Bilateria;Deuterostomia;Chordata;Craniata;Vertebrata;Gnathostomata;Teleostomi;Euteleostomi;Sarcopterygii;Dipnotetrapodomorpha;Tetrapoda;Amniota;Mammalia;Theria;Eutheria;Boreoeutheria;Laurasiatheria;Artiodactyla;Ruminantia;Pecora;Bovidae;Bovinae;Bos;Bos taurus   131567;2759;33154;33208;6072;33213;33511;7711;89593;7742;7776;117570;117571;8287;1338369;32523;32524;40674;32525;9347;1437010;314145;91561;9845;35500;9895;27592;9903;9913
        376619  cellular organisms;Bacteria;Proteobacteria;Gammaproteobacteria;Thiotrichales;Francisellaceae;Francisella;Francisella tularensis;Francisella tularensis subsp. holarctica;Francisella tularensis subsp. holarctica LVS     131567;2;1224;1236;72273;34064;262;263;119857;376619
        349741  cellular organisms;Bacteria;PVC group;Verrucomicrobia;Verrucomicrobiae;Verrucomicrobiales;Akkermansiaceae;Akkermansia;Akkermansia muciniphila;Akkermansia muciniphila ATCC BAA-835      131567;2;1783257;74201;203494;48461;1647988;239934;239935;349741
        239935  cellular organisms;Bacteria;PVC group;Verrucomicrobia;Verrucomicrobiae;Verrucomicrobiales;Akkermansiaceae;Akkermansia;Akkermansia muciniphila   131567;2;1783257;74201;203494;48461;1647988;239934;239935
        314101  cellular organisms;Bacteria;environmental samples;uncultured murine large bowel bacterium BAC 54B       131567;2;48479;314101
        11932   Viruses;Riboviria;Pararnavirae;Artverviricota;Revtraviricetes;Ortervirales;Retroviridae;unclassified Retroviridae;Intracisternal A-particles;Mouse Intracisternal A-particle    10239;2559587;2732397;2732409;2732514;2169561;11632;35276;11749;11932
        1327037 Viruses;Duplodnaviria;Heunggongvirae;Uroviricota;Caudoviricetes;Caudovirales;Siphoviridae;unclassified Siphoviridae;Croceibacter phage P2559Y   10239;2731341;2731360;2731618;2731619;28883;10699;196894;1327037
        123124124
        3
        92489   cellular organisms;Bacteria;Proteobacteria;Gammaproteobacteria;Enterobacterales;Erwiniaceae;Erwinia;Erwinia oleae       131567;2;1224;1236;91347;1903409;551;796334
        1458427 cellular organisms;Bacteria;Proteobacteria;Betaproteobacteria;Burkholderiales;Comamonadaceae;Serpentinomonas;Serpentinomonas raichei    131567;2;1224;28216;80840;80864;2490452;1458425

    or read taxids from STDIN:

        $ cat taxids.txt | taxonkit lineage

1. And ranks of all nodes:

        $ echo 2697049 \
            | taxonkit lineage -t -R \
            | csvtk transpose -Ht
        2697049
        Viruses;Riboviria;Orthornavirae;Pisuviricota;Pisoniviricetes;Nidovirales;Cornidovirineae;Coronaviridae;Orthocoronavirinae;Betacoronavirus;Sarbecovirus;Severe acute respiratory syndrome-related coronavirus;Severe acute respiratory syndrome coronavirus 2
        10239;2559587;2732396;2732408;2732506;76804;2499399;11118;2501931;694002;2509511;694009;2697049
        superkingdom;clade;kingdom;phylum;class;order;suborder;family;subfamily;genus;subgenus;species;no rank

    Another way to show lineage detail of a taxID
    
        $ echo 2697049 \
            | taxonkit lineage -t \
            | csvtk cut -Ht -f 3 \
            | csvtk unfold -Ht -f 1 -s ";" \
            | taxonkit lineage -r -n -L \
            | csvtk cut -Ht -f 1,3,2 \
            | csvtk pretty -t 
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
        
## reformat

Usage

```text
Reformat lineage in canonical ranks

Output format can be formated by flag --format, available placeholders:

    {k}: superkingdom
    {p}: phylum
    {c}: class
    {o}: order
    {f}: family
    {g}: genus
    {s}: species
    {t}: subspecies/strain
    
    {S}: subspecies
    {T}: strain

When these's no nodes of rank "subspecies" nor "stain",
you can switch -S/--pseudo-strain to use the node with lowest rank
as subspecies/strain name, if which rank is lower than "species". 
This flag affects {t}, {S}, {T}.

Output format can contains some escape charactors like "\t".

This command appends reformated lineage to the input line.
The corresponding taxIDs of reformated lineage can be provided as another
column by flag "-t/--show-lineage-taxids".

Usage:
  taxonkit reformat [flags]

Flags:
  -P, --add-prefix                     add prefixes for all ranks, single prefix for a rank is defined by flag --prefix-X
  -d, --delimiter string               field delimiter in input lineage (default ";")
  -F, --fill-miss-rank                 fill missing rank with lineage information of the next higher rank
  -f, --format string                  output format, placeholders of rank are needed (default "{k};{p};{c};{o};{f};{g};{s}")
  -h, --help                           help for reformat
  -i, --lineage-field int              field index of lineage. data should be tab-separated (default 2)
  -r, --miss-rank-repl string          replacement string for missing rank
  -p, --miss-rank-repl-prefix string   prefix for estimated taxon level (default "unclassified ")
  -R, --miss-taxid-repl string         replacement string for missing taxid
      --prefix-S string                prefix for subspecies, used along with flag -P/--add-prefix (default "S__")
      --prefix-T string                prefix for strain, used along with flag -P/--add-prefix (default "T__")
      --prefix-c string                prefix for class, used along with flag -P/--add-prefix (default "c__")
      --prefix-f string                prefix for family, used along with flag -P/--add-prefix (default "f__")
      --prefix-g string                prefix for genus, used along with flag -P/--add-prefix (default "g__")
      --prefix-k string                prefix for superkingdom, used along with flag -P/--add-prefix (default "k__")
      --prefix-o string                prefix for order, used along with flag -P/--add-prefix (default "o__")
      --prefix-p string                prefix for phylum, used along with flag -P/--add-prefix (default "p__")
      --prefix-s string                prefix for species, used along with flag -P/--add-prefix (default "s__")
      --prefix-t string                prefix for subspecies/strain, used along with flag -P/--add-prefix (default "t__")
  -S, --pseudo-strain                  use the node with lowest rank as strain name, only if which rank is lower than "species" and not "subpecies" nor "strain". This flag affects {t}, {S}, {T}
  -t, --show-lineage-taxids            show corresponding taxids of reformated lineage
  -T, --trim                           do not fill missing rank lower than current rank

```

Examples:

1. Example lineage (produced by: `taxonkit lineage taxids.txt | awk '$2!=""' > lineage.txt`).

        $ cat lineage.txt
        9606    cellular organisms;Eukaryota;Opisthokonta;Metazoa;Eumetazoa;Bilateria;Deuterostomia;Chordata;Craniata;Vertebrata;Gnathostomata;Teleostomi;Euteleostomi;Sarcopterygii;Dipnotetrapodomorpha;Tetrapoda;Amniota;Mammalia;Theria;Eutheria;Boreoeutheria;Euarchontoglires;Primates;Haplorrhini;Simiiformes;Catarrhini;Hominoidea;Hominidae;Homininae;Homo;Homo sapiens
        9913    cellular organisms;Eukaryota;Opisthokonta;Metazoa;Eumetazoa;Bilateria;Deuterostomia;Chordata;Craniata;Vertebrata;Gnathostomata;Teleostomi;Euteleostomi;Sarcopterygii;Dipnotetrapodomorpha;Tetrapoda;Amniota;Mammalia;Theria;Eutheria;Boreoeutheria;Laurasiatheria;Artiodactyla;Ruminantia;Pecora;Bovidae;Bovinae;Bos;Bos taurus
        376619  cellular organisms;Bacteria;Proteobacteria;Gammaproteobacteria;Thiotrichales;Francisellaceae;Francisella;Francisella tularensis;Francisella tularensis subsp. holarctica;Francisella tularensis subsp. holarctica LVS
        349741  cellular organisms;Bacteria;PVC group;Verrucomicrobia;Verrucomicrobiae;Verrucomicrobiales;Akkermansiaceae;Akkermansia;Akkermansia muciniphila;Akkermansia muciniphila ATCC BAA-835
        239935  cellular organisms;Bacteria;PVC group;Verrucomicrobia;Verrucomicrobiae;Verrucomicrobiales;Akkermansiaceae;Akkermansia;Akkermansia muciniphila
        314101  cellular organisms;Bacteria;environmental samples;uncultured murine large bowel bacterium BAC 54B
        11932   Viruses;Riboviria;Pararnavirae;Artverviricota;Revtraviricetes;Ortervirales;Retroviridae;unclassified Retroviridae;Intracisternal A-particles;Mouse Intracisternal A-particle
        1327037 Viruses;Duplodnaviria;Heunggongvirae;Uroviricota;Caudoviricetes;Caudovirales;Siphoviridae;unclassified Siphoviridae;Croceibacter phage P2559Y
        92489   cellular organisms;Bacteria;Proteobacteria;Gammaproteobacteria;Enterobacterales;Erwiniaceae;Erwinia;Erwinia oleae
        1458427 cellular organisms;Bacteria;Proteobacteria;Betaproteobacteria;Burkholderiales;Comamonadaceae;Serpentinomonas;Serpentinomonas raichei

1. Default output format (`"{k};{p};{c};{o};{f};{g};{s}"`).

        # reformated lineages are appended to the input data
        $ taxonkit reformat lineage.txt 
        ...
        239935  cellular organisms;Bacteria;PVC group;Verrucomicrobia;Verrucomicrobiae;Verrucomicrobiales;Akkermansiaceae;Akkermansia;Akkermansia muciniphila     Bacteria;Verrucomicrobia;Verrucomicrobiae;Verrucomicrobiales;Akkermansiaceae;Akkermansia;Akkermansia muciniphila
        ...

        $ 
        $ taxonkit reformat lineage.txt | tee lineage.txt.reformat
        
        $ cut -f 1,3 lineage.txt.reformat
        9606    Eukaryota;Chordata;Mammalia;Primates;Hominidae;Homo;Homo sapiens
        9913    Eukaryota;Chordata;Mammalia;Artiodactyla;Bovidae;Bos;Bos taurus
        376619  Bacteria;Proteobacteria;Gammaproteobacteria;Thiotrichales;Francisellaceae;Francisella;Francisella tularensis
        349741  Bacteria;Verrucomicrobia;Verrucomicrobiae;Verrucomicrobiales;Akkermansiaceae;Akkermansia;Akkermansia muciniphila
        239935  Bacteria;Verrucomicrobia;Verrucomicrobiae;Verrucomicrobiales;Akkermansiaceae;Akkermansia;Akkermansia muciniphila
        314101  Bacteria;;;;;;uncultured murine large bowel bacterium BAC 54B
        11932   Viruses;Artverviricota;Revtraviricetes;Ortervirales;Retroviridae;Intracisternal A-particles;Mouse Intracisternal A-particle
        1327037 Viruses;Uroviricota;Caudoviricetes;Caudovirales;Siphoviridae;;Croceibacter phage P2559Y
        92489   Bacteria;Proteobacteria;Gammaproteobacteria;Enterobacterales;Erwiniaceae;Erwinia;Erwinia oleae
        1458427 Bacteria;Proteobacteria;Betaproteobacteria;Burkholderiales;Comamonadaceae;Serpentinomonas;Serpentinomonas raichei
            
        # aligned 
        $ cat lineage.txt \
            | taxonkit reformat \
            | csvtk -H -t cut -f 1,3 \
            | csvtk -H -t sep -f 2 -s ';' -R \
            | csvtk add-header -t -n taxid,kindom,phylum,class,order,family,genus,species \
            | csvtk pretty -t
            
            taxid     kindom      phylum            class                 order                family            genus                        species
            9606      Eukaryota   Chordata          Mammalia              Primates             Hominidae         Homo                         Homo sapiens
            9913      Eukaryota   Chordata          Mammalia              Artiodactyla         Bovidae           Bos                          Bos taurus
            376619    Bacteria    Proteobacteria    Gammaproteobacteria   Thiotrichales        Francisellaceae   Francisella                  Francisella tularensis
            349741    Bacteria    Verrucomicrobia   Verrucomicrobiae      Verrucomicrobiales   Akkermansiaceae   Akkermansia                  Akkermansia muciniphila
            239935    Bacteria    Verrucomicrobia   Verrucomicrobiae      Verrucomicrobiales   Akkermansiaceae   Akkermansia                  Akkermansia muciniphila
            314101    Bacteria                                                                                                                uncultured murine large bowel bacterium BAC 54B
            11932     Viruses     Artverviricota    Revtraviricetes       Ortervirales         Retroviridae      Intracisternal A-particles   Mouse Intracisternal A-particle
            1327037   Viruses     Uroviricota       Caudoviricetes        Caudovirales         Siphoviridae                                   Croceibacter phage P2559Y
            92489     Bacteria    Proteobacteria    Gammaproteobacteria   Enterobacterales     Erwiniaceae       Erwinia                      Erwinia oleae
            1458427   Bacteria    Proteobacteria    Betaproteobacteria    Burkholderiales      Comamonadaceae    Serpentinomonas              Serpentinomonas raichei

1. And `subspecies/stain` (`{t}`), `subspecies` (`{S}`), and `strain` (`{T}`) are also available.


        # default operation
        $ echo -ne "239935\n83333\n1408252\n2697049\n2605619\n" \
            | taxonkit lineage -n -r \
            | taxonkit reformat -f '{t};{S};{T}' \
            | csvtk -H -t cut -f 1,4,3,5 \
            | csvtk -H -t sep -f 4 -s ';' -R \
            | csvtk -H -t add-header -n "taxid,rank,name,subspecies/stain,subspecies,strain" \
            | csvtk pretty -t
        taxid     rank         name                                              subspecies/stain        subspecies              strain
        239935    species      Akkermansia muciniphila                                                                           
        83333     strain       Escherichia coli K-12                             Escherichia coli K-12                           Escherichia coli K-12
        1408252   subspecies   Escherichia coli R178                             Escherichia coli R178   Escherichia coli R178   
        2697049   no rank      Severe acute respiratory syndrome coronavirus 2                                                   
        2605619   no rank      Escherichia coli O16:H48

        
        # fill missing ranks
        # see example below for -F/--fill-miss-rank
        #
        $ echo -ne "239935\n83333\n1408252\n2697049\n2605619\n" \
            | taxonkit lineage -n -r \
            | taxonkit reformat -f '{t};{S};{T}' --fill-miss-rank \
            | csvtk -H -t cut -f 1,4,3,5 \
            | csvtk -H -t sep -f 4 -s ';' -R \
            | csvtk -H -t add-header -n "taxid,rank,name,subspecies/stain,subspecies,strain" \
            | csvtk pretty -t
        taxid     rank         name                                              subspecies/stain                                                                       subspecies                                                                      strain
        239935    species      Akkermansia muciniphila                           unclassified Akkermansia muciniphila subspecies/strain                                 unclassified Akkermansia muciniphila subspecies                                 unclassified Akkermansia muciniphila strain
        83333     strain       Escherichia coli K-12                             Escherichia coli K-12                                                                  unclassified Escherichia coli subspecies                                        Escherichia coli K-12
        1408252   subspecies   Escherichia coli R178                             Escherichia coli R178                                                                  Escherichia coli R178                                                           unclassified Escherichia coli R178 strain
        2697049   no rank      Severe acute respiratory syndrome coronavirus 2   unclassified Severe acute respiratory syndrome-related coronavirus subspecies/strain   unclassified Severe acute respiratory syndrome-related coronavirus subspecies   unclassified Severe acute respiratory syndrome-related coronavirus strain
        2605619   no rank      Escherichia coli O16:H48                          unclassified Escherichia coli subspecies/strain                                        unclassified Escherichia coli subspecies                                        unclassified Escherichia coli strain

        # When these's no nodes of rank "subspecies" nor "stain",
        # you can switch -S/--pseudo-strain to use the node with lowest rank
        # as subspecies/strain name, if which rank is lower than "species"
        #
        $ echo -ne "239935\n83333\n1408252\n2697049\n2605619\n" \
            | taxonkit lineage -n -r \
            | taxonkit reformat -f '{t};{S};{T}' --fill-miss-rank  --pseudo-strain \
            | csvtk -H -t cut -f 1,4,3,5 \
            | csvtk -H -t sep -f 4 -s ';' -R \
            | csvtk -H -t add-header -n "taxid,rank,name,subspecies/stain,subspecies,strain" \
            | csvtk pretty -t
        taxid     rank         name                                              subspecies/stain                                         subspecies                                        strain
        239935    species      Akkermansia muciniphila                           unclassified Akkermansia muciniphila subspecies/strain   unclassified Akkermansia muciniphila subspecies   unclassified Akkermansia muciniphila strain
        83333     strain       Escherichia coli K-12                             Escherichia coli K-12                                    unclassified Escherichia coli subspecies          Escherichia coli K-12
        1408252   subspecies   Escherichia coli R178                             Escherichia coli R178                                    Escherichia coli R178                             unclassified Escherichia coli R178 strain
        2697049   no rank      Severe acute respiratory syndrome coronavirus 2   Severe acute respiratory syndrome coronavirus 2          Severe acute respiratory syndrome coronavirus 2   Severe acute respiratory syndrome coronavirus 2
        2605619   no rank      Escherichia coli O16:H48                          Escherichia coli O16:H48                                 Escherichia coli O16:H48                          Escherichia coli O16:H48


1. Add prefix (`-P/--add-prefix`).

        $ cat lineage.txt \
            | taxonkit reformat -P \
            | csvtk -H -t cut -f 1,3 
        9606    k__Eukaryota;p__Chordata;c__Mammalia;o__Primates;f__Hominidae;g__Homo;s__Homo sapiens
        9913    k__Eukaryota;p__Chordata;c__Mammalia;o__Artiodactyla;f__Bovidae;g__Bos;s__Bos taurus
        376619  k__Bacteria;p__Proteobacteria;c__Gammaproteobacteria;o__Thiotrichales;f__Francisellaceae;g__Francisella;s__Francisella tularensis
        349741  k__Bacteria;p__Verrucomicrobia;c__Verrucomicrobiae;o__Verrucomicrobiales;f__Akkermansiaceae;g__Akkermansia;s__Akkermansia muciniphila
        239935  k__Bacteria;p__Verrucomicrobia;c__Verrucomicrobiae;o__Verrucomicrobiales;f__Akkermansiaceae;g__Akkermansia;s__Akkermansia muciniphila
        314101  k__Bacteria;p__;c__;o__;f__;g__;s__uncultured murine large bowel bacterium BAC 54B
        11932   k__Viruses;p__Artverviricota;c__Revtraviricetes;o__Ortervirales;f__Retroviridae;g__Intracisternal A-particles;s__Mouse Intracisternal A-particle
        1327037 k__Viruses;p__Uroviricota;c__Caudoviricetes;o__Caudovirales;f__Siphoviridae;g__;s__Croceibacter phage P2559Y
        92489   k__Bacteria;p__Proteobacteria;c__Gammaproteobacteria;o__Enterobacterales;f__Erwiniaceae;g__Erwinia;s__Erwinia oleae
        1458427 k__Bacteria;p__Proteobacteria;c__Betaproteobacteria;o__Burkholderiales;f__Comamonadaceae;g__Serpentinomonas;s__Serpentinomonas raichei

        
1. Show corresponding taxids of reformated lineage (flag `-t/--show-lineage-taxids`)

        $ cat lineage.txt \
            | taxonkit reformat -t \
            | csvtk -H -t cut -f 1,4 \
            | csvtk -H -t sep -f 2 -s ';' -R \
            | csvtk add-header -t -n taxid,kindom,phylum,class,order,family,genus,species \
            | csvtk pretty -t   
            
        taxid     kindom   phylum    class     order     family    genus     species
        9606      2759     7711      40674     9443      9604      9605      9606
        9913      2759     7711      40674     91561     9895      9903      9913
        376619    2        1224      1236      72273     34064     262       263
        349741    2        74201     203494    48461     1647988   239934    239935
        239935    2        74201     203494    48461     1647988   239934    239935
        314101    2                                                          314101
        11932     10239    2732409   2732514   2169561   11632     11749     11932
        1327037   10239    2731618   2731619   28883     10699               1327037
        92489     2        1224      1236      91347     1903409   551       796334
        1458427   2        1224      28216     80840     80864     2490452   1458425


1. Use custom symbols for unclassfied ranks (`-r/--miss-rank-repl`)

        $ taxonkit reformat lineage.txt -r "__"  | cut -f 3
        Eukaryota;Chordata;Mammalia;Primates;Hominidae;Homo;Homo sapiens
        Eukaryota;Chordata;Mammalia;Artiodactyla;Bovidae;Bos;Bos taurus
        Bacteria;Proteobacteria;Gammaproteobacteria;Thiotrichales;Francisellaceae;Francisella;Francisella tularensis
        Bacteria;Verrucomicrobia;Verrucomicrobiae;Verrucomicrobiales;Akkermansiaceae;Akkermansia;Akkermansia muciniphila
        Bacteria;Verrucomicrobia;Verrucomicrobiae;Verrucomicrobiales;Akkermansiaceae;Akkermansia;Akkermansia muciniphila
        Bacteria;__;__;__;__;__;uncultured murine large bowel bacterium BAC 54B
        Viruses;Artverviricota;Revtraviricetes;Ortervirales;Retroviridae;Intracisternal A-particles;Mouse Intracisternal A-particle
        Viruses;Uroviricota;Caudoviricetes;Caudovirales;Siphoviridae;__;Croceibacter phage P2559Y
        Bacteria;Proteobacteria;Gammaproteobacteria;Enterobacterales;Erwiniaceae;Erwinia;Erwinia oleae
        Bacteria;Proteobacteria;Betaproteobacteria;Burkholderiales;Comamonadaceae;Serpentinomonas;Serpentinomonas raichei
        
        $ taxonkit reformat lineage.txt -r Unassigned | cut -f 3
        Eukaryota;Chordata;Mammalia;Primates;Hominidae;Homo;Homo sapiens
        Eukaryota;Chordata;Mammalia;Artiodactyla;Bovidae;Bos;Bos taurus
        Bacteria;Proteobacteria;Gammaproteobacteria;Thiotrichales;Francisellaceae;Francisella;Francisella tularensis
        Bacteria;Verrucomicrobia;Verrucomicrobiae;Verrucomicrobiales;Akkermansiaceae;Akkermansia;Akkermansia muciniphila
        Bacteria;Verrucomicrobia;Verrucomicrobiae;Verrucomicrobiales;Akkermansiaceae;Akkermansia;Akkermansia muciniphila
        Bacteria;Unassigned;Unassigned;Unassigned;Unassigned;Unassigned;uncultured murine large bowel bacterium BAC 54B
        Viruses;Artverviricota;Revtraviricetes;Ortervirales;Retroviridae;Intracisternal A-particles;Mouse Intracisternal A-particle
        Viruses;Uroviricota;Caudoviricetes;Caudovirales;Siphoviridae;Unassigned;Croceibacter phage P2559Y
        Bacteria;Proteobacteria;Gammaproteobacteria;Enterobacterales;Erwiniaceae;Erwinia;Erwinia oleae
        Bacteria;Proteobacteria;Betaproteobacteria;Burkholderiales;Comamonadaceae;Serpentinomonas;Serpentinomonas raichei
            

1. **Estimate and fill missing rank with original lineage information**
  (`-F, --fill-miss-rank`, very useful for formating input data for
  [LEfSe](https://bitbucket.org/biobakery/biobakery/wiki/lefse)).
  You can change the prefix "unclassified" using flag `-p/--miss-rank-repl-prefix`.

        $ cat lineage.txt \
            | taxonkit reformat -F \
            | csvtk -H -t cut -f 1,3 \
            | csvtk -H -t sep -f 2 -s ';' -R \
            | csvtk add-header -t -n taxid,kindom,phylum,class,order,family,genus,species \
            | csvtk pretty -t
            
        taxid     kindom      phylum                         class                         order                         family                         genus                             species
        9606      Eukaryota   Chordata                       Mammalia                      Primates                      Hominidae                      Homo                              Homo sapiens
        9913      Eukaryota   Chordata                       Mammalia                      Artiodactyla                  Bovidae                        Bos                               Bos taurus
        376619    Bacteria    Proteobacteria                 Gammaproteobacteria           Thiotrichales                 Francisellaceae                Francisella                       Francisella tularensis
        349741    Bacteria    Verrucomicrobia                Verrucomicrobiae              Verrucomicrobiales            Akkermansiaceae                Akkermansia                       Akkermansia muciniphila
        239935    Bacteria    Verrucomicrobia                Verrucomicrobiae              Verrucomicrobiales            Akkermansiaceae                Akkermansia                       Akkermansia muciniphila
        314101    Bacteria    unclassified Bacteria phylum   unclassified Bacteria class   unclassified Bacteria order   unclassified Bacteria family   unclassified Bacteria genus       uncultured murine large bowel bacterium BAC 54B
        11932     Viruses     Artverviricota                 Revtraviricetes               Ortervirales                  Retroviridae                   Intracisternal A-particles        Mouse Intracisternal A-particle
        1327037   Viruses     Uroviricota                    Caudoviricetes                Caudovirales                  Siphoviridae                   unclassified Siphoviridae genus   Croceibacter phage P2559Y
        92489     Bacteria    Proteobacteria                 Gammaproteobacteria           Enterobacterales              Erwiniaceae                    Erwinia                           Erwinia oleae
        1458427   Bacteria    Proteobacteria                 Betaproteobacteria            Burkholderiales               Comamonadaceae                 Serpentinomonas                   Serpentinomonas raichei

1. Only some ranks.

        $ cat lineage.txt \
            | taxonkit reformat -F -f "{s};{p}"\
            | csvtk -H -t cut -f 1,3 \
            | csvtk -H -t sep -f 2 -s ';' -R \
            | csvtk add-header -t -n taxid,species,phylum \
            | csvtk pretty -t
            
        taxid     species                                           phylum
        9606      Homo sapiens                                      Chordata
        9913      Bos taurus                                        Chordata
        376619    Francisella tularensis                            Proteobacteria
        349741    Akkermansia muciniphila                           Verrucomicrobia
        239935    Akkermansia muciniphila                           Verrucomicrobia
        314101    uncultured murine large bowel bacterium BAC 54B   unclassified Bacteria phylum
        11932     Mouse Intracisternal A-particle                   Artverviricota
        1327037   Croceibacter phage P2559Y                         Uroviricota
        92489     Erwinia oleae                                     Proteobacteria
        1458427   Serpentinomonas raichei                           Proteobacteria

1. **For some taxids which rank is higher than the lowest rank in `-f/--format`, use `-T/--trim` to avoid fill missing rank lower than current rank**.

        $ echo -ne "2\n239934\n239935\n" \
            | taxonkit lineage \
            | taxonkit reformat -T \
            | sed -r "s/;+$//" \
            | csvtk -H -t cut -f 1,3 
        2       Bacteria
        239934  Bacteria;Verrucomicrobia;Verrucomicrobiae;Verrucomicrobiales;Akkermansiaceae;Akkermansia
        239935  Bacteria;Verrucomicrobia;Verrucomicrobiae;Verrucomicrobiales;Akkermansiaceae;Akkermansia;Akkermansia muciniphila

1. Support tab in format string

        $ echo 9606 \
            | taxonkit lineage \
            | taxonkit reformat -f "{k}\t{p}\t{c}\t{o}\t{f}\t{g}\t{s}\t{S}" \
            | csvtk cut -t -f -2
        9606    Eukaryota       Chordata        Mammalia        Primates        Hominidae       Homo    Homo sapiens

1. List seven-level lineage for all taxIDs.

        # replace empty taxon with "Unassigned"
        $ taxonkit list --ids 1 \
            | taxonkit lineage \
            | taxonkit reformat -r Unassigned 
            | gzip -c > all.lineage.tsv.gz
        
        # tab-delimited seven-levels
        $ taxonkit list --ids 1 \
            | taxonkit lineage \
            | taxonkit reformat -r Unassigned -f "{k}\t{p}\t{c}\t{o}\t{f}\t{g}\t{s}" \
            | csvtk cut -H -t -f -2 \
            | head -n 5 \
            | csvtk pretty -t
        
        # 8-level
        $ taxonkit list --ids 1 \
            | taxonkit lineage \
            | taxonkit reformat -r Unassigned -f "{k}\t{p}\t{c}\t{o}\t{f}\t{g}\t{s}\t{t}" \
            | csvtk cut -H -t -f -2 \
            | head -n 5 \
            | csvtk pretty -t
        
        # Fill and trim
        $ memusg -t -s ' taxonkit list --ids 1 \
            | taxonkit lineage \
            | taxonkit reformat -F -T \
            | sed -r "s/;+$//" \
            | gzip -c > all.lineage.tsv.gz '
        
        elapsed time: 19.930s
        peak rss: 6.25 GB


1. From taxid to 7-ranks lineage:

        $ cat taxids.txt | taxonkit lineage | taxonkit reformat

## name2taxid

Usage

```text
Convert scientific names to taxIDs

Attention:

  1. Some taxIDs share the same scientific names, e.g, Drosophila.
     These input lines are duplicated with multiple taxIDs.

    $ echo Drosophila | taxonkit name2taxid | taxonkit lineage -i 2 -r -L
    Drosophila      7215    genus
    Drosophila      32281   subgenus
    Drosophila      2081351 genus

Usage:
  taxonkit name2taxid [flags]

Flags:
  -h, --help             help for name2taxid
  -i, --name-field int   field index of name. data should be tab-separated (default 1)
  -s, --sci-name         only searching scientific names
  -r, --show-rank        show rank

```

Examples

Example data

    $ cat names.txt
    Homo sapiens
    Akkermansia muciniphila ATCC BAA-835
    Akkermansia muciniphila
    Mouse Intracisternal A-particle
    Wei Shen
    uncultured murine large bowel bacterium BAC 54B
    Croceibacter phage P2559Y

1. Default.

        # taxonkit name2taxid names.txt
        $ cat names.txt | taxonkit name2taxid | csvtk pretty -t
        Homo sapiens                                      9606
        Akkermansia muciniphila ATCC BAA-835              349741
        Akkermansia muciniphila                           239935
        Mouse Intracisternal A-particle                   11932
        Wei Shen                                          
        uncultured murine large bowel bacterium BAC 54B   314101
        Croceibacter phage P2559Y                         1327037

1. Show rank.

        $ cat names.txt | taxonkit name2taxid --show-rank | csvtk pretty -t
        Homo sapiens                                      9606      species
        Akkermansia muciniphila ATCC BAA-835              349741    strain
        Akkermansia muciniphila                           239935    species
        Mouse Intracisternal A-particle                   11932     species
        Wei Shen                                                    
        uncultured murine large bowel bacterium BAC 54B   314101    species
        Croceibacter phage P2559Y                         1327037   species

1. From name to lineage.

        $ cat names.txt | taxonkit name2taxid | taxonkit lineage --taxid-field 2
        Homo sapiens    9606    cellular organisms;Eukaryota;Opisthokonta;Metazoa;Eumetazoa;Bilateria;Deuterostomia;Chordata;Craniata;Vertebrata;Gnathostomata;Teleostomi;Euteleostomi;Sarcopterygii;Dipnotetrapodomorpha;Tetrapoda;Amniota;Mammalia;Theria;Eutheria;Boreoeutheria;Euarchontoglires;Primates;Haplorrhini;Simiiformes;Catarrhini;Hominoidea;Hominidae;Homininae;Homo;Homo sapiens
        Akkermansia muciniphila ATCC BAA-835    349741  cellular organisms;Bacteria;PVC group;Verrucomicrobia;Verrucomicrobiae;Verrucomicrobiales;Akkermansiaceae;Akkermansia;Akkermansia muciniphila;Akkermansia muciniphila ATCC BAA-835
        Akkermansia muciniphila 239935  cellular organisms;Bacteria;PVC group;Verrucomicrobia;Verrucomicrobiae;Verrucomicrobiales;Akkermansiaceae;Akkermansia;Akkermansia muciniphila
        Mouse Intracisternal A-particle 11932   Viruses;Ortervirales;Retroviridae;unclassified Retroviridae;Intracisternal A-particles;Mouse Intracisternal A-particle
        Wei Shen
        uncultured murine large bowel bacterium BAC 54B 314101  cellular organisms;Bacteria;environmental samples;uncultured murine large bowel bacterium BAC 54B
        Croceibacter phage P2559Y       1327037 Viruses;Caudovirales;Siphoviridae;unclassified Siphoviridae;Croceibacter phage P2559Y

1. **Some taxIDs share the same scientific names**, e.g, Drosophila.

        $ echo Drosophila \
            | taxonkit name2taxid \
            | taxonkit lineage -i 2 -r \
            | taxonkit reformat -i 3 \
            | csvtk cut -H -t -f 1,2,4,5 \
            | csvtk pretty -t
        Drosophila   7215      genus      Eukaryota;Arthropoda;Insecta;Diptera;Drosophilidae;Drosophila;
        Drosophila   32281     subgenus   Eukaryota;Arthropoda;Insecta;Diptera;Drosophilidae;Drosophila;
        Drosophila   2081351   genus      Eukaryota;Basidiomycota;Agaricomycetes;Agaricales;Psathyrellaceae;Drosophila;

## filter

Usage

```text
Filter taxIDs by taxonomic rank range

Attentions:

  1. Flag -L/--lower-than and -H/--higher-than are exclusive, and can be
     used along with -E/--equal-to which values can be different.
  2. A list of pre-ordered ranks is in ~/.taxonkit/ranks.txt, you can use
     your list by -r/--rank-file, the format specification is below.
  3. All ranks in taxonomy database should be defined in rank file.
  4. Ranks can be removed with black list via -B/--black-list.
  5. TaxIDs with no rank can be optionally discarded by -N/--discard-noranks.
  6. But when filtering with -L/--lower-than, you can use
    -n/--save-predictable-norank to save some special ranks without order,
    where rank of the closest higher node is still lower than rank cutoff.

Rank file:

  1. Blank lines or lines starting with "#" are ignored.
  2. Ranks are in decending order and case ignored.
  3. Ranks with same order should be in one line separated with comma (",", no space).
  4. Ranks without order should be assigned a prefix symbol "!" for each rank.

Usage:
  taxonkit filter [flags]

Flags:
  -B, --black-list strings        black list of ranks to discard, e.g., '-B "no rank" -B "clade"
  -N, --discard-noranks           discard all ranks without order, type "taxonkit filter --help" for details
  -R, --discard-root              discard root taxid, defined by --root-taxid
  -E, --equal-to strings          output taxIDs with rank equal to some ranks, multiple values can be separated with comma "," (e.g., -E "genus,species"), or give multiple times (e.g., -E genus -E species)
  -h, --help                      help for filter
  -H, --higher-than string        output taxIDs with rank higher than a rank, exclusive with --lower-than
      --list-order                list user defined ranks in order, from "$HOME/.taxonkit/ranks.txt"
      --list-ranks                list ordered ranks in taxonomy database, sorted in user defined order
  -L, --lower-than string         output taxIDs with rank lower than a rank, exclusive with --higher-than
  -r, --rank-file string          user-defined ordered taxonomic ranks, type "taxonkit filter --help" for details
      --root-taxid uint32         root taxid (default 1)
  -n, --save-predictable-norank   do not discard some special ranks without order when using -L, where rank of the closest higher node is still lower than rank cutoff
  -i, --taxid-field int           field index of taxid. input data should be tab-separated (default 1)

```

Examples

1. Example data

        $ echo 349741 | taxonkit lineage -t | cut -f 3 | sed 's/;/\n/g' > taxids2.txt
        
        $ cat taxids2.txt
        131567
        2
        1783257
        74201
        203494
        48461
        1647988
        239934
        239935
        349741
        
        
        $ cat taxids2.txt  | taxonkit lineage -r | csvtk -Ht cut -f 1,3,2 | csvtk pretty -t
        131567    no rank        cellular organisms
        2         superkingdom   cellular organisms;Bacteria
        1783257   clade          cellular organisms;Bacteria;PVC group
        74201     phylum         cellular organisms;Bacteria;PVC group;Verrucomicrobia
        203494    class          cellular organisms;Bacteria;PVC group;Verrucomicrobia;Verrucomicrobiae
        48461     order          cellular organisms;Bacteria;PVC group;Verrucomicrobia;Verrucomicrobiae;Verrucomicrobiales
        1647988   family         cellular organisms;Bacteria;PVC group;Verrucomicrobia;Verrucomicrobiae;Verrucomicrobiales;Akkermansiaceae
        239934    genus          cellular organisms;Bacteria;PVC group;Verrucomicrobia;Verrucomicrobiae;Verrucomicrobiales;Akkermansiaceae;Akkermansia
        239935    species        cellular organisms;Bacteria;PVC group;Verrucomicrobia;Verrucomicrobiae;Verrucomicrobiales;Akkermansiaceae;Akkermansia;Akkermansia muciniphila
        349741    strain         cellular organisms;Bacteria;PVC group;Verrucomicrobia;Verrucomicrobiae;Verrucomicrobiales;Akkermansiaceae;Akkermansia;Akkermansia muciniphila;Akkermansia muciniphila ATCC BAA-835

1. Equal to certain rank(s) (`-E/--equal-to`)

        $ cat taxids2.txt \
            | taxonkit filter -E Phylum -E Class \
            | taxonkit lineage -r \
            | csvtk -Ht cut -f 1,3,2 \
            | csvtk pretty -t
        74201    phylum   cellular organisms;Bacteria;PVC group;Verrucomicrobia
        203494   class    cellular organisms;Bacteria;PVC group;Verrucomicrobia;Verrucomicrobiae
            
1. Lower than a rank (`-L/--lower-than`)

        $ cat taxids2.txt \
            | taxonkit filter -L genus \
            | taxonkit lineage -r -n -L \
            | csvtk -Ht cut -f 1,3,2 \
            | csvtk pretty -t
        239935   species   Akkermansia muciniphila
        349741   strain    Akkermansia muciniphila ATCC BAA-835

1. Higher than a rank (`-H/--higher-than`)

        $ cat taxids2.txt \
            | taxonkit filter -H phylum \
            | taxonkit lineage -r -n -L \
            | csvtk -Ht cut -f 1,3,2 \
            | csvtk pretty -t
        2   superkingdom   Bacteria

1. "no rank" and "clade" have no rank and can be filter out via `-N/--discard-noranks`.
  Futher ranks can be removed with black list via `-B/--black-list`

        # 562 is taxID of Escherichia coli
        $ taxonkit list --ids 562 \
            | taxonkit filter -L species \
            | taxonkit lineage -r -n -L \
            | csvtk cut -Ht -f 1,3,2 \
            | csvtk freq -Ht -f 2 -nr \
            | csvtk pretty -t
        strain       2950
        no rank      149
        serotype     141
        serogroup    95
        isolate      1
        subspecies   1
        
        $ taxonkit list --ids 562 \
            | taxonkit filter -L species -N -B strain \
            | taxonkit lineage -r -n -L \
            | csvtk cut -Ht -f 1,3,2 \
            | csvtk freq -Ht -f 2 -nr \
            | csvtk pretty -t
        serotype     141
        serogroup    95
        isolate      1
        subspecies   1
        
1. Combine of `-L/-H` with `-E`.

        $ cat taxids2.txt \
            | taxonkit filter -L genus -E genus  \
            | taxonkit lineage -r -n -L \
            | csvtk -Ht cut -f 1,3,2 \
            | csvtk pretty -t
        239934   genus     Akkermansia
        239935   species   Akkermansia muciniphila
        349741   strain    Akkermansia muciniphila ATCC BAA-835

1. **Special cases of "no rank"**. (`-n/--save-predictable-norank`).
    When filtering with `-L/--lower-than`, you can use
    `-n/--save-predictable-norank` to save some special ranks without order,
    where rank of the closest higher node is still lower than rank cutoff.

        $ echo -ne "2605619\n1327037\n" \
            | taxonkit lineage -t \
            | csvtk cut -Ht -f 3 \
            | csvtk unfold -Ht -f 1 -s ";" \
            | taxonkit lineage -r -n -L \
            | csvtk cut -Ht -f 1,3,2 \
            | csvtk pretty -t 
        131567    no rank        cellular organisms
        2         superkingdom   Bacteria
        1224      phylum         Proteobacteria
        1236      class          Gammaproteobacteria
        91347     order          Enterobacterales
        543       family         Enterobacteriaceae
        561       genus          Escherichia
        562       species        Escherichia coli
        2605619   no rank        Escherichia coli O16:H48
        
        10239     superkingdom   Viruses
        2731341   clade          Duplodnaviria
        2731360   clade          Heunggongvirae
        2731618   phylum         Uroviricota
        2731619   class          Caudoviricetes
        28883     order          Caudovirales
        10699     family         Siphoviridae
        196894    no rank        unclassified Siphoviridae
        1327037   species        Croceibacter phage P2559Y
    
        # save taxids
        $  echo -ne "2605619\n1327037\n" \
            | taxonkit lineage -t \
            | csvtk cut -Ht -f 3 \
            | csvtk unfold -Ht -f 1 -s ";" \
            | tee taxids4.txt
        131567
        2
        1224
        1236
        91347
        543
        561
        562
        2605619
        10239
        2731341
        2731360
        2731618
        2731619
        28883
        10699
        196894
        1327037

    Now, filter nodes of rank <= species.
        
        $ cat taxids4.txt \
            | taxonkit filter -L species -E species -N -n \
            | taxonkit lineage -r -n -L \
            | csvtk cut -Ht -f 1,3,2 \
            | csvtk pretty -t
        562       species   Escherichia coli
        2605619   no rank   Escherichia coli O16:H48
        1327037   species   Croceibacter phage P2559Y

    Note that 2605619 (no rank) is saved because its parent node 562 is <= species.

## lca

Usage

```text
Compute lowest common ancestor (LCA) for taxIDs

Attention:

  1. This command computes LCA taxID for a list of taxIDs 
     in a field ("-i/--taxids-field) of tab-delimited file or STDIN.
  2. TaxIDs should have the same separater ("-s/--separater"),
     single charactor separater is prefered.
  3. Empty lines or lines without valid taxIDs in the field are omitted.
  4. If some taxIDs are not found in database, it returns 0.
  
Examples:

    $ echo 239934, 239935, 349741 | taxonkit lca  -s ", "
    239934, 239935, 349741  239934

    $ echo 239934  239935  349741 9606  | taxonkit lca
    239934 239935 349741 9606       131567

Usage:
  taxonkit lca [flags]

Flags:
  -h, --help               help for lca
  -s, --separater string   separater for taxIDs (default " ")
  -D, --skip-deleted       skip deleted taxIDs and compute with left ones
  -U, --skip-unfound       skip unfound taxIDs and compute with left ones
  -i, --taxids-field int   field index of taxid. input data should be tab-separated (default 1)

```

Examples:

1. Example data

        $ taxonkit list --ids 9605 -nr --indent "    "
        9605 [genus] Homo
            9606 [species] Homo sapiens
                63221 [subspecies] Homo sapiens neanderthalensis
                741158 [subspecies] Homo sapiens subsp. 'Denisova'
            1425170 [species] Homo heidelbergensis
            2665952 [no rank] environmental samples
                2665953 [species] Homo sapiens environmental sample

1. Simple one

        $ echo 63221 2665953 | taxonkit lca
        63221 2665953   9605
        
1. Custom field (`-i/--taxids-field`) and separater (`-s/--separater`).

        $ echo -ne "a\t63221,2665953\nb\t63221, 741158\n"
        a       63221,2665953
        b       63221, 741158
        
        $ echo -ne "a\t63221,2665953\nb\t63221, 741158\n" \
            | taxonkit lca -i 2 -s ","
        a       63221,2665953   9605
        b       63221, 741158   9606
        
1. Merged taxIDs.

        # merged
        $ echo 92487 92488 92489 | taxonkit lca
        10:08:26.578 [WARN] taxid 92489 was merged into 796334
        92487 92488 92489       1236
        
1. Deleted taxIDs, you can ommit theses and continue compute with left onces with (`-D/--skip-deleted`).

        $ echo 1 2 3 | taxonkit lca 
        10:30:17.678 [WARN] taxid 3 not found
        1 2 3   0
        
        $ time  echo 1 2 3 | taxonkit lca -D
        10:29:31.828 [WARN] taxid 3 was deleted
        1 2 3   1

1. TaxIDs not found in database, you can ommit theses and continue compute with left onces with (`-U/--skip-unfound`).

        $ echo 61021 61022 11111111 | taxonkit lca
        10:31:44.929 [WARN] taxid 11111111 not found
        61021 61022 11111111    0
        
        $ echo 61021 61022 11111111 | taxonkit lca -U
        10:32:02.772 [WARN] taxid 11111111 not found
        61021 61022 11111111    2628496

## taxid-changelog

Usage

```text
Create taxID changelog from dump archives

Steps:

    # dependencies:
    #   rush - https://github.com/shenwei356/rush/

    mkdir -p archive; cd archive;

    # --------- download ---------

    # option 1
    # for fast network connection
    wget ftp://ftp.ncbi.nlm.nih.gov/pub/taxonomy/taxdump_archive/taxdmp*.zip

    # option 2
    # for slow network connection
    url=https://ftp.ncbi.nlm.nih.gov/pub/taxonomy/taxdump_archive/
    wget $url -O - -o /dev/null \
        | grep taxdmp | perl -ne '/(taxdmp_.+?.zip)/; print "$1\n";' \
        | rush -j 2 -v url=$url 'axel -n 5 {url}/{}' \
            --immediate-output  -c -C download.rush

    # --------- unzip ---------

    ls taxdmp*.zip | rush -j 1 'unzip {} names.dmp nodes.dmp merged.dmp delnodes.dmp -d {@_(.+)\.}'

    # --------- create log ---------

    cd ..
    taxonkit taxid-changelog -i archive -o taxid-changelog.csv.gz --verbose

Output format (CSV):

    # fields        comments
    taxid           # taxid
    version         # version / time of archive, e.g, 2019-07-01
    change          # change, values:
                    #   NEW             newly added
                    #   REUSE_DEL       deleted taxids being reused
                    #   REUSE_MER       merged taxids being reused
                    #   DELETE          deleted
                    #   MERGE           merged into another taxid
                    #   ABSORB          other taxids merged into this one
                    #   CHANGE_NAME     scientific name changed
                    #   CHANGE_RANK     rank changed
                    #   CHANGE_LIN_LIN  lineage taxids remain but lineage remain
                    #   CHANGE_LIN_TAX  lineage taxids changed
                    #   CHANGE_LIN_LEN  lineage length changed
    change-value    # variable values for changes: 
                    #   1) new taxid for MERGE
                    #   2) merged taxids for ABSORB
                    #   3) empty for others
    name            # scientific name
    rank            # rank
    lineage         # full lineage of the taxid
    lineage-taxids  # taxids of the lineage

    # you can use csvtk to investigate them. e.g.,
    csvtk grep -f taxid -p 1390515 taxid-changelog.csv.gz

Usage:
  taxonkit taxid-changelog [flags]

Flags:
  -i, --archive string   directory containing decompressed dumped archives
  -h, --help             help for taxid-changelog
  
```

[Details](https://github.com/shenwei356/taxid-changelog)

1. Example 1  (*E.coli* with taxid `562`)

        $ pigz -cd taxid-changelog.csv.gz \
            | csvtk grep -f taxid -p 562 \
            | csvtk pretty
        taxid   version      change           change-value    name               rank      lineage                                                                                                                            lineage-taxids
        562     2014-08-01   NEW                              Escherichia coli   species   cellular organisms;Bacteria;Proteobacteria;Gammaproteobacteria;Enterobacteriales;Enterobacteriaceae;Escherichia;Escherichia coli   131567;2;1224;1236;91347;543;561;562
        562     2014-08-01   ABSORB           662101;662104   Escherichia coli   species   cellular organisms;Bacteria;Proteobacteria;Gammaproteobacteria;Enterobacteriales;Enterobacteriaceae;Escherichia;Escherichia coli   131567;2;1224;1236;91347;543;561;562
        562     2015-11-01   ABSORB           1637691         Escherichia coli   species   cellular organisms;Bacteria;Proteobacteria;Gammaproteobacteria;Enterobacteriales;Enterobacteriaceae;Escherichia;Escherichia coli   131567;2;1224;1236;91347;543;561;562
        562     2016-10-01   CHANGE_LIN_LIN                   Escherichia coli   species   cellular organisms;Bacteria;Proteobacteria;Gammaproteobacteria;Enterobacterales;Enterobacteriaceae;Escherichia;Escherichia coli    131567;2;1224;1236;91347;543;561;562
        562     2018-06-01   ABSORB           469598          Escherichia coli   species   cellular organisms;Bacteria;Proteobacteria;Gammaproteobacteria;Enterobacterales;Enterobacteriaceae;Escherichia;Escherichia coli    131567;2;1224;1236;91347;543;561;562

        # merged taxids
        $ pigz -cd taxid-changelog.csv.gz \
            | csvtk grep -f taxid -p 662101,662104,1637691,469598 \
            | csvtk pretty
        taxid     version      change           change-value   name                        rank      lineage                                                                                                                                     lineage-taxids
        469598    2014-08-01   NEW                             Escherichia sp. 3_2_53FAA   species   cellular organisms;Bacteria;Proteobacteria;Gammaproteobacteria;Enterobacteriales;Enterobacteriaceae;Escherichia;Escherichia sp. 3_2_53FAA   131567;2;1224;1236;91347;543;561;469598
        469598    2016-10-01   CHANGE_LIN_LIN                  Escherichia sp. 3_2_53FAA   species   cellular organisms;Bacteria;Proteobacteria;Gammaproteobacteria;Enterobacterales;Enterobacteriaceae;Escherichia;Escherichia sp. 3_2_53FAA    131567;2;1224;1236;91347;543;561;469598
        469598    2018-06-01   MERGE            562            Escherichia sp. 3_2_53FAA   species   cellular organisms;Bacteria;Proteobacteria;Gammaproteobacteria;Enterobacterales;Enterobacteriaceae;Escherichia;Escherichia sp. 3_2_53FAA    131567;2;1224;1236;91347;543;561;469598
        662101    2014-08-01   MERGE            562                                                                                                                                                                                              
        662104    2014-08-01   MERGE            562                                                                                                                                                                                              
        1637691   2015-04-01   DELETE                                                                                                                                                                                                            
        1637691   2015-05-01   REUSE_DEL                       Escherichia sp. MAR         species   cellular organisms;Bacteria;Proteobacteria;Gammaproteobacteria;Enterobacteriales;Enterobacteriaceae;Escherichia;Escherichia sp. MAR         131567;2;1224;1236;91347;543;561;1637691
        1637691   2015-11-01   MERGE            562            Escherichia sp. MAR         species   cellular organisms;Bacteria;Proteobacteria;Gammaproteobacteria;Enterobacteriales;Enterobacteriaceae;Escherichia;Escherichia sp. MAR         131567;2;1224;1236;91347;543;561;1637691

1. Example 2 (SARS-CoV-2).

        $ time pigz -cd taxid-changelog.csv.gz \
                | csvtk grep -f taxid -p 2697049 \
                | csvtk pretty
        taxid     version      change           change-value   name                                              rank      lineage                                                                                                                                                                                                                                                        lineage-taxids
        2697049   2020-02-01   NEW                             Wuhan seafood market pneumonia virus              species   Viruses;Riboviria;Nidovirales;Cornidovirineae;Coronaviridae;Orthocoronavirinae;Betacoronavirus;unclassified Betacoronavirus;Wuhan seafood market pneumonia virus                                                                                               10239;2559587;76804;2499399;11118;2501931;694002;696098;2697049
        2697049   2020-03-01   CHANGE_NAME                     Severe acute respiratory syndrome coronavirus 2   no rank   Viruses;Riboviria;Nidovirales;Cornidovirineae;Coronaviridae;Orthocoronavirinae;Betacoronavirus;Sarbecovirus;Severe acute respiratory syndrome-related coronavirus;Severe acute respiratory syndrome coronavirus 2                                              10239;2559587;76804;2499399;11118;2501931;694002;2509511;694009;2697049
        2697049   2020-03-01   CHANGE_RANK                     Severe acute respiratory syndrome coronavirus 2   no rank   Viruses;Riboviria;Nidovirales;Cornidovirineae;Coronaviridae;Orthocoronavirinae;Betacoronavirus;Sarbecovirus;Severe acute respiratory syndrome-related coronavirus;Severe acute respiratory syndrome coronavirus 2                                              10239;2559587;76804;2499399;11118;2501931;694002;2509511;694009;2697049
        2697049   2020-03-01   CHANGE_LIN_LEN                  Severe acute respiratory syndrome coronavirus 2   no rank   Viruses;Riboviria;Nidovirales;Cornidovirineae;Coronaviridae;Orthocoronavirinae;Betacoronavirus;Sarbecovirus;Severe acute respiratory syndrome-related coronavirus;Severe acute respiratory syndrome coronavirus 2                                              10239;2559587;76804;2499399;11118;2501931;694002;2509511;694009;2697049
        2697049   2020-06-01   CHANGE_LIN_LEN                  Severe acute respiratory syndrome coronavirus 2   no rank   Viruses;Riboviria;Orthornavirae;Pisuviricota;Pisoniviricetes;Nidovirales;Cornidovirineae;Coronaviridae;Orthocoronavirinae;Betacoronavirus;Sarbecovirus;Severe acute respiratory syndrome-related coronavirus;Severe acute respiratory syndrome coronavirus 2   10239;2559587;2732396;2732408;2732506;76804;2499399;11118;2501931;694002;2509511;694009;2697049
        2697049   2020-07-01   CHANGE_RANK                     Severe acute respiratory syndrome coronavirus 2   isolate   Viruses;Riboviria;Orthornavirae;Pisuviricota;Pisoniviricetes;Nidovirales;Cornidovirineae;Coronaviridae;Orthocoronavirinae;Betacoronavirus;Sarbecovirus;Severe acute respiratory syndrome-related coronavirus;Severe acute respiratory syndrome coronavirus 2   10239;2559587;2732396;2732408;2732506;76804;2499399;11118;2501931;694002;2509511;694009;2697049
        2697049   2020-08-01   CHANGE_RANK                     Severe acute respiratory syndrome coronavirus 2   no rank   Viruses;Riboviria;Orthornavirae;Pisuviricota;Pisoniviricetes;Nidovirales;Cornidovirineae;Coronaviridae;Orthocoronavirinae;Betacoronavirus;Sarbecovirus;Severe acute respiratory syndrome-related coronavirus;Severe acute respiratory syndrome coronavirus 2   10239;2559587;2732396;2732408;2732506;76804;2499399;11118;2501931;694002;2509511;694009;2697049
        
        real    0m7.644s
        user    0m16.749s
        sys     0m3.985s

2. Example 3 (All subspecies and strain in *Akkermansia muciniphila* 239935)

        # species in Akkermansia
        $ taxonkit list --show-rank --show-name --indent "    " --ids 239935
        239935 [species] Akkermansia muciniphila
            349741 [strain] Akkermansia muciniphila ATCC BAA-835
        
        # check them all  
        $ pigz -cd taxid-changelog.csv.gz \
            | csvtk grep -f taxid -P <(taxonkit list --indent "" --ids 239935) \
            | csvtk pretty                                                                                                                                                   lineage-taxids
        taxid    version      change           change-value   name                                   rank      lineage                                                                                                                                                                                                         lineage-taxids
        239935   2014-08-01   NEW                             Akkermansia muciniphila                species   cellular organisms;Bacteria;Chlamydiae/Verrucomicrobia group;Verrucomicrobia;Verrucomicrobiae;Verrucomicrobiales;Verrucomicrobiaceae;Akkermansia;Akkermansia muciniphila                                        131567;2;51290;74201;203494;48461;203557;239934;239935
        239935   2015-05-01   CHANGE_LIN_TAX                  Akkermansia muciniphila                species   cellular organisms;Bacteria;Chlamydiae/Verrucomicrobia group;Verrucomicrobia;Verrucomicrobiae;Verrucomicrobiales;Akkermansiaceae;Akkermansia;Akkermansia muciniphila                                            131567;2;51290;74201;203494;48461;1647988;239934;239935
        239935   2016-03-01   CHANGE_LIN_TAX                  Akkermansia muciniphila                species   cellular organisms;Bacteria;PVC group;Verrucomicrobia;Verrucomicrobiae;Verrucomicrobiales;Akkermansiaceae;Akkermansia;Akkermansia muciniphila                                                                   131567;2;1783257;74201;203494;48461;1647988;239934;239935
        239935   2016-05-01   ABSORB           1834199        Akkermansia muciniphila                species   cellular organisms;Bacteria;PVC group;Verrucomicrobia;Verrucomicrobiae;Verrucomicrobiales;Akkermansiaceae;Akkermansia;Akkermansia muciniphila                                                                   131567;2;1783257;74201;203494;48461;1647988;239934;239935
        349741   2014-08-01   NEW                             Akkermansia muciniphila ATCC BAA-835   no rank   cellular organisms;Bacteria;Chlamydiae/Verrucomicrobia group;Verrucomicrobia;Verrucomicrobiae;Verrucomicrobiales;Verrucomicrobiaceae;Akkermansia;Akkermansia muciniphila;Akkermansia muciniphila ATCC BAA-835   131567;2;51290;74201;203494;48461;203557;239934;239935;349741
        349741   2015-05-01   CHANGE_LIN_TAX                  Akkermansia muciniphila ATCC BAA-835   no rank   cellular organisms;Bacteria;Chlamydiae/Verrucomicrobia group;Verrucomicrobia;Verrucomicrobiae;Verrucomicrobiales;Akkermansiaceae;Akkermansia;Akkermansia muciniphila;Akkermansia muciniphila ATCC BAA-835       131567;2;51290;74201;203494;48461;1647988;239934;239935;349741
        349741   2016-03-01   CHANGE_LIN_TAX                  Akkermansia muciniphila ATCC BAA-835   no rank   cellular organisms;Bacteria;PVC group;Verrucomicrobia;Verrucomicrobiae;Verrucomicrobiales;Akkermansiaceae;Akkermansia;Akkermansia muciniphila;Akkermansia muciniphila ATCC BAA-835                              131567;2;1783257;74201;203494;48461;1647988;239934;239935;349741
        349741   2020-07-01   CHANGE_RANK                     Akkermansia muciniphila ATCC BAA-835   strain    cellular organisms;Bacteria;PVC group;Verrucomicrobia;Verrucomicrobiae;Verrucomicrobiales;Akkermansiaceae;Akkermansia;Akkermansia muciniphila;Akkermansia muciniphila ATCC BAA-835                              131567;2;1783257;74201;203494;48461;1647988;239934;239935;349741

[More](https://github.com/shenwei356/taxid-changelog)
    
## genautocomplete

Usage

```text
Generate shell autocompletion script

Note: The current version supports Bash only.
This should work for *nix systems with Bash installed.

Howto:

1. run: taxonkit genautocomplete

2. create and edit ~/.bash_completion file if you don't have it.

        nano ~/.bash_completion

   add the following:

        for bcfile in ~/.bash_completion.d/* ; do
          . $bcfile
        done

Usage:
  taxonkit genautocomplete [flags]

Flags:
      --file string   autocompletion file (default "/home/shenwei/.bash_completion.d/taxonkit.sh")
  -h, --help          help for genautocomplete
      --type string   autocompletion type (currently only bash supported) (default "bash")

```

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

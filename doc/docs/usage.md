# Usage and Examples

## Before use

Please download and uncompress these files:

- [ftp://ftp.ncbi.nih.gov/pub/taxonomy/taxdump.tar.gz](ftp://ftp.ncbi.nih.gov/pub/taxonomy/taxdump.tar.gz)

And copy "names.dmp" and "nodes.dmp" to data directory: "$HOME/.taxonkit".

## taxonkit

```
TaxonKit - Cross-platform and Efficient NCBI Taxonomy Toolkit

Version: 0.2.4

Author: Wei Shen <shenwei356@gmail.com>

Source code: https://github.com/shenwei356/taxonkit
Documents  : http://bioinf.shenwei.me/taxonkit

Dataset:

    Please download and decompress "taxdump.tar.gz":
    ftp://ftp.ncbi.nih.gov/pub/taxonomy/taxdump.tar.gz

    and copy "names.dmp" and "nodes.dmp" to data directory:
    "/home/shenwei/.taxonkit"

Usage:
  taxonkit [command]

Available Commands:
  genautocomplete generate shell autocompletion script
  help            Help about any command
  lineage         query lineage of given taxids
  list            list taxon tree of given taxids
  name2taxid      query taxid by taxon scientific name
  reformat        reformat lineage
  version         print version information and check for update

Flags:
  -h, --help                help for taxonkit
      --names-file string   names.dmp file (default "/home/shenwei/.taxonkit/names.dmp")
      --nodes-file string   nodes.dmp file (default "/home/shenwei/.taxonkit/nodes.dmp")
  -o, --out-file string     out file ("-" for stdout, suffix .gz for gzipped out) (default "-")
  -j, --threads int         number of CPUs. (default value: 1 for single-CPU PC, 2 for others) (default 2)
      --verbose             print verbose information

Use "taxonkit [command] --help" for more information about a command.

```

## list

Usage

```
list taxon tree of given taxids

Usage:
  taxonkit list [flags]

Flags:
      --ids string      taxid(s), multiple values should be separated by comma
      --indent string   indent (default "  ")
      --json            output in JSON format. you can save the result in file with suffix ".json" and open with modern text editor
      --show-name       output scientific name
      --show-rank       output rank

```

Examples

1. Default usage

        $ taxonkit list --ids 9605,239934
        9605
          9606
            63221
            741158
          1425170

        239934
          239935
            349741
          512293
            512294
            1131822
            1262691
            1263034
          1131336
          1574264
          1574265
          1638783
          1679444
          1755639
          1896967

1. Removing indent. The list could be used to extract sequences from BLAST database with `blastdbcmd` (see [tutorial](http://bioinf.shenwei.me/taxonkit/tutorial/))

        $ taxonkit list --ids 9605,239934 --indent ""
        9605
        9606
        63221
        741158
        1425170

        239934
        239935
        349741
        512293
        512294
        1131822
        1262691
        1263034
        1131336
        1574264
        1574265
        1638783
        1679444
        1755639
        1896967


    **Performance:** Time and memory usage for whole taxon tree:

        $ # emptying the buffers cache
        $ su -c "free && sync && echo 3 > /proc/sys/vm/drop_caches && free"

        $ memusg -t taxonkit list --ids 1 --indent "" > t0.txt
        elapsed time: 3.093s
        peak rss: 82.79 MB

1. Adding names

        $ taxonkit list --show-rank --show-name --ids 9605,239934
        9605 [genus] Homo
          9606 [species] Homo sapiens
            63221 [subspecies] Homo sapiens neanderthalensis
            741158 [subspecies] Homo sapiens ssp. Denisova
          1425170 [species] Homo heidelbergensis

        239934 [genus] Akkermansia
          239935 [species] Akkermansia muciniphila
            349741 [no rank] Akkermansia muciniphila ATCC BAA-835
          512293 [no rank] environmental samples
            512294 [species] uncultured Akkermansia sp.
            1131822 [species] uncultured Akkermansia sp. SMG25
            1262691 [species] Akkermansia sp. CAG:344
            1263034 [species] Akkermansia muciniphila CAG:154
          1131336 [species] Akkermansia sp. KLE1605
          1574264 [species] Akkermansia sp. KLE1797
          1574265 [species] Akkermansia sp. KLE1798
          1638783 [species] Akkermansia sp. UNK.MGS-1
          1679444 [species] Akkermansia glycaniphila
          1755639 [species] Akkermansia sp. MC_55
          1896967 [species] Akkermansia sp. 54_46

    **Performance:** Time and memory usage for whole taxon tree:

        $ # emptying the buffers cache
        $ su -c "free && sync && echo 3 > /proc/sys/vm/drop_caches && free"

        $ memusg -t taxonkit list --show-rank --show-name --ids 1 > t1.txt
        elapsed time: 9.825s
        peak rss: 648.65 MB

1. Output in JSON format, so you can easily collapse and uncollapse taxonomy tree in modern text editor.

        $ taxonkit list --show-rank --show-name --ids 9605,239934 --json
        {
          "9605 [genus] Homo": {
            "9606 [species] Homo sapiens": {
              "63221 [subspecies] Homo sapiens neanderthalensis": {},
              "741158 [subspecies] Homo sapiens ssp. Denisova": {}
            }
            "1425170 [species] Homo heidelbergensis": {}
          },
          "239934 [genus] Akkermansia": {
            "239935 [species] Akkermansia muciniphila": {
              "349741 [no rank] Akkermansia muciniphila ATCC BAA-835": {}
            }
            "512293 [no rank] environmental samples": {
              "512294 [species] uncultured Akkermansia sp.": {},
              "1131822 [species] uncultured Akkermansia sp. SMG25": {},
              "1262691 [species] Akkermansia sp. CAG:344": {},
              "1263034 [species] Akkermansia muciniphila CAG:154": {}
            }
            "1131336 [species] Akkermansia sp. KLE1605": {},
            "1574264 [species] Akkermansia sp. KLE1797": {},
            "1574265 [species] Akkermansia sp. KLE1798": {},
            "1638783 [species] Akkermansia sp. UNK.MGS-1": {},
            "1679444 [species] Akkermansia glycaniphila": {},
            "1755639 [species] Akkermansia sp. MC_55": {},
            "1896967 [species] Akkermansia sp. 54_46": {}
          }
        }

    Snapshot of taxonomy (taxid 1) in kate:
    ![taxon.json.png](files/taxon.json.png)


## lineage

Usage

```
query lineage of given taxids

Usage:
  taxonkit lineage [flags]

Flags:
  -d, --delimiter string      field delimiter in lineage (default ";")
  -h, --help                  help for lineage
  -t, --show-lineage-taxids   show lineage consisting of taxids
  -i, --taxid-field int       field index of taxid. data should be tab-separated (default 1)

```

Examples

1. Full lineage:

        # note that 10000000 is a fake taxid
        $ cat taxids.txt
        9606
        349741
        239935
        11932
        314101
        1327037
        10000000


        $ taxonkit lineage -t taxids.txt
        9606    cellular organisms;Eukaryota;Opisthokonta;Metazoa;Eumetazoa;Bilateria;Deuterostomia;Chordata;Craniata;Vertebrata;Gnathostomata;Teleostomi;Euteleostomi;Sarcopterygii;Dipnotetrapodomorpha;Tetrapoda;Amniota;Mammalia;Theria;Eutheria;Boreoeutheria;Euarchontoglires;Primates;Haplorrhini;Simiiformes;Catarrhini;Hominoidea;Hominidae;Homininae;Homo;Homo sapiens
        349741  cellular organisms;Bacteria;PVC group;Verrucomicrobia;Verrucomicrobiae;Verrucomicrobiales;Akkermansiaceae;Akkermansia;Akkermansia muciniphila;Akkermansia muciniphila ATCC BAA-835
        239935  cellular organisms;Bacteria;PVC group;Verrucomicrobia;Verrucomicrobiae;Verrucomicrobiales;Akkermansiaceae;Akkermansia;Akkermansia muciniphila
        11932   Viruses;Retro-transcribing viruses;Retroviridae;unclassified Retroviridae;Intracisternal A-particles;Mouse Intracisternal A-particle
        314101  cellular organisms;Bacteria;environmental samples;uncultured murine large bowel bacterium BAC 54B
        1327037 Viruses;dsDNA viruses, no RNA stage;Caudovirales;Siphoviridae;unclassified Siphoviridae;Croceibacter phage P2559Y
        10000000

1. show lineage consisting of taxids:

        $ taxonkit lineage -t taxids.txt
        9606    cellular organisms;Eukaryota;Opisthokonta;Metazoa;Eumetazoa;Bilateria;Deuterostomia;Chordata;Craniata;Vertebrata;Gnathostomata;Teleostomi;Euteleostomi;Sarcopterygii;Dipnotetrapodomorpha;Tetrapoda;Amniota;Mammalia;Theria;Eutheria;Boreoeutheria;Euarchontoglires;Primates;Haplorrhini;Simiiformes;Catarrhini;Hominoidea;Hominidae;Homininae;Homo;Homo sapiens    131567;2759;33154;33208;6072;33213;33511;7711;89593;7742;7776;117570;117571;8287;1338369;32523;32524;40674;32525;9347;1437010;314146;9443;376913;314293;9526;314295;9604;207598;9605;9606
        349741  cellular organisms;Bacteria;PVC group;Verrucomicrobia;Verrucomicrobiae;Verrucomicrobiales;Akkermansiaceae;Akkermansia;Akkermansia muciniphila;Akkermansia muciniphila ATCC BAA-835  131567;2;1783257;74201;203494;48461;1647988;239934;239935;349741
        239935  cellular organisms;Bacteria;PVC group;Verrucomicrobia;Verrucomicrobiae;Verrucomicrobiales;Akkermansiaceae;Akkermansia;Akkermansia muciniphila       131567;2;1783257;74201;203494;48461;1647988;239934;239935
        11932   Viruses;Retro-transcribing viruses;Retroviridae;unclassified Retroviridae;Intracisternal A-particles;Mouse Intracisternal A-particle        10239;35268;11632;35276;11749;11932
        314101  cellular organisms;Bacteria;environmental samples;uncultured murine large bowel bacterium BAC 54B   131567;2;48479;314101
        1327037 Viruses;dsDNA viruses, no RNA stage;Caudovirales;Siphoviridae;unclassified Siphoviridae;Croceibacter phage P2559Y   10239;35237;28883;10699;196894;1327037
        10000000

    or read taxids from STDIN:

        $ cat taxids.txt | taxonkit lineage

## reformat

Usage

```
reformat lineage

Output format can be formated by flag --format, available placeholders:

    {k}: superkingdom
    {p}: phylum
    {c}: class
    {o}: order
    {f}: family
    {g}: genus
    {s}: species
    {S}: subspecies

Output format can contains some escape charactors like "\t".

This command appends reformated lineage to the input line.
The corresponding taxids of reformated lineage can be provided as another
column by flag "-t/--show-lineage-taxids".

Usage:
  taxonkit reformat [flags]

Flags:
  -d, --delimiter string               field delimiter in input lineage (default ";")
  -F, --fill-miss-rank                 estimate and fill missing rank with original lineage information (recommended)
  -f, --format string                  output format, placeholders of rank are needed (default "{k};{p};{c};{o};{f};{g};{s}")
  -h, --help                           help for reformat
  -i, --lineage-field int              field index of lineage. data should be tab-separated (default 2)
  -r, --miss-rank-repl string          replacement string for missing rank, if given "", "unclassified xxx xxx" will used, where "unclassified " is settable by flag -p/--miss-rank-repl-prefix
  -p, --miss-rank-repl-prefix string   prefix for estimated taxon level (default "unclassified ")
  -R, --miss-taxid-repl string         replacement string for missing taxid
  -t, --show-lineage-taxids            show corresponding taxids of reformated lineage

```

Examples:

1. Example lineage (produced by: `taxonkit lineage taxids.txt > lineage.txt`)

        $ cat lineage.txt
        9606    cellular organisms;Eukaryota;Opisthokonta;Metazoa;Eumetazoa;Bilateria;Deuterostomia;Chordata;Craniata;Vertebrata;Gnathostomata;Teleostomi;Euteleostomi;Sarcopterygii;Dipnotetrapodomorpha;Tetrapoda;Amniota;Mammalia;Theria;Eutheria;Boreoeutheria;Euarchontoglires;Primates;Haplorrhini;Simiiformes;Catarrhini;Hominoidea;Hominidae;Homininae;Homo;Homo sapiens
        349741  cellular organisms;Bacteria;PVC group;Verrucomicrobia;Verrucomicrobiae;Verrucomicrobiales;Akkermansiaceae;Akkermansia;Akkermansia muciniphila;Akkermansia muciniphila ATCC BAA-835
        239935  cellular organisms;Bacteria;PVC group;Verrucomicrobia;Verrucomicrobiae;Verrucomicrobiales;Akkermansiaceae;Akkermansia;Akkermansia muciniphila
        11932   Viruses;Retro-transcribing viruses;Retroviridae;unclassified Retroviridae;Intracisternal A-particles;Mouse Intracisternal A-particle
        314101  cellular organisms;Bacteria;environmental samples;uncultured murine large bowel bacterium BAC 54B
        1327037 Viruses;dsDNA viruses, no RNA stage;Caudovirales;Siphoviridae;unclassified Siphoviridae;Croceibacter phage P2559Y

1. Default output format (`"{k};{p};{c};{o};{f};{g};{s}"`)

        # reformated lineages are appended to the input data
        $ taxonkit reformat --lineage-field 2 lineage.txt
        $ taxonkit reformat lineage.txt
        9606    cellular organisms;Eukaryota;Opisthokonta;Metazoa;Eumetazoa;Bilateria;Deuterostomia;Chordata;Craniata;Vertebrata;Gnathostomata;Teleostomi;Euteleostomi;Sarcopterygii;Dipnotetrapodomorpha;Tetrapoda;Amniota;Mammalia;Theria;Eutheria;Boreoeutheria;Euarchontoglires;Primates;Haplorrhini;Simiiformes;Catarrhini;Hominoidea;Hominidae;Homininae;Homo;Homo sapiens  Eukaryota;Chordata;Mammalia;Primates;Hominidae;Homo;Homo sapiens
        ...

        $ taxonkit reformat lineage.txt | cut -f 3
        Eukaryota;Chordata;Mammalia;Primates;Hominidae;Homo;Homo sapiens
        Bacteria;Verrucomicrobia;Verrucomicrobiae;Verrucomicrobiales;Akkermansiaceae;Akkermansia;Akkermansia muciniphila
        Bacteria;Verrucomicrobia;Verrucomicrobiae;Verrucomicrobiales;Akkermansiaceae;Akkermansia;Akkermansia muciniphila
        Viruses;;;;Retroviridae;Intracisternal A-particles;Mouse Intracisternal A-particle
        Bacteria;;;;;;uncultured murine large bowel bacterium BAC 54B
        Viruses;;;Caudovirales;Siphoviridae;;Croceibacter phage P2559Y

1. Show corresponding taxids of reformated lineage

        $ cat lineage.txt | taxonkit reformat -t
        9606    cellular organisms;Eukaryota;Opisthokonta;Metazoa;Eumetazoa;Bilateria;Deuterostomia;Chordata;Craniata;Vertebrata;Gnathostomata;Teleostomi;Euteleostomi;Sarcopterygii;Dipnotetrapodomorpha;Tetrapoda;Amniota;Mammalia;Theria;Eutheria;Boreoeutheria;Euarchontoglires;Primates;Haplorrhini;Simiiformes;Catarrhini;Hominoidea;Hominidae;Homininae;Homo;Homo sapiens       Eukaryota;Chordata;Mammalia;Primates;Hominidae;Homo;Homo sapiens     2759;7711;40674;9443;9604;9605;9606
        9913    cellular organisms;Eukaryota;Opisthokonta;Metazoa;Eumetazoa;Bilateria;Deuterostomia;Chordata;Craniata;Vertebrata;Gnathostomata;Teleostomi;Euteleostomi;Sarcopterygii;Dipnotetrapodomorpha;Tetrapoda;Amniota;Mammalia;Theria;Eutheria;Boreoeutheria;Laurasiatheria;Cetartiodactyla;Ruminantia;Pecora;Bovidae;Bovinae;Bos;Bos taurus     Eukaryota;Chordata;Mammalia;;Bovidae;Bos;Bos taurus     2759;7711;40674;;9895;9903;9913
        376619  cellular organisms;Bacteria;Proteobacteria;Gammaproteobacteria;Thiotrichales;Francisellaceae;Francisella;Francisella tularensis;Francisella tularensis subsp. holarctica;Francisella tularensis subsp. holarctica LVS     Bacteria;Proteobacteria;Gammaproteobacteria;Thiotrichales;Francisellaceae;Francisella;Francisella tularensis 2;1224;1236;72273;34064;262;263
        349741  cellular organisms;Bacteria;PVC group;Verrucomicrobia;Verrucomicrobiae;Verrucomicrobiales;Akkermansiaceae;Akkermansia;Akkermansia muciniphila;Akkermansia muciniphila ATCC BAA-835   Bacteria;Verrucomicrobia;Verrucomicrobiae;Verrucomicrobiales;Akkermansiaceae;Akkermansia;Akkermansia muciniphila 2;74201;203494;48461;1647988;239934;239935
        239935  cellular organisms;Bacteria;PVC group;Verrucomicrobia;Verrucomicrobiae;Verrucomicrobiales;Akkermansiaceae;Akkermansia;Akkermansia muciniphila        Bacteria;Verrucomicrobia;Verrucomicrobiae;Verrucomicrobiales;Akkermansiaceae;Akkermansia;Akkermansia muciniphila     2;74201;203494;48461;1647988;239934;239935
        314101  cellular organisms;Bacteria;environmental samples;uncultured murine large bowel bacterium BAC 54B    Bacteria;;;;;;uncultured murine large bowel bacterium BAC 54B   2;;;;;;314101
        11932   Viruses;Retro-transcribing viruses;Retroviridae;unclassified Retroviridae;Intracisternal A-particles;Mouse Intracisternal A-particle Viruses;;;;Retroviridae;Intracisternal A-particles;Mouse Intracisternal A-particle   10239;;;;11632;11749;11932
        1327037 Viruses;dsDNA viruses, no RNA stage;Caudovirales;Siphoviridae;unclassified Siphoviridae;Croceibacter phage P2559Y    Viruses;;;Caudovirales;Siphoviridae;;Croceibacter phage P2559Y  10239;;;28883;10699;;1327037


1. Use custom strings for unclassfied ranks

        $ taxonkit reformat lineage.txt --miss-rank-repl "__" | cut -f 3
        Eukaryota;Chordata;Mammalia;Primates;Hominidae;Homo;Homo sapiens
        Eukaryota;Chordata;Mammalia;__;Bovidae;Bos;Bos taurus
        Bacteria;Proteobacteria;Gammaproteobacteria;Thiotrichales;Francisellaceae;Francisella;Francisella tularensis
        Bacteria;Verrucomicrobia;Verrucomicrobiae;Verrucomicrobiales;Akkermansiaceae;Akkermansia;Akkermansia muciniphila
        Bacteria;Verrucomicrobia;Verrucomicrobiae;Verrucomicrobiales;Akkermansiaceae;Akkermansia;Akkermansia muciniphila
        Bacteria;__;__;__;__;__;uncultured murine large bowel bacterium BAC 54B
        Viruses;__;__;__;Retroviridae;Intracisternal A-particles;Mouse Intracisternal A-particle
        Viruses;__;__;Caudovirales;Siphoviridae;__;Croceibacter phage P2559Y


1. Estimate and fill missing rank with original lineage information
   (**recommended**, very useful for formating input data for
   [LEfSe](https://bitbucket.org/biobakery/biobakery/wiki/lefse))

       $ cat lineage.txt | taxonkit reformat -t -F > lineage.txt.reformat.fill
       $ cat lineage.txt.reformat.fill \
           | perl -pe 's/^/Taxid   : /; \
               s/\t/\nLineage : /; \
               s/\t/\nReformat: /; \
               s/\t/\nTaxids  : /; \
               print "\n";'

       Taxid   : 9606
       Lineage : cellular organisms;Eukaryota;Opisthokonta;Metazoa;Eumetazoa;Bilateria;Deuterostomia;Chordata;Craniata;Vertebrata;Gnathostomata;Teleostomi;Euteleostomi;Sarcopterygii;Dipnotetrapodomorpha;Tetrapoda;Amniota;Mammalia;Theria;Eutheria;Boreoeutheria;Euarchontoglires;Primates;Haplorrhini;Simiiformes;Catarrhini;Hominoidea;Hominidae;Homininae;Homo;Homo sapiens
       Reformat: Eukaryota;Chordata;Mammalia;Primates;Hominidae;Homo;Homo sapiens
       Taxids  : 2759;7711;40674;9443;9604;9605;9606

       Taxid   : 9913
       Lineage : cellular organisms;Eukaryota;Opisthokonta;Metazoa;Eumetazoa;Bilateria;Deuterostomia;Chordata;Craniata;Vertebrata;Gnathostomata;Teleostomi;Euteleostomi;Sarcopterygii;Dipnotetrapodomorpha;Tetrapoda;Amniota;Mammalia;Theria;Eutheria;Boreoeutheria;Laurasiatheria;Cetartiodactyla;Ruminantia;Pecora;Bovidae;Bovinae;Bos;Bos taurus
       Reformat: Eukaryota;Chordata;Mammalia;unclassified Mammalia order;Bovidae;Bos;Bos taurus
       Taxids  : 2759;7711;40674;;9895;9903;9913

       Taxid   : 376619
       Lineage : cellular organisms;Bacteria;Proteobacteria;Gammaproteobacteria;Thiotrichales;Francisellaceae;Francisella;Francisella tularensis;Francisella tularensis subsp. holarctica;Francisella tularensis subsp. holarctica LVS
       Reformat: Bacteria;Proteobacteria;Gammaproteobacteria;Thiotrichales;Francisellaceae;Francisella;Francisella tularensis
       Taxids  : 2;1224;1236;72273;34064;262;263

       Taxid   : 349741
       Lineage : cellular organisms;Bacteria;PVC group;Verrucomicrobia;Verrucomicrobiae;Verrucomicrobiales;Akkermansiaceae;Akkermansia;Akkermansia muciniphila;Akkermansia muciniphila ATCC BAA-835
       Reformat: Bacteria;Verrucomicrobia;Verrucomicrobiae;Verrucomicrobiales;Akkermansiaceae;Akkermansia;Akkermansia muciniphila
       Taxids  : 2;74201;203494;48461;1647988;239934;239935

       Taxid   : 239935
       Lineage : cellular organisms;Bacteria;PVC group;Verrucomicrobia;Verrucomicrobiae;Verrucomicrobiales;Akkermansiaceae;Akkermansia;Akkermansia muciniphila
       Reformat: Bacteria;Verrucomicrobia;Verrucomicrobiae;Verrucomicrobiales;Akkermansiaceae;Akkermansia;Akkermansia muciniphila
       Taxids  : 2;74201;203494;48461;1647988;239934;239935

       Taxid   : 314101
       Lineage : cellular organisms;Bacteria;environmental samples;uncultured murine large bowel bacterium BAC 54B
       Reformat: Bacteria;unclassified Bacteria phylum;unclassified Bacteria class;unclassified Bacteria order;unclassified Bacteria family;unclassified Bacteria genus;uncultured murine large bowel bacterium BAC 54B
       Taxids  : 2;;;;;;314101

       Taxid   : 11932
       Lineage : Viruses;Retro-transcribing viruses;Retroviridae;unclassified Retroviridae;Intracisternal A-particles;Mouse Intracisternal A-particle
       Reformat: Viruses;unclassified Viruses phylum;unclassified Viruses class;unclassified Viruses order;Retroviridae;Intracisternal A-particles;Mouse Intracisternal A-particle
       Taxids  : 10239;;;;11632;11749;11932

       Taxid   : 1327037
       Lineage : Viruses;dsDNA viruses, no RNA stage;Caudovirales;Siphoviridae;unclassified Siphoviridae;Croceibacter phage P2559Y
       Reformat: Viruses;unclassified Viruses phylum;unclassified Viruses class;Caudovirales;Siphoviridae;unclassified Siphoviridae genus;Croceibacter phage P2559Y
       Taxids  : 10239;;;28883;10699;;1327037

1. Support tab in format string

        $ echo 9606 | taxonkit lineage | taxonkit reformat -f "{k}\t{p}\t{c}\t{o}\t{f}\t{g}\t{s}\t{S}"
        9606    cellular organisms;Eukaryota;Opisthokonta;Metazoa;Eumetazoa;Bilateria;Deuterostomia;Chordata;Craniata;Vertebrata;Gnathostomata;Teleostomi;Euteleostomi;Sarcopterygii;Dipnotetrapodomorpha;Tetrapoda;Amniota;Mammalia;Theria;Eutheria;Boreoeutheria;Euarchontoglires;Primates;Haplorrhini;Simiiformes;Catarrhini;Hominoidea;Hominidae;Homininae;Homo;Homo sapiens Eukaryota       Chordata        Mammalia        Primates        Hominidae  Homo    Homo sapiens


1. from taxid to 7-ranks lineage:

        $ cat taxids.txt | taxonkit lineage | taxonkit reformat

## name2taxid

Usage

```
query taxid by taxon scientific name

Usage:
  taxonkit name2taxid [flags]

Flags:
  -i, --name-field int   field index of name. data should be tab-separated (default 1)
      --show-rank        show rank

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

1. default

        # taxonkit name2taxid names.txt
        $ cat names.txt | taxonkit name2taxid
        Homo sapiens    9606
        Akkermansia muciniphila ATCC BAA-835    349741
        Akkermansia muciniphila 239935
        Mouse Intracisternal A-particle 11932
        Wei Shen
        uncultured murine large bowel bacterium BAC 54B 314101
        Croceibacter phage P2559Y       1327037

1. show rank

        $ cat names.txt | taxonkit name2taxid --show-rank
        Homo sapiens    9606    species
        Akkermansia muciniphila ATCC BAA-835    349741  no rank
        Akkermansia muciniphila 239935  species
        Mouse Intracisternal A-particle 11932   species
        Wei Shen
        uncultured murine large bowel bacterium BAC 54B 314101  species
        Croceibacter phage P2559Y       1327037 species

1. from name to lineage

        $ cat names.txt | taxonkit name2taxid | taxonkit lineage --taxid-field 2
        Homo sapiens    9606    cellular organisms;Eukaryota;Opisthokonta;Metazoa;Eumetazoa;Bilateria;Deuterostomia;Chordata;Craniata;Vertebrata;Gnathostomata;Teleostomi;Euteleostomi;Sarcopterygii;Dipnotetrapodomorpha;Tetrapoda;Amniota;Mammalia;Theria;Eutheria;Boreoeutheria;Euarchontoglires;Primates;Haplorrhini;Simiiformes;Catarrhini;Hominoidea;Hominidae;Homininae;Homo;Homo sapiens
        Akkermansia muciniphila ATCC BAA-835    349741  cellular organisms;Bacteria;PVC group;Verrucomicrobia;Verrucomicrobiae;Verrucomicrobiales;Akkermansiaceae;Akkermansia;Akkermansia muciniphila;Akkermansia muciniphila ATCC BAA-835
        Akkermansia muciniphila 239935  cellular organisms;Bacteria;PVC group;Verrucomicrobia;Verrucomicrobiae;Verrucomicrobiales;Akkermansiaceae;Akkermansia;Akkermansia muciniphila
        Mouse Intracisternal A-particle 11932   Viruses;Retro-transcribing viruses;Retroviridae;unclassified Retroviridae;Intracisternal A-particles;Mouse Intracisternal A-particle
        Wei Shen
        uncultured murine large bowel bacterium BAC 54B 314101  cellular organisms;Bacteria;environmental samples;uncultured murine large bowel bacterium BAC 54B
        Croceibacter phage P2559Y       1327037 Viruses;dsDNA viruses, no RNA stage;Caudovirales;Siphoviridae;unclassified Siphoviridae;Croceibacter phage P2559Y

        $ cat names.txt | taxonkit name2taxid | taxonkit lineage --taxid-field 2 --show-lineage-taxids
        Homo sapiens    9606    cellular organisms;Eukaryota;Opisthokonta;Metazoa;Eumetazoa;Bilateria;Deuterostomia;Chordata;Craniata;Vertebrata;Gnathostomata;Teleostomi;Euteleostomi;Sarcopterygii;Dipnotetrapodomorpha;Tetrapoda;Amniota;Mammalia;Theria;Eutheria;Boreoeutheria;Euarchontoglires;Primates;Haplorrhini;Simiiformes;Catarrhini;Hominoidea;Hominidae;Homininae;Homo;Homo sapiens        131567;2759;33154;33208;6072;33213;33511;7711;89593;7742;7776;117570;117571;8287;1338369;32523;32524;40674;32525;9347;1437010;314146;9443;376913;314293;9526;314295;9604;207598;9605;9606
        ...

## genautocomplete

Usage

```
generate shell autocompletion script

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

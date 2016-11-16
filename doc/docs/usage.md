# Usage and Examples

## Before use

Please download and uncompress these files:

- [ftp://ftp.ncbi.nih.gov/pub/taxonomy/taxdump.tar.gz](ftp://ftp.ncbi.nih.gov/pub/taxonomy/taxdump.tar.gz)


## taxonkit

Usage

```
TaxonKit - Cross-platform and Efficient NCBI Taxonomy Toolkit

Version: 0.1.3

Author: Wei Shen <shenwei356@gmail.com>

Source code: https://github.com/shenwei356/taxonkit
Documents  : http://bioinf.shenwei.me/taxonkit

Dataset:

    Please download and uncompress these files:
    ftp://ftp.ncbi.nih.gov/pub/taxonomy/taxdump.tar.gz

Usage:
  taxonkit [command]

Available Commands:
  lineage     query lineage of given taxids from file
  list        list taxon tree of given taxids
  reformat    reformat lineage
  version     print version information and check for update

Flags:
  -o, --out-file string   out file ("-" for stdout, suffix .gz for gzipped out) (default "-")
  -j, --threads int       number of CPUs. (default value: 1 for single-CPU PC, 2 for others) (default 2)

Use "taxonkit [command] --help" for more information about a command.

```

## list

Usage

```
list taxon tree of given taxids

Usage:
  taxonkit list [flags]

Flags:
      --ids string      taxid(s), multiple values should be seperated by comma (default "1")
      --indent string   indent (default "  ")
      --json            output in JSON format. you can save the result in file with suffix ".json" and open with modern text editor
      --names string    names.dmp file, when it given taxid will be followed by its scientific name
      --nodes string    nodes.dmp file (default "nodes.dmp")

```

Examples

1. Default usage

        $ taxonkit list --nodes nodes.dmp --ids 9605,239934
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

        $ taxonkit list --nodes nodes.dmp --ids 9605,239934 --indent ""
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

        $ memusg -t taxonkit list --nodes nodes.dmp --ids 1 --indent "" > t0.txt
        elapsed time: 2.987s
        peak rss: 82.79 MB

1. Adding names

        $ taxonkit list --nodes nodes.dmp --names names.dmp --ids 9605,239934
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

        $ memusg -t taxonkit list --nodes nodes.dmp --names names.dmp --ids 1 > t1.txt
        elapsed time: 9.825s
        peak rss: 648.65 MB

1. Output in JSON format, so you can easily collapse and uncollapse taxonomy tree in modern text editor.

        $ taxonkit list --nodes nodes.dmp --names names.dmp --ids 9605,239934 --json
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
query lineage of given taxids from file

Usage:
  taxonkit lineage [flags]

Flags:
  -f, --formated-rank   show formated rank
      --names string    names.dmp file (default "names.dmp")
      --nodes string    nodes.dmp file (default "nodes.dmp")

```

Examples

1. Full lineage:

        $ taxonkit lineage --nodes nodes.dmp --names names.dmp  t.taxid
        349741  cellular organisms;cellular organisms;Bacteria;PVC group;Verrucomicrobia;Verrucomicrobiae;Verrucomicrobiales;Akkermansiaceae;Akkermansia;Akkermansia muciniphila;Akkermansia muciniphila ATCC BAA-835
        834     cellular organisms;cellular organisms;Bacteria;FCB group;Fibrobacteres;Fibrobacteria;Fibrobacterales;Fibrobacteraceae;Fibrobacter;Fibrobacter succinogenes;Fibrobacter succinogenes subsp. succinogenes

2. Formated rank:

        $ taxonkit lineage --nodes nodes.dmp --names names.dmp -f t.taxid
        349741  k__Bacteria;p__Verrucomicrobia;c__Verrucomicrobiae;o__Verrucomicrobiales;f__Akkermansiaceae;g__Akkermansia;s__Akkermansia muciniphila
        834     k__Bacteria;p__Fibrobacteres;c__Fibrobacteria;o__Fibrobacterales;f__Fibrobacteraceae;g__Fibrobacter;s__Fibrobacter succinogenes;S__Fibrobacter succinogenes subsp. succinogenes

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

Usage:
  taxonkit reformat [flags]

Flags:
  -b, --blank string       blank string for missing level (default "__")
  -d, --delimiter string   field delimiter in input lineage (default ";")
  -f, --format string      output format, placeholder of is need (default "{k};{p};{c};{o};{f};{g};{s}")
      --names string       names.dmp file (default "names.dmp")
      --nodes string       nodes.dmp file (default "nodes.dmp")

```

Examples:

Example lineage list:

    $ cat lineage.txt
    root;cellular organisms;Bacteria;environmental samples <Bacteria>;uncultured murine large bowel bacterium BAC 54B;
    root;cellular organisms;Bacteria;PVC group;Verrucomicrobia;Verrucomicrobiae;Verrucomicrobiales;Akkermansiaceae;Akkermansia;Akkermansia muciniphila;
    root;Viruses;dsDNA viruses, no RNA stage;Caudovirales;Siphoviridae;unclassified Siphoviridae;Croceibacter phage P2559Y;
    root;Viruses;Retro-transcribing viruses;Retroviridae;unclassified Retroviridae;Intracisternal A-particles;Mouse Intracisternal A-particle;

1. Default output format ("{k};{p};{c};{o};{f};{g};{s}")

        $ taxonkit reformat lineage.txt | cut -f 2
        Bacteria;__;__;__;__;__;uncultured murine large bowel bacterium BAC 54B
        Bacteria;Verrucomicrobia;Verrucomicrobiae;Verrucomicrobiales;Akkermansiaceae;Akkermansia;Akkermansia muciniphila
        Viruses;__;__;Caudovirales;Siphoviridae;__;Croceibacter phage P2559Y
        Viruses;__;__;__;Retroviridae;Intracisternal A-particles;Mouse Intracisternal A-particle

2. Extracting species

        $ taxonkit reformat lineage.txt -f "{s}" | cut -f 2
        uncultured murine large bowel bacterium BAC 54B
        Akkermansia muciniphila
        Croceibacter phage P2559Y
        Mouse Intracisternal A-particle


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

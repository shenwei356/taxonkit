# GTDB taxonomy taxdump file with trackable TaxIds

## Results

## Method

### Generation of TaxIds

We just hash the name (in lower case) of each taxon node to `uint64`
using [xxhash](https://github.com/cespare/xxhash/) and convert it to `uint32`.

For these missing some taxon nodes, GTDB uses names of parent nodes
e.g., [GCA_018897955.1](https://gtdb.ecogenomic.org/genome?gid=GCA_018897955.1).
So in these cases, TaxIds keep distinct.

    GB_GCA_018897955.1      d__Archaea;p__EX4484-52;c__EX4484-52;o__EX4484-52;f__LFW-46;g__LFW-46;s__LFW-46 sp018897155

We also detect duplicate names with different ranks, e.g., [GCA_003663585.1](https://gtdb.ecogenomic.org/genome?gid=GCA_003663585.1).
The Class and Genus have the same name `B47-G6`, while the Order and Family between them have different names.
In this case, we reassign a new TaxId by increasing the TaxId of name at lower rank until it being distinct.

    GB_GCA_003663585.1      d__Archaea;p__Thermoplasmatota;c__B47-G6;o__B47-G6B;f__47-G6;g__B47-G6;s__B47-G6 sp003663585

### Taxonomic hierarchy

We add an extra subspecies rank ("no rank") to the seven ranks of GTDB taxonomy,
with the assembly accession (without version number) as the name of the taxon node.

A GTDB species cluster contains >=1 assemblies, each can be treated as a strain.
So we can assign each assembly a TaxId with the rank of "no rank" below the species rank.
Therefore, we can also track the changes of these assemblies via the TaxId later.

### Data and tools

GTDB taxnomy files are download from https://data.gtdb.ecogenomic.org/releases/, and organized as:

    $ tree taxonomy/
    taxonomy/
    ├── R080
    │   └── bac_taxonomy_r80.tsv
    ├── R083
    │   └── bac_taxonomy_r83.tsv
    ├── R086
    │   ├── ar122_taxonomy_r86.2.tsv
    │   └── bac120_taxonomy_r86.2.tsv
    ├── R089
    │   ├── ar122_taxonomy_r89.tsv
    │   └── bac120_taxonomy_r89.tsv
    ├── R095
    │   ├── ar122_taxonomy_r95.tsv.gz
    │   └── bac120_taxonomy_r95.tsv.gz
    ├── R202
    │   ├── ar122_taxonomy_r202.tsv.gz
    │   └── bac120_taxonomy_r202.tsv.gz
    └── R207
        ├── ar53_taxonomy_r207.tsv.gz
        └── bac120_taxonomy_r207.tsv.gz

[TaxonKit](https://github.com/shenwei356/taxonkit) v0.11.0 or later version is needed.

### Steps
    
1. Generating taxdump files for the first version r80:

        $ taxonkit create-taxdump taxonomy/R080/*.tsv* --gtdb --out-dir gtdb-taxdump/R080
        15:19:59.816 [WARN] --gtdb-re-subs failed to extract ID for subspecies, the origninal value is used instead. e.g., UBA11420
        15:19:59.964 [INFO] 94759 records saved to gtdb-taxdump/R080/taxid.map
        15:20:00.011 [INFO] 110345 records saved to gtdb-taxdump/R080/nodes.dmp
        15:20:00.048 [INFO] 110345 records saved to gtdb-taxdump/R080/names.dmp
        15:20:00.048 [INFO] 0 records saved to gtdb-taxdump/R080/merged.dmp
        15:20:00.048 [INFO] 0 records saved to gtdb-taxdump/R080/delnodes.dmp
    
1. For later versions, we need the taxdump files of the first version to track merged and deleted nodes.




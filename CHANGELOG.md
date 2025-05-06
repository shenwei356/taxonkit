- [TaxonKit v0.20.0](https://github.com/shenwei356/taxonkit/releases/tag/v0.20.0)
[![Github Releases (by Release)](https://img.shields.io/github/downloads/shenwei356/taxonkit/v0.20.0/total.svg)](https://github.com/shenwei356/taxonkit/releases/tag/v0.20.0)
    
    This version is mainly for maintaining compatibility with NCBI's recent changes([1](https://ncbiinsights.ncbi.nlm.nih.gov/2025/02/27/new-ranks-ncbi-taxonomy/), [2](https://ncbiinsights.ncbi.nlm.nih.gov/2025/04/25/ncbi-taxonomy-updates-virus-classification-april-2025/)).    
    **Please remove the `ranks.txt` file in `~/.taxonkit/` or other directories containing taxdump files (nodes.dmp)**.
    
    - `taxonkit reformat`:
        - Tried to make it compatible with NCBI's recent changes but failed.
        - **Update doc to recommend using `taxonkit reformat2`**.
    - `taxonkit reformat2`:
        - Fix handling of deleted TaxIds. [#117](https://github.com/shenwei356/taxonkit/issues/117)
        - **Change the default output format to be compatible with NCBI's recent changes**. [#115](https://github.com/shenwei356/taxonkit/issues/115), [#116](https://github.com/shenwei356/taxonkit/issues/116))
    - `taxonkit filter`:
        - Support new ranks: `acellular root` and `cellular root`. 
          **Please remove the `ranks.txt` file in `~/.taxonkit/` or other directories containing taxdump files (nodes.dmp)**.
        - Automaticlly switch on `-N/--discard-noranks` when only `-E/--equal-to` is given and the value is not one of ranks without order ("no rank", "clade").
    - `taxonkit name2taxid`:
        - Remove duplicated records caused by capitalization.
- [TaxonKit v0.19.0](https://github.com/shenwei356/taxonkit/releases/tag/v0.19.0)
[![Github Releases (by Release)](https://img.shields.io/github/downloads/shenwei356/taxonkit/v0.19.0/total.svg)](https://github.com/shenwei356/taxonkit/releases/tag/v0.19.0)
    - new command `taxonkit reformat2`: Reformat lineage in chosen ranks, allowing more ranks than 'reformat'
    - `taxonkit reformat`:
        - Fix `-T/--trim` which did not work for `-r/--miss-rank-repl`. [#106](https://github.com/shenwei356/taxonkit/issues/106)
- [TaxonKit v0.18.0](https://github.com/shenwei356/taxonkit/releases/tag/v0.18.0)
[![Github Releases (by Release)](https://img.shields.io/github/downloads/shenwei356/taxonkit/v0.18.0/total.svg)](https://github.com/shenwei356/taxonkit/releases/tag/v0.18.0)
    - `taxonkit reformat`:
        - Add a placeholder for rank "realm", "{r}", which is common in Virus taxonomy like [ictv](https://github.com/shenwei356/ictv-taxdump). [#102](https://github.com/shenwei356/taxonkit/issues/102)
    - `taxonkit name2taxid`:
        - Show warning for names with multiple taxids. [#103](https://github.com/shenwei356/taxonkit/issues/103)
- [TaxonKit v0.17.0](https://github.com/shenwei356/taxonkit/releases/tag/v0.17.0)
[![Github Releases (by Release)](https://img.shields.io/github/downloads/shenwei356/taxonkit/v0.17.0/total.svg)](https://github.com/shenwei356/taxonkit/releases/tag/v0.17.0)
    - `taxonkit filter`:
        - Fix keeping taxids with a rank of "no rank" or "clade". [#97](https://github.com/shenwei356/taxonkit/issues/97)
    - `taxonkit list`:
        - Accept input TaxIds from stdin/files, not just the flag `--ids`. [#93](https://github.com/shenwei356/taxonkit/issues/93)
    - `taxonkit name2taxid`:
        - Add support of fuzzy match with `-f/--fuzzy`. [#88](https://github.com/shenwei356/taxonkit/issues/88)
    - `taxonkit profile2cami`:
        - Add a new flag `-S/--no-sum-up` to disable summing up abundances. [#99](https://github.com/shenwei356/taxonkit/issues/99)
- [TaxonKit v0.16.0](https://github.com/shenwei356/taxonkit/releases/tag/v0.16.0)
[![Github Releases (by Release)](https://img.shields.io/github/downloads/shenwei356/taxonkit/v0.16.0/total.svg)](https://github.com/shenwei356/taxonkit/releases/tag/v0.16.0)
    - `taxonkit create-taxdump`:
        - TaxIds are generated from the hash value of "rank+taxon_name", this allowed duplicated names with different ranks (common in GTDB). [#92](https://github.com/shenwei356/taxonkit/issues/92).
    - `taxonkit taxid-changelog/create-taxdump`:
        - Add notes about the possible error records in taxid-changelog built from taxonomic data created by `create-taxdump`. [#91](https://github.com/shenwei356/taxonkit/issues/91)
        - Dont't worry, a single version of taxonomic data created by `taxonkit create-taxdump` has no problem, it's just the changelog might not be perfect.
    - `taxonkit lca`:
        - new flag `-K/--keep-invalid`: print the query even if no single valid taxid left. [#89](https://github.com/shenwei356/taxonkit/issues/89)
- [TaxonKit v0.15.1](https://github.com/shenwei356/taxonkit/releases/tag/v0.15.1)
[![Github Releases (by Release)](https://img.shields.io/github/downloads/shenwei356/taxonkit/v0.15.1/total.svg)](https://github.com/shenwei356/taxonkit/releases/tag/v0.15.1)
    - `taxonkit name2taxid`:
        - remove the restriction of name types. [#87](https://github.com/shenwei356/taxonkit/issues/87)
- [TaxonKit v0.15.0](https://github.com/shenwei356/taxonkit/releases/tag/v0.15.0)
[![Github Releases (by Release)](https://img.shields.io/github/downloads/shenwei356/taxonkit/v0.15.0/total.svg)](https://github.com/shenwei356/taxonkit/releases/tag/v0.15.0)
    - `taxonkit reformat`:
        - For lineages with more than one node, if it fails to query TaxId with the parent-child pair, use the last child only. [#82](https://github.com/shenwei356/taxonkit/issues/82)
        - The flag `-T/--trim` also does not add the prefix for missing ranks lower than the current rank. [#82](https://github.com/shenwei356/taxonkit/issues/82)
        - New flag `-s/--miss-rank-repl-suffix` to set the suffix for estimated taxon names. [#85](https://github.com/shenwei356/taxonkit/issues/85)
- [TaxonKit v0.14.2](https://github.com/shenwei356/taxonkit/releases/tag/v0.14.2)
[![Github Releases (by Release)](https://img.shields.io/github/downloads/shenwei356/taxonkit/v0.14.2/total.svg)](https://github.com/shenwei356/taxonkit/releases/tag/v0.14.2)
    - `taxonkit filter`:
        - fix checking merged/deleted/not-found taxids. [#80](https://github.com/shenwei356/taxonkit/issues/80)
    - `taxonkit lca`:
        - add a new flag `-b/--buffer-size` to set the size of the line buffer. [#75](https://github.com/shenwei356/taxonkit/issues/75)
        - fix typos: `--separater` -> `--separater`, the former is still available for backward compatibility.
    - `taxonkit reformat`:
        - output compatible format for TaxIds not found in the database. [#79](https://github.com/shenwei356/taxonkit/issues/79)
    - `taxonkit taxid-changelog`:
        - support gzip-compressed taxdump files for saving space. [#78](https://github.com/shenwei356/taxonkit/issues/78)
- [TaxonKit v0.14.1](https://github.com/shenwei356/taxonkit/releases/tag/v0.14.1)
[![Github Releases (by Release)](https://img.shields.io/github/downloads/shenwei356/taxonkit/v0.14.1/total.svg)](https://github.com/shenwei356/taxonkit/releases/tag/v0.14.1)
    - `taxonkit reformat`:
        - The flag `-S/--pseudo-strain` does not require `-F/--fill-miss-rank` now.
        - For taxa of rank >= species, `{t}`, `{S}`, and `T` outputs nothing when using `-S/--pseudo-strain`.
- [TaxonKit v0.14.0](https://github.com/shenwei356/taxonkit/releases/tag/v0.14.0)
[![Github Releases (by Release)](https://img.shields.io/github/downloads/shenwei356/taxonkit/v0.14.0/total.svg)](https://github.com/shenwei356/taxonkit/releases/tag/v0.14.0)
    - `taxonkit create-taxdump`:
        - save taxIds in `int32` instead of `uint32`, as BLAST and DIAMOND do. [#70](https://github.com/shenwei356/taxonkit/issues/70)
    - `taxonkit list`:
        - do not skip visited subtrees when some of give taxids are descendants of others. [#68](https://github.com/shenwei356/taxonkit/issues/68)
    - `taxonkit`:
        - When environment variable `TAXONKIT_DB` is set, explicitly setting `--data-dir` will overide the value of `TAXONKIT_DB`.
- [TaxonKit v0.13.0](https://github.com/shenwei356/taxonkit/releases/tag/v0.13.0)
[![Github Releases (by Release)](https://img.shields.io/github/downloads/shenwei356/taxonkit/v0.13.0/total.svg)](https://github.com/shenwei356/taxonkit/releases/tag/v0.13.0)
    - `taxonkit reformat`:
        - add a new placeholder `{K}` for rank `kingdom`. [#64](https://github.com/shenwei356/taxonkit/issues/64)
        - do not panic for invalid TaxIds, e.g., the column name, when using `-I--taxid-field`.
    - `taxonkit create-taxdump`: 
        - fix merged.dmp and delnodes.dmp. Thanks to @apcamargo ! [gtdb-taxdump/issues/2](https://github.com/shenwei356/gtdb-taxdump/issues/2).
        - fix bug of handling non-GTDB data when using `-A/--field-accession` and no rank names given:
          the colname of the accession column would be treated as one of the ranks, which messed up all the ranks.
        - fix the default option value of `--field-accession-re` which wrongly remove prefix like `Sp_`. [#65](https://github.com/shenwei356/taxonkit/issues/65)
    - `taxonkit list`:
        - fix warning message of merged taxids.
- [TaxonKit v0.12.0](https://github.com/shenwei356/taxonkit/releases/tag/v0.12.0)
[![Github Releases (by Release)](https://img.shields.io/github/downloads/shenwei356/taxonkit/v0.12.0/total.svg)](https://github.com/shenwei356/taxonkit/releases/tag/v0.12.0)
    - `taxonkit create-taxdump`: 
        - accepts arbitrary ranks [#60](https://github.com/shenwei356/taxonkit/issues/60)
        - better handle of taxa with same names.
        - many flags changed.
- [TaxonKit v0.11.1](https://github.com/shenwei356/taxonkit/releases/tag/v0.11.1)
[![Github Releases (by Release)](https://img.shields.io/github/downloads/shenwei356/taxonkit/v0.11.1/total.svg)](https://github.com/shenwei356/taxonkit/releases/tag/v0.11.1)
    - `taxonkit create-taxdump`: fix bug of missing Class rank, contributed by @apcamargo. The flag `--gtdb` was not effected. [#57](https://github.com/shenwei356/taxonkit/issues/57)
- [TaxonKit v0.11.0](https://github.com/shenwei356/taxonkit/releases/tag/v0.11.0)
[![Github Releases (by Release)](https://img.shields.io/github/downloads/shenwei356/taxonkit/v0.11.0/total.svg)](https://github.com/shenwei356/taxonkit/releases/tag/v0.11.0)
    - new command `taxonkit create-taxdump`: Create NCBI-style taxdump files for custom taxonomy, e.g., GTDB and ICTV. [#56](https://github.com/shenwei356/taxonkit/issues/56)
- [TaxonKit v0.10.1](https://github.com/shenwei356/taxonkit/releases/tag/v0.10.1)
[![Github Releases (by Release)](https://img.shields.io/github/downloads/shenwei356/taxonkit/v0.10.1/total.svg)](https://github.com/shenwei356/taxonkit/releases/tag/v0.10.1)
    - `taxonkit cami2-filter`: fix option `--show-rank` which did not work in v0.10.0.
- [TaxonKit v0.10.0](https://github.com/shenwei356/taxonkit/releases/tag/v0.10.0)
[![Github Releases (by Release)](https://img.shields.io/github/downloads/shenwei356/taxonkit/v0.10.0/total.svg)](https://github.com/shenwei356/taxonkit/releases/tag/v0.10.0)
    - new command `taxonkit cami2-filter`: Remove taxa of given TaxIds and their descendants in CAMI metagenomic profile
    - `taxonkit reformat`: fix panic for deleted taxid using `-F/--fill-miss-rank `. [#55](https://github.com/shenwei356/taxonkit/issues/55)
- [TaxonKit v0.9.0](https://github.com/shenwei356/taxonkit/releases/tag/v0.9.0)
[![Github Releases (by Release)](https://img.shields.io/github/downloads/shenwei356/taxonkit/v0.9.0/total.svg)](https://github.com/shenwei356/taxonkit/releases/tag/v0.9.0)
    - new command `taxonkit profile2cami`: converting metagenomic profile table to CAMI format
- [TaxonKit v0.8.0](https://github.com/shenwei356/taxonkit/releases/tag/v0.8.0)
[![Github Releases (by Release)](https://img.shields.io/github/downloads/shenwei356/taxonkit/v0.8.0/total.svg)](https://github.com/shenwei356/taxonkit/releases/tag/v0.8.0)
    - `taxonkit reformat`:
        - **accept input of TaxIds via flag `-I/--taxid-field`**.
        - **accept single taxonomy names**.
        - *show warning message for TaxIds with the same lineage*. [#42](https://github.com/shenwei356/taxonkit/issues/42)
        - better flag checking. [#40](https://github.com/shenwei356/taxonkit/issues/40)
    - `taxonkit lca`:
        - slightly speedup.
    - `taxonkit genautocomplete`:
        - support bash|zsh|fish/powershell
- [TaxonKit v0.7.2](https://github.com/shenwei356/taxonkit/releases/tag/v0.7.2)
[![Github Releases (by Release)](https://img.shields.io/github/downloads/shenwei356/taxonkit/v0.7.2/total.svg)](https://github.com/shenwei356/taxonkit/releases/tag/v0.7.2)
    - `taxonkit lineage`:
        - new flag `-R/--show-lineage-ranks` for appending ranks of all levels.
        - reduce memory occupation and slightly speedup.
    - `taxonkit filter`:
        - flag `-E/--equal-to` supports multiple values.
        - new flag `-n/--save-predictable-norank`: do not discard some special ranks without order when using -L, where rank of the closest higher node is still lower than rank cutoff.
    - `taxonkit reformat`:
        - **new placeholder `{t}` for `subspecies/strain`, `{T}` for `strain`**. Thanks @wqssf102 for feedback.
        - new flag `-S/--pseudo-strain` for using the node with lowest rank as strain name, only if which rank is lower than "species" and not "subpecies" nor "strain". 
- [TaxonKit v0.7.1](https://github.com/shenwei356/taxonkit/releases/tag/v0.7.1)
[![Github Releases (by Release)](https://img.shields.io/github/downloads/shenwei356/taxonkit/v0.7.1/total.svg)](https://github.com/shenwei356/taxonkit/releases/tag/v0.7.1)
    - `taxonkit filter`: 
        - disable unnecessary stdin check when using flag `--list-order` or `--list-ranks`. [#36](https://github.com/shenwei356/taxonkit/issues/36)
        - better handling of black list, empty default value: "no rank" and "clade". And you need use `-N/--discard-noranks` to explicitly filter out "no rank", "clade". [#37](https://github.com/shenwei356/taxonkit/issues/37)
        - update help message. Thanks @standage for improve this command! [#38](https://github.com/shenwei356/taxonkit/issues/38)
- [TaxonKit v0.7.0](https://github.com/shenwei356/taxonkit/releases/tag/v0.7.0)
[![Github Releases (by Release)](https://img.shields.io/github/downloads/shenwei356/taxonkit/v0.7.0/total.svg)](https://github.com/shenwei356/taxonkit/releases/tag/v0.7.0)
    - `taxonkit`: **2-3X faster taxonomy data loading**.
    - new command `taxonkit filter`: **filtering TaxIds by taxonomic rank range**. [#32](https://github.com/shenwei356/taxonkit/issues/32)
    - new command `taxonkit lca`: Computing lowest common ancestor (LCA) for TaxIds.
    - `taxonkit reformat`:
        - new flag `-P/--add-prefix`: **add prefixes for all ranks**, single prefix for a rank is defined by flag `--prefix-X`, where `X` may be `k`, `p`, `c`, `o`, `f`, `s`, `S`.
        - new flag `-T/--trim`: do not fill missing rank lower than current rank.
    - `taxonkit list`: do not duplicate root node.
- [TaxonKit v0.6.2](https://github.com/shenwei356/taxonkit/releases/tag/v0.6.2)
[![Github Releases (by Release)](https://img.shields.io/github/downloads/shenwei356/taxonkit/v0.6.2/total.svg)](https://github.com/shenwei356/taxonkit/releases/tag/v0.6.2)
    - `taxonkit reformat -F`: fix taxids of abbreviated lineage containing names shared by different taxids. [#35](https://github.com/shenwei356/taxonkit/issues/35)
- [TaxonKit v0.6.1](https://github.com/shenwei356/taxonkit/releases/tag/v0.6.1)
[![Github Releases (by Release)](https://img.shields.io/github/downloads/shenwei356/taxonkit/v0.6.1/total.svg)](https://github.com/shenwei356/taxonkit/releases/tag/v0.6.1)
    - `taxonkit lineage`: 
        - new flag `-n/--show-name` for appending scientific name.
        - new flag `-L/--no-lineage` for hide lineage, this is for fast retrieving names or/and ranks.
    - `taxonkit reformat`:
        - fix flag `-F/--fill-miss-rank`.
        - discard order restriction of rank symbols.
- [TaxonKit v0.6.0](https://github.com/shenwei356/taxonkit/releases/tag/v0.6.0)
[![Github Releases (by Release)](https://img.shields.io/github/downloads/shenwei356/taxonkit/v0.6.0/total.svg)](https://github.com/shenwei356/taxonkit/releases/tag/v0.6.0)
    - `taxonkit list`:
        - check merged and deleted taxids.
        - fix bug of json output. [#30](https://github.com/shenwei356/taxonkit/issues/30)
    - `taxonkit name2taxid`: new flag `-s/--sci-name` for limiting to searching scientific names. [#29](https://github.com/shenwei356/taxonkit/issues/29)
    - `taxonkit version`: make checking update optional
- [TaxonKit v0.5.0](https://github.com/shenwei356/taxonkit/releases/tag/v0.5.0)
[![Github Releases (by Release)](https://img.shields.io/github/downloads/shenwei356/taxonkit/v0.5.0/total.svg)](https://github.com/shenwei356/taxonkit/releases/tag/v0.5.0)
    - `taxonkit`: requiring delnodes.dmp and merged.dmp.
    - `taxonkit lineage`: detect deleted and merged taxids now. [#19](https://github.com/shenwei356/taxonkit/issues/19)
    - `taxonkit list/name2taxid`: add short flag `-r` for `--show-rank`, `-n` for `--show-name`.
- [TaxonKit v0.4.3](https://github.com/shenwei356/taxonkit/releases/tag/v0.4.3)
[![Github Releases (by Release)](https://img.shields.io/github/downloads/shenwei356/taxonkit/v0.4.3/total.svg)](https://github.com/shenwei356/taxonkit/releases/tag/v0.4.3)
    - `taxonkit taxid-changelog`: rewrite logic, fix bug and add more change types
- [TaxonKit v0.4.2](https://github.com/shenwei356/taxonkit/releases/tag/v0.4.2)
[![Github Releases (by Release)](https://img.shields.io/github/downloads/shenwei356/taxonkit/v0.4.2/total.svg)](https://github.com/shenwei356/taxonkit/releases/tag/v0.4.2)
    - `taxonkit taxid-changelog`: change output of `ABSORB`, do not merged into one record for changes in different versions.
- [TaxonKit v0.4.1](https://github.com/shenwei356/taxonkit/releases/tag/v0.4.1)
[![Github Releases (by Release)](https://img.shields.io/github/downloads/shenwei356/taxonkit/v0.4.1/total.svg)](https://github.com/shenwei356/taxonkit/releases/tag/v0.4.1)
    - `taxonkit taxid-changelog`: 
        - add fields: `name` and `rank`.
        - and fix sorting bug.
        - detailed lineage change status
- [TaxonKit v0.4.0](https://github.com/shenwei356/taxonkit/releases/tag/v0.4.0)
[![Github Releases (by Release)](https://img.shields.io/github/downloads/shenwei356/taxonkit/v0.4.0/total.svg)](https://github.com/shenwei356/taxonkit/releases/tag/v0.4.0)
    - new command: `taxonkit taxid-changelog`: for creating taxid changelog from dump archive
- [TaxonKit v0.3.0](https://github.com/shenwei356/taxonkit/releases/tag/v0.3.0)
[![Github Releases (by Release)](https://img.shields.io/github/downloads/shenwei356/taxonkit/v0.3.0/total.svg)](https://github.com/shenwei356/taxonkit/releases/tag/v0.3.0)
    - this version is almost the same as v0.2.5
- [TaxonKit v0.2.5](https://github.com/shenwei356/taxonkit/releases/tag/v0.2.5)
[![Github Releases (by Release)](https://img.shields.io/github/downloads/shenwei356/taxonkit/v0.2.5/total.svg)](https://github.com/shenwei356/taxonkit/releases/tag/v0.2.5)
    - add global flag: `--line-buffered` to disable output buffer. [#11](https://github.com/shenwei356/taxonkit/issues/11)
    - replace global flags `--names-file` and `--nodes-file` with `--data-dir`, also support environment variable `TAXONKIT_DB`. [#17](https://github.com/shenwei356/taxonkit/issues/17)
    - `taxonkit reformat`: detects lineages containing unofficial taxon name and won't show panic message.
    - `taxonkit name2taxid`: supports synonyms names. [#9](https://github.com/shenwei356/taxonkit/commit/d4fac1c1138a571957f52eb431ff0d85c03852a8)
    - `taxokit lineage`: add flag `-r/--show-rank` to print rank at another new column.
- [TaxonKit v0.2.4](https://github.com/shenwei356/taxonkit/releases/tag/v0.2.4)
[![Github Releases (by Release)](https://img.shields.io/github/downloads/shenwei356/taxonkit/v0.2.4/total.svg)](https://github.com/shenwei356/taxonkit/releases/tag/v0.2.4)
    - `taxonkit reformat`:
        - more accurate result when using flag `-F/--fill-miss-rank` to estimate and fill missing rank with original lineage information
        - supporting escape strings like `\t`, `\n`, [#5](https://github.com/shenwei356/taxonkit/issues/5)
        - outputting corresponding taxids for reformated lineage. [#8](https://github.com/shenwei356/taxonkit/issues/8)
    - `taxonkit lineage`:
        - fix bug for taxid `1` [#7](https://github.com/shenwei356/taxonkit/issues/7)
        - add flag `-d/--delimiter`.
- [TaxonKit v0.2.3](https://github.com/shenwei356/taxonkit/releases/tag/v0.2.3)
[![Github Releases (by Release)](https://img.shields.io/github/downloads/shenwei356/taxonkit/v0.2.3/total.svg)](https://github.com/shenwei356/taxonkit/releases/tag/v0.2.3)
    - fix bug brought in v0.2.1
- [TaxonKit v0.2.2](https://github.com/shenwei356/taxonkit/releases/tag/v0.2.2)
[![Github Releases (by Release)](https://img.shields.io/github/downloads/shenwei356/taxonkit/v0.2.2/total.svg)](https://github.com/shenwei356/taxonkit/releases/tag/v0.2.2)
    - make verbose information optional [#4](https://github.com/shenwei356/taxonkit/issues/4)
- [TaxonKit v0.2.1](https://github.com/shenwei356/taxonkit/releases/tag/v0.2.1)
[![Github Releases (by Release)](https://img.shields.io/github/downloads/shenwei356/taxonkit/v0.2.1/total.svg)](https://github.com/shenwei356/taxonkit/releases/tag/v0.2.1)
    - `taxonkit list`: fix bug of no output for leaf nodes of the taxonomic tree.
    [#4](https://github.com/shenwei356/taxonkit/issues/4)
    - add new command `genautocomplete` to generate **shell autocompletion** script!
- [TaxonKit v0.2.0](https://github.com/shenwei356/taxonkit/releases/tag/v0.2.0)
[![Github Releases (by Release)](https://img.shields.io/github/downloads/shenwei356/taxonkit/v0.2.0/total.svg)](https://github.com/shenwei356/taxonkit/releases/tag/v0.2.0)
    - add command `name2taxid` to query taxid by taxon scientific name.
    - `lineage`, `reformat`: ***changed flags and default operations***,
 check the [usage](http://bioinf.shenwei.me/taxonkit/usage/).
- [TaxonKit v0.1.8](https://github.com/shenwei356/taxonkit/releases/tag/v0.1.8)
[![Github Releases (by Release)](https://img.shields.io/github/downloads/shenwei356/taxonkit/v0.1.8/total.svg)](https://github.com/shenwei356/taxonkit/releases/tag/v0.1.8)
    - `taxonkit lineage`, add an extra column of lineage in Taxid.
     [#3](https://github.com/shenwei356/taxonkit/issues/3). e.g.,
    - fix colorful output in windows.
- [TaxonKit v0.1.7](https://github.com/shenwei356/taxonkit/releases/tag/v0.1.7)
[![Github Releases (by Release)](https://img.shields.io/github/downloads/shenwei356/taxonkit/v0.1.7/total.svg)](https://github.com/shenwei356/taxonkit/releases/tag/v0.1.7)
    - `taxonkit reformat`: supports reading stdin from output of `taxonkit lineage`,
 reformated lineages are appended to input data.
- [TaxonKit v0.1.6](https://github.com/shenwei356/taxonkit/releases/tag/v0.1.6)
[![Github Releases (by Release)](https://img.shields.io/github/downloads/shenwei356/taxonkit/v0.1.6/total.svg)](https://github.com/shenwei356/taxonkit/releases/tag/v0.1.6)
    - remove flag `-f/--formated-rank` from `taxonkit lineage`,
      using `taxonkit reformat` can archieve same result.
- [TaxonKit v0.1.5](https://github.com/shenwei356/taxonkit/releases/tag/v0.1.5)
[![Github Releases (by Release)](https://img.shields.io/github/downloads/shenwei356/taxonkit/v0.1.5/total.svg)](https://github.com/shenwei356/taxonkit/releases/tag/v0.1.5)
    - reorganize code and flags
- [TaxonKit v0.1.4](https://github.com/shenwei356/taxonkit/releases/tag/v0.1.4)
[![Github Releases (by Release)](https://img.shields.io/github/downloads/shenwei356/taxonkit/v0.1.4/total.svg)](https://github.com/shenwei356/taxonkit/releases/tag/v0.1.4)
    - add flag `--fill` for `taxonkit reformat`, which estimates and fills missing rank with original lineage information
- [TaxonKit v0.1.3](https://github.com/shenwei356/taxonkit/releases/tag/v0.1.3)
[![Github Releases (by Release)](https://img.shields.io/github/downloads/shenwei356/taxonkit/v0.1.3/total.svg)](https://github.com/shenwei356/taxonkit/releases/tag/v0.1.3)
    - add command of `taxonkit reformat` which reformats full lineage to custom format
- [TaxonKit v0.1.2](https://github.com/shenwei356/taxonkit/releases/tag/v0.1.2)
[![Github Releases (by Release)](https://img.shields.io/github/downloads/shenwei356/taxonkit/v0.1.2/total.svg)](https://github.com/shenwei356/taxonkit/releases/tag/v0.1.2)
    - add command of `taxonkit lineage`, users can query lineage of given taxon IDs from file
- [TaxonKit v0.1.1](https://github.com/shenwei356/taxonkit/releases/tag/v0.1.1)
[![Github Releases (by Release)](https://img.shields.io/github/downloads/shenwei356/taxonkit/v0.1.1/total.svg)](https://github.com/shenwei356/taxonkit/releases/tag/v0.1.1)
    - add feature of `taxonkit list`, users can choose output in readable JSON
 format by flag `--json` so the taxonomy tree could be collapse and
 uncollapse in modern text editor.
- [TaxonKit v0.1](https://github.com/shenwei356/taxonkit/releases/tag/v0.1)
[![Github Releases (by Release)](https://img.shields.io/github/downloads/shenwei356/taxonkit/v0.1/total.svg)](https://github.com/shenwei356/taxonkit/releases/tag/v0.1)
    - first release

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

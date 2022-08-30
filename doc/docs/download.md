# Download

`TaxonKit` is implemented in [Go](https://golang.org/) programming language,
 executable binary files **for most popular operating systems** are freely available
  in [release](https://github.com/shenwei356/taxonkit/releases) page.

## Current Version

- [TaxonKit v0.13.0](https://github.com/shenwei356/taxonkit/releases/tag/v0.13.0) [![Github Releases (by Release)](https://img.shields.io/github/downloads/shenwei356/taxonkit/v0.13.0/total.svg)](https://github.com/shenwei356/taxonkit/releases/tag/v0.13.0)
    - `taxonkit reformat`:
        - add a new placeholder `{K}` for rank `kingdom`. [#64](https://github.com/shenwei356/taxonkit/issues/64)
        - do not panic for invalid TaxIds, e.g., the column name, when using `-I--taxid-field`.
    - `taxonkit create-taxdump`: 
        - fix merged.dmp and delnodes.dmp. Thanks to @apcamargo ! [gtdb-taxdump/issues/2](https://github.com/shenwei356/gtdb-taxdump/issues/2).
        - fix bug of handling non-GTDB data when using `-A/--field-accession` and no rank names given.

### Please cite

> Shen, W., Ren, H., TaxonKit: a practical and efficient NCBI Taxonomy toolkit,
> Journal of Genetics and Genomics, [https://doi.org/10.1016/j.jgg.2021.03.006](https://www.sciencedirect.com/science/article/pii/S1673852721000837)

### Links

***Tips***

- run `taxonkit version` to check update !!!
- run `taxonkit genautocomplete` to update Bash completion !!!


OS     |Arch      |File, 中国镜像                                                                                                                                                                              |Download Count
:------|:---------|:-------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|:-----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------
Linux  |**64-bit**|[**taxonkit_linux_amd64.tar.gz**](https://github.com/shenwei356/taxonkit/releases/download/v0.13.0/taxonkit_linux_amd64.tar.gz),<br/> [中国镜像](http://app.shenwei.me/data/taxonkit/taxonkit_linux_amd64.tar.gz)                  |[![Github Releases (by Asset)](https://img.shields.io/github/downloads/shenwei356/taxonkit/latest/taxonkit_linux_amd64.tar.gz.svg?maxAge=3600)](https://github.com/shenwei356/taxonkit/releases/download/v0.13.0/taxonkit_linux_amd64.tar.gz)
Linux  |**arm64**|[**taxonkit_linux_arm64.tar.gz**](https://github.com/shenwei356/taxonkit/releases/download/v0.13.0/taxonkit_linux_arm64.tar.gz),<br/> [中国镜像](http://app.shenwei.me/data/taxonkit/taxonkit_linux_arm64.tar.gz)                  |[![Github Releases (by Asset)](https://img.shields.io/github/downloads/shenwei356/taxonkit/latest/taxonkit_linux_arm64.tar.gz.svg?maxAge=3600)](https://github.com/shenwei356/taxonkit/releases/download/v0.13.0/taxonkit_linux_arm64.tar.gz)
macOS  |**64-bit**|[**taxonkit_darwin_amd64.tar.gz**](https://github.com/shenwei356/taxonkit/releases/download/v0.13.0/taxonkit_darwin_amd64.tar.gz),<br/> [中国镜像](http://app.shenwei.me/data/taxonkit/taxonkit_darwin_amd64.tar.gz)               |[![Github Releases (by Asset)](https://img.shields.io/github/downloads/shenwei356/taxonkit/latest/taxonkit_darwin_amd64.tar.gz.svg?maxAge=3600)](https://github.com/shenwei356/taxonkit/releases/download/v0.13.0/taxonkit_darwin_amd64.tar.gz)
macOS  |**arm64** |[**taxonkit_darwin_arm64.tar.gz**](https://github.com/shenwei356/taxonkit/releases/download/v0.13.0/taxonkit_darwin_arm64.tar.gz),<br/> [中国镜像](http://app.shenwei.me/data/taxonkit/taxonkit_darwin_arm64.tar.gz)               |[![Github Releases (by Asset)](https://img.shields.io/github/downloads/shenwei356/taxonkit/latest/taxonkit_darwin_arm64.tar.gz.svg?maxAge=3600)](https://github.com/shenwei356/taxonkit/releases/download/v0.13.0/taxonkit_darwin_arm64.tar.gz)
Windows|**64-bit**|[**taxonkit_windows_amd64.exe.tar.gz**](https://github.com/shenwei356/taxonkit/releases/download/v0.13.0/taxonkit_windows_amd64.exe.tar.gz),<br/> [中国镜像](http://app.shenwei.me/data/taxonkit/taxonkit_windows_amd64.exe.tar.gz)|[![Github Releases (by Asset)](https://img.shields.io/github/downloads/shenwei356/taxonkit/latest/taxonkit_windows_amd64.exe.tar.gz.svg?maxAge=3600)](https://github.com/shenwei356/taxonkit/releases/download/v0.13.0/taxonkit_windows_amd64.exe.tar.gz)

## Installation

[Download Page](https://github.com/shenwei356/taxonkit/releases)

`TaxonKit` is implemented in [Go](https://golang.org/) programming language,
 executable binary files **for most popular operating systems** are freely available
  in [release](https://github.com/shenwei356/taxonkit/releases) page.

#### Method 1: Download binaries (latest stable/dev version)

Just [download](https://github.com/shenwei356/taxonkit/releases) compressed
executable file of your operating system,
and uncompress it with `tar -zxvf *.tar.gz` command or other tools.
And then:

1. **For Linux-like systems**
    1. If you have root privilege simply copy it to `/usr/local/bin`:

            sudo cp taxonkit /usr/local/bin/

    1. Or copy to anywhere in the environment variable `PATH`:

            mkdir -p $HOME/bin/; cp taxonkit $HOME/bin/

1. **For windows**, just copy `taxonkit.exe` to `C:\WINDOWS\system32`.

#### Method 2: Install via conda  (latest stable version) [![Install-with-conda](https://anaconda.org/bioconda/taxonkit/badges/installer/conda.svg)](http://bioinf.shenwei.me/taxonkit/download/) [![Anaconda Cloud](https://anaconda.org/bioconda/taxonkit/badges/version.svg)](https://anaconda.org/bioconda/taxonkit) [![downloads](https://anaconda.org/bioconda/taxonkit/badges/downloads.svg)](https://anaconda.org/bioconda/taxonkit)

    conda install -c bioconda taxonkit

#### Method 3: Install via homebrew (may not the lastest version)

    brew install brewsci/bio/taxonkit

#### Method 4: Compile from source (latest stable/dev version)

1. [Install go](https://go.dev/doc/install)

        wget https://go.dev/dl/go1.17.13.linux-amd64.tar.gz

        tar -zxf go1.17.13.linux-amd64.tar.gz -C $HOME/

        # or 
        #   echo "export PATH=$PATH:$HOME/go/bin" >> ~/.bashrc
        #   source ~/.bashrc
        export PATH=$PATH:$HOME/go/bin

2. Compile TaxonKit

        # ------------- the latest stable version -------------

        go get -v -u github.com/shenwei356/taxonkit/taxonkit

        # The executable binary file is located in:
        #   ~/go/bin/taxonkit
        # You can also move it to anywhere in the $PATH
        mkdir -p $HOME/bin
        cp ~/go/bin/taxonkit $HOME/bin/


        # --------------- the development version --------------

        git clone https://github.com/shenwei356/taxonkit
        cd taxonkit/taxonkit/
        go build

        # The executable binary file is located in:
        #   ./taxonkit
        # You can also move it to anywhere in the $PATH
        mkdir -p $HOME/bin
        cp ./taxonkit $HOME/bin/

## Bash-completion

Supported shell: bash|zsh|fish|powershell

Bash:

    # generate completion shell
    taxonkit genautocomplete --shell bash

    # configure if never did.
    # install bash-completion if the "complete" command is not found.
    echo "for bcfile in ~/.bash_completion.d/* ; do source \$bcfile; done" >> ~/.bash_completion
    echo "source ~/.bash_completion" >> ~/.bashrc

Zsh:

    # generate completion shell
    taxonkit genautocomplete --shell zsh --file ~/.zfunc/_taxonkit

    # configure if never did
    echo 'fpath=( ~/.zfunc "${fpath[@]}" )' >> ~/.zshrc
    echo "autoload -U compinit; compinit" >> ~/.zshrc

fish:

    taxonkit genautocomplete --shell fish --file ~/.config/fish/completions/taxonkit.fish

## Dataset

1. Download and uncompress `taxdump.tar.gz`: ftp://ftp.ncbi.nih.gov/pub/taxonomy/taxdump.tar.gz 
2. Copy `names.dmp`, `nodes.dmp`, `delnodes.dmp` and `merged.dmp` to data directory: `$HOME/.taxonkit`,
e.g., `/home/shenwei/.taxonkit` ,
3. Optionally copy to some other directories, and later you can refer to using flag `--data-dir`,
or environment variable `TAXONKIT_DB`.

All-in-one command:

    wget -c ftp://ftp.ncbi.nih.gov/pub/taxonomy/taxdump.tar.gz 
    tar -zxvf taxdump.tar.gz
    
    mkdir -p $HOME/.taxonkit
    cp names.dmp nodes.dmp delnodes.dmp merged.dmp $HOME/.taxonkit
    
**Update dataset**: Simply re-download the taxdump files, uncompress and override old ones.

## Release history

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
[![Github Releases (by Release)](https://img.shields.io/github/downloads/shenwei356/taxonkit/v0.6.2/total.svg)](https://github.com/shenwei356/taxonkit/releases/tag/v0.6.1)
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

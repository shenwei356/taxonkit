# TaxonKit - A Practical and Efficient NCBI Taxonomy Toolkit

<img src="taxonkit.jpg" alt="taxonkit" width="800" align="center" />

- **Documents:** [https://bioinf.shenwei.me/taxonkit](https://bioinf.shenwei.me/taxonkit)
([**Usage&Examples**](https://bioinf.shenwei.me/taxonkit/usage/),
[**Tutorial**](https://bioinf.shenwei.me/taxonkit/tutorial/), [**中文介绍**](https://bioinf.shenwei.me/taxonkit/chinese/))
- **Source code:** [https://github.com/shenwei356/taxonkit](https://github.com/shenwei356/taxonkit)
[![GitHub stars](https://img.shields.io/github/stars/shenwei356/taxonkit.svg?style=social&label=Star&?maxAge=2592000)](https://github.com/shenwei356/taxonkit)
[![license](https://img.shields.io/github/license/shenwei356/taxonkit.svg?maxAge=2592000)](https://github.com/shenwei356/taxonkit/blob/master/LICENSE)
[![Built with GoLang](https://img.shields.io/badge/powered_by-go-6362c2.svg?style=flat)](https://golang.org)
- **Latest version:** [![Latest Version](https://img.shields.io/github/release/shenwei356/taxonkit.svg?style=flat?maxAge=86400)](https://github.com/shenwei356/taxonkit/releases)
[![Github Releases](https://img.shields.io/github/downloads/shenwei356/taxonkit/latest/total.svg?maxAge=3600)](https://bioinf.shenwei.me/taxonkit/download/)
[![Cross-platform](https://img.shields.io/badge/platform-any-ec2eb4.svg?style=flat)](https://bioinf.shenwei.me/taxonkit/download/)
[![Anaconda Cloud](	https://anaconda.org/bioconda/taxonkit/badges/version.svg)](https://anaconda.org/bioconda/taxonkit)
- [**Taxid-Changelog**](https://github.com/shenwei356/taxid-changelog): tracking all changes of TaxIds, including deletion, new adding, merge, reuse, and rank/name changes.
- **[Please cite](#citation):** [https://doi.org/10.1016/j.jgg.2021.03.006](https://www.sciencedirect.com/science/article/pii/S1673852721000837)
- [pytaxonkit](https://github.com/bioforensics/pytaxonkit), Python bindings for TaxonKit.



## Table of Contents

<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->

- [Features](#features)
- [Subcommands](#subcommands)
- [Benchmark](#benchmark)
- [Dataset](#dataset)
- [Installation](#installation)
- [Command-line completion](#command-line-completion)
- [Citation](#citation)
- [Contact](#contact)
- [License](#license)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->


## Features

- **Easy to install** ([download](http://bioinf.shenwei.me/taxonkit/download/))
    - Statically linked executable binaries for multiple platforms (Linux/Windows/macOS, amd64/arm64) 
    - Light weight and out-of-the-box, no dependencies, no compilation, no configuration
    - No database building, just download [NCBI taxonomy data](https://ftp.ncbi.nih.gov/pub/taxonomy/taxdump.tar.gz) and uncompress to `$HOME/.taxonkit`
- **Easy to use** ([usages and examples](http://bioinf.shenwei.me/taxonkit/usage/))
    - Supporting [bash-completion](#bash-completion)
    - Fast (see [benchmark](#benchmark)), multiple-CPUs supported, most operations cost 2-10s.
    - Detailed usages and examples
    - Supporting STDIN and (gzipped) input/output file, easily integrated in pipe
- **Versatile commands** 
    - [Usage and examples](http://bioinf.shenwei.me/taxonkit/usage/)
    - Featured command: [tracking monthly changelog of all TaxIds](https://github.com/shenwei356/taxid-changelog)
    - Featured command: [reformating lineage into format of seven-level ("superkingdom, phylum, class, order, family, genus, species"](https://bioinf.shenwei.me/taxonkit/usage/#reformat)
    - Featured command: [filtering taxiDs by a rank range](http://bioinf.shenwei.me/taxonkit/usage/#filter), e.g., at or below genus rank.

## Subcommands

Subcommand                                                                    |Function
:-----------------------------------------------------------------------------|:----------------------------------------------
[`list`](https://bioinf.shenwei.me/taxonkit/usage/#list)                      |List taxonomic subtrees (TaxIds) bellow given TaxIds
[`lineage`](https://bioinf.shenwei.me/taxonkit/usage/#lineage)                |Query taxonomic lineage of given TaxIds
[`reformat`](https://bioinf.shenwei.me/taxonkit/usage/#reformat)              |Reformat lineage in canonical ranks
[`name2taxid`](https://bioinf.shenwei.me/taxonkit/usage/#name2taxid)          |Convert scientific names to TaxIds
[`filter`](https://bioinf.shenwei.me/taxonkit/usage/#filter)                  |Filter TaxIds by taxonomic rank range
[`lca`](https://bioinf.shenwei.me/taxonkit/usage/#lca)                        |Compute lowest common ancestor (LCA) for TaxIds
[`taxid-changelog`](https://bioinf.shenwei.me/taxonkit/usage/#taxid-changelog)|Create TaxId changelog from dump archives
[`profile2cami`](https://bioinf.shenwei.me/taxonkit/usage/#profile2cami)      |Convert metagenomic profile table to CAMI format 
`version`                                                                     |Print version information and check for update
`genautocomplete`                                                             |Generate shell autocompletion script

## Benchmark

1. Getting complete lineage for given TaxIds

   <img src="bench/bench.get_lineage.reformat.tsv.png" alt="" width="600" align="center" />

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

## Installation

Go to [Download Page](https://bioinf.shenwei.me/taxonkit/download) for more download options and changelogs.

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

1. **For Windows**, just copy `taxonkit.exe` to `C:\WINDOWS\system32`.

#### Method 2: Install via conda  (latest stable version) [![Install-with-conda](https://anaconda.org/bioconda/taxonkit/badges/installer/conda.svg)](https://bioinf.shenwei.me/taxonkit/download/) [![Anaconda Cloud](https://anaconda.org/bioconda/taxonkit/badges/version.svg)](https://anaconda.org/bioconda/taxonkit) [![downloads](https://anaconda.org/bioconda/taxonkit/badges/downloads.svg)](https://anaconda.org/bioconda/taxonkit)

    conda install -c bioconda taxonkit

#### Method 3: Install via homebrew (latest stable version)

    brew install brewsci/bio/taxonkit

#### Method 4: For Go developer (latest stable/dev version)

    go get -u github.com/shenwei356/taxonkit/taxonkit

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

## Citation

If you use TaxonKit in your work, please cite:

> Shen, W., Ren, H., TaxonKit: a practical and efficient NCBI Taxonomy toolkit,
> Journal of Genetics and Genomics, [https://doi.org/10.1016/j.jgg.2021.03.006](https://www.sciencedirect.com/science/article/pii/S1673852721000837)

## Contact

[Create an issue](https://github.com/shenwei356/taxonkit/issues) to report bugs,
propose new functions or ask for help.

## License

[MIT License](https://github.com/shenwei356/taxonkit/blob/master/LICENSE)

## Starchart

<img src="https://starchart.cc/shenwei356/taxonkit.svg" alt="Stargazers over time" style="max-width: 100%">

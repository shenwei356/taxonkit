# TaxonKit - A Cross-platform and Efficient NCBI Taxonomy Toolkit

- **Documents:** [https://bioinf.shenwei.me/taxonkit](https://bioinf.shenwei.me/taxonkit)
([**Usage**](https://bioinf.shenwei.me/taxonkit/usage/),
[**Tutorial**](https://bioinf.shenwei.me/taxonkit/tutorial/))
- **Source code:** [https://github.com/shenwei356/taxonkit](https://github.com/shenwei356/taxonkit)
[![GitHub stars](https://img.shields.io/github/stars/shenwei356/taxonkit.svg?style=social&label=Star&?maxAge=2592000)](https://github.com/shenwei356/taxonkit)
[![license](https://img.shields.io/github/license/shenwei356/taxonkit.svg?maxAge=2592000)](https://github.com/shenwei356/taxonkit/blob/master/LICENSE)
[![Built with GoLang](https://img.shields.io/badge/powered_by-go-6362c2.svg?style=flat)](https://golang.org)
- **Latest version:** [![Latest Version](https://img.shields.io/github/release/shenwei356/taxonkit.svg?style=flat?maxAge=86400)](https://github.com/shenwei356/taxonkit/releases)
[![Github Releases](https://img.shields.io/github/downloads/shenwei356/taxonkit/latest/total.svg?maxAge=3600)](https://bioinf.shenwei.me/taxonkit/download/)
[![Cross-platform](https://img.shields.io/badge/platform-any-ec2eb4.svg?style=flat)](https://bioinf.shenwei.me/taxonkit/download/)
[![Anaconda Cloud](	https://anaconda.org/bioconda/taxonkit/badges/version.svg)](https://anaconda.org/bioconda/taxonkit)
- [**Taxid-Changelog**](https://github.com/shenwei356/taxid-changelog): tracking all changes of taxids, including taxids deletion, new adding, merge, reuse, and rank/name changes.
- **[Please cite](#citation):** [![doi](https://img.shields.io/badge/doi-10.1101%2F513523-blue.svg?style=flat)](https://doi.org/10.1101/513523)

## Features

1. **Easy to install**
   - **Statically linked executable binaries for multiple platforms** (Linux/Windows/macOS,x86/amd64/arm64)
     ([download](http://bioinf.shenwei.me/taxonkit/download/))
   - **Light weight and out-of-the-box, no dependencies, no compilation, no configuration**
     ([download](http://bioinf.shenwei.me/taxonkit/download/))
2. **Easy to use**   
   - **Supporting [Bash-completion](#bash-completion)**
   - **Fast** (see [benchmark](#benchmark)), **multiple-CPUs supported**
   - **Well documented usage and examples** (detailed [usage](http://bioinf.shenwei.me/taxonkit/usage/))
   - **Supporting STDIN and gzipped input/output file, easy being used in pipe**
3. **Versatile commands**   
   - [Usage and examples](http://bioinf.shenwei.me/taxonkit/usage/)

## Subcommands

- [`list`](https://bioinf.shenwei.me/taxonkit/usage/#list)    List taxonomic tree of given taxIDs
- [`lineage`](https://bioinf.shenwei.me/taxonkit/usage/#lineage) Query taxonomic lineage of given taxIDs
- [`reformat`](https://bioinf.shenwei.me/taxonkit/usage/#reformat) Reformat lineage in canonical ranks
- [`name2taxid`](https://bioinf.shenwei.me/taxonkit/usage/#name2taxid) Convert scientific names to taxIDs
- [`filter`](https://bioinf.shenwei.me/taxonkit/usage/#filter) Filter taxIDs by taxonomic rank range
- [`taxid-changelog`](https://bioinf.shenwei.me/taxonkit/usage/#taxid-changelog) Create taxID changelog from dump archives

Misc

- `version`   print version information and check for update
- `genautocomplete` generate shell autocompletion script

## Installation

Go to [Download Page](https://bioinf.shenwei.me/taxonkit/download) for more download options and changelogs.

`TaxonKit` is implemented in [Go](https://golang.org/) programming language,
 executable binary files **for most popular operating systems** are freely available
  in [release](https://github.com/shenwei356/taxonkit/releases) page.

#### Method 1: Download binaries (latest stable/dev version)

Just [download](https://github.com/shenwei356/taxonkit/releases) compressed
executable file of your operating system,
and decompress it with `tar -zxvf *.tar.gz` command or other tools.
And then:

1. **For Linux-like systems**
    1. If you have root privilege simply copy it to `/usr/local/bin`:

            sudo cp taxonkit /usr/local/bin/

    1. Or copy to anywhere in the environment variable `PATH`:

            mkdir -p $HOME/bin/; cp taxonkit $HOME/bin/

1. **For windows**, just copy `taxonkit.exe` to `C:\WINDOWS\system32`.

#### Method 2: Install via conda  (latest stable version) [![Install-with-conda](https://anaconda.org/bioconda/taxonkit/badges/installer/conda.svg)](https://bioinf.shenwei.me/taxonkit/download/) [![Anaconda Cloud](https://anaconda.org/bioconda/taxonkit/badges/version.svg)](https://anaconda.org/bioconda/taxonkit) [![downloads](https://anaconda.org/bioconda/taxonkit/badges/downloads.svg)](https://anaconda.org/bioconda/taxonkit)

    conda install -c bioconda taxonkit

#### Method 3: Install via homebrew (latest stable version)

    brew install brewsci/bio/taxonkit

#### Method 4: For Go developer (latest stable/dev version)

    go get -u github.com/shenwei356/taxonkit/taxonkit

## Bash-completion

Note: The current version supports Bash only.
This should work for *nix systems with Bash installed.

Howto:

1. run: `taxonkit genautocomplete`

2. create and edit `~/.bash_completion` file if you don't have it.

        nano ~/.bash_completion

    add the following:

        for bcfile in ~/.bash_completion.d/* ; do
          . $bcfile
        done

## Citation

If you use TaxonKit in your work, please cite [the preprint](https://www.biorxiv.org/content/early/2019/01/08/513523).

## Contact

Email me for any problem when using TaxonKit. shenwei356(at)gmail.com

[Create an issue](https://github.com/shenwei356/taxonkit/issues) to report bugs,
propose new functions or ask for help.

## License

[MIT License](https://github.com/shenwei356/taxonkit/blob/master/LICENSE)

## Starchart

<img src="https://starchart.cc/shenwei356/taxonkit.svg" alt="Stargazers over time" style="max-width: 100%">

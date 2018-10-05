# TaxonKit - Cross-platform and Efficient NCBI Taxonomy Toolkit

- **Documents:** [http://bioinf.shenwei.me/taxonkit](http://bioinf.shenwei.me/taxonkit)
([**Usage**](http://bioinf.shenwei.me/taxonkit/usage/),
[**Tutorial**](http://bioinf.shenwei.me/taxonkit/tutorial/))
- **Source code:** [https://github.com/shenwei356/taxonkit](https://github.com/shenwei356/taxonkit)
[![GitHub stars](https://img.shields.io/github/stars/shenwei356/taxonkit.svg?style=social&label=Star&?maxAge=2592000)](https://github.com/shenwei356/taxonkit)
[![license](https://img.shields.io/github/license/shenwei356/taxonkit.svg?maxAge=2592000)](https://github.com/shenwei356/taxonkit/blob/master/LICENSE)
[![Built with GoLang](https://img.shields.io/badge/powered_by-go-6362c2.svg?style=flat)](https://golang.org)
[![Go Report Card](https://goreportcard.com/badge/github.com/shenwei356/taxonkit)](https://goreportcard.com/report/github.com/shenwei356/taxonkit)
- **Latest version:** [![Latest Version](https://img.shields.io/github/release/shenwei356/taxonkit.svg?style=flat?maxAge=86400)](https://github.com/shenwei356/taxonkit/releases)
[![Github Releases](https://img.shields.io/github/downloads/shenwei356/taxonkit/latest/total.svg?maxAge=3600)](http://bioinf.shenwei.me/taxonkit/download/)
[![Cross-platform](https://img.shields.io/badge/platform-any-ec2eb4.svg?style=flat)](http://bioinf.shenwei.me/taxonkit/download/)
[![Install-with-conda](	https://anaconda.org/bioconda/taxonkit/badges/installer/conda.svg)](http://bioinf.shenwei.me/taxonkit/download/)
[![Anaconda Cloud](	https://anaconda.org/bioconda/taxonkit/badges/version.svg)](https://anaconda.org/bioconda/taxonkit)


## Subcommands

- `list`    list taxon tree of given taxids
- `lineage` query lineage of given taxids (supporting STDIN)
- `reformat` reformat lineage (supporting STDIN)
- `name2taxid` query taxid by taxon scientific name (supporting STDIN)

Misc

- `version`   print version information and check for update
- `genautocomplete` generate shell autocompletion script

## Installation

Go to [Download Page](http://bioinf.shenwei.me/taxonkit/download) for more download options and changelogs.

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

#### Method 2: Install via conda  (latest stable version) [![Install-with-conda](https://anaconda.org/bioconda/taxonkit/badges/installer/conda.svg)](http://bioinf.shenwei.me/taxonkit/download/) [![Anaconda Cloud](https://anaconda.org/bioconda/taxonkit/badges/version.svg)](https://anaconda.org/bioconda/taxonkit) [![downloads](https://anaconda.org/bioconda/taxonkit/badges/downloads.svg)](https://anaconda.org/bioconda/taxonkit)

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

## Contact

Email me for any problem when using taxonkit. shenwei356(at)gmail.com

[Create an issue](https://github.com/shenwei356/taxonkit/issues) to report bugs,
propose new functions or ask for help.

## License

[MIT License](https://github.com/shenwei356/taxonkit/blob/master/LICENSE)

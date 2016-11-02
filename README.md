# TaxonKit - Crossplatform and Efficient NCBI Taxonomy Toolkit

**Documents:** [http://bioinf.shenwei.me/taxonkit](http://bioinf.shenwei.me/taxonkit)
([**Usage**](http://bioinf.shenwei.me/taxonkit/usage/))

**Source code:** [https://github.com/shenwei356/taxonkit](https://github.com/shenwei356/taxonkit)
[![GitHub stars](https://img.shields.io/github/stars/shenwei356/taxonkit.svg?style=social&label=Star&?maxAge=2592000)](https://github.com/shenwei356/taxonkit)
[![license](https://img.shields.io/github/license/shenwei356/taxonkit.svg?maxAge=2592000)](https://github.com/shenwei356/taxonkit/blob/master/LICENSE)

**Latest version:** [![Latest Version](https://img.shields.io/github/release/shenwei356/taxonkit.svg?style=flat?maxAge=86400)](https://github.com/shenwei356/taxonkit/releases)
[![Github Releases](https://img.shields.io/github/downloads/shenwei356/taxonkit/latest/total.svg?maxAge=3600)](http://bioinf.shenwei.me/taxonkit/download/)


## Subcommands

- `list`    list taxon tree of given taxon IDs


## Installation

Go to [Download Page](http://bioinf.shenwei.me/taxonkit/download) for more download options and changelogs.

`taxonkit` is implemented in [Golang](https://golang.org/) programming language,
 executable binary files **for most popular operating systems** are freely available
  in [release](https://github.com/shenwei356/taxonkit/releases) page.

Just [download](https://github.com/shenwei356/taxonkit/releases) compressed
executable file of your operating system,
and uncompress it with `tar -zxvf *.tar.gz` command or other tools.
And then:

1. **For Linux-like systems**
    1. If you have root privilege simply copy it to `/usr/local/bin`:

            sudo cp taxonkit /usr/local/bin/

    1. Or add the directory of the executable file to environment variable
    `PATH`:

            echo export PATH=\$PATH:/PATH/OF/taxonkit >> ~/.bashrc


1. **For windows**, just copy `taxonkit.exe` to `C:\WINDOWS\system32`.

For Go developer, just one command:

    go get -u github.com/shenwei356/taxonkit/taxonkit

## Contact

Email me for any problem when using taxonkit. shenwei356(at)gmail.com

[Create an issue](https://github.com/shenwei356/taxonkit/issues) to report bugs,
propose new functions or ask for help.

## License

[MIT License](https://github.com/shenwei356/taxonkit/blob/master/LICENSE)

# Download

`TaxonKit` is implemented in [Go](https://golang.org/) programming language,
 executable binary files **for most popular operating system** are freely available
  in [release](https://github.com/shenwei356/taxonkit/releases) page.

## Current Version

[TaxonKit v0.1.1](https://github.com/shenwei356/taxonkit/releases/tag/v0.1.1)
[![Github Releases (by Release)](https://img.shields.io/github/downloads/shenwei356/taxonkit/v0.1.1/total.svg)](https://github.com/shenwei356/taxonkit/releases/tag/v0.1.1)

- add feature of `taxonkit list`, users can choose output in readable JSON
 format by flag `--json` so the taxonomy tree could be collapse and
 uncollapse in modern text editor.

Links:

- **Linux**
    - [![Github Releases (by Asset)](https://img.shields.io/github/downloads/shenwei356/taxonkit/latest/taxonkit_linux_386.tar.gz.svg?maxAge=3600)](https://github.com/shenwei356/taxonkit/releases/download/v0.1/taxonkit_linux_386.tar.gz)
    [taxonkit_linux_386.tar.gz](https://github.com/shenwei356/taxonkit/releases/download/v0.1/taxonkit_linux_386.tar.gz)
    - [![Github Releases (by Asset)](https://img.shields.io/github/downloads/shenwei356/taxonkit/latest/taxonkit_linux_amd64.tar.gz.svg?maxAge=3600)](https://github.com/shenwei356/taxonkit/releases/download/v0.1/taxonkit_linux_amd64.tar.gz)
    [taxonkit_linux_amd64.tar.gz](https://github.com/shenwei356/taxonkit/releases/download/v0.1/taxonkit_linux_amd64.tar.gz)
- **Mac OS X**
    - [![Github Releases (by Asset)](https://img.shields.io/github/downloads/shenwei356/taxonkit/latest/taxonkit_darwin_386.tar.gz.svg?maxAge=3600)](https://github.com/shenwei356/taxonkit/releases/download/v0.1/taxonkit_darwin_386.tar.gz)
      [taxonkit_darwin_386.tar.gz](https://github.com/shenwei356/taxonkit/releases/download/v0.1/taxonkit_darwin_386.tar.gz)
    - [![Github Releases (by Asset)](https://img.shields.io/github/downloads/shenwei356/taxonkit/latest/taxonkit_darwin_amd64.tar.gz.svg?maxAge=3600)](https://github.com/shenwei356/taxonkit/releases/download/v0.1/taxonkit_darwin_amd64.tar.gz)
      [taxonkit_darwin_amd64.tar.gz](https://github.com/shenwei356/taxonkit/releases/download/v0.1/taxonkit_darwin_amd64.tar.gz)
- **Windows**
    - [![Github Releases (by Asset)](https://img.shields.io/github/downloads/shenwei356/taxonkit/latest/taxonkit_windows_386.exe.tar.gz.svg?maxAge=3600)](https://github.com/shenwei356/taxonkit/releases/download/v0.1/taxonkit_windows_386.exe.tar.gz)
    [taxonkit_windows_386.exe.tar.gz](https://github.com/shenwei356/taxonkit/releases/download/v0.1/taxonkit_windows_386.exe.tar.gz)
    - [![Github Releases (by Asset)](https://img.shields.io/github/downloads/shenwei356/taxonkit/latest/taxonkit_windows_amd64.exe.tar.gz.svg?maxAge=3600)](https://github.com/shenwei356/taxonkit/releases/download/v0.1/taxonkit_windows_amd64.exe.tar.gz)
    [taxonkit_windows_amd64.exe.tar.gz](https://github.com/shenwei356/taxonkit/releases/download/v0.1/taxonkit_windows_amd64.exe.tar.gz)

## Installation

[Download Page](https://github.com/shenwei356/taxonkit/releases)

`TaxonKit` is implemented in [Go](https://golang.org/) programming language,
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

## Previous Versions

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

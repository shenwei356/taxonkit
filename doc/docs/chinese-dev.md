## 现有工具比较

想要从NCBI获取生物的谱系信息，可以在 NCBI Taxonomy网站上用TaxID或者名称查询。
比如可以用`Homo sapiens`或`9606`搜索“人”的分类学信息，以及密码子表，Entrez记录统计等。

同时也可以通过NCBI的官方工具包 [E-utilities](https://www.ncbi.nlm.nih.gov/books/NBK179288/)
([ftp](https://ftp.ncbi.nlm.nih.gov/entrez/entrezdirect/))。

    $ esearch -db taxonomy -query "txid9606 [Organism]" \
        | efetch -format xml \
        | xtract -pattern Lineage -element Lineage

此外也有一些工具提供类似的功能，部分软件：

工具       |编程语言      |数据获取方式      |使用方式        |备注
:----------|:-------------|:-----------------|:---------------|:-------------------------------------------------------
E-utilities|shell/Perl/C++|远程Web调用       |命令行          |官方程序，Taxonomy操作仅为其部分功能
BioPython  |Python        |远程Web调用       |脚本            |包装entrez接口，Taxonomy操作仅为其部分功能
ETE Toolkit|Python        |本地数据库        |脚本/命令行     |Taxonomy操作仅为其部分功能
Taxize     |R             |远程Web调用         |脚本            |ropensci；支持多种数据库；功能较丰富
Taxopy     |Python        |本地数据文件        |脚本/命令行     | 仅基本功能

选择工具一般考虑几个方面：

1. 是否满足功能需求。大多工具仅有基本的查询、获取完整谱系的功能，都没法将完整谱系格式化为"界门纲目科属种"的格式；
1. 软件安装便利性。上述工具都不需要手动编译安装，除了E-utilities的部分组件需要手动完成，其它基本都能用对应编程语言的包管理工具安装；
1. 配置便利性。部分建立本地数据库的软件则需要先构建数据库，不过基本都是嵌入式的sqlite，比较简单快捷，空间占用也能接受；
1. 使用便利性。提供命令行接口的工具实用较为便捷，也便于整合到分析流程；
   而仅提供包/库的工具，需要使用者在语言终端或编写脚本进行调用，灵活但需要一定编程基础。
1. 计算效率。通过网络调用的软件受网络状态影响大，且在大批量调用的时候速度较慢；实用本地数据库则较为高效。

最初我想要的功能只是根据获取"界门纲目科属种"格式的谱系，发现没有现成工具，而后又有新的需求无法满足，即获取某个类别所有的TaxID。
故开始编写工具来实现，并逐步扩展其功能。

其实最简单的方法就是自己下载数据文件进行解析。

## NCBI Taxonomy 数据文件
    
NCBI Taxonomy数据库将所有生物的**分类学关系**组织为一棵“有根树”（rooted tree）,
与进化树（Phylogenetic tree）不同: 进化树是按**进化关系**”组织，且可以为“无根树”(unrooted tree)。

NCBI Taxonomy公开数据格式有两种，旧的名称为 `taxdump.tar.gz` ，文件大小约50Mb，内含以下文件。

    nodes.dmp       # [当前版本] 节点信息
                    #    重要内容： tax_id, parent tax_id, rank
    names.dmp       # [当前版本] 名称信息
                    #    重要内容： tax_id, name_txt
    merged.dmp      # [目前为止] 被合并的节点信息
                    #    重要内容： old_tax_id, new_tax_id
    delnodes.dmp    # [目前为止] 被删除的nodes信息
                    #    重要内容： tax_id
                    
    citations.dmp   # 引用信息
    division.dmp    # division信息
    gencode.dmp     # 遗传编码信息
    gc.prt          # 遗传编码表
    readme.txt      # 说明文档
    
其中最主要的是前4个文件：

1. `nodes.dmp` 主要包含当前版本的所有分类学单元节点（taxon）
的唯一标识符（taxonomic identifier, 简称TaxId, taxid, tax_id)，
分类学水平(rank），及其父节点的TaxID。
2. `names.dmp` 主要包含包含当前版本的所有TaxID及其统一科学名称（scientific name）和别名。
3. `merged.dmp` 包含了到当前版本为止，所有被合并的TaxID与合并到的新TaxID。
4. `delnodes.dmp` 包含了到当前版本为止，所有被删除的TaxID。

在2018年2月的时候，[推出了新的格式](https://ncbiinsights.ncbi.nlm.nih.gov/2018/02/22/new-taxonomy-files-available-with-lineage-type-and-host-information/)，
额外包含了谱系（lineage），类型（type）和宿主（host）信息。 
文件名称为`new_taxdump.tar.gz`，文件大小约110Mb。
相对旧版，新版本文件数量和内容更多，主要是因为增加了lineage和类型信息。
事实上lineage是可以从`nodes.dmp`和`names.dmp`计算而来。
新版格式所含文件如下：

    nodes.dmp
    names.dmp
    merged.dmp
    delnodes.dmp
    
    fullnamelineage.dmp
    TaxIDlineage.dmp
    rankedlineage.dmp
    
    host.dmp
    typeoftype.dmp
    typematerial.dmp
    
    citations.dmp
    division.dmp
    gencode.dmp
    readme.txt

NCBI Taxonomy数据每天都在更新，每月初（大多为1号）的数据作为存档保存在 `taxdump_archive/` 目录，
旧版本最早数据到2014年8月，新版本只到2018年12月。

## TaxonKit 开发思路

大家应该都有安装生物信息软件的痛苦回忆，在conda出现之前，很多软件都需要手动安装依赖、再编译安装。
不同操作系统，操作系统版本，编译器版本给软件安装带来了巨大的困难。
如果开发者没注意软件的跨平台、可移植性更是如此。

好的软件一定要考虑以下几个方面：

1. 安装便利性。
    1. 尽可能简化安装步骤，甚至一键/一条命令安装。
    1. 减少对外部软件/包的依赖。
    1. 对多平台（windows/linux）的兼容性。
    1. 尽量提供编译好的 静态链接可执行程序（Statically linked executable binaries）。
1. 配置便利性。
    1. 尽可能简化配置，自动化配置，甚至零配置。
1. 使用便利性。
    1. 丰富的文档：安装，使用，例子。
    1. 软件结构合理，模块化。
    1. 友好的报错信息，指出详细的错误原因，而不是只报segmentation fault，或扔出一堆错误信息。
    1. 丰富的命令行参数，满足不同功能需求。
    1. 支持标准输入/输出，从而便于整合到分析流程。
    1. 可选支持shell补全，便于快速调用子命令和参数。
1. 计算效率。
    1. 尽可能占用低内存、低存储。
    1. 尽量减少计算时间，充分利用多CPU。
1. 持续的支持。
    1. 根据用户需求修复bug、增加新功能。
    1. 定期更新发布新版本。

在实现TaxonKit的时候，我已经开始编写seqkit和csvtk软件，有了一定的经验，也基本能达到上述所有要求。

TaxonKit使用Go语言编写，这样可以轻松编译出支持Linux, Windows, 
macOS等操作系统的不同架构（x86/arm）的可执行二进制文件。
由于Go是编译型语言，在运行效率上也有保证。
至于配置、使用等便利性则依赖于开发者。

分类学数据使用NCBI taxonomy的公开数据。
数据访问方式的选择：通过网络访问官方Web接口的方式太慢，只考虑本地访问。
本地访问有几种方式：

1. 直接访问数据库：又分嵌入式数据库如SQLite，第三方数据库入MySQL。后者不考虑，配置太麻烦。
1. Client-Server模式：
    1. Web接口：服务端启动守护进程，长期保持数据库连接，对外提供Web（RESTful）接口，
   客户端本地或远程调用。先前已经开发了一个原型（https://github.com/shenwei356/gtaxon），
   但通过RESTful接口（HTTP）大批量调用，访问速度较慢。
    1. Socket接口：与Web借口类似，因为没有使用http协议，速度应该会高一些。但没有尝试。

最后测试发现，直接解析数据文件的速度也很快，5秒左右（存储为NVMe SSD），完全满足要求。
完全不用搭建数据库，配置更容易，也便于更新数据。
近日又进一步优化到2秒左右，非常快速。内存也在500Mb-1.5G左右，完全可以接受。

TaxonKit为命令行工具，采用子命令的方式来执行不同功能，大多数子命令支持标准输入/输出，便于使用命令行管道进行流水作业。

## 局限性

- 分类学数据库有很多，TaxonKit目前只支持应用最广泛的NCBI Taxonomy。
- 对于GTDB Taxonomy，可以通过现有工具，如[gtdb_to_taxdump](https://github.com/nick-youngblut/gtdb_to_taxdump)，
  将其数据转换为NCBI taxdump文件。
 

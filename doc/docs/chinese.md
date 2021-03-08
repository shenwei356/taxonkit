# TaxonKit: 小巧、高效、实用的NCBI分类学数据命令行工具集

## NCBI Taxonomy 数据库

从事生物多样性的研究者对NCBI Taxonomy数据库一定不会陌生，
它包含了NCBI所有核酸和蛋白序列数据库中每条序列对应的物种名称与分类学信息。
大多数生态学研究对物种组成的描述都是基于NCBI Taxonomy数据库，
当然目前也开始使用其他数据库，如GTDB等。

NCBI Taxonomy数据库始于1991年，一直随着Entrez数据库和其他数据库更新，
1996年推出网页版。NCBI Taxonomy数据库官方地址为 https://www.ncbi.nlm.nih.gov/taxonomy ，
公开数据下载地址为 https://ftp.ncbi.nih.gov/pub/taxonomy/  ，
数据每小时更新，每个月初生成一份数据归档存于 taxdump_archive 目录，最早可追溯到2014年8月。

## TaxonKit 使用

TaxonKit是采用Go语言编写的命令行工具，
提供Linux, Windows, macOS操作系统不同架构（x86-64/arm64）的静态编译的可执行二进制文件。
发布的压缩包不足3Mb，除了Github托管外，还提供国内镜像供下载，同时还支持conda和homebrew安装。
用户只需要**下载、解压，开箱即用，无需配置**，仅需下载解压NCBI Taxonomy数据文件解压到指定目录即可。

- 源代码 https://github.com/shenwei356/taxonkit ，
- 文档 http://bioinf.shenwei.me/taxonkit （介绍、使用说明、例子、教程）

选择系统对应的版本下载最新版 https://github.com/shenwei356/taxonkit/releases ，解压后添加环境变量即可使用。或可选conda安装

    conda install taxonkit -c bioconda -y
    # 表格数据处理，推荐使用 csvtk 更高效
    conda install csvtk -c bioconda -y

测试数据下载可直接 https://github.com/shenwei356/taxonkit 下载项目压缩包，或使用git clone下载项目文件夹，其中的example为测试数据

    git clone https://github.com/shenwei356/taxonkit

TaxonKit为命令行工具，采用子命令的方式来执行不同功能，
大多数子命令支持标准输入/输出，便于使用命令行管道进行流水作业，
轻松整合进分析流程中。

子命令                                                                         |功能
:-----------------------------------------------------------------------------|:----------------------------------------------
[`list`](https://bioinf.shenwei.me/taxonkit/usage/#list)                      |列出指定taxID下所有子单元的的TaxID
[`lineage`](https://bioinf.shenwei.me/taxonkit/usage/#lineage)                |根据TaxID获取完整谱系（lineage）
[`reformat`](https://bioinf.shenwei.me/taxonkit/usage/#reformat)              |将完整谱系转化为“界门纲目科属种株"的自定义格式
[`name2taxid`](https://bioinf.shenwei.me/taxonkit/usage/#name2taxid)          |将分类单元名称转化为TaxID
[`filter`](https://bioinf.shenwei.me/taxonkit/usage/#filter)                  |按分类学水平范围过滤TaxIDs
[`lca`](https://bioinf.shenwei.me/taxonkit/usage/#lca)                        |计算最低公共祖先(LCA)
[`taxid-changelog`](https://bioinf.shenwei.me/taxonkit/usage/#taxid-changelog)|追踪TaxID变更记录
`version`                                                                     |显示版本信息、检测新版本
`genautocomplete`                                                             |生成shell自动补全配置脚本

备注：

- 输出：
    - 所有命令输出中包含输入数据内容，在此基础上增加列。
    - 所有命令默认输出到标准输出（stdout），可通过重定向（`>`）写入文件。
    - 或通过全局参数`-o`或`--out-file`指定输出文件，且可自动识别输出文件后缀（`.gz`）输出gzip格式。
- 输入：
    - 除了`list`与`taxid-changelog`之外，`lineage`, `reformat`, `name2taxid`, `filter` 与 `lca`
      均可从标准输入（stdin）读取输入数据，也可通过位置参数（positional arguments）输入，即命令后面不带
      任何flag的参数，如 `taxonkit lineage taxids.txt`
    - 输入格式为单列，或者制表符分隔的格式，输入数据所在列用`-i`或`--taxid-field`指定。

TaxonKit直接解析NCBI Taxonomy数据文件（2秒左右），配置更容易，也便于更新数据，占用内存在500Mb-1.5G左右。
数据下载：

    # 有时下载失败，可多试几次；或尝试浏览器下载此链接
    wget -c https://ftp.ncbi.nih.gov/pub/taxonomy/taxdump.tar.gz 
    tar -zxvf taxdump.tar.gz
    
    # 解压文件存于家目录中.taxonkit/，程序默认数据库默认目录
    mkdir -p $HOME/.taxonkit
    cp names.dmp nodes.dmp delnodes.dmp merged.dmp $HOME/.taxonkit


### list 列出指定taxID所在子树的所有TaxID

`taxonkit list`用于列出指定TaxID所在分类学单元（taxon）的子树（subtree）的所有taxon的TaxID，可选显示名称和分类学水平。
此功能与NCBI Taxonomy网页版类似。

如，

    # 以人属(9605)和肠道中著名的Akk菌属(239934)为例
    $ taxonkit list --show-rank --show-name --indent "    " --ids 9605,239934
    9605 [genus] Homo
        9606 [species] Homo sapiens
            63221 [subspecies] Homo sapiens neanderthalensis
            741158 [subspecies] Homo sapiens subsp. 'Denisova'
        1425170 [species] Homo heidelbergensis
        2665952 [no rank] environmental samples
            2665953 [species] Homo sapiens environmental sample

    239934 [genus] Akkermansia
        239935 [species] Akkermansia muciniphila
            349741 [strain] Akkermansia muciniphila ATCC BAA-835
        512293 [no rank] environmental samples
            512294 [species] uncultured Akkermansia sp.
            1131822 [species] uncultured Akkermansia sp. SMG25
            1262691 [species] Akkermansia sp. CAG:344
            1263034 [species] Akkermansia muciniphila CAG:154
        1679444 [species] Akkermansia glycaniphila
        2608915 [no rank] unclassified Akkermansia
            1131336 [species] Akkermansia sp. KLE1605
        ...


list使用最广泛的的功能是获取某个类别（比如细菌、病毒、某个属等）下所有的TaxID，
用来从NCBI nt/nr中获取对应的核酸/蛋白序列，从而搭建特异性的BLAST数据库。
官网提供了相应的详细步骤： http://bioinf.shenwei.me/taxonkit/tutorial 。

    # 所有细菌的TaxID
    $ taxonkit list --show-rank --show-name --ids 2 > /dev/null


### lineage 根据TaxID获取完整谱系

分类学数据相关最常见的功能就是根据TaxID获取完整谱系。
TaxonKit可根据输入文件提供的TaxID列表快速计算lineage，并可选提供名称，分类学水平，以及谱系对应的TaxID。

值得注意的是，随着Taxonomy数据的频繁更新，有的TaxID可能被删除、或合并（merge）到其它TaxID中，
TaxonKit会自动识别，并进行提示，对于被合并的TaxID，TaxonKit会按新TaxID进行计算。

    # 使用example中的测试数据
    $ head taxids.txt
    9606
    9913
    376619
    # 查找指定taxids列表的物种信息，tee可输出屏幕并写入文件
    $ taxonkit lineage taxids.txt | tee lineage.txt 
    19:22:13.077 [WARN] taxid 92489 was merged into 796334
    19:22:13.077 [WARN] taxid 1458427 was merged into 1458425
    19:22:13.077 [WARN] taxid 123124124 not found
    19:22:13.077 [WARN] taxid 3 was deleted
    9606    cellular organisms;Eukaryota;Opisthokonta;Metazoa;Eumetazoa;Bilateria;Deuterostomia;Chordata;Craniata;Vertebrata;Gnathostomata;Teleostomi;Euteleostomi;Sarcopterygii;Dipnotetrapodomorpha;Tetrapoda;Amniota;Mammalia;Theria;Eutheria;Boreoeutheria;Euarchontoglires;Primates;Haplorrhini;Simiiformes;Catarrhini;Hominoidea;Hominidae;Homininae;Homo;Homo sapiens
    9913    cellular organisms;Eukaryota;Opisthokonta;Metazoa;Eumetazoa;Bilateria;Deuterostomia;Chordata;Craniata;Vertebrata;Gnathostomata;Teleostomi;Euteleostomi;Sarcopterygii;Dipnotetrapodomorpha;Tetrapoda;Amniota;Mammalia;Theria;Eutheria;Boreoeutheria;Laurasiatheria;Artiodactyla;Ruminantia;Pecora;Bovidae;Bovinae;Bos;Bos taurus
    376619  cellular organisms;Bacteria;Proteobacteria;Gammaproteobacteria;Thiotrichales;Francisellaceae;Francisella;Francisella tularensis;Francisella tularensis subsp. holarctica;Francisella tularensis subsp. holarctica LVS
    349741  cellular organisms;Bacteria;PVC group;Verrucomicrobia;Verrucomicrobiae;Verrucomicrobiales;Akkermansiaceae;Akkermansia;Akkermansia muciniphila;Akkermansia muciniphila ATCC BAA-835
    239935  cellular organisms;Bacteria;PVC group;Verrucomicrobia;Verrucomicrobiae;Verrucomicrobiales;Akkermansiaceae;Akkermansia;Akkermansia muciniphila
    314101  cellular organisms;Bacteria;environmental samples;uncultured murine large bowel bacterium BAC 54B
    11932   Viruses;Riboviria;Pararnavirae;Artverviricota;Revtraviricetes;Ortervirales;Retroviridae;unclassified Retroviridae;Intracisternal A-particles;Mouse Intracisternal A-particle
    1327037 Viruses;Duplodnaviria;Heunggongvirae;Uroviricota;Caudoviricetes;Caudovirales;Siphoviridae;unclassified Siphoviridae;Croceibacter phage P2559Y
    123124124
    3
    92489   cellular organisms;Bacteria;Proteobacteria;Gammaproteobacteria;Enterobacterales;Erwiniaceae;Erwinia;Erwinia oleae
    1458427 cellular organisms;Bacteria;Proteobacteria;Betaproteobacteria;Burkholderiales;Comamonadaceae;Serpentinomonas;Serpentinomonas raicheisms;Bacteria;Proteobacteria;Betaproteobacteria;Burkholderiales;Comamonadaceae;Serpentinomonas;Serpentinomonas raichei

与其它软件的性能相比，当查询数量较少时ETE较快，数量较多时则TaxonKit更快。
在不同数据量规模上 TaxonKit速度一直很稳定，均为2-3秒，时间主要花在解析Taxonomy数据文件上。


列出lineage每个分类学单元的的taxID和rank和名称，比如SARS-COV-2。

    # lineage提取SARS-COV-2的世系
    $ echo "2697049" \
        | taxonkit lineage -t -R \
        | sed "s/\t/\n/g"
    2697049
    Viruses;Riboviria;Orthornavirae;Pisuviricota;Pisoniviricetes;Nidovirales;Cornidovirineae;Coronaviridae;Orthocoronavirinae;Betacoronavirus;Sarbecovirus;Severe acute respiratory syndrome-related coronavirus;Severe acute respiratory syndrome coronavirus 2
    10239;2559587;2732396;2732408;2732506;76804;2499399;11118;2501931;694002;2509511;694009;2697049
    superkingdom;clade;kingdom;phylum;class;order;suborder;family;subfamily;genus;subgenus;species;no rank

### reformat 生成标准层级物种注释 

有时候，我们并不需要完整的分类学谱系（full lineage），因为很多级别即不常用，而且不完整。通常只想保留界门纲目科属种。

值得注意的是，**不是所有物种都有完整的界门纲目科属种水平，特别是病毒以及一些环境样品**。
TaxonKit可以用自定义内容替代缺失的分类单元，如用“__”替代。
更<s>厉害</s>有用的是，**TaxonKit还可以用更高层级的分类单元信息来补齐缺失的层级** (`-F/--fill-miss-rank`)，比如

    # 没有genus的病毒
    $ echo 1327037 | taxonkit lineage | taxonkit reformat | cut -f 1,3
    1327037 Viruses;Uroviricota;Caudoviricetes;Caudovirales;Siphoviridae;;Croceibacter phage P2559Y
    
    # -F参数会用family信息来补齐genus信息
    $ echo 1327037 | taxonkit lineage | taxonkit reformat -F | cut -f 1,3
    1327037 Viruses;Uroviricota;Caudoviricetes;Caudovirales;Siphoviridae;unclassified Siphoviridae genus;Croceibacter phage P2559Y

输出格式可选只输出部分分类学水平，还支持制表符（`"\t"`），再配合作者的另一个工具csvtk，可以输出漂亮的结果。

其它有用的选项：

- `-P/--add-prefix`：给每个分类学水平添加前缀，比如`s__species`。
- `-t/--show-lineage-taxids`：输出分类学单元对应的TaxID。
- `-r/--miss-rank-repl`: 替代没有对应rank的taxon名称
- `-S/--pseudo-strain`: 对于低于species且rank既不是subspecies也不是stain的taxid，使用水平最低taxon名称做为菌株名称。

例，

    $ echo -ne "349741\n1327037"\
        | taxonkit lineage \
        | taxonkit reformat -f "{k}\t{p}\t{c}\t{o}\t{f}\t{g}\t{s}" -F -P \
        | csvtk cut -t -f -2 \
        | csvtk add-header -t -n taxid,kindom,phylum,class,order,family,genus,species \
        | csvtk pretty -t
    
    taxid     kindom        phylum               class                 order                   family               genus                                species
    349741    k__Bacteria   p__Verrucomicrobia   c__Verrucomicrobiae   o__Verrucomicrobiales   f__Akkermansiaceae   g__Akkermansia                       s__Akkermansia muciniphila
    1327037   k__Viruses    p__Uroviricota       c__Caudoviricetes     o__Caudovirales         f__Siphoviridae      g__unclassified Siphoviridae genus   s__Croceibacter phage P2559Y

    # 便于小屏幕查看，用csvtk进行转置
    $ echo -ne "349741\n1327037"\
        | taxonkit lineage \
        | taxonkit reformat -f "{k}\t{p}\t{c}\t{o}\t{f}\t{g}\t{s}" -F -P \
        | csvtk cut -t -f -2 \
        | csvtk add-header -t -n taxid,kindom,phylum,class,order,family,genus,species \
        | csvtk transpose -t \
        | csvtk pretty -t
    
    taxid     349741                       1327037
    kindom    k__Bacteria                  k__Viruses
    phylum    p__Verrucomicrobia           p__Uroviricota
    class     c__Verrucomicrobiae          c__Caudoviricetes
    order     o__Verrucomicrobiales        o__Caudovirales
    family    f__Akkermansiaceae           f__Siphoviridae
    genus     g__Akkermansia               g__unclassified Siphoviridae genus
    species   s__Akkermansia muciniphila   s__Croceibacter phage P2559Y
    
    # 到株水平，以sars-cov-2为例
    $ echo -ne "2697049"\
        | taxonkit lineage \
        | taxonkit reformat -f "{k}\t{p}\t{c}\t{o}\t{f}\t{g}\t{s}\t{t}" -F -P -S \
        | csvtk cut -t -f -2 \
        | csvtk add-header -t -n taxid,kindom,phylum,class,order,family,genus,species,strain \
        | csvtk transpose -t \
        | csvtk pretty -t
    
    taxid     2697049
    kindom    k__Viruses
    phylum    p__Pisuviricota
    class     c__Pisoniviricetes
    order     o__Nidovirales
    family    f__Coronaviridae
    genus     g__Betacoronavirus
    species   s__Severe acute respiratory syndrome-related coronavirus
    strain    t__Severe acute respiratory syndrome coronavirus 2

### name2taxid 将分类单元名称转化为TaxID

将分类单元名称转化为TaxID非常容易理解，唯一要注意的是**某些taxID对应相同的名称**，比如

    # -i指定列，-r显示级别，-L不显示世系
    $ echo Drosophila | taxonkit name2taxid | taxonkit lineage -i 2 -r -L
    Drosophila      7215    genus
    Drosophila      32281   subgenus
    Drosophila      2081351 genus

获取TaxID之后，可以立即传给taxonkit进行后续操作，但要注意用`-i`指定taxID所在列。
    
### filter 按分类学水平范围过滤TaxIDs

filter可以按**分类学水平范围**过滤TaxIDs，注意，不仅仅是特定的Rank，而是一个**范围**。
比如genus及以下的分类学水平，用`-L genus -E genus`，类似于 `<= genus`。

    $ cat taxids2.txt \
        | taxonkit filter -L genus -E genus  \
        | taxonkit lineage -r -n -L \
        | csvtk -Ht cut -f 1,3,2 \
        | csvtk pretty -t
    239934   genus     Akkermansia
    239935   species   Akkermansia muciniphila
    349741   strain    Akkermansia muciniphila ATCC BAA-835

### lca 计算最低公共祖先(LCA)

比如人属的例子

    $ taxonkit list --ids 9605 -nr --indent "    "    
    9605 [genus] Homo
        9606 [species] Homo sapiens
            63221 [subspecies] Homo sapiens neanderthalensis
            741158 [subspecies] Homo sapiens subsp. 'Denisova'
        1425170 [species] Homo heidelbergensis
        2665952 [no rank] environmental samples
            2665953 [species] Homo sapiens environmental sample

TaxID的分隔符可用`-s/--separater`指定，默认为" "。

    # 计算两个物种的最近共同祖先，以上面尼安德特人亚种和海德堡人种
    $ echo 63221 2665953 | taxonkit lca
    63221 2665953   9605

    # 其它分隔符，且不小心多了空格
    $ echo -ne "a\t63221,2665953\nb\t63221, 741158\n"
    a       63221,2665953
    b       63221, 741158
    
    $ echo -ne "a\t63221,2665953\nb\t63221, 741158\n" \
        | taxonkit lca -i 2 -s ","
    a       63221,2665953   9605
    b       63221, 741158   9606
        
## TaxID changelog 追踪TaxID变更记录

NCBI Taxonomy数据每天都在更新，每月初（大多为1号）的数据作为存档保存在 `taxdump_archive/` 目录，
旧版本最早数据到2014年8月，新版本只到2018年12月。

TaxonKit可以追踪所有TaxID每个月的变化，输出到csv文件中，可以通过命令行工具进行查询。
数据和文档单独托管在 https://github.com/shenwei356/taxid-changelog 。

除了简单的增加、删除、合并之外，作者将TaxID改变做了细分。输出格式如下


    # 列            备注
    taxid           # taxid
    version         # version / time of archive, e.g, 2019-07-01
    change          # change, values:
                    #   NEW             新增
                    #   REUSE_DEL       前期被删除，现在又重新加入
                    #   REUSE_MER       前期被合并，现在又重新加入
                    #   DELETE          删除
                    #   MERGE           合并到另一个TaxID
                    #   ABSORB          其他TaxID合并到当前TaxID
                    #   CHANGE_NAME     名称改变
                    #   CHANGE_RANK     分类学水平改变
                    #   CHANGE_LIN_LIN  谱系的TaxID没有变化，谱系改变（某些TaxID的名称变了）
                    #   CHANGE_LIN_TAX  谱系的TaxID改变
                    #   CHANGE_LIN_LEN  谱系的长度/深度发生变化
    change-value    # variable values for changes: 
                    #   1) new taxid for MERGE
                    #   2) merged taxids for ABSORB
                    #   3) empty for others
    name            # scientific name
    rank            # rank
    lineage         # full lineage of the taxid
    lineage-taxids  # taxids of the lineage

数据文件可以在前面网站上下载，`taxid-changelog.csv.gz`，130M左右，解压后2.2G，因为是gzip格式，完全不需要解压即可分析。
下文使用了`pigz`代替`zcat`和`gzip`提高解压速度。

例1 superkingdom也能消失 ，比如类病毒(Viroids)在2019年5月被删除了。
作者是在某一天无意中发现此事，所以决定刨根问底，开发了这个子命令。

    # 下载
    wget -c https://github.com/shenwei356/taxid-changelog/releases/download/v2021.01/taxid-changelog.csv.gz
    # 安装多线程解压索软件。或者用gzip替换。
    conda install pigz
    
    $ pigz -cd taxid-changelog.csv.gz \
        | csvtk grep -f rank -p superkingdom \
        | csvtk pretty 
    taxid   version      change   change-value   name        rank           lineage                        lineage-taxids
    2       2014-08-01   NEW                     Bacteria    superkingdom   cellular organisms;Bacteria    131567;2
    2157    2014-08-01   NEW                     Archaea     superkingdom   cellular organisms;Archaea     131567;2157
    2759    2014-08-01   NEW                     Eukaryota   superkingdom   cellular organisms;Eukaryota   131567;2759
    10239   2014-08-01   NEW                     Viruses     superkingdom   Viruses                        10239
    12884   2014-08-01   NEW                     Viroids     superkingdom   Viroids                        12884
    12884   2019-05-01   DELETE                  Viroids     superkingdom   Viroids                        12884

例2 SARS-CoV-2 。可见新冠病毒在2020年2月加入，随后3月和6月份改了名称，谱系等信息。查询速度也很快。

    # 本例子只显示了部分列。
    $ time pigz -cd taxid-changelog.csv.gz \
        | csvtk grep -f taxid -p 2697049 \
        | csvtk cut -f version,change,name,rank \
        | csvtk pretty

    version      change           name                                              rank
    2020-02-01   NEW              Wuhan seafood market pneumonia virus              species
    2020-03-01   CHANGE_NAME      Severe acute respiratory syndrome coronavirus 2   no rank
    2020-03-01   CHANGE_RANK      Severe acute respiratory syndrome coronavirus 2   no rank
    2020-03-01   CHANGE_LIN_LEN   Severe acute respiratory syndrome coronavirus 2   no rank
    2020-06-01   CHANGE_LIN_LEN   Severe acute respiratory syndrome coronavirus 2   no rank
    2020-07-01   CHANGE_RANK      Severe acute respiratory syndrome coronavirus 2   isolate
    2020-08-01   CHANGE_RANK      Severe acute respiratory syndrome coronavirus 2   no rank

    real    0m7.644s
    user    0m16.749s
    sys     0m3.985s

更多有意思的发现详见[taxid-changelog](https://github.com/shenwei356/taxid-changelog)

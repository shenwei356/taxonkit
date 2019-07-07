#!/usr/bin/env python3

import sys
from ete3 import NCBITaxa

ncbi = NCBITaxa()

for taxid in sys.stdin:
    taxid = taxid.strip()
    
    lineage = ncbi.get_lineage(taxid)
    names = ncbi.get_taxid_translator(lineage)
    lineage_s = "; ".join([names[taxid] for taxid in lineage[1:]])
    
    print("{}\t{}".format(taxid, lineage_s))

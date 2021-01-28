#!/usr/bin/env python3

import sys
from os.path import expanduser
import taxopy

dirhome = expanduser("~")

taxdb = taxopy.TaxDb(nodes_dmp=dirhome+"/.taxopy/nodes.dmp", names_dmp=dirhome+"/.taxopy/names.dmp", keep_files=True)

for taxid in sys.stdin:
    taxid = taxid.strip()
    
    taxon = taxopy.Taxon(taxid, taxdb)
    
    lineage = taxon.name_lineage
    
    lineage_s = ""
    if len(lineage) > 0:
        lineage.reverse()
        lineage_s = "; ".join(lineage[1:])
    
    print("{}\t{}".format(taxid, lineage_s))

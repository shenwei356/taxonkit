#!/usr/bin/env python3

import sys
from os.path import expanduser
from taxadb.taxid import TaxID


handler = TaxID(dbtype="sqlite", dbname=expanduser("~") + "/.taxadb/taxadb.sqlite")

for taxid in sys.stdin:
    taxid = taxid.strip()
    
    lineage = handler.lineage_name(taxid, reverse=True)
    
    lineage_s = ""
    if len(lineage) > 0:
        lineage_s = "; ".join(lineage)
    
    print("{}\t{}".format(taxid, lineage_s))

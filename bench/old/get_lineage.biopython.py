#!/usr/bin/env python3

import sys
from Bio import Entrez

Entrez.email = "anonymous@test.com"

for taxid in sys.stdin:
    taxid = taxid.strip()
    
    handle = Entrez.efetch(db="Taxonomy", id=taxid, retmode="xml")
    records = Entrez.read(handle)
    lineage_s = records[0]["Lineage"] + "; " + records[0]["ScientificName"]
    
    print("{}\t{}".format(taxid, lineage_s))

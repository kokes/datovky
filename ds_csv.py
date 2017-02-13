"""
Naber vsechna XML z https://portal.gov.cz/portal/obcan/rejstriky/ogd/x-sds.html
a vysypej z nich CSV.

Vsechno dej do slozky ./data a pust tohle
"""
from lxml import etree
import gzip
from glob import glob
import os
import csv

fns = glob('data/*.xml.gz')

with open('data/datafile_ds.xsd') as f:
    spec = etree.parse(f).getroot()

# najdem specifikaci, je treba ji rekurzivne projit

els = spec.find('.//{%s}element[@name="box"]' % spec.nsmap['xs'])

def get_spec(els):
    cls = []
    for el in els.find('.//{%s}sequence' % els.nsmap['xs']).getchildren():
        nm, tp = el.attrib['name'], el.attrib['type']

        # komplexni typ
        if not tp.startswith('xs:'):
            dt = spec.find('.//{%s}complexType[@name="%s"]' % (spec.nsmap['xs'], tp[tp.find(':')+1:]))
            for ee in get_spec(dt):
                cls.append('%s_%s' % (nm, ee))
            
            continue # aby se nepridal ten root
            
        cls.append(nm)
    
    return cls

# seznam properties v tech xml
props = get_spec(els)

if not os.path.isdir('./csv'):
    os.mkdir('csv')

for fn in fns:
    print('Parsuju %s' % fn)

    tfn = os.path.join('csv', '%s.csv' % os.path.split(fn)[-1].split('.')[0])

    with gzip.open(fn, 'rt') as gf:
        et = etree.parse(gf).getroot()

        print('Mame %d zaznamu' % len(et.getchildren()))

        fl = open(tfn, 'w')

        cr = csv.writer(fl)
        cr.writerow(props)

        for el in et.iterchildren():
            vals = []
            for pr in props:
                tg = './/{%s}' % el.nsmap[None] + ('/{%s}' % el.nsmap[None]).join(pr.split('_'))
                dt = el.find(tg)
                vl = None
                if dt is not None:
                    vl = dt.text
                vals.append(vl)

            cr.writerow(vals)

        fl.close()

MV zveřejnilo spoustu stovek tisíc majitelů datových schránek, mezi nimi velkou část soukromých držitelů. Vznikla tak nejspíš největší databáze právnických osob vč. základních údajů. Jelikož jde o velká XML, napsal jsem si parsery na zpracování do CSV.

Pythoní verze závisí na `lxml` a standardních knihovnách. Za běhu načítá XSD specifikaci, takže by měla být rezilientní vůči změnám. Go verze má zakodovanou současnou podobu specifikace, takže se možná rozsype. Je ale násobně rychlejší než Pythoní verze. Vybírej dle potřeby.

Parsery nejsou nějak echt ošetřené, za správnost výstupů neručím.

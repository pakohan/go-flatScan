package main

const (
	base             string = "http://kleinanzeigen.ebay.de"
	searchSite       string = "%s/anzeigen/s-wohnung-mieten/berlin/anzeige:angebote/seite:%d/c203l3331"
	entitiyFlatOffer string = "OFFER"
	zipEntity        string = "ZIP"
	userEntitiy      string = "USER"
	email            string = `
You have a new offer:
{{if gt .RentN 0.0}}Rent: {{.RentN}}
{{end}}Adresse: {{if gt (len .Street) 0}}{{.Street}}
{{end}}{{if gt (len .District) 0}}{{.District}}
{{end}}{{if gt .Zip 0}}{{.Zip}} {{end}}Berlin
Rooms: {{.Rooms}}
Size: {{.Size}}
Url: http://kleinanzeigen.ebay.de{{.Url}}

Remove ZIP: http://flat-scan.appspot.com/removeZip?ID={{.Zip}}
Set as invalid: http://flat-scan.appspot.com/toggleOffer?ID={{md5 .Url}}&valid=false

Description: {{.Description}}`
)

var plz map[int]string = map[int]string{
	10115: "Mitte",
	10117: "Mitte",
	10119: "Mitte / Prenzlauer Berg",
	10178: "Mitte",
	10179: "Mitte",
	10243: "Friedrichshain",
	10245: "Friedrichshain",
	10247: "Friedrichshain",
	10249: "Friedrichshain",
	10315: "Friedrichsfelde",
	10317: "Rummelsburg",
	10318: "Karlshorst",
	10319: "Friedrichsfelde",
	10365: "Lichtenberg",
	10367: "Lichtenberg",
	10369: "Lichtenberg",
	10405: "Prenzlauer Berg",
	10407: "Prenzlauer Berg",
	10409: "Prenzlauer Berg",
	10435: "Mitte / Prenzlauer Berg",
	10437: "Prenzlauer Berg",
	10439: "Prenzlauer Berg",
	10551: "Moabit",
	10553: "Tiergarten",
	10555: "Tiergarten",
	10557: "Tiergarten",
	10559: "Tiergarten",
	10585: "Charlottenburg",
	10587: "Charlottenburg",
	10589: "Charlottenburg",
	10623: "BhfZoo",
	10625: "Charlottenburg",
	10627: "Charlottenburg",
	10629: "Charlottenburg",
	10707: "Wilmersdorf",
	10709: "Wilmersdorf",
	10711: "Halensee",
	10713: "Wilmersdorf",
	10715: "Wilmersdorf",
	10717: "Wilmersdorf",
	10719: "Wilmersdorf",
	10777: "Wilmersdorf / Schöneberg",
	10779: "Wilmersdorf / Schöneberg",
	10781: "Schöneberg",
	10783: "Schöneberg",
	10785: "Tiergarten",
	10787: "Tiergarten",
	10789: "Charlottenburg / Schöneberg",
	10823: "Schöneberg",
	10825: "Schöneberg",
	10827: "Schöneberg",
	10829: "Schöneberg",
	10961: "Kreuzberg",
	10963: "Kreuzberg",
	10965: "Kreuzberg / Neukölln / Tempelhof",
	10967: "Kreuzberg",
	10969: "Kreuzberg",
	10997: "Kreuzberg",
	10999: "Kreuzberg",
	12043: "Neukölln",
	12045: "Neukölln",
	12047: "Neukölln",
	12049: "Neukölln",
	12051: "Neukölln",
	12053: "Neukölln",
	12055: "Neukölln",
	12057: "Neukölln",
	12059: "Neukölln",
	12099: "Tempelhof",
	12101: "FH Tempelhof",
	12103: "Schöneberg / Tempelhof",
	12105: "Mariendorf",
	12107: "Mariendorf",
	12109: "Mariendorf",
	12157: "Schöneberg",
	12159: "Friedenau",
	12161: "Friedenau",
	12163: "Steglitz",
	12165: "Steglitz",
	12167: "Steglitz",
	12169: "Steglitz",
	12203: "Lichterfelde",
	12205: "Lichterfelde",
	12207: "Lichterfelde",
	12209: "Lichterfelde",
	12247: "Lankwitz",
	12249: "Lankwitz",
	12277: "Marienfelde",
	12279: "Marienfelde",
	12305: "Lichtenrade",
	12307: "Lichtenrade",
	12309: "Lichtenrade",
	12347: "Britz",
	12349: "Buckow",
	12351: "Buckow",
	12353: "Buckow",
	12355: "Rudow",
	12357: "Buckow / Rudow",
	12359: "Britz",
	12435: "Treptow",
	12437: "Baumschulenweg",
	12439: "Niederschöneweide",
	12459: "Oberschöneweide",
	12487: "Johannisthal",
	12489: "Adlershof",
	12524: "Altglienicke",
	12526: "Bohnsdorf",
	12527: "Grünau / FH Schönefeld",
	12529: "Schönefeld",
	12555: "Köpenick",
	12557: "Köpenick",
	12559: "Müggelheim",
	12587: "Friedrichshagen",
	12589: "Wilhelmshagen",
	12619: "Kaulsdorf",
	12621: "Kaulsdorf",
	12623: "Mahlsdorf",
	12625: "Waldesruh",
	12627: "Hellersdorf",
	12629: "Hellersdorf",
	12679: "Marzahn",
	12681: "Marzahn",
	12683: "Biesdorf",
	12685: "Marzahn",
	12687: "Marzahn",
	12689: "Ahrensfelde",
	13051: "Malchow",
	13053: "Hohenschönhausen",
	13055: "Hohenschönhausen",
	13057: "Falkenberg",
	13059: "Wartenberg",
	13086: "Weißensee",
	13088: "Weißensee",
	13089: "Heinersdorf",
	13125: "Buch / Karow",
	13127: "Buchholz",
	13129: "Blankenburg",
	13156: "Niederschönhausen",
	13158: "Rosenthal",
	13159: "Blankenfelde",
	13187: "Pankow",
	13189: "Pankow",
	13347: "Wedding",
	13349: "Wedding",
	13351: "Wedding",
	13353: "Wedding",
	13355: "Wedding",
	13357: "Wedding",
	13359: "Wedding",
	13403: "Reinickendorf",
	13405: "FH Tegel",
	13407: "Reinickendorf",
	13409: "Reinickendorf",
	13435: "Wittenau",
	13437: "Wittenau",
	13439: "Märkisches Viertel",
	13465: "Frohnau",
	13467: "Hermsdorf",
	13469: "Waidmannslust",
	13503: "Heiligensee",
	13505: "Konradshöhe",
	13507: "Tegel",
	13509: "Borsigwalde",
	13581: "Klosterfelde / Wilhelmstadt",
	13583: "Spandau",
	13585: "Spandau",
	13587: "Hakenfelde",
	13589: "Spandau",
	13591: "Staaken",
	13593: "Pichelsdorf",
	13595: "Pichelsdorf",
	13597: "Charlottenburg",
	13599: "Haselhorst",
	13627: "Spandau",
	13629: "Siemensstadt",
	14050: "Charlottenburg",
	14052: "Charlottenburg",
	14053: "Charlottenburg",
	14055: "Charlottenburg / Wilmersdorf",
	14057: "Charlottenburg",
	14059: "Charlottenburg",
	14089: "Gatow / Kladow",
	14109: "Wannsee",
	14129: "Nikolassee",
	14163: "Zehlendorf",
	14165: "Zehlendorf",
	14167: "Zehlendorf",
	14169: "Zehlendorf",
	14193: "Grunewald",
	14195: "Dahlem",
	14197: "Wilmersdorf",
	14199: "Schmargendorf",
}

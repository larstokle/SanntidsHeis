# SanntidsHeis

##Avoid losing orders:
	-Lokal backupprogram på alle maskiner kjørende
		- kommuniserer ved fil i tilfelle restart ol.
		- fjerner/skriver til fil om forrige sesjon avsluttet riktig
	-remote programstartup når PC "våkner" igjen?
	-Network check: kjøre selv om nettverk er borte
	-Timestamp for sync.

##HW feil:
	- feil kjøretid mellom etasjer.
	- feil sensor går høy under kjøring.
	- ingen sensorsignal.
	- kanppe flickering
	- si ifra om ingen mulighet for å behandle order...

##Heisoppførsel (Cost function):
 - behandle interne først etter retning
 - ta eksterne på veien
	- eksterne FIFO, nærmeste ledig tar den.
		- Ingen ledig: Heis spør etter nye ordre når den blir ledig.



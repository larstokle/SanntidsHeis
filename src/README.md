# SanntidsHeis
- Ingen master altså P2P
- Transaksjonsmanager på hvem som tar bestillinger (voting)
	- Telle hvor mange ganger en bommer på en vote: restart etter nok bom.
	- Feil voting: prøv å sende og kalkuler på ny.
- Hartbeat for å telle antall heiser på nettverk
	- inneholder id og timestamp siden sist hartbeat til alle heiser som er detektert for feilsjekk. fungerer også som teller for hvor mange heiser som er til stede på nettverk.
- Sender knappetrykk og venter på "vote" for sync av kø.

##Moduler/funksjonaliteter
- [ ] main
	- [ ] quit
	- [ ] start backup
	- [ ] init
	- [ ] switch to local master
	- [ ] start elevator manager and network manager
- [ ] ElevatorManager
	- [x] HW
		- [x] Kjør motor
		- [x] Sett bestillingslys
		- [x] sett etasjelys
		- [x] etasjesensor
		- [x] knappetrykk
	- [ ] Cost 
	- [ ] LocalEventmanager
		- [x] knappetrykk behandling
		- [ ] feilevents
	- [ ] Que
		- [x] new order
		- [x] neste ordre
		- [x] fjern ordre
		- [ ] syncronize
	- [x] FSM
		- [x] new event
		- [x] get state
- [ ] Log Manager
	- [ ] write order
	- [ ] read order
	- [ ] clear file
	- [ ] read file
- [ ] Transaction manager
	- [ ] send order
	- [ ] recieve order
	- [ ] send response
	- [ ] handle voting
	- [ ] restart transaction
	- [ ] recieve sync
- [ ] Networking
	- [ ] send
	- [ ] recieve
	- [ ] error check


##Avoid losing orders:
- Lokal backupprogram på alle maskiner kjørende
	- kommuniserer ved fil i tilfelle restart ol.
		- knapp -> nettverk -> kø -> log -> oppdater lys.
	- fjerner/skriver til fil om forrige sesjon avsluttet riktig
- Network check: kjøre selv om nettverk er borte
- Timestamp for sync.

##Heisoppførsel (Cost function):
- behandle interne først etter retning (mulig vi må ta høyde for kukunger..? cost for å skifte retning?)
- ta eksterne på veien
- eksterne FIFO, nærmeste ledig tar den.
	- Ingen ledig: Heis spør etter nye ordre når den blir ledig.

##HW feil:
- feil kjøretid mellom etasjer.
- feil sensor går høy under kjøring.
- ingen sensorsignal.
- knappe flickering
- si ifra om ingen mulighet for å behandle order...



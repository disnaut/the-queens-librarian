# the-queens-librarian
## Requirements from Client
- [ ] Robust ability to search for cards
    - [X] Manacost
    - [X] Colors
    - [X] Types
    - [X] Keywords
    - [X] Rarity
    - [X] Name
    - [X] Artist
- [ ] Ability for users to create their own colloqulialisms that can then be turned into valid search terms.
    - E.g. "Power Nine", "Mana Rocks", "Aristocrats"
- [ ] Ability to catalogue cards as part of a **Personal** collection
- [ ] Ability to create decks based on various formats
- [ ] When creating a deck, app will inform the user of what cards they do, and do not, have based on their personal collection.
- [ ] Update pricing information to give the user an estimate of how much the rest of a deck would cost
**Requirements Subject To Change**
# Current Happenings
- Currently working Unit Tests for what is currently implemented
- Working on getting a "set up" guide working
- Working on getting a .config file working
- Working on re-factoring the Card.go to utilize go routines


## 02/19/2023
- Created some basic query strings for URL's, and added `mongo.go` with some stubbed out code.
## 02/23/2023
- Removed the example files and updated functionality to interact with the entire database. Going to start creating branches to extend functionality of the webAPI to encompass more search functions, as well as the creation of an overall collection from the user, as well as decks that the user can create.
- Updated README to include requirements gathered from client.
## 03/06/2023
- Updated search function to include everything in requirements but mana cost.

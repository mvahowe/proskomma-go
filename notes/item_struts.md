# ITEM STRUTS

Use Go enums for fixed types, eg itemType

ITEMS = []ITEM*

ITEM = TOKEN | SCOPE | GRAFT

ITEMINTERFACE = type int, subType int

## Option 1: Keep succinct representation

- more compact
- faster to search
- more latency to access


TOKEN = ITEMINTERFACE, charsIndex int

SCOPE = ITEMINTERFACE, labelIndexes []int

GRAFT = ITEMINTERFACE, seqId int

## Option 2: Unpack succinct representation

- quicker to access
- easier to use

TOKEN = ITEMINTERFACE, chars string

SCOPE = ITEMINTERFACE, label string

GRAFT = ITEMINTERFACE, seqId string

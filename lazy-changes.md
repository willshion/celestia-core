Tracking changes:


- make Tx of the form (namespace || data)
  - this makes certain assumptions about the (length and structure) of the transactions

Hash / Merkle root of all Tx:

- add row and column roots to the block/header to enable data availability
  - DataHash (classical merkle root of all messages/Txs) becomes obsolete but could be a simpler way for full-nodes
  storing the whole block anyways to validate a block
  - if we delete DataHash (which becomes obsolete) this breaks the assumption of the header structure
  used allover the tendermint code base ...
  - for a first version: we kept the structure as it is and simply add additional fields
  but later we need to adapt the the DataHash to consist of the following:
  > We nevertheless represent all of the row and column roots as a a single dataRoot_i to allow ‘super-light’ clients
  > which do not download the row and column roots, but these clients cannot be assured of data availability and thus
  > do not fully benefit from the increased security of allowing fraud proofs.

https://arxiv.org/abs/1809.09044

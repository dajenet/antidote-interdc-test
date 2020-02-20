# Antidote InterDc Test

AntidoteDB stress test for inter dc communication. Creates two datacenters and increments a single counter concurrently on both dcs. Also applies package loss between both dcs.
`run.sh` contains all steps to run. `stop.sh` can be used to clean up after the test.

# apriori-miner

A pure-Go implementation of the Apriori algorithm for mining frequent itemsets
and generating association rules from transaction data. The tool also computes
rule quality metrics (lift, leverage, conviction, cosine, Jaccard, Kulczynski)
and persists mining results to disk with SHA-256 integrity verification.

## Features

- Apriori frequent itemset mining with configurable minimum support
- Association rule generation with minimum confidence threshold
- Six interestingness measures for rule quality assessment
- Flexible rule filtering (by metrics, items, length)
- Result persistence with checksum verification (corruption detection)
- Incremental merge of results from separate mining batches
- Result validation for integrity constraints

## Build & Test

```bash
export GOTOOLCHAIN=local CGO_ENABLED=0
go build ./...
go test ./...
```

## Usage

```bash
# Mine frequent itemsets and print rules
apriori-miner -input transactions.txt -min-support 0.2 -min-confidence 0.6

# Input format: one transaction per line, items separated by whitespace
# Example transactions.txt:
#   bread milk eggs
#   bread butter
#   milk diaper beer
```

## Directory Structure

```
internal/apriori    Apriori algorithm (frequent sets + rule generation)
internal/metrics    Interestingness measures (lift, leverage, conviction, cosine, Jaccard, Kulczynski)
internal/filter     Rule filtering engine (thresholds, item constraints, length)
internal/persist    Result persistence (JSON + SHA-256 checksum, merge, validate)
```

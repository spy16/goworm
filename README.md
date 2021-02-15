# 🐛 GoWorm

GoWorm is a simulation system for spiking neural networks. 

`goworm` is generic enough to run different models defined using simple
CSV files. A full connectome of C.Elegans is provided.

## Usage 

1. Install the tool using `go get -u -v github.com/spy16/goworm`

2. Run the model using:

```shell
$ goworm -model c_elegans.csv -tick 100ms -threshold 30

# -tick represents the simulation step interval (use <1 second).
# -threshold represents the cell threshold for firing.
```

Note: Refer [c_elegans.csv](./c_elegans.csv) for model format.


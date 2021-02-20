# 🐛 GoWorm

GoWorm is a simulation system for spiking neural networks.

`goworm` is generic enough to run different models defined using simple JSON files. A full connectome of C.Elegans is
provided.

## Usage

1. Install the console version using `go get -u -v github.com/spy16/goworm`

2. Run the model using:

    ```shell
    $ goworm -model c_elegans.json -tick 100ms -addr :8081
    
    # -model configures the model file to read from (can be "-" to read from standard input)
    # -tick configures the simulation step interval (use <1 second).
    # -addr configures the
    ```

3. Poke different cell groups using `curl http://localhost:8081/poke?group=nose`

Note: Refer [c_elegans.json](./c_elegans.json) for model format.

## References:

1. <https://github.com/Flowx08/Celegans-simulation>
2. <https://github.com/Connectome/GoPiGo>

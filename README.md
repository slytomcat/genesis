
# It is simple genetic simulation of mortal and immortal populations.

## The simulation assumes that:
1. a population consists of some amount of creturee and it lives into some environment for some amount of time periods, lets call them years.
2. the environment has a number permanently changing (each year) factors and some permanent capacity.
3. a creature:
  - has age
  - has some amount of chromosomes and they have different lengths
  - has a fertility age after which it can reproduce a new creature (every year) with some probability. During the reproduction the random partner is selected form population and the new creature gets a random combitation of the partners chomosomes. With some probability the chomosomes of a new creature can mutate: they change its length.
  - with some probability it can die due to combination of 3 factor:
    - a. the compatibility to the changing factor of environment: it is determined as average of minimum distances between any of its chromosome and each changing environment factor
    - b. the factor that depends on population size and environment capacyty (it is groving by exponential rule: low for small population and becomming much much bigger for populatin that bigger than the environment capacity)
    - c. the age factor: it is increasing each year on some value
  - creatures on immortal population always has the death factor 'c' equal to zero

The simulation utility is written on golang and requires go v.23.1 or newer for building the binary.

## Building

    go build

## Usage

First review and change the settings into [`settings.yaml`](settings.yaml) file. Then you may run the utility:

### start simulation within the random environment

    ./genesis random
### generate and save a new random environment

    ./genesis store   
### run simulation within the random environment

    ./genesis stored  
    

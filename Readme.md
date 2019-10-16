# GoCodeVis
The Golang Code Visualizer

This project reads a Golang source tree and generates text files using PlantUML's DSL.  These files can be processed by [PlantUML](http://plantuml.com/) to generate PNGs which contain the UML class digrams of the Go source tree.

Now, you may be saying "but Go isn't an real OO language".  I think that receiver function declarations (pointer or value) seem awfully similar to how we defined classes in C++ and that non-receiver function declarations, which are bound to a scope, seem awfully similar to static methods of a class named after the binding scope.  Thus, it is close enough that we can use UML class diagrams at the least to have interesting discussions about the coupling and cohesive characteristics of the code base (the two most important design characteristics of any code base).

I've included the diagrams of this project in the [projectdiagrams folder](https://github.com/drp2179/gocodevis/tree/master/prodjectdiagrams/) as examples of what this project can do.

## Usage 
gocodevis --help
gocodevis -output <outputpath> -root <sourceroot> -verbose



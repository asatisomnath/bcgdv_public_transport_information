## bcgdv public transport information

To build:

```
$ go install github.com/asatisomnath/bcgdv_public_transport_information
```

To run:

```
$ go run main.go
```


APIs Params:

```

Find a vehicle for a given time and X & Y coordinates: if there exist a stop for given x and y then it finds the next line name arriving otherwise firstly it will get the nearest stop and return the same. 
http://localhost:8081/arriving/10:00:00/1/1
http://localhost:8081/arriving/10:10:00/2/8
http://localhost:8081/arriving/10:52:00/0/0

Return the vehicle arriving next at a given stop : it returns the next line name arriving at the stops otherwise return doesn't exist.
http://localhost:8081/arriving/1
http://localhost:8081/arriving/5

Indicate if a given line is currently delayed: it returns the delay value for the line name otherwise it return Line is on time :)
http://localhost:8081/delay/200
http://localhost:8081/delay/M4
```


- ## LSMA
  - length: 20
  - offset: 8
  - source: low

- ## Liniear Regression
 - upperDeviation: 2
 - lowerDeviation: -2
 - length: 120
 - source: close
 - description: Calculate a linreg on {length} candle before and add it to the stdDev * upper/lowerDeviation

 - ## BB %B:
  - length: 200
  - source: close
  - stdDev: 1
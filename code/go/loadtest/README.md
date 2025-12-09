# Morning Rush Hour Car Inflow: Cars per Second

Estimated number of cars entering major cities per second during morning rush:

| City | Vehicles Entering (Morning Rush) | Cars per Second |
|------|---------------------------------|----------------|
| Tokyo | ~1,000,000 | ~93 |
| Mexico City | ~650,000 | ~60 |
| NYC (Manhattan) | ~275,000 | ~25 |
| Paris | ~200,000 | ~19 |

## Testing billing API

### Result test for 1min with 200 cars/s

| Metric                     | Value           |
|----------------------------|-----------------|
| Total requests sent        | 11,970          |
| Successful requests        | 11,970          |
| Failed requests            | 0               |
| Total test duration        | 1m0.000449125s  |
| Average request latency    | 1.924132ms      |
| Minimum request latency    | 853.708µs       |
| Maximum request latency    | 79.724542ms     |
| Achieved RPS               | 199.50          |

### Result test for 10min with 300 cars/s

| Metric                     | Value            |
|----------------------------|------------------|
| Total requests sent        | 177,928          |
| Successful requests        | 177,928          |
| Failed requests            | 0                |
| Total test duration        | 10m 0.0015935s   |
| Average request latency    | 3.313398 ms      |
| Minimum request latency    | 817.75 µs        |
| Maximum request latency    | 475.874083 ms    |
| Achieved RPS               | 296.55           |

### Result test for 60min with 300 cars/s

| Metric                     | Value             |
|----------------------------|-------------------|
| Total requests sent        | 1,070,590         |
| Successful requests        | 1,070,590         |
| Failed requests            | 0                 |
| Total test duration        | 1h 0m 0.042680625s|
| Average request latency    | 3.459595 ms       |
| Minimum request latency    | 793.167 µs        |
| Maximum request latency    | 816.322583 ms     |
| Achieved RPS               | 297.38            |

## Testing the limits:

| Metric                     | Value          |
|----------------------------|----------------|
| Total requests sent        | 35,712         |
| Successful requests        | 35,712         |
| Failed requests            | 0              |
| Total test duration        | 1m 0.000026s   |
| Average request latency    | 14.005705 ms   |
| Minimum request latency    | 778.667 µs     |
| Maximum request latency    | 532.184792 ms  |
| Achieved RPS               | 595.20         |

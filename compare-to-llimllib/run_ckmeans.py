import sys
import ckmeans

data = [line for line in sys.stdin]

nclusters = int(data[0])
vals = [float(val) for val in data[1:]]

result = ckmeans.ckmeans(vals, nclusters)
for i, cluster in enumerate(result):
    for val in cluster:
        print(i, val)

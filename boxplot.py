import matplotlib.pyplot as plt
from collections import defaultdict
import collections
import numpy as np
import os

script_path = os.path.dirname(__file__)
data = defaultdict(list)

for i in range(10,16):
  rel_path = "nur100/ids_ed25519_gt_"+str(i)

  f = open(os.path.join(script_path, rel_path), "r")
  lines = f.readlines()

  for line in lines:
    line = line.split()
    data[i].append(int(line[1]))

od = collections.OrderedDict(sorted(data.items()))

keys = []
values = []

for key in od:
  if len(data[key]) < 100:
    continue

  values.append(data[key])
  keys.append(key)

# print(od.keys())

plt.boxplot(values)
plt.xticks(np.arange(1, len(keys)+1), keys)

plt.show()
import numpy as np
import os
import matplotlib.pyplot as plt

dirname = 'results'

plt.title('Recall-Query per second(1/s) tradeoff')
plt.xlabel('Recall')
plt.ylabel('Query per second (1/s)')
plt.yscale('log')
plt.xlim(0.0, 1.03)

markers = ['x', '.', '^', 'o', 's', 'D']

for i, f in enumerate(os.listdir(dirname)):
    data = np.loadtxt(dirname + '/' + f, delimiter=',')
    if len(data.shape) == 1:
        data = data.reshape(-1, 2)
    data = data[np.argsort(data[:,1])]
    last_x = float('-inf')
    comparator = \
      (lambda xv, lx: xv > lx) if last_x < 0 else (lambda xv, lx: xv < lx)

    xs, ys = [], []
    for xv, yv in data:
        if comparator(xv, last_x):
            last_x = xv
            xs.append(xv)
            ys.append(yv)
    plt.plot(np.array(xs), 1.0/np.array(ys), marker=markers[i], label=f.replace('.csv', ''))

plt.legend(bbox_to_anchor=(1.05, 1), loc='upper left', borderaxespad=0)
plt.subplots_adjust(right=0.7)

plt.savefig('out.png')

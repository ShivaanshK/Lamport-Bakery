# Lamport-Bakery

This repository implements Lamport's bakery algorithm for mutual exclusion.

To run a demo, do the following.

```bash
git@github.com:ShivaanshK/Lamport-Bakery.git
cd Lamport-Bakery
```
n = Number of processes randomly competing for access to the criticial section

i = Number of iterations each process will run the algorithm for
```bash
go run main.go -numProcesses=n numIterations=i
```

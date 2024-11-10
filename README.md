# MapReduce

TODO:
1. start the servers in go routines instead of main routine (done)
2. implement map and reduce functions in pkg/utils.go (done)
3. execute  tasks on worker side (done)
4. assign next after worker is finished with current task
5. notify master of finished tasks after map (done)
6. simulate failures
7. handle failures for workers (done , not tested , point 6)
8. change reduce function input from list of key value pair to key, list of value pair (involves changing how the mapper will send data) (done)
9. debug kernel case (done)

## Instructions to run
1. At seperate terminals run the workers

```bash
go run cmd/worker/main.go -port 707(0/1/2/3/)
```
2. Run the master

```bash
go run cmd/master/main.go
```
3. To generate new test cases, run `gen.py`

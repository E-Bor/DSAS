## Data Store As Service Project

This project design and implement as prototype of production ready solution. It contain a futures:
- gRPC call to start some task
- A queue with auto-scheduling, which is based on average task completion time and estimated date.
- Pipeline based on independent components
- High performance for I/O bound tasks
- It's easy to do integrations. You can simply write a separate application in Go, and then compile it into a module. It will automatically be available from the DSAS system.
- Easily scalable
- Bash scripts for auto compile integrations

## Architecture 
![image](https://github.com/user-attachments/assets/835960b3-c86d-4ac8-b7bc-5ae14630c2e3)

## Requered
Proto contracts [here](https://github.com/E-Bor/DSAS_Proto)

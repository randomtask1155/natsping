## Build

```
go get github.com/nats-io/go-nats
go get github.com/randomtask1155/natsping
```

## get usage

```
Usage:

./natsping -s [ ip:port | ip:port,ip:port ] -u username -p password -sub subject -m message -t <timeout value in seconds default=10>

Use Case 1
	bosh director is reporting agent ping timeout after <X> seconds even though the agent is up and running on the remove vm. In this case you may want to ping the remote bosh agent manualy
	Get the agent id from the bosh vms --details you want ping.  The agent id will be used in the request subject as agent.<agent id>.  If your agent id is 05e3468d-72e1-4796-a871-0c143c25013a then you could send a ping with the following command.

./natsping -s "10.193.67.11:4222" -u nats -p "password" -sub "agent.05e3468d-72e1-4796-a871-0c143c25013a" -m '{"method":"ping","arguments":[], "reply_to": "agent.reply_to_natsping"}'
```

## Example Output

```
natsping -s "10.193.67.11:4222" -u nats -p "d_Nn8fsDbO5zVdCW0CieljAGFvgBn1cg" -sub "agent.05e3468d-72e1-4796-a871-0c143c25013a" -m '{"method":"ping","arguments":[], "reply_to": "agent.reply_to_natsping"}'
Using url nats://nats:d_Nn8fsDbO5zVdCW0CieljAGFvgBn1cg@10.193.67.11:4222
Published [agent.05e3468d-72e1-4796-a871-0c143c25013a] : '{"method":"ping","arguments":[], "reply_to": "agent.reply_to_natsping"}'
[2017-06-16 18:22:46.172985421 -0500 CDT]Received on subject [agent.05e3468d-72e1-4796-a871-0c143c25013a]: '{"method":"ping","arguments":[], "reply_to": "agent.reply_to_natsping"}'
[2017-06-16 18:22:46.173591726 -0500 CDT]Received on subject [agent.reply_to_natsping]: '{"value":"pong"}'
Reply received successfully!
```

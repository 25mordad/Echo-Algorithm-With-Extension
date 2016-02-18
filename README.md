# Echo algorithm with extinction
An election algorithm for undirected networks
At any time, each process takes part in at most one wave. Suppose a process p that is participating in a wave tagged with q is hit by a wave tagged with r.
• If q<r, then p makes the sender its parent, changes to the wave tagged with r (it abandons all the wave messages it received earlier), and treats the incoming message accordingly.
• If q>r, then p continues with the wave tagged with q (it purges the incoming message).
• If q = r, then p treats the incoming message according to the echo algorithm of the wave tagged with q.
If the wave tagged with p completes, by executing a decide event at p, then p becomesthe leader.
- Distributed Algorithms by Wan Fokkink-Page79

# configuration file:
- First line is the name of root for example: 127.0.0.1:8082, it means the ip address is 127.0.0.1 and the port is 8082
- After first line you should write the neighbor in each line with ip and port for example if the node has two neighbors, we write down:
127.0.0.1:8083
127.0.0.1:8084
- If the node is initiator, we should write: initiator:[nodeIP]:[nodePort], for example: initiator:127.0.0.1:8082

# node
How can I add a node?
- Just copy the node folder and rename it to what you want. Then change the configuration file as you need.
- All node files (node.go) are the same.

#Graph for Sample:
	84 --------- 81 ---------- 83 -------- 85
				  \           /
				   \         /
				    \       /
				     \     /
				      \   /
				        82 
			




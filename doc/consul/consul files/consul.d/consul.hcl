# The datacenter in which the agent is running
datacenter = "dc1"
# The data directory for the agent to store state
data_dir = "/opt/consul"
# Specifies the secret key to use for encryption of Consul network traffic
encrypt = "Luj2FZWwlt8475wD1WtwUQ=="
# Address of another agent to join upon starting up
retry_join = ["192.168.1.89"]
performance {
  raft_multiplier = 1
}

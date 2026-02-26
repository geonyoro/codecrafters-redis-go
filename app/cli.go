package main

type CliArgs struct {
	Port      int
	Host      string
	ReplicaOf string
}

func (c *CliArgs) WithPort(port int) *CliArgs {
	c.Port = port
	return c
}

func (c *CliArgs) WithHost(host string) *CliArgs {
	c.Host = host
	return c
}

func (c *CliArgs) WithReplicaOf(replicaOf string) *CliArgs {
	c.ReplicaOf = replicaOf
	return c
}

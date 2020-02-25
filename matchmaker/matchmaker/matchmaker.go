package matchmaker

import (
	"fmt"

	nomad "github.com/hashicorp/nomad/api"
)

func nomadClient() (*nomad.Jobs, error) {
	n, err := nomad.NewClient(nomad.DefaultConfig())
	if err != nil {
		return &nomad.Jobs{}, err
	}
	return n.Jobs(), nil
}

// CreateGameServer func ask Nomad for the creation of a game server
func CreateGameServer() {
	n, err := nomadClient()
	if err != nil {
		panic(err)
	}

	// Check for minecraft-server-register job:
	jobs, _, err := n.List(nil)
	if err != nil {
		panic(err)
	}

	for j := range jobs {
		if j.Name == "minecraft-server-register" {
			fmt.Println("Found")
		} else {
			fmt.Println("Register")
		}
	}

}

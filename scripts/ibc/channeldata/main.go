package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

func main() {
	var channelsResponse ChannelsResponse
	res, err := http.Get("https://api.sifchain.finance/ibc/core/channel/v1/channels")
	if err != nil {
		panic(err)
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}
	res.Body.Close()
	err = json.Unmarshal(body, &channelsResponse)
	if err != nil {
		panic(err)
	}

	var connectionsResponse ConnectionsResponse
	conRes, err := http.Get("https://api.sifchain.finance/ibc/core/connection/v1/connections")
	if err != nil {
		panic(err)
	}
	body, err = ioutil.ReadAll(conRes.Body)
	if err != nil {
		panic(err)
	}
	conRes.Body.Close()
	err = json.Unmarshal(body, &connectionsResponse)
	if err != nil {
		panic(err)
	}

	for _, channel := range channelsResponse.Channels {
		var clientID string
		for _, connection := range connectionsResponse.Connections {
			if connection.ID == channel.ConnectionsHops[0] {
				clientID = connection.ClientID
				break
			}
		}
		clientResponse := GetClientState(clientID)
		fmt.Printf("%s,%s,%s,%s,%s\n", channel.ChannelID, channel.Counterparty.ChannelID, channel.ConnectionsHops[0], clientID, clientResponse.ClientState.ChainID)
	}
}

type ChannelsResponse struct {
	Channels []struct {
		ChannelID    string `json:"channel_id"`
		Counterparty struct {
			ChannelID string `json:"channel_id"`
		} `json:"counterparty"`
		ConnectionsHops []string `json:"connection_hops"`
	}
}

type ConnectionsResponse struct {
	Connections []struct {
		ID       string `json:"id"`
		ClientID string `json:"client_id"`
	} `json:"connections"`
}

type ClientResponse struct {
	ClientState struct {
		ChainID string `json:"chain_id"`
	} `json:"client_state"`
}

func GetClientState(clientID string) ClientResponse {
	var clientResponse ClientResponse
	clientsRes, err := http.Get("https://api.sifchain.finance/ibc/core/client/v1/client_states/" + clientID)
	if err != nil {
		panic(err)
	}
	body, err := ioutil.ReadAll(clientsRes.Body)
	if err != nil {
		panic(err)
	}
	clientsRes.Body.Close()
	err = json.Unmarshal(body, &clientResponse)
	if err != nil {
		panic(err)
	}

	return clientResponse
}

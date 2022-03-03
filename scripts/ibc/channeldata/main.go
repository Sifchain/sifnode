package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
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
	err = json.Unmarshal(body, &connectionsResponse)
	if err != nil {
		panic(err)
	}

	for _, channel := range channelsResponse.Channels {
		var clientId string
		for _, connection := range connectionsResponse.Connections {
			if connection.Id == channel.ConnectionsHops[0] {
				clientId = connection.ClientId
				break
			}
		}
		clientResponse := GetClientState(clientId)
		log.Printf("%s,%s,%s,%s,%s", channel.ChannelId, channel.Counterparty.ChannelId, channel.ConnectionsHops[0], clientId, clientResponse.ClientState.ChainId)
	}
}

type ChannelsResponse struct {
	Channels []struct {
		ChannelId    string `json:"channel_id"`
		Counterparty struct {
			ChannelId string `json:"channel_id"`
		} `json:"counterparty"`
		ConnectionsHops []string `json:"connection_hops"`
	}
}

type ConnectionsResponse struct {
	Connections []struct {
		Id       string `json:"id"`
		ClientId string `json:"client_id"`
	} `json:"connections"`
}

type ClientsResponse struct {
	ClientStates []struct {
		ClientId    string `json:"client_id"`
		ClientState struct {
			ChainId string `json:"chain_id"`
		} `json:"client_state"`
	} `json:"client_states"`
}

type ClientResponse struct {
	ClientState struct {
		ChainId string `json:"chain_id"`
	} `json:"client_state"`
}

func GetClientState(clientId string) ClientResponse {
	var clientResponse ClientResponse
	clientsRes, err := http.Get("https://api.sifchain.finance/ibc/core/client/v1/client_states/" + clientId)
	if err != nil {
		panic(err)
	}
	body, err := ioutil.ReadAll(clientsRes.Body)
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(body, &clientResponse)
	if err != nil {
		panic(err)
	}

	return clientResponse
}

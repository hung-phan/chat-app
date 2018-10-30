package infrastructure

import (
	"sync"
)

type baseClient struct {
	ID            string
	Hub           *ClientHub
	channels      []DataChannel
	channelsMutex sync.RWMutex
}

func (bc *baseClient) GetID() string {
	return bc.ID
}

func (bc *baseClient) GetHub() *ClientHub {
	return bc.Hub
}

func (bc *baseClient) AddListener(ch DataChannel) {
	bc.channelsMutex.Lock()
	defer bc.channelsMutex.Unlock()

	bc.channels = append(bc.channels, ch)
}

func (bc *baseClient) RemoveListener(ch DataChannel) {
	bc.channelsMutex.Lock()
	defer bc.channelsMutex.Unlock()

	for index := len(bc.channels) - 1; index >= 0; index-- {
		channel := bc.channels[index]

		if channel == ch {
			bc.channels = append(bc.channels[:index], bc.channels[index+1:]...)
		}
	}
}

func (bc *baseClient) emit(data []byte) {
	bc.channelsMutex.RLock()
	defer bc.channelsMutex.RUnlock()

	for _, channel := range bc.channels {
		go func(channel DataChannel, data []byte) {
			channel <- data
		}(channel, data)
	}
}
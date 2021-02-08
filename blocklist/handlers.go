package blocklist

import (
	"errors"
	"fmt"
	"github.com/brharrelldev/BlockListAPI/graph/model"
	"sync"
)

func (bl *BlockList) GetIPDetails(ip string) (*model.IPAddress, error) {

	result, found, err := bl.db.RetrieveByIp(ip)
	if err != nil {
		return nil, fmt.Errorf("error retriving ip %v", err)
	}

	if !found {
		return nil, errors.New("no results returned")
	}

	var ipDetails []*model.IPAddress
	for _, res := range result {
		ipd := &model.IPAddress{
			ID:           res.Id,
			Response:     res.Response,
			IPAddress:    res.IPAddress,
			ResponseCode: res.ResponseCode,
			CreatedAt:    res.CreatedAt,
			UpdatedAt:    res.UpdatedAt,
		}

		ipDetails = append(ipDetails, ipd)

	}

	if ipDetails == nil {
		return nil, errors.New("unknown error occured when rendering results")
	}

	return ipDetails[0], nil
}

func (bl *BlockList) Enqueue(ips []*string) ([]*model.Status, error) {

	var err error
	statusChan := make(chan *model.Status)

	var statuses []*model.Status

	wg := &sync.WaitGroup{}
	for _, ip := range ips {

		wg.Add(1)
		go bl.run(wg, *ip, statusChan)
		fmt.Println(err)

		select {
		case status := <-statusChan:
			fmt.Println("received status")
			statuses = append(statuses, status)
			
		}

	}
	wg.Wait()

	return statuses, nil

}

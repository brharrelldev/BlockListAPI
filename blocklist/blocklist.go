package blocklist

import (
	"errors"
	"fmt"
	"github.com/brharrelldev/BlockListAPI/database"
	"github.com/brharrelldev/BlockListAPI/graph/model"
	"github.com/dgryski/trifles/uuid"
	"log"
	"net"
	"strings"
	"sync"
	"time"
)

type BlockList struct {
	db *database.BlockListDB
}

type Result struct {
	Id           string
	IPAddress    string
	ResponseCode string
	Response     string
	UpdatedAt    string
	CreatedAt    string
}

func NewBlockList(db *database.BlockListDB) (*BlockList, error) {

	if db == nil {
		return nil, errors.New("empty db instance passed")
	}

	return &BlockList{db: db}, nil

}

func (bl *BlockList) reverse(ip string) (string, error) {
	var octectSlice []string
	octets := strings.Split(ip, ".")

	if len(octets) < 4 {
		return "", errors.New("invalid ip address format")
	}

	octectSlice = append(octectSlice, octets[3])
	octectSlice = append(octectSlice, octets[2])
	octectSlice = append(octectSlice, octets[1])
	octectSlice = append(octectSlice, octets[0])

	for _, octet := range octectSlice {

		if len(octet) > 3 || len(octet) < 1 {
			return "", errors.New("invalid ip address")
		}
	}

	reverseIp := strings.Join(octectSlice, ".")

	return reverseIp, nil

}

func (bl *BlockList) lookupIp(ip string) (string, string, error) {

	reverse, err := bl.reverse(ip)
	if err != nil {
		return "", "", fmt.Errorf("error attempting to reverse ip %v", err)
	}

	reverseLookup := fmt.Sprintf("%s.zen.spamhaus.org", reverse)
	origLookup := fmt.Sprintf("%s.zen.spamhaus.org", ip)

	resp, err := net.LookupIP(reverseLookup)
	if err != nil {
		if !strings.Contains(err.Error(), "no such host") {
			return "", "", fmt.Errorf("error looking up address %v", err)
		}

		fmt.Println(origLookup)
		resp, err = net.LookupIP(origLookup)
		if err != nil {
			fmt.Println(err)
			return "", "", fmt.Errorf("original IP lookError failed, bailing %v", err)
		}

		respCode := bl.lookError(resp[0].String())

		return resp[0].String(), respCode.Error(), nil
	}

	respCode := bl.lookError(resp[0].String())

	return resp[0].String(), respCode.Error(), nil

}

func (bl *BlockList) run(wg *sync.WaitGroup, ip string,statusChan chan *model.Status) {

	defer wg.Done()


	_, err := bl.GetIPDetails(ip)
	if err != nil {
		if strings.Contains(err.Error(), "no results") {
			resp, respCode, err := bl.lookupIp(ip)
			if err != nil {
				statusChan <- &model.Status{
					IP: ip,
					Message:  err.Error(),
				}


				return
			}

			id := uuid.UUIDv4()
			dbIPInput := database.IPDetail{
				Id:           id,
				IPAddress:    ip,
				UpdatedAt:    time.Now().String(),
				CreatedAt:    time.Now().String(),
				Response:     resp,
				ResponseCode: respCode,
			}

			if err := bl.db.CreateIpDetail(dbIPInput); err != nil {
				statusChan <- &model.Status{
					IP: ip,
					Message:  err.Error(),
				}
				return
			}

			status := &model.Status{
				IP: ip,
				Message: "added to database",
			}

			statusChan <- status


		}

		statusChan <- &model.Status{
			IP:      ip,
			Message: err.Error(),
		}

		return
	}

	log.Printf("updating ip=%s with timestamp=%s", ip, time.Now().String())
	if err := bl.db.UpdateTimestamp(ip, time.Now().String()); err != nil{
		statusChan <- &model.Status{
			IP: ip,
			Message:  err.Error(),
		}
		return
	}

	status := &model.Status{
		IP: ip,
		Message: "timestamp updated",
	}


	statusChan <- status

	return

}

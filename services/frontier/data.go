package main

import (
	"context"
	"fmt"
	"io"
	"sync"
	"time"

	"github.com/stellar/go/historyarchive"
	"github.com/stellar/go/ingest"
	"github.com/stellar/go/xdr"
)

type Data struct {
	archive *historyarchive.Archive

	mutex sync.RWMutex

	rootHAS  historyarchive.HistoryArchiveState
	accounts map[string]Account
}

type Account struct {
	AccountEntry     xdr.AccountEntry
	TrustLineEntries []xdr.TrustLineEntry
}

func NewData(archiveURL string) (*Data, error) {
	archive, err := historyarchive.Connect(
		archiveURL,
		historyarchive.ConnectOptions{},
	)
	if err != nil {
		return nil, err
	}
	data := &Data{
		archive: archive,
		accounts: map[string]Account{},
	}
	return data, nil
}

func (d *Data) Update() error {
	rootHAS, err := d.archive.GetRootHAS()
	if err != nil {
		return err
	}

	if rootHAS == d.rootHAS {
		return nil
	}

	reader, err := ingest.NewCheckpointChangeReader(context.Background(), d.archive, d.rootHAS.CurrentLedger, rootHAS.CurrentLedger)
	if err != nil {
		panic(err)
	}

	accounts := map[string]Account{}

	ts := time.Now()
	fmt.Println("Update start:", ts)

	for {
		entry, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			panic(err)
		}

		switch entry.Type {
		case xdr.LedgerEntryTypeAccount:
			a := entry.Post.Data.Account
			if a == nil {
				break
			}
			da := accounts[a.AccountId.Address()]
			da.AccountEntry = *a
			accounts[a.AccountId.Address()] = da
		case xdr.LedgerEntryTypeTrustline:
			t := entry.Post.Data.TrustLine
			if t == nil {
				break
			}
			da := accounts[t.AccountId.Address()]
			da.TrustLineEntries = append(da.TrustLineEntries, *t)
			accounts[t.AccountId.Address()] = da
		}
	}

	tf := time.Now()
	td := tf.Sub(ts)
	fmt.Println("Update finish:", tf)
	fmt.Println("Update time:", td)

	d.mutex.Lock()
	defer d.mutex.Unlock()

	d.rootHAS = rootHAS
	for k, v := range accounts {
		d.accounts[k] = v
	}

	return nil
}

func (d *Data) RootHAS() historyarchive.HistoryArchiveState {
	d.mutex.RLock()
	defer d.mutex.RUnlock()

	return d.rootHAS
}

func (d *Data) Accounts() map[string]Account {
	d.mutex.RLock()
	defer d.mutex.RUnlock()

	return d.accounts
}

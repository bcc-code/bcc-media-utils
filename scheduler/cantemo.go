package main

import (
	"github.com/bcc-code/bccm-utils/scheduler/cantemo"
	"github.com/bcc-code/mediabank-bridge/log"
)

func getItemIDsFromSearchQuery(client *cantemo.Client, query string) ([]string, error) {
	res, err := client.Search().Put(query, 1)
	if err != nil {
		return nil, err
	}

	var ids []string
	for _, r := range res.Results {
		ids = append(ids, r.VidispineId)
	}

	for res.Page < res.Pages {
		res, err = client.Search().Put(query, res.Page+1)
		if err != nil {
			log.L.Error().Err(err).Send()
		}
		for _, r := range res.Results {
			ids = append(ids, r.VidispineId)
		}
	}

	return ids, nil
}

func iterateThroughIDs(client *cantemo.Client, ids []string) {
	for _, id := range ids {
		info, err := client.Items().Get(id)
		if err != nil {
			log.L.Error().Err(err)
			continue
		}

		log.L.Debug().Str("id", info.Id).Msg("Retrieved id")
	}
}

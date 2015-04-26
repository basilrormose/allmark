// Copyright 2014 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package orchestrator

import (
	"allmark.io/modules/model"
	"allmark.io/modules/web/view/viewmodel"
	"fmt"
)

type FeedOrchestrator struct {
	*Orchestrator
}

func (orchestrator *FeedOrchestrator) GetRootEntry(hostname string) viewmodel.FeedEntry {

	rootItem := orchestrator.rootItem()
	if rootItem == nil {
		orchestrator.logger.Fatal("No root item found.")
	}

	return orchestrator.createFeedEntryModel(hostname, rootItem)
}

func (orchestrator *FeedOrchestrator) GetEntries(hostname string, itemsPerPage, page int) (entries []viewmodel.FeedEntry, found bool) {

	// validate page number
	if page < 1 {
		orchestrator.logger.Fatal("Invalid page number (%v).", page)
	}

	rootItem := orchestrator.rootItem()
	if rootItem == nil {
		orchestrator.logger.Fatal("No root item found")
	}

	feedEntries := make([]viewmodel.FeedEntry, 0)

	latestItems, found := pagedItems(orchestrator.getLatestItems(rootItem.Route()), itemsPerPage, page)
	if !found {
		return feedEntries, false
	}

	for _, item := range latestItems {
		feedEntries = append(feedEntries, orchestrator.createFeedEntryModel(hostname, item))
	}

	return feedEntries, true
}

func (orchestrator *FeedOrchestrator) createFeedEntryModel(hostname string, item *model.Item) viewmodel.FeedEntry {

	addressPrefix := fmt.Sprintf("http://%s/%s/", hostname, item.Route().Value())
	pathProvider := orchestrator.absolutePather(addressPrefix)

	// item location
	location := pathProvider.Path(item.Route().Value())

	// content
	content, err := orchestrator.converter.Convert(orchestrator.getItemByAlias, pathProvider, item)
	if err != nil {
		content = err.Error()
	}

	// creation date
	creationDate := item.MetaData.CreationDate.Format("2006-01-02")

	return viewmodel.FeedEntry{
		Title:       item.Title,
		Description: content,
		Link:        location,
		PubDate:     creationDate,
	}
}
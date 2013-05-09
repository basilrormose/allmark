// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package mapper

import (
	"fmt"
	"github.com/andreaskoch/allmark/converter"
	"github.com/andreaskoch/allmark/path"
	"github.com/andreaskoch/allmark/repository"
	"github.com/andreaskoch/allmark/types"
	"github.com/andreaskoch/allmark/view"
)

func Map(item *repository.Item, repositoryPathProvider *path.Provider) *view.Model {

	// map childs first
	// for _, child := range item.Childs {
	// 	Map(child, repositoryPathProvider) // recurse
	// }

	// paths
	relativePath := item.PathProvider().GetWebRoute(item)
	absolutePath := repositoryPathProvider.GetWebRoute(item)

	// convert the item
	parsedItem, err := converter.Convert(item)
	if err != nil {
		return view.Error(fmt.Sprintf("%s", err), relativePath, absolutePath)
	}

	var model *view.Model

	switch itemType := parsedItem.MetaData.ItemType; itemType {
	case types.DocumentItemType:
		model = createDocumentMapperFunc(parsedItem, relativePath, absolutePath)

	case types.RepositoryItemType, types.CollectionItemType:
		model = createDocumentMapperFunc(parsedItem, relativePath, absolutePath)
		model.SubEntries = getSubModels(item, repositoryPathProvider)

	case types.MessageItemType:
		model = createMessageMapperFunc(parsedItem, relativePath, absolutePath)

	default:
		model = view.Error(fmt.Sprintf("There is no mapper available for items of type %q", itemType), relativePath, absolutePath)
	}

	// assign the model to the item
	item.Model = model

	return model
}

func getSubModels(item *repository.Item, repositoryPathProvider *path.Provider) []*view.Model {

	items := item.Childs
	models := make([]*view.Model, 0)

	for _, child := range items {
		models = append(models, Map(child, repositoryPathProvider))
	}

	return models
}

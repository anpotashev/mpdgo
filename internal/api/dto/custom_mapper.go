package dto

import (
	"github.com/anpotashev/mpdgo/pkg/mpdapi"
)

func MapMpdTree(root mpdapi.DirectoryItem) DirectoryItem {
	result := MapDirectoryItem(root)
	for _, child := range root.Children {
		var childItem TreeItem
		switch child.(type) {
		case *mpdapi.DirectoryItem:
			childItem = MapMpdTree(*child.(*mpdapi.DirectoryItem))
		case *mpdapi.FileItem:
			childItem = MapFileItem(*child.(*mpdapi.FileItem))
		}
		result.Children = append(result.Children, childItem)
	}
	return result
}

func MapSlice[T any, K any](in []T, mapFunc func(T) K) []K {
	var result []K
	for _, elem := range in {
		result = append(result, mapFunc(elem))
	}
	return result
}

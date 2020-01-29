package web

import (
	"gitlab.com/kabestan/repo/baseapp/internal/svc"
)

type (
	AuthCache interface {
		SetService(svc *svc.Service)
		PathPermissionTags(path string) (tags []string, err error)
	}
)

type (
	Cache struct {
		cache   map[string][]string
		service *svc.Service
	}
)

func NewCache() *Cache {
	return &Cache{
		cache: make(map[string][]string),
	}
}

func (c *Cache) SetService(s *svc.Service) {
	c.service = s
}

func (c *Cache) PathPermissionTags(path string) (tags []string, err error) {
	tags, ok := c.cache[path]
	if ok {
		return tags, nil
	}

	tags, err = c.service.ResourcePermissionTagsByPath(path)
	if err != nil {
		return tags, err
	}

	c.cache[path] = tags

	return tags, nil

}

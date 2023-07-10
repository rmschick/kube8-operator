/*
Copyright 2023 The Kubernetes collector-controller Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Code generated by lister-gen. DO NOT EDIT.

package v1

import (
	v1 "github.com/FishtechCSOC/terminal-poc-deployment/pkg/apis/collector/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/tools/cache"
)

// CollectorLister helps list Collectors.
// All objects returned here must be treated as read-only.
type CollectorLister interface {
	// List lists all Collectors in the indexer.
	// Objects returned here must be treated as read-only.
	List(selector labels.Selector) (ret []*v1.Collector, err error)
	// Collectors returns an object that can list and get Collectors.
	Collectors(namespace string) CollectorNamespaceLister
	CollectorListerExpansion
}

// collectorLister implements the CollectorLister interface.
type collectorLister struct {
	indexer cache.Indexer
}

// NewCollectorLister returns a new CollectorLister.
func NewCollectorLister(indexer cache.Indexer) CollectorLister {
	return &collectorLister{indexer: indexer}
}

// List lists all Collectors in the indexer.
func (s *collectorLister) List(selector labels.Selector) (ret []*v1.Collector, err error) {
	err = cache.ListAll(s.indexer, selector, func(m interface{}) {
		ret = append(ret, m.(*v1.Collector))
	})
	return ret, err
}

// Collectors returns an object that can list and get Collectors.
func (s *collectorLister) Collectors(namespace string) CollectorNamespaceLister {
	return collectorNamespaceLister{indexer: s.indexer, namespace: namespace}
}

// CollectorNamespaceLister helps list and get Collectors.
// All objects returned here must be treated as read-only.
type CollectorNamespaceLister interface {
	// List lists all Collectors in the indexer for a given namespace.
	// Objects returned here must be treated as read-only.
	List(selector labels.Selector) (ret []*v1.Collector, err error)
	// Get retrieves the Collector from the indexer for a given namespace and name.
	// Objects returned here must be treated as read-only.
	Get(name string) (*v1.Collector, error)
	CollectorNamespaceListerExpansion
}

// collectorNamespaceLister implements the CollectorNamespaceLister
// interface.
type collectorNamespaceLister struct {
	indexer   cache.Indexer
	namespace string
}

// List lists all Collectors in the indexer for a given namespace.
func (s collectorNamespaceLister) List(selector labels.Selector) (ret []*v1.Collector, err error) {
	err = cache.ListAllByNamespace(s.indexer, s.namespace, selector, func(m interface{}) {
		ret = append(ret, m.(*v1.Collector))
	})
	return ret, err
}

// Get retrieves the Collector from the indexer for a given namespace and name.
func (s collectorNamespaceLister) Get(name string) (*v1.Collector, error) {
	obj, exists, err := s.indexer.GetByKey(s.namespace + "/" + name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.NewNotFound(v1.Resource("collector"), name)
	}
	return obj.(*v1.Collector), nil
}

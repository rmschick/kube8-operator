/*
Package server manages the lifecycle of multiple HTTP/gRPC servers uniformly. This is meant to create HTTP servers and
handle lifecycle management for all of them in a uniform manor while also allowing for controls over what paths are
available in each HTTP server (so that not all paths are available on every port). Currently only HTTP 1.x entrypoints
are supported, because integrations should not need a public endpoint (that we know of).
*/
package server

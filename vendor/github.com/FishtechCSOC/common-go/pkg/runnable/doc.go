/*
Package runnable is meant to simplify managing workers and signals. Runnable is just a worker group that passes a
cancelable context to its workers so they can gracefully shutdown once an acceptable graceful shutdown signal is sent
(ie SIGTERM or SIGINT).
*/
package runnable

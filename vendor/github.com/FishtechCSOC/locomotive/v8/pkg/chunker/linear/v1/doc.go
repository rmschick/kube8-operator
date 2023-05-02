/*
Package linear provides an implementation of chunk.Chunker that process entries until it fills a chunk at which case
it pushed along

linear is meant to handle really large batches with creating large call-stacks (and memory usage) with the trade-off of
being slower
*/
package linear

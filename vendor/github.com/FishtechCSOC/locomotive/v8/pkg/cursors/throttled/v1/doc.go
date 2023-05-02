/*
Package throttled provides an implementation of monorail.Cursor that is meant to wrap an existing checkpoint with throttling logic.

throttled is meant to ensure that calls to underlying checkpoint are limited by throttling logic.
*/
package throttled

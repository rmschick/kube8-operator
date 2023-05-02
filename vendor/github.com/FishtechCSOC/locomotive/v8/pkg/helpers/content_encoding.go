package helpers

import (
	"bytes"
	"compress/gzip"
	"context"
	"fmt"
	"io"
	"strings"
)

const (
	xgzipApplicationHeader = "application/x-gzip"
	gzipApplicationHeader  = "application/gzip"
	xgzipHeader            = "x-gzip"
	gzipHeader             = "gzip"
	xbz2ApplicationHeader  = "application/x-bzip2"
	bz2ApplicationHeader   = "application/bzip2"
	xbz2Header             = "x-bzip2"
	bz2Header              = "bzip2"
)

func Gunzip(ctx context.Context, gzData []byte) ([]byte, error) {
	reader, err := gzip.NewReader(bytes.NewReader(gzData))
	if err != nil {
		return nil, fmt.Errorf("failed to initialize GZip reader: %w", err)
	}

	data, err := io.ReadAll(reader)
	if err != nil {
		return nil, fmt.Errorf("IO error while reading GZip data: %w", err)
	}

	return data, nil
}

func CheckForGZip(contentEncoding, contentType, objectKey string) bool {
	switch {
	case checkForGZipHeader(contentEncoding) || checkForGZipHeader(contentType):
		return true
	case contentEncoding != "" || contentType != "":
		return false
	default:
		return strings.HasSuffix(objectKey, ".gz")
	}
}

func CheckForGZipWithOverride(contentEncoding, contentType, objectKey, encodingOverride string) bool {
	switch {
	case checkForGZipHeader(encodingOverride):
		return true
	case encodingOverride != "":
		return false
	default:
		return CheckForGZip(contentEncoding, contentType, objectKey)
	}
}

func checkForGZipHeader(header string) bool {
	if header == xgzipApplicationHeader ||
		header == gzipApplicationHeader ||
		header == xgzipHeader ||
		header == gzipHeader {
		return true
	}

	return false
}

func CheckForBZ2(contentEncoding, contentType, objectKey string) bool {
	switch {
	case checkForBZ2Header(contentEncoding) || checkForBZ2Header(contentType):
		return true
	case contentEncoding != "" || contentType != "":
		return false
	default:
		return strings.HasSuffix(objectKey, ".bz2")
	}
}

func CheckForBZ2WithOverride(contentEncoding, contentType, objectKey, encodingOverride string) bool {
	switch {
	case checkForBZ2Header(encodingOverride):
		return true
	case encodingOverride != "":
		return false
	default:
		return CheckForBZ2(contentEncoding, contentType, objectKey)
	}
}

func checkForBZ2Header(header string) bool {
	if header == xbz2ApplicationHeader ||
		header == bz2ApplicationHeader ||
		header == xbz2Header ||
		header == bz2Header {
		return true
	}

	return false
}

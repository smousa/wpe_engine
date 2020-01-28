package account

import (
	"context"
	"encoding/csv"
	"io"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
	"golang.org/x/sync/semaphore"
)

var (
	// ErrInvalidHeader is an error generated when the csv file includes an
	// invalid header
	ErrInvalidHeader = errors.New("invalid header")

	outHeader = []string{
		"Account ID",
		"First Name",
		"Created On",
		"Status",
		"Status Set On",
	}
)

// WPStreamerOption is an option that can be passed into the WPStreamer
type WPStreamerOption func(s *WPStreamer)

// WithMaxConcurrentRequests returns a WPStreamerOption that limits the number
// of max concurrent requests made by the client
func WithMaxConcurrentRequests(n int64) WPStreamerOption {
	return func(s *WPStreamer) {
		s.maxConcurrentRequests = n
	}
}

// WPStreamer is the mechanism which we will transform the results.  It will
// read from the input, look up the record and dump the output.
type WPStreamer struct {
	client                Client
	maxConcurrentRequests int64
}

// NewWPStreamer instantiates a new WPStreamer
func NewWPStreamer(client Client, ops ...WPStreamerOption) *WPStreamer {
	s := &WPStreamer{
		client:                client,
		maxConcurrentRequests: 10, // Default
	}
	for _, op := range ops {
		op(s)
	}
	return s
}

// Stream transforms the input into the desired out format
func (s *WPStreamer) Stream(ctx context.Context, r io.Reader, w io.Writer) error {
	log := logrus.WithContext(ctx)

	cr := csv.NewReader(r)
	cr.FieldsPerRecord = 4

	cw := csv.NewWriter(w)

	// read the header line
	record, err := cr.Read()
	if err != nil {
		return errors.Wrap(err, "could not read header")
	}
	if err := validateInHeader(record); err != nil {
		return errors.Wrap(err, "could not read header")
	}

	// write the header line
	if err := cw.Write(outHeader); err != nil {
		// Errors won't typically show up here because the csv writer requires a
		// call to Flush before data gets written to the outfile.  So, I guess
		// the primary case is if we run out of memory.  If that is an often
		// enough use case, then we can consider flushing after writing a set
		// number of rows
		return errors.Wrap(err, "could not write header")
	}

	// There are a few ways we can implement this next bit, so I will go over
	// the options:
	//
	// 1. Read all the records from the stream and do a GetBulk call to get the
	// accounts relevant to the request.
	// - We can't implement this because there is no GetBulk endpoint, but I
	// like this idea because we don't have to make so many http requests per
	// file.  If the file is huge, then we can chunk the number of rows read
	// in to keep memory usage manageable.
	//
	// 2. Do 1, except pull all the accounts from the /accounts endpoint
	// - This is okay as long as the number of records that come back from
	// /accounts is small.  We could store the results in memory and could
	// stream the input one at a time.
	//
	// 3. Look up accounts row by row.
	// - This is my least favorite implementation because of the overhead of
	// performing multiple http requests, and therefore wouldn't scale well if
	// we have millions of accounts to look up.  I could add a caching layer to
	// the client if I thought that a file would have the same account listed
	// multiple times, but intuitively that doesn't seem like a common enough
	// use case to be worthwhile (unless this were implemented as a server).
	// The upside to this approach is that it is pretty easy to implement and
	// easy things can be worthwhile to do when you know you are going to have
	// to change it later, but don't know how yet.

	g, gctx := errgroup.WithContext(ctx)
	sem := semaphore.NewWeighted(s.maxConcurrentRequests)

	for {
		// read the record from the input
		raw, err := cr.Read()
		if err == io.EOF {
			// wait for all pending processes to finish
			if err := g.Wait(); err != nil {
				return errors.Wrap(err, "could not process data")
			}

			// flush the write buffer and make sure everything is a-ok
			cw.Flush()
			if err := cw.Error(); err != nil {
				return errors.Wrap(err, "could not flush to writer")
			}

			return nil
		} else if err != nil {
			// Wait for all pending processes to finish
			if err := g.Wait(); err != nil {
				// log any additional errors that may otherwise be swallowed
				log.WithError(err).Error("Could not process data")
			}
			return errors.Wrap(err, "could not read row")
		}

		inRecord := InRecord(raw)

		// Acquire a resource that will permit the creation of a http request.
		// We use a semaphore here in order to throttle the number of
		// concurrent requests made to the server.
		err = sem.Acquire(gctx, 1)
		if err != nil {
			return errors.Wrap(g.Wait(), "could not process data")
		}

		g.Go(func() error {
			defer sem.Release(1)

			// prepare the output
			outRecord := []string{
				inRecord.AccountId(),
				inRecord.FirstName(),
				inRecord.CreatedOn(),
				"",
				"",
			}

			// Get the account from the server
			resp, err := s.client.GetAccount(gctx, &GetAccountRequest{
				AccountId: inRecord.AccountId(),
			})
			if err != nil {
				// log an error if there is a problem with a record
				log.WithError(err).WithField("account_id", inRecord.AccountId()).Error("Could not look up account id")
			} else {
				// populate the record with the response from the server
				outRecord[3] = resp.Status
				outRecord[4] = resp.CreatedOn
			}

			// dump the output
			if err := cw.Write(outRecord); err != nil {
				// Again, errors shouldn't really appear here since the real
				// magic doesn't happen until we do a call to Flush
				return errors.Wrap(err, "could not write row")
			}
			return nil
		})
	}
}

func validateInHeader(record []string) error {
	if len(record) != len(inHeader) {
		return ErrInvalidHeader
	}
	for i, header := range inHeader {
		if record[i] != header {
			return ErrInvalidHeader
		}
	}
	return nil
}

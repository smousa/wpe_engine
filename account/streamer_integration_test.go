//+build integration

package account_test

import (
	"bytes"
	"context"
	"encoding/csv"
	"io"

	"github.com/pkg/errors"
	. "github.com/wpe_merge/wpe_merge/account"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("StreamerIntegration", func() {
	var (
		streamer *WPStreamer

		ctx    context.Context
		cancel context.CancelFunc
		client *WPClient
	)

	BeforeEach(func() {
		client = NewWPClient(emulator.URL())
		streamer = NewWPStreamer(client)
		ctx, cancel = context.WithCancel(context.Background())
	})

	AfterEach(func() {
		cancel()
		emulator.ResetData()
	})

	getReader := func(records [][]string) io.Reader {
		// create dummy file
		var dummy bytes.Buffer
		err := csv.NewWriter(&dummy).WriteAll(records)
		Ω(err).ShouldNot(HaveOccurred())
		return bytes.NewReader(dummy.Bytes())
	}

	assertWriter := func(actual []byte, expected [][]string) {
		cr := csv.NewReader(bytes.NewReader(actual))

		// ensure that the header is the first line
		header, err := cr.Read()
		Ω(err).ShouldNot(HaveOccurred())
		Ω(header).Should(Equal([]string{
			"Account ID",
			"First Name",
			"Created On",
			"Status",
			"Status Set On",
		}))

		// match all of the expected rows
		records, err := cr.ReadAll()
		Ω(err).ShouldNot(HaveOccurred())
		Ω(records).Should(ConsistOf(expected))
	}

	It("should return an error if the reader returns an error", func() {
		var (
			r = ErrReader{}
			w = &bytes.Buffer{}
		)
		err := streamer.Stream(ctx, r, w)
		Ω(err).Should(HaveOccurred())
		Ω(w.Len()).Should(BeZero())
	})

	It("should return an error if the reader is missing a header", func() {
		var (
			r = getReader([][]string{
				{"1", "jdoe", "Jane", "2020-01-01"},
				{"2", "bdole", "Bob", "2020-02-02"},
				{"4", "gknight", "Gladys", "2020-03-03"},
			})
			w = &bytes.Buffer{}
		)

		// make a request
		err := streamer.Stream(ctx, r, w)
		Ω(errors.Cause(err)).Should(Equal(ErrInvalidHeader))
		Ω(w.Len()).Should(BeZero())
	})

	It("should skip the row if the server returns an error", func() {
		var (
			r = getReader([][]string{
				{"Account ID", "Account Name", "First Name", "Created On"},
				{"1", "jdoe", "Jane", "2020-01-01"},
			})

			w = &bytes.Buffer{}
		)

		err := streamer.Stream(ctx, r, w)
		Ω(err).ShouldNot(HaveOccurred())

		assertWriter(w.Bytes(), [][]string{
			{"1", "Jane", "2020-01-01", "", ""},
		})
	})

	It("should return an error if the writer returns an error", func() {
		emulator.LoadData(
			&Account{
				AccountId: 1,
				Status:    "good",
				CreatedOn: "2019-12-12",
			},
		)
		var (
			r = getReader([][]string{
				{"Account ID", "Account Name", "First Name", "Created On"},
				{"1", "jdoe", "Jane", "2020-01-01"},
			})

			w = ErrWriter{}
		)

		err := streamer.Stream(ctx, r, w)
		Ω(err).Should(HaveOccurred())

	})

	It("should process all accounts", func() {
		accounts := []*Account{
			{
				AccountId: 1,
				Status:    "good",
				CreatedOn: "2019-12-12",
			}, {
				AccountId: 2,
				Status:    "great",
				CreatedOn: "2019-11-11",
			}, {
				AccountId: 4,
				Status:    "grape",
				CreatedOn: "2019-10-10",
			},
		}
		emulator.LoadData(accounts...)

		var (
			r = getReader([][]string{
				{"Account ID", "Account Name", "First Name", "Created On"},
				{"1", "jdoe", "Jane", "2020-01-01"},
				{"2", "bdole", "Bob", "2020-02-02"},
				{"4", "gknight", "Gladys", "2020-03-03"},
			})

			w = &bytes.Buffer{}
		)

		err := streamer.Stream(ctx, r, w)
		Ω(err).ShouldNot(HaveOccurred())

		assertWriter(w.Bytes(), [][]string{
			{"1", "Jane", "2020-01-01", "good", "2019-12-12"},
			{"2", "Bob", "2020-02-02", "great", "2019-11-11"},
			{"4", "Gladys", "2020-03-03", "grape", "2019-10-10"},
		})
	})
})

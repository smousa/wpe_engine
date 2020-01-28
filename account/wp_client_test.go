package account_test

import (
	"context"

	. "github.com/smousa/wpe_engine/account"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("WpClient", func() {
	var (
		client *WPClient

		ctx    context.Context
		cancel context.CancelFunc
	)

	BeforeEach(func() {
		client = NewWPClient(emulator.URL())
		ctx, cancel = context.WithCancel(context.Background())
	})

	AfterEach(func() {
		cancel()
		emulator.ResetData()
	})

	Context("without accounts", func() {
		It("should return an empty response from accounts", func() {
			resp, err := client.GetAccounts(ctx, &GetAccountsRequest{})
			Ω(err).ShouldNot(HaveOccurred())
			Ω(resp).ShouldNot(BeNil())
			Ω(resp.Results).Should(BeEmpty())
		})

		It("should return an error if looking up a non-existant record", func() {
			resp, err := client.GetAccount(ctx, &GetAccountRequest{
				AccountId: "29",
			})
			Ω(err).Should(HaveOccurred())
			Ω(resp).Should(BeNil())
		})

		It("should return an error if looking up an invalid account id", func() {
			resp, err := client.GetAccount(ctx, &GetAccountRequest{
				AccountId: "aa",
			})
			Ω(err).Should(HaveOccurred())
			Ω(resp).Should(BeNil())
		})
	})

	Context("with accounts", func() {
		var accounts = []*Account{
			{
				AccountId: 1,
				Status:    "good",
				CreatedOn: "01/22/2020",
			}, {
				AccountId: 2,
				Status:    "bad",
				CreatedOn: "01/19/2019",
			},
		}

		BeforeEach(func() {
			emulator.LoadData(accounts...)
		})

		It("should return accounts", func() {
			resp, err := client.GetAccounts(ctx, &GetAccountsRequest{})
			Ω(err).ShouldNot(HaveOccurred())
			Ω(resp).ShouldNot(BeNil())
			Ω(resp.Results).Should(ConsistOf(accounts))
		})

		It("should return an existing account", func() {
			resp, err := client.GetAccount(ctx, &GetAccountRequest{
				AccountId: "1",
			})
			Ω(err).ShouldNot(HaveOccurred())
			Ω(resp).Should(Equal(accounts[0]))
		})

		It("should return an error if an account doesn't exist", func() {
			resp, err := client.GetAccount(ctx, &GetAccountRequest{
				AccountId: "29",
			})
			Ω(err).Should(HaveOccurred())
			Ω(resp).Should(BeNil())
		})
	})
})

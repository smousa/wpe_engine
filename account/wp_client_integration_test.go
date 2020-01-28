//+build integration

package account_test

import (
	"context"
	"fmt"

	. "github.com/wpe_merge/wpe_merge/account"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("WpClientIntegration", func() {
	var (
		client *WPClient

		ctx    context.Context
		cancel context.CancelFunc
	)

	BeforeEach(func() {
		client = NewWPClient("http://interview.wpengine.io")
		ctx, cancel = context.WithCancel(context.Background())
	})

	AfterEach(func() {
		cancel()
	})

	It("should retrieve all accounts", func() {
		resp, err := client.GetAccounts(ctx, &GetAccountsRequest{})
		Ω(err).ShouldNot(HaveOccurred())
		Ω(resp).ShouldNot(BeNil())
		Ω(resp.Results).ShouldNot(BeEmpty())
	})

	Context("finding a particular account", func() {
		var accounts []*Account

		BeforeEach(func() {
			resp, err := client.GetAccounts(ctx, &GetAccountsRequest{})
			Ω(err).ShouldNot(HaveOccurred())
			Ω(resp).ShouldNot(BeNil())
			Ω(resp.Results).ShouldNot(BeEmpty())
			accounts = resp.Results
		})

		It("should return the matching account", func() {
			resp, err := client.GetAccount(ctx, &GetAccountRequest{
				AccountId: fmt.Sprintf("%d", accounts[0].AccountId),
			})
			Ω(err).ShouldNot(HaveOccurred())
			Ω(resp).Should(Equal(accounts[0]))
		})

		It("should return an error if the account is not valid", func() {
			resp, err := client.GetAccount(ctx, &GetAccountRequest{
				AccountId: "haha",
			})
			Ω(err).Should(HaveOccurred())
			Ω(resp).Should(BeNil())
		})

		It("should return an error if the account id doesn't exist", func() {
			resp, err := client.GetAccount(ctx, &GetAccountRequest{
				AccountId: "0",
			})
			Ω(err).Should(HaveOccurred())
			Ω(resp).Should(BeNil())
		})
	})
})

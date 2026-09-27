// Copyright (c) 2025-2026 The Pearl Research Labs
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package main

import (
	"testing"

	"github.com/pearl-research-labs/pearl/node/btcjson"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// newestFirstPage builds a page exactly the way Wallet.ListTransactions
// returns it: newest (fewest confirmations) first. See
// wallet/wallet/wallet.go: "Return newer results first by starting at
// mempool height and working down to the genesis block."
func newestFirstPage() []btcjson.ListTransactionsResult {
	return []btcjson.ListTransactionsResult{
		{TxID: "pending-tx", Category: "send", Amount: -1.5, Confirmations: 0},
		{TxID: "one-conf-tx", Category: "receive", Amount: 2.0, Confirmations: 1},
		{TxID: "two-conf-tx", Category: "generate", Amount: 3229.64, Confirmations: 2},
	}
}

// Guards pearl-research-labs/pearl#321: the Transactions screen must keep the
// wallet's newest-first page order instead of reversing it (which buried the
// pending send at the bottom of the list).
func TestTxPageOptionsPreservesNewestFirst(t *testing.T) {
	page := newestFirstPage()

	opts := txPageOptions(page)

	require.Len(t, opts, len(page))
	for i := range page {
		assert.Equal(t, page[i].TxID, opts[i].Value,
			"option %d must select the %d-th newest transaction", i, i)
	}
	assert.Equal(t, "pending-tx", opts[0].Value,
		"the pending (newest) transaction must be the first row")
}

// Guards pearl-research-labs/pearl#321: the Overview "Recent activity" list
// must show the newest activity at the top.
func TestRecentActivityRowsPreservesNewestFirst(t *testing.T) {
	page := newestFirstPage()

	rows := recentActivityRows(page)

	require.Len(t, rows, len(page))
	for i := range page {
		assert.Equal(t, txRow(page[i]), rows[i],
			"row %d must render the %d-th newest transaction", i, i)
	}
	assert.Equal(t, txRow(page[0]), rows[0],
		"the first row must render the newest transaction")
}

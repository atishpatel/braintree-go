package braintree

import "testing"

func TestTransactionCreate(t *testing.T) {
	tx, err := testGateway.Transaction().Create(&Transaction{
		Type:   "sale",
		Amount: 100.00,
		CreditCard: &CreditCard{
			Number:         testCreditCards["visa"].Number,
			ExpirationDate: "05/14",
		},
	})

	if err != nil {
		t.Fatal(err)
	}
	if tx.Id == "" {
		t.Fatal("Received invalid ID on new transaction")
	}
}

// This test will fail unless you set up your Braintree sandbox account correctly. See TESTING.md for details.
func TestTransactionCreateWhenGatewayRejected(t *testing.T) {
	_, err := testGateway.Transaction().Create(&Transaction{
		Type:   "sale",
		Amount: 2010.00,
		CreditCard: &CreditCard{
			Number:         testCreditCards["visa"].Number,
			ExpirationDate: "05/14",
		},
	})

	if err == nil {
		t.Fatal("Did not receive error when creating invalid transaction")
	}
	if err.Error() != "Card Issuer Declined CVV" {
		t.Fatal(err)
	}
}

func TestFindTransaction(t *testing.T) {
	createdTransaction, err := testGateway.Transaction().Create(&Transaction{
		Type:   "sale",
		Amount: 100.00,
		CreditCard: &CreditCard{
			Number:         testCreditCards["mastercard"].Number,
			ExpirationDate: "05/14",
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	foundTransaction, err := testGateway.Transaction().Find(createdTransaction.Id)
	if err != nil {
		t.Fatal(err)
	}

	if createdTransaction.Id != foundTransaction.Id {
		t.Fatal("transaction ids do not match")
	}
}

func TestFindNonExistantTransaction(t *testing.T) {
	_, err := testGateway.Transaction().Find("bad_transaction_id")
	if err == nil {
		t.Fatal("Did not receive error when finding an invalid tx ID")
	}
	if err.Error() != "Not Found (404)" {
		t.Fatal(err)
	}
}

func TestAllTransactionFields(t *testing.T) {
	tx := &Transaction{
		Type:    "sale",
		Amount:  100.00,
		OrderId: "my_custom_order",
		CreditCard: &CreditCard{
			Number:         testCreditCards["visa"].Number,
			ExpirationDate: "05/14",
			CVV:            "100",
		},
		Customer: &Customer{
			FirstName: "Lionel",
		},
		BillingAddress: &Address{
			StreetAddress: "1 E Main St",
			Locality:      "Chicago",
			Region:        "IL",
			PostalCode:    "60637",
		},
		ShippingAddress: &Address{
			StreetAddress: "1 E Main St",
			Locality:      "Chicago",
			Region:        "IL",
			PostalCode:    "60637",
		},
		Options: &TransactionOptions{
			SubmitForSettlement:              true,
			StoreInVault:                     true,
			AddBillingAddressToPaymentMethod: true,
			StoreShippingAddressInVault:      true,
		},
	}

	tx2, err := testGateway.Transaction().Create(tx)
	if err != nil {
		t.Fatal(err)
	}

	if tx2.Type != tx.Type {
		t.Fail()
	}
	if tx2.Amount != tx.Amount {
		t.Fail()
	}
	if tx2.OrderId != tx.OrderId {
		t.Fail()
	}
	if tx2.Customer.FirstName != tx.Customer.FirstName {
		t.Fail()
	}
	if tx2.BillingAddress.StreetAddress != tx.BillingAddress.StreetAddress {
		t.Fail()
	}
	if tx2.BillingAddress.Locality != tx.BillingAddress.Locality {
		t.Fail()
	}
	if tx2.BillingAddress.Region != tx.BillingAddress.Region {
		t.Fail()
	}
	if tx2.BillingAddress.PostalCode != tx.BillingAddress.PostalCode {
		t.Fail()
	}
	if tx2.ShippingAddress.StreetAddress != tx.ShippingAddress.StreetAddress {
		t.Fail()
	}
	if tx2.ShippingAddress.Locality != tx.ShippingAddress.Locality {
		t.Fail()
	}
	if tx2.ShippingAddress.Region != tx.ShippingAddress.Region {
		t.Fail()
	}
	if tx2.ShippingAddress.PostalCode != tx.ShippingAddress.PostalCode {
		t.Fail()
	}
	if tx2.CreditCard.Token == "" {
		t.Fail()
	}
	if tx2.Customer.Id == "" {
		t.Fail()
	}
	if tx2.Status != "submitted_for_settlement" {
		t.Fail()
	}
}

func TestTransactionCreateFromPaymentMethodCode(t *testing.T) {
	customer, err := testGateway.Customer().Create(&Customer{
		CreditCard: &CreditCard{
			Number:         testCreditCards["discover"].Number,
			ExpirationDate: "05/14",
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	if customer.CreditCards.CreditCard[0].Token == "" {
		t.Fatal("invalid token")
	}

	tx, err := testGateway.Transaction().Create(&Transaction{
		Type:               "sale",
		CustomerID:         customer.Id,
		Amount:             100,
		PaymentMethodToken: customer.CreditCards.CreditCard[0].Token,
	})

	if err != nil {
		t.Fatal(err)
	}
	if tx.Id == "" {
		t.Fatal("invalid tx id")
	}
}

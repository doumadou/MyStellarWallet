package stellar

import (
	"encoding/json"
	"net/http"
	// "io/ioutil"
	"encoding/base64"
	"fmt"
	b "github.com/stellar/go/build"
	"github.com/stellar/go/clients/horizon"
	"github.com/stellar/go/keypair"
	"github.com/stellar/go/xdr"
	"log"
	"strconv"
	"strings"
)

//2017/03/23 22:32:10 SCHGA46CCBGDGYLY5BG5QBQ2SWKZPBDZLMRBITBITI5LFZS43GRROQTJ
//2017/03/23 22:32:10 GB64V42R7ZNZWQRGDKG5VT2C4NROTNWI4DCISJGJSAVDAFXJFL7COM44

var (
	address    = "GB64V42R7ZNZWQRGDKG5VT2C4NROTNWI4DCISJGJSAVDAFXJFL7COM44"
	secret_key = "SCHGA46CCBGDGYLY5BG5QBQ2SWKZPBDZLMRBITBITI5LFZS43GRROQTJ"

	//address = "GDYTWJAQDHU6SWHFYRJP7Q6SF6DKFFUI54YFDDR2N7K77G74YIZ3KDQX"
	//secret_key = "SC4DDNUB2562RZRQ3QGN3JONHFI7XVAE5U4AW5QTQSOJ43NUSFFRZQJ2"

	dest_address = "GDYTWJAQDHU6SWHFYRJP7Q6SF6DKFFUI54YFDDR2N7K77G74YIZ3KDQX"
)

func CreateAccount() (address string, secret_key string) {
	pair, err := keypair.Random()
	if err != nil {
		log.Fatal(err)
	}

	//log.Println(pair.secret())
	address = pair.Address()
	secret_key = pair.Seed()
	return address, secret_key
}

func GetAccoutInfo(address string) (horizon.Account, error) {
	account, err := horizon.DefaultPublicNetClient.LoadAccount(address)
	if err != nil {
		fmt.Println(err)
	}
	return account, err
}

func SendCurreny(from_secret_key string, dest_address string, send_amount string) error {
	from := from_secret_key

	to := dest_address

	tx := b.Transaction(
		b.SourceAccount{from},
		b.AutoSequence{horizon.DefaultPublicNetClient},
		b.PublicNetwork,
		b.Payment(
			b.Destination{to},
			b.NativeAmount{send_amount},
		),
	)

	txe := tx.Sign(from)
	fmt.Println(txe)
	txeB64, err := txe.Base64()

	if err != nil {
		//panic(err)
		return err
	}

	fmt.Println("tx base64: %s", txeB64)

	resp, err := horizon.DefaultPublicNetClient.SubmitTransaction(txeB64)
	if err != nil {
		fmt.Println(resp)
		//panic(err)
		return err
	}

	fmt.Println("transaction posted in ledger:", resp.Ledger)
	return nil
}

func GetBookOffer(issuer string, Code string) (horizon.OrderBookSummary, error) {
	orderBook, err := horizon.DefaultPublicNetClient.LoadOrderBook(horizon.Asset{Type: "native"}, horizon.Asset{"credit_alphanum4", Code, issuer})
	return orderBook, err
}

func TrustIssue(secret_key string, issuer string, code string) error {

	seed := secret_key
	tx := b.Transaction(
		b.SourceAccount{seed},
		b.AutoSequence{horizon.DefaultPublicNetClient},
		b.PublicNetwork,
		b.Trust(code, issuer,
			b.Limit("100000000000"),
			//b.SourceAccount{address},
		),
	)

	txe := tx.Sign(seed)
	txeB64, err := txe.Base64()
	if err != nil {
		fmt.Println(err)
		return err
	}

	fmt.Println("%v", txe)
	fmt.Printf("tx base64: %s\n", txeB64)
	resp, err := horizon.DefaultPublicNetClient.SubmitTransaction(txeB64)
	if err != nil {
		fmt.Println(err)
		return err
	}

	fmt.Println("transaction posted in ledger:", resp.Ledger)
	return nil
}

func TrustIssueChange(secret_key string, issue string, code string) error {
	seed := secret_key
	tx := b.Transaction(
		b.SourceAccount{seed},
		b.AutoSequence{horizon.DefaultPublicNetClient},
		b.PublicNetwork,
		//b.Trust(code, issue, b.Limit("100.000")),
		b.Trust(code, issue),
	)

	txe := tx.Sign(seed)
	txeB64, _ := txe.Base64()

	fmt.Printf("tx base64: %s", txeB64)
	resp, err := horizon.DefaultPublicNetClient.SubmitTransaction(txeB64)
	if err != nil {
		fmt.Println(err)
		return err
	}

	fmt.Println("transaction posted in ledger:", resp.Ledger)
	return nil
}

func RemoveTrustIssue(secret_key, issue string, code string) error {
	//operationSource := "GCVJCNUHSGKOTBBSXZJ7JJZNOSE2YDNGRLIDPMQDUEQWJQSE6QZSDPNU"
	tx := b.Transaction(
		b.SourceAccount{secret_key},
		b.AutoSequence{horizon.DefaultPublicNetClient},
		b.PublicNetwork,
		b.RemoveTrust(
			code,
			issue,
			//SourceAccount{operationSource},
		),
	)

	txe := tx.Sign(secret_key)
	txeB64, _ := txe.Base64()

	fmt.Printf("tx base64: %s", txeB64)
	resp, err := horizon.DefaultPublicNetClient.SubmitTransaction(txeB64)
	if err != nil {
		fmt.Println(err)
		return err
	}

	fmt.Println("transaction posted in ledger:", resp.Ledger)
	return nil
}
func CreateOffer_1(issuser string, code string, price string, amount string) {
	rate := b.Rate{
		Selling: b.NativeAsset(),
		Buying:  b.CreditAsset("CNY", "GAREELUB43IRHWEASCFBLKHURCGMHE5IF6XSE7EXDLACYHGRHM43RFOX"),
		Price:   b.Price("0.0161"),
	}

	seed := secret_key
	tx := b.Transaction(
		b.SourceAccount{seed},
		b.AutoSequence{horizon.DefaultPublicNetClient},
		b.PublicNetwork,
		//b.CreateOffer(rate, "20"),
		//UpdateOffer(rate, "40", OfferID(2)),
		b.DeleteOffer(rate, b.OfferID(2753)),
	)

	txe := tx.Sign(seed)
	txeB64, err := txe.Base64()
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("%v", txe)
	fmt.Printf("tx base64: %s\n", txeB64)
	resp, err := horizon.DefaultPublicNetClient.SubmitTransaction(txeB64)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("transaction posted in ledger:", resp.Ledger)
}

func CancelOffer(secret_key string, issuser string, code string, price string, offerid int) error {
	rate := b.Rate{
		Selling: b.NativeAsset(),
		Buying:  b.CreditAsset(code, issuser),
		Price:   b.Price(price),
	}

	seed := secret_key
	tx := b.Transaction(
		b.SourceAccount{seed},
		b.AutoSequence{horizon.DefaultPublicNetClient},
		b.PublicNetwork,
		b.DeleteOffer(rate, b.OfferID(offerid)),
	)

	txe := tx.Sign(seed)
	txeB64, err := txe.Base64()
	if err != nil {
		fmt.Println(err)
		return err
	}

	fmt.Println("%v", txe)
	fmt.Printf("tx base64: %s\n", txeB64)
	resp, err := horizon.DefaultPublicNetClient.SubmitTransaction(txeB64)
	if err != nil {
		fmt.Println(err)
		return err
	}

	fmt.Println("transaction posted in ledger:", resp.Ledger)
	return nil
}

func CreateOffer(secret_key string, issuser string, code string, price string, amount string, trade_type string) error {

	var rate b.Rate
	var a b.Amount
	if trade_type == "BUY" {
		f_price, err := strconv.ParseFloat(price, 32)
		if err != nil {
			return err
		}

		point_prec := strings.IndexByte(price, '.')
		if point_prec < 0 {
			point_prec = 0
		} else {
			point_prec = len(price) - point_prec
		}

		rate = b.Rate{
			Buying:  b.NativeAsset(),
			Selling: b.CreditAsset(code, issuser),
			Price:   b.Price(strconv.FormatFloat(1/f_price, 'f', point_prec, 32)),
		}
		f_amount, err := strconv.ParseFloat(amount, 32)
		if err != nil {
			return err
		}
		s_amount := strconv.FormatFloat(f_amount*f_price, 'f', 8, 32)

		fmt.Printf("s_amount: %s\n", s_amount)

		a = b.Amount(s_amount)
	} else {
		rate = b.Rate{
			Buying:  b.CreditAsset(code, issuser),
			Selling: b.NativeAsset(),
			Price:   b.Price(price),
		}
		a = b.Amount(amount)
	}

	seed := secret_key
	tx := b.Transaction(
		b.SourceAccount{seed},
		b.AutoSequence{horizon.DefaultPublicNetClient},
		b.PublicNetwork,
		b.CreateOffer(rate, a),
	)

	txe := tx.Sign(seed)
	txeB64, err := txe.Base64()
	if err != nil {
		fmt.Println(err)
		return err
	}

	fmt.Println("%v", txe)
	fmt.Printf("tx base64: %s\n", txeB64)
	resp, err := horizon.DefaultPublicNetClient.SubmitTransaction(txeB64)
	if err != nil {
		fmt.Println(err)
		return err
	}

	fmt.Println("transaction posted in ledger:", resp.Ledger)
	return nil
}

func GetAccoutOffers(address string) (horizon.OffersPage, error) {
	res, err := horizon.DefaultPublicNetClient.LoadAccountOffers(address)
	if err != nil {
		fmt.Println(err)
		return res, err
	}
	for i := 0; i < len(res.Embedded.Records); i++ {
		fmt.Println(res.Embedded.Records[i].ID)
		fmt.Println(res.Embedded.Records[i].PT)
		fmt.Println(res.Embedded.Records[i].Seller)
		fmt.Println(res.Embedded.Records[i].Selling)
		fmt.Println(res.Embedded.Records[i].Buying)
		fmt.Println(res.Embedded.Records[i].Amount)
		fmt.Println(res.Embedded.Records[i].Price)
	}
	return res, err
}

func main_t1() {

	if address == "" {

		pair, err := keypair.Random()
		if err != nil {
			log.Fatal(err)
		}

		//log.Println(pair.secret())
		address = pair.Address()
		secret_key = pair.Seed()
	}

	log.Println(address)
	log.Println(secret_key)

	//horizon.DefaultPublicNetClient.LoadAccount()

	account, err := horizon.DefaultPublicNetClient.LoadAccount(address)
	if err != nil {
		log.Println(err)
		//log.Fatal(err)
	}

	fmt.Println("Balances for account:", address)

	fmt.Println(account)

	for _, balance := range account.Balances {
		log.Println(balance)
		log.Println(balance.Balance)
		log.Println(balance.Limit)
		log.Println(balance.Asset.Type)
		log.Println(balance.Asset.Code)
		log.Println(balance.Asset.Issuer)
	}

	from := secret_key

	to := dest_address

	tx := b.Transaction(
		b.SourceAccount{from},
		b.AutoSequence{horizon.DefaultPublicNetClient},
		b.PublicNetwork,
		b.Payment(
			b.Destination{to},
			b.NativeAmount{"0.1"},
		),
	)

	//SourceAccount{seed},
	//Sequence{1},
	//PublicNetwork,
	//Payment(
	//    Destination{"GAWSI2JO2CF36Z43UGMUJCDQ2IMR5B3P5TMS7XM7NUTU3JHG3YJUDQXA"},
	//    NativeAmount{"50"},
	//),

	txe := tx.Sign(from)
	fmt.Println(txe)
	txeB64, err := txe.Base64()

	if err != nil {
		panic(err)
	}

	fmt.Println("tx base64: %s", txeB64)

	resp, err := horizon.DefaultPublicNetClient.SubmitTransaction(txeB64)
	if err != nil {
		fmt.Println(resp)
		panic(err)
	}

	fmt.Println("transaction posted in ledger:", resp.Ledger)

	//resp, err := http.Get("https://horizon-testnet.stellar.org/friendbot?addr=" + address)
	//if err != nil {
	//    log.Println(err)
	//    //log.Fatal(err)
	//}

	//defer resp.Body.Close()
	//body, err := ioutil.ReadAll(resp.Body)
	//if err != nil {
	//    log.Fatal(err)
	//}

	//fmt.Println(string(body))
}

type Payment struct {
	Id               string
	Paging_token     string
	Source_account   string
	Account          string
	Starting_balance string
	Type             string
	Asset_type       string
	From             string
	To               string
	Amount           string
}

type EmbeddedResult struct {
	Records []Payment
}

type PaymentResult struct {
	Embedded EmbeddedResult `json:"_embedded"`
}

func decodeResponse(resp *http.Response, object interface{}) (err error) {
	defer resp.Body.Close()
	decoder := json.NewDecoder(resp.Body)

	if !(resp.StatusCode >= 200 && resp.StatusCode < 300) {
		horizonError := &horizon.Error{
			Response: resp,
		}
		decodeError := decoder.Decode(&horizonError.Problem)
		if decodeError != nil {
			return decodeError
		}
		return horizonError
	}

	err = decoder.Decode(&object)
	if err != nil {
		return err
	}
	return nil
}

func GetHistroyPayments(address string) ([]Payment, error) {
	resp, err := horizon.DefaultPublicNetClient.HTTP.Get(horizon.DefaultPublicNetClient.URL + "/accounts/" + address + "/payments")
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	var f PaymentResult

	err = decodeResponse(resp, &f)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	fmt.Println(f)
	for i := 0; i < len(f.Embedded.Records); i++ {
		fmt.Println(f.Embedded.Records[i])
	}

	return f.Embedded.Records, nil
}

type Transaction struct {
	Id             string
	Paging_token   string
	Created_at     string
	Source_account string
	Envelope_xdr   string
	Memo_type      string
	//Tx             xdr.TransactionEnvelope
	Op_name          string
	Dest_account     string
	Amount           string
	Price            string
	Asset            string
	Starting_balance int64
}

type TxEmbeddedResult struct {
	Records []Transaction
}

type TransactionResult struct {
	Embedded TxEmbeddedResult `json:"_embedded"`
}

//func parseData(data string) {
//	rawr := strings.NewReader(data)
//	b64r := base64.NewDecoder(base64.StdEncoding, rawr)
//
//	var tx xdr.TransactionEnvelope
//	bytesRead, err := xdr.Unmarshal(b64r, &tx)
//
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	fmt.Printf("This tx has %d operations\n", len(tx.Tx.Operations))
//	fmt.Println(tx.Tx.Operations[0].SourceAccount)
//	fmt.Println(tx.Tx.Operations[0].Body)
//	fmt.Println(tx.Tx.Operations[0].Body.Type)
//	//fmt.Println(tx.Tx.Operations[0].Body.CreateAccountOp      )
//	if tx.Tx.Operations[0].Body.CreateAccountOp != nil {
//		fmt.Println(tx.Tx.Operations[0].Body.CreateAccountOp.Destination)
//		fmt.Println(tx.Tx.Operations[0].Body.CreateAccountOp.StartingBalance)
//	}
//	//fmt.Println(tx.Tx.Operations[0].Body.PaymentOp            )
//	if tx.Tx.Operations[0].Body.PaymentOp != nil {
//		fmt.Println(tx.Tx.Operations[0].Body.PaymentOp.Destination)
//		fmt.Println(tx.Tx.Operations[0].Body.PaymentOp.Destination.Type)
//		fmt.Println(tx.Tx.Operations[0].Body.PaymentOp.Destination.Ed25519)
//		//fmt.Println(base64.NewDecoder(base64.StdEncoding, tx.Tx.Operations[0].Body.PaymentOp.Destination.Ed25519))
//		fmt.Println(tx.Tx.Operations[0].Body.PaymentOp.Asset)
//		fmt.Println(tx.Tx.Operations[0].Body.PaymentOp.Amount)
//	}
//	//fmt.Println(tx.Tx.Operations[0].Body.PathPaymentOp        )
//	if tx.Tx.Operations[0].Body.PathPaymentOp != nil {
//		fmt.Println(tx.Tx.Operations[0].Body.PathPaymentOp.SendAsset)
//		fmt.Println(tx.Tx.Operations[0].Body.PathPaymentOp.SendMax)
//		fmt.Println(tx.Tx.Operations[0].Body.PathPaymentOp.Destination)
//		fmt.Println(tx.Tx.Operations[0].Body.PathPaymentOp.DestAsset)
//		fmt.Println(tx.Tx.Operations[0].Body.PathPaymentOp.DestAmount)
//		fmt.Println(tx.Tx.Operations[0].Body.PathPaymentOp.Path)
//	}
//	//fmt.Println(tx.Tx.Operations[0].Body.ManageOfferOp        )
//	if tx.Tx.Operations[0].Body.ManageOfferOp != nil {
//		fmt.Println(tx.Tx.Operations[0].Body.ManageOfferOp.Selling)
//		fmt.Println(tx.Tx.Operations[0].Body.ManageOfferOp.Buying)
//		fmt.Println(tx.Tx.Operations[0].Body.ManageOfferOp.Amount)
//		fmt.Println(tx.Tx.Operations[0].Body.ManageOfferOp.Price)
//		fmt.Println(tx.Tx.Operations[0].Body.ManageOfferOp.OfferId)
//		// OfferId == 0 撤单
//	}
//	//fmt.Println(tx.Tx.Operations[0].Body.CreatePassiveOfferOp )
//	if tx.Tx.Operations[0].Body.CreatePassiveOfferOp != nil {
//		fmt.Println(tx.Tx.Operations[0].Body.CreatePassiveOfferOp.Selling)
//		fmt.Println(tx.Tx.Operations[0].Body.CreatePassiveOfferOp.Buying)
//		fmt.Println(tx.Tx.Operations[0].Body.CreatePassiveOfferOp.Amount)
//		fmt.Println(tx.Tx.Operations[0].Body.CreatePassiveOfferOp.Price)
//	}
//	//fmt.Println(tx.Tx.Operations[0].Body.SetOptionsOp         )
//	//	fmt.Println(tx.Tx.Operations[0].Body.ChangeTrustOp        )
//	if tx.Tx.Operations[0].Body.ChangeTrustOp != nil {
//		fmt.Println(tx.Tx.Operations[0].Body.ChangeTrustOp.Line)
//		fmt.Println(tx.Tx.Operations[0].Body.ChangeTrustOp.Limit)
//	}
//	// fmt.Println(tx.Tx.Operations[0].Body.AllowTrustOp         )
//	if tx.Tx.Operations[0].Body.AllowTrustOp != nil {
//		fmt.Println(tx.Tx.Operations[0].Body.AllowTrustOp.Trustor)
//		fmt.Println(tx.Tx.Operations[0].Body.AllowTrustOp.Asset)
//		fmt.Println(tx.Tx.Operations[0].Body.AllowTrustOp.Authorize)
//	}
//	// fmt.Println(tx.Tx.Operations[0].Body.Destination          )
//	//fmt.Println(tx.Tx.Operations[0].Body.ManageDataOp         )
//	if tx.Tx.Operations[0].Body.ManageDataOp != nil {
//		fmt.Println(tx.Tx.Operations[0].Body.ManageDataOp.DataName)
//		fmt.Println(tx.Tx.Operations[0].Body.ManageDataOp.DataValue)
//	}
//	// Output: read 192 bytes
//	// This tx has 1 operations
//}

type OrderInfo struct {
	Created_at     string
	Selling_code   string
	Selling_issuer string
	Buying_code    string
	Buying_issuer  string
	Amount         float64
	Price          float64
	OrderId        int64
}

func GetHistroyOps(address string) ([]OrderInfo, error) {
	var result []OrderInfo
	var f TransactionResult
	resp, err := horizon.DefaultPublicNetClient.HTTP.Get(horizon.DefaultPublicNetClient.URL + "/accounts/" + address + "/transactions" + "?cursor=&limit=100&order=desc")
	//resp, err := horizon.DefaultPublicNetClient.HTTP.Get(horizon.DefaultPublicNetClient.URL + "/accounts/" + address + "/transactions")
	if err != nil {
		fmt.Println(err)
		return result, err
	}

	err = decodeResponse(resp, &f)
	if err != nil {
		fmt.Println(err)
		return result, err
	}
	fmt.Println(f)
	for i := 0; i < len(f.Embedded.Records); i++ {
		rawr := strings.NewReader(f.Embedded.Records[i].Envelope_xdr)
		b64r := base64.NewDecoder(base64.StdEncoding, rawr)

		var tx xdr.TransactionEnvelope
		_, err := xdr.Unmarshal(b64r, &tx)

		if err != nil {
			log.Fatal(err)
		} else {
			if tx.Tx.Operations[0].Body.ManageOfferOp != nil {
				var info OrderInfo
				info.Created_at = f.Embedded.Records[i].Created_at
				s_selling := fmt.Sprintf("%v", tx.Tx.Operations[0].Body.ManageOfferOp.Selling)
				l_s := strings.Split(s_selling, "/")
				if len(l_s) == 3 {
					info.Selling_code = l_s[1]
					info.Selling_issuer = l_s[2]
				} else {
					info.Selling_code = "XML"
					info.Selling_issuer = s_selling
				}
				s_buying := fmt.Sprintf("%v", tx.Tx.Operations[0].Body.ManageOfferOp.Buying)
				l_b := strings.Split(s_buying, "/")
				if len(l_b) == 3 {
					info.Buying_code = l_b[1]
					info.Buying_issuer = l_b[2]
				} else {
					info.Buying_code = "XML"
					info.Buying_issuer = s_selling
				}
				info.Amount = float64(int64(tx.Tx.Operations[0].Body.ManageOfferOp.Amount) / 10000000)
				info.Price = float64(tx.Tx.Operations[0].Body.ManageOfferOp.Price.N) / float64(tx.Tx.Operations[0].Body.ManageOfferOp.Price.D)
				info.OrderId = int64(tx.Tx.Operations[0].Body.ManageOfferOp.OfferId)
				// 0=create a new offer, otherwise edit an existing offer
				result = append(result, info)
			}
		}
	}
	return result, err
}

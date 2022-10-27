package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/devAlvinSyahbana/golang-rfq/graph/generated"
	"github.com/devAlvinSyahbana/golang-rfq/graph/model"
	"github.com/devAlvinSyahbana/golang-rfq/service"
	"github.com/lib/pq"
	"gopkg.in/gomail.v2"
)

// CreateRfq is the resolver for the createRFQ field.
func (r *mutationResolver) CreateRfq(ctx context.Context, input model.NewRfq) (*model.Rfq, error) {
	sqlStatement := `INSERT INTO rfq.header(
		"CompanyName", "CompanyAddress", "CompanyWebsite", "QuotationDate", "QuotationNo", "QuotationExpires", "MadeForName", "MadeForAddress", "MadeForPhone", "SentToName", "SentToAddress", "SentToPhone", "Disc", "Tax", "Interest", "SNK")
		VALUES ($1, $2, $3, $4, $5,$6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16) RETURNING *;`
	response := &model.Rfq{}
	items := []*model.Item{}

	err := r.DB.QueryRow(sqlStatement, input.CompanyName, input.CompanyAddress, input.CompanyWebsite, input.QuotationDate, input.QuotationNo, input.QuotationExpires, input.MadeForName, input.MadeForAddress, input.MadeForPhone, input.SentToName, input.SentToAddress, input.SentToPhone, input.Disc, input.Tax, input.Interest, pq.Array(input.Snk)).Scan(
		&response.CompanyName,
		&response.CompanyAddress,
		&response.CompanyWebsite,
		&response.QuotationDate,
		&response.QuotationNo,
		&response.QuotationExpires,
		&response.MadeForName,
		&response.MadeForAddress,
		&response.MadeForPhone,
		&response.SentToName,
		&response.SentToAddress,
		&response.SentToPhone,
		&response.Disc,
		&response.Tax,
		&response.Interest,
		pq.Array(&response.Snk),
		&response.ID)
	for _, item := range input.Items {
		newItem := &model.Item{}

		println(response.ID)
		sqlStatementItem := `INSERT INTO rfq.items(
			"HeaderID", "Nama", "Harga", "Qty")
			VALUES ($1, $2, $3, $4) RETURNING *;`
		err := r.DB.QueryRow(sqlStatementItem, response.ID, item.Nama, item.Harga, item.Qty).Scan(&newItem.HeaderID, &newItem.Nama, &newItem.Harga, &newItem.Qty)
		if err != nil {
			panic(err)
		}
		items = append(items, newItem)
	}
	m := gomail.NewMessage()
	m.SetHeader("From", "noreply@tripatra.com")
	m.SetHeader("To", "alvin9b.tik@gmail.com")
	m.SetHeader("Subject", "New Request For Quotation !")
	m.SetBody("text/html", "New quotation has created with quotation no : <b>"+response.QuotationNo+"</b>")

	d := gomail.NewDialer("smtp.gmail.com", 587, "dev.alvin.syahbana@gmail.com", "ggthhfkoanmhgxuh")

	// Send the email to Bob, Cora and Dan.
	if err := d.DialAndSend(m); err != nil {
		panic(err)
	}

	if err != nil {
		panic(err)
	}
	response.Items = items
	return response, err
}

// Login is the resolver for the login field.
func (r *mutationResolver) Login(ctx context.Context, input model.Login) (*model.LoginResponse, error) {
	user := &model.Login{}
	err := r.DB.QueryRow("SELECT * FROM rfq.user WHERE email = ($1) AND password = ($2)", input.Email, input.Password).Scan(&user.Email, &user.Password)

	if err != nil {
		return nil, fmt.Errorf("Invalid login")
	}
	token, _ := service.JwtGenerate(user.Email)
	response := model.LoginResponse{Token: token}
	return &response, nil
}

// Rfq is the resolver for the RFQ field.
func (r *mutationResolver) Rfq(ctx context.Context, input model.RFQInput) (*model.Rfq, error) {
	response := &model.Rfq{}

	err := r.DB.QueryRow(`SELECT * FROM rfq.header WHERE "ID" = ($1)`, input.ID).Scan(
		&response.CompanyName,
		&response.CompanyAddress,
		&response.CompanyWebsite,
		&response.QuotationDate,
		&response.QuotationNo,
		&response.QuotationExpires,
		&response.MadeForName,
		&response.MadeForAddress,
		&response.MadeForPhone,
		&response.SentToName,
		&response.SentToAddress,
		&response.SentToPhone,
		&response.Disc,
		&response.Tax,
		&response.Interest,
		pq.Array(&response.Snk),
		&response.ID,
	)

	rows, err := r.DB.Query(`SELECT * FROM rfq.items WHERE "HeaderID" = ($1)`, input.ID)
	if err != nil {
		panic(err)
	}
	responseArray := []*model.Item{}
	for rows.Next() {
		responseItem := &model.Item{}
		rows.Scan(&responseItem.HeaderID,
			&responseItem.Nama,
			&responseItem.Harga,
			&responseItem.Qty)
		responseArray = append(responseArray, responseItem)
	}
	response.Items = responseArray
	if err != nil {
		panic(err)
	}
	return response, nil
}

// RFQList is the resolver for the RFQList field.
func (r *queryResolver) RFQList(ctx context.Context) ([]*model.RFQList, error) {
	rows, err := r.DB.Query(`SELECT "ID", "CompanyName","QuotationNo" FROM rfq.header`)
	if err != nil {
		panic(err)
	}
	responseArray := []*model.RFQList{}
	for rows.Next() {
		response := &model.RFQList{}
		rows.Scan(&response.ID,
			&response.CompanyName,
			&response.QuotationNo)
		responseArray = append(responseArray, response)
	}
	return responseArray, nil
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }

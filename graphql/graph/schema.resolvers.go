package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"graphql/graph/auth"
	"graphql/graph/generated"
	"graphql/graph/model"
	"math/rand"
	"strconv"
	"time"
)

func (r *mutationResolver) RequestSignInCode(ctx context.Context, input model.RequestSignInCodeInput) (*model.ErrorPayload, error) {
	if input.Phone == "" {
		return &model.ErrorPayload{Message: "Phone should not be empty"}, nil
	}

	var userId int64
	err := r.DB.NewSelect().
		Column("id").
		Table("users").
		Where("phone = ?", input.Phone).
		Scan(ctx, &userId)

	if err != nil {
		return &model.ErrorPayload{Message: fmt.Sprintf("User with phone %s doesn't exist", input.Phone)}, nil
	}

	seed := rand.NewSource(time.Now().UnixNano())
	code := rand.New(seed).Intn(10000)

	fmt.Println(code)

	values := map[string]interface{}{
		"user_id": strconv.FormatInt(userId, 10),
		"code":    strconv.Itoa(code),
	}

	_, err = r.DB.NewInsert().
		Table("codes").
		Model(&values).
		Exec(ctx)

	return nil, err
}

func (r *mutationResolver) SignInByCode(ctx context.Context, input model.SignInByCodeInput) (model.SignInOrErrorPayload, error) {
	_, err := strconv.Atoi(input.Code)

	if err != nil {
		return &model.ErrorPayload{Message: fmt.Sprintf("Code %s should be a number", input.Code)}, nil
	}

	var data map[string]interface{}
	err = r.DB.NewSelect().
		Table("codes").
		ColumnExpr("users.id as id").
		ColumnExpr("users.phone as phone").
		Join("JOIN users ON users.id = codes.user_id").
		Where("users.phone = ?", input.Phone).
		Where("codes.code = ?", input.Code).
		Scan(ctx, &data)

	if err != nil {
		return &model.ErrorPayload{Message: fmt.Sprintf("Unable to sign in by phone %s and code %s", input.Phone, input.Code)}, nil
	}

	user := model.User{ID: int(data["id"].(int64)), Phone: data["phone"].(string)}

	token, err := uuid.NewUUID()
	values := map[string]interface{}{
		"user_id": strconv.Itoa(user.ID),
		"token":   token.String(),
	}

	_, err = r.DB.NewInsert().
		Table("tokens").
		Model(&values).
		Exec(ctx)

	if err != nil {
		return nil, err
	}

	return model.SignInPayload{Token: token.String(), Viewer: &model.Viewer{User: &user}}, nil
}

func (r *queryResolver) Products(ctx context.Context) ([]*model.Product, error) {
	var products []*model.Product
	var rows []map[string]interface{}

	err := r.DB.NewSelect().
		Table("products").
		Scan(ctx, &rows)

	for _, m := range rows {
		id := m["id"].(int64)
		name := m["name"].(string)
		products = append(products, &model.Product{ID: int(id), Name: name})
	}

	return products, err
}

func (r *queryResolver) Viewer(ctx context.Context) (*model.Viewer, error) {
	token := auth.ForContext(ctx)
	if token == "" {
		return nil, errors.New("token shouldn't be empty")
	}

	var data map[string]interface{}
	err := r.DB.NewSelect().
		ColumnExpr("users.id as id").
		ColumnExpr("users.phone as phone").
		Join("JOIN users ON users.id = tokens.user_id").
		Table("tokens").
		Where("tokens.token = ?", token).
		Scan(ctx, &data)

	if err != nil {
		return nil, errors.New(fmt.Sprintf("authorization failed for token %s", token))
	}
	user := model.User{ID: int(data["id"].(int64)), Phone: data["phone"].(string)}
	return &model.Viewer{User: &user}, nil
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }

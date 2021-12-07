package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"
	"graphql/graph/generated"
	"graphql/graph/model"
	"math/rand"
	"strconv"
	"time"
)

func (r *mutationResolver) RequestSignInCode(ctx context.Context, input model.RequestSignInCodeInput) (*model.ErrorPayload, error) {
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
	panic(fmt.Errorf("not implemented"))
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
	panic(fmt.Errorf("not implemented"))
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }

package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"os"
	"shop-graphql/auth"
	generated1 "shop-graphql/graph/generated"
	"shop-graphql/graph/model"
	"shop-graphql/sms"
	"strconv"
	"time"

	"github.com/google/uuid"
)

func (r *mutationResolver) RequestSignInCode(ctx context.Context, input model.RequestSignInCodeInput) (*model.ErrorPayload, error) {
	if input.Phone == "" {
		return &model.ErrorPayload{Message: "Phone should not be empty"}, nil
	}

	userId, err := r.Repository.GetUserIdByPhone(ctx, input.Phone)

	if err != nil {
		return &model.ErrorPayload{Message: fmt.Sprintf("User with phone %s doesn't exist", input.Phone)}, nil
	}

	seed := rand.NewSource(time.Now().UnixNano())
	code := (rand.New(seed).Intn(10000-1000) + 1000) % 9999

	err = sms.SendSMS(input.Phone, fmt.Sprintf("Your code is %d", code))

	if err != nil {
		fmt.Fprintf(os.Stderr, "Cannot send SMS with code %d: %s\n", code, err)
	}

	err = r.Repository.AddCodeForUser(ctx, code, userId)
	return nil, err
}

func (r *mutationResolver) SignInByCode(ctx context.Context, input model.SignInByCodeInput) (model.SignInOrErrorPayload, error) {
	_, err := strconv.Atoi(input.Code)

	if err != nil {
		return &model.ErrorPayload{Message: fmt.Sprintf("Code %s should be a number", input.Code)}, nil
	}

	user, err := r.Repository.GetUserByCodeAndPhone(ctx, input.Code, input.Phone)

	if err != nil {
		return &model.ErrorPayload{Message: fmt.Sprintf("Unable to sign in by phone %s and code %s", input.Phone, input.Code)}, nil
	}

	token, err := uuid.NewUUID()

	if err != nil {
		return nil, err
	}

	err = r.Repository.AddTokenForUser(ctx, token.String(), user.ID)

	if err != nil {
		return nil, err
	}

	return model.SignInPayload{Token: token.String(), Viewer: &model.Viewer{User: user}}, nil
}

func (r *queryResolver) Products(ctx context.Context) ([]*model.Product, error) {
	return r.Repository.GetAllProducts(ctx)
}

func (r *queryResolver) Viewer(ctx context.Context) (*model.Viewer, error) {
	token := auth.ForContext(ctx)
	if token == "" {
		return nil, errors.New("token shouldn't be empty")
	}

	user, err := r.Repository.GetUserByToken(ctx, token)

	if err != nil {
		return nil, errors.New(fmt.Sprintf("authorization failed for token %s", token))
	}

	return &model.Viewer{User: user}, nil
}

// Mutation returns generated1.MutationResolver implementation.
func (r *Resolver) Mutation() generated1.MutationResolver { return &mutationResolver{r} }

// Query returns generated1.QueryResolver implementation.
func (r *Resolver) Query() generated1.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }

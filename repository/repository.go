package repository

import (
	"context"
	"github.com/uptrace/bun"
	"shop-graphql/graph/model"
	"strconv"
)

type Repository struct {
	DB *bun.DB
}

func (r *Repository) GetUserIdByPhone(ctx context.Context, phone string) (int64, error) {
	var userId int64
	err := r.DB.NewSelect().
		Column("id").
		Table("users").
		Where("phone = ?", phone).
		Scan(ctx, &userId)

	return userId, err
}

func (r *Repository) AddCodeForUser(ctx context.Context, code int, userId int64) error {
	values := map[string]interface{}{
		"user_id": strconv.FormatInt(userId, 10),
		"code":    strconv.Itoa(code),
	}

	_, err := r.DB.NewInsert().
		Table("codes").
		Model(&values).
		Exec(ctx)

	return err
}

func (r *Repository) GetUserByCodeAndPhone(ctx context.Context, code, phone string) (*model.User, error) {
	var data map[string]interface{}
	err := r.DB.NewSelect().
		Table("codes").
		ColumnExpr("users.id as id").
		ColumnExpr("users.phone as phone").
		Join("JOIN users ON users.id = codes.user_id").
		Where("users.phone = ?", phone).
		Where("codes.code = ?", code).
		Scan(ctx, &data)

	if err != nil {
		return nil, err
	}

	user := model.User{ID: int(data["id"].(int64)), Phone: data["phone"].(string)}
	return &user, nil
}

func (r *Repository) AddTokenForUser(ctx context.Context, token string, userId int) error {
	values := map[string]interface{}{
		"user_id": userId,
		"token":   token,
	}

	_, err := r.DB.NewInsert().
		Table("tokens").
		Model(&values).
		Exec(ctx)

	return err
}

func (r *Repository) GetAllProducts(ctx context.Context) ([]*model.Product, error) {
	var products []*model.Product
	var rows []map[string]interface{}

	err := r.DB.NewSelect().
		Table("products").
		Scan(ctx, &rows)

	if err != nil {
		return nil, err
	}

	for _, m := range rows {
		id := m["id"].(int64)
		name := m["name"].(string)
		products = append(products, &model.Product{ID: int(id), Name: name})
	}

	return products, nil
}

func (r *Repository) GetUserByToken(ctx context.Context, token string) (*model.User, error) {
	var data map[string]interface{}
	err := r.DB.NewSelect().
		ColumnExpr("users.id as id").
		ColumnExpr("users.phone as phone").
		Join("JOIN users ON users.id = tokens.user_id").
		Table("tokens").
		Where("tokens.token = ?", token).
		Scan(ctx, &data)

	if err != nil {
		return nil, err
	}

	user := model.User{ID: int(data["id"].(int64)), Phone: data["phone"].(string)}
	return &user, nil
}

package repositories

import (
	"context"
	"goozinshe/models"
	"goozinshe/logger"

	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

type UsersRepository struct {
	db *pgxpool.Pool
}

func NewUsersRepository(conn *pgxpool.Pool) *UsersRepository {
	return &UsersRepository{db: conn}
}

func (r *UsersRepository) FindAll(c context.Context) ([]models.User, error) {
	logger := logger.GetLogger()
	logger.Info("Fetching all users")

	rows, err := r.db.Query(c, "select id, name, email, password_hash from users order by id")
	if err != nil {
		logger.Error("Could not fetch users", zap.Error(err))
		return nil, err
	}

defer rows.Close()

	users := make([]models.User, 0)
	for rows.Next() {
		var user models.User
		if err := rows.Scan(&user.Id, &user.Name, &user.Email, &user.PasswordHash); err != nil {
			logger.Error("Could not scan user row", zap.Error(err))
			return nil, err
		}
		users = append(users, user)
	}

	if rows.Err() != nil {
		logger.Error("Error occurred during rows iteration", zap.Error(err))
		return nil, err
	}

	logger.Info("Successfully fetched users", zap.Int("count", len(users)))
	return users, nil
}

func (r *UsersRepository) FindById(c context.Context, id int) (models.User, error) {
	logger := logger.GetLogger()
	logger.Info("Fetching user by ID", zap.Int("user_id", id))

	var user models.User
	row := r.db.QueryRow(c, "select id, name, email from users where id = $1", id)
	if err := row.Scan(&user.Id, &user.Name, &user.Email); err != nil {
		logger.Error("Could not fetch user", zap.Error(err))
		return models.User{}, err
	}

	logger.Info("Successfully fetched user", zap.Int("user_id", id))
	return user, nil
}

func (r *UsersRepository) FindByEmail(c context.Context, email string) (models.User, error) {
	logger := logger.GetLogger()
	logger.Info("Fetching user by email", zap.String("email", email))

	var user models.User
	row := r.db.QueryRow(c, "select id, name, email, password_hash from users where email = $1", email)
	if err := row.Scan(&user.Id, &user.Name, &user.Email, &user.PasswordHash); err != nil {
		logger.Error("Could not fetch user by email", zap.Error(err))
		return models.User{}, err
	}

	logger.Info("Successfully fetched user by email", zap.String("email", email))
	return user, nil
}

func (r *UsersRepository) Create(c context.Context, user models.User) (int, error) {
	logger := logger.GetLogger()
	logger.Info("Creating new user", zap.String("email", user.Email))

	var id int
	err := r.db.QueryRow(c, "insert into users(name, email, password_hash) values($1, $2, $3) returning id", 
		user.Name, user.Email, user.PasswordHash).Scan(&id)

	if err != nil {
		logger.Error("Could not create user", zap.Error(err))
		return 0, err
	}

	logger.Info("Successfully created user", zap.Int("user_id", id))
	return id, nil
}

func (r *UsersRepository) Update(c context.Context, id int, user models.User) error {
	logger := logger.GetLogger()
	logger.Info("Updating user", zap.Int("user_id", id))

	_, err := r.db.Exec(c, "update users set name=$1, email=$2 where id=$3", user.Name, user.Email, id)
	if err != nil {
		logger.Error("Could not update user", zap.Error(err))
		return err
	}

	logger.Info("Successfully updated user", zap.Int("user_id", id))
	return nil
}

func (r *UsersRepository) ChangePasswordHash(c context.Context, id int, password string) error {
	logger := logger.GetLogger()
	logger.Info("Updating user password", zap.Int("user_id", id))

	_, err := r.db.Exec(c, "update users set password_hash=$1 where id=$2", password, id)
	if err != nil {
		logger.Error("Could not update user password", zap.Error(err))
		return err
	}

	logger.Info("Successfully updated user password", zap.Int("user_id", id))
	return nil
}

func (r *UsersRepository) Delete(c context.Context, id int) error {
	logger := logger.GetLogger()
	logger.Info("Deleting user", zap.Int("user_id", id))

	_, err := r.db.Exec(c, "delete from users where id=$1", id)
	if err != nil {
		logger.Error("Could not delete user", zap.Error(err))
		return err
	}

	logger.Info("Successfully deleted user", zap.Int("user_id", id))
	return nil
}



package repositories

import (
	"context"
	"goozinshe/models"

	"github.com/jackc/pgx/v5/pgxpool"
)

type UsersRepository struct {
	db *pgxpool.Pool
}

func NewUsersRepository(conn *pgxpool.Pool) *UsersRepository {
	return &UsersRepository{db: conn}
}

func (r *UsersRepository) FindAll(c context.Context) ([]models.User, error) {
	rows, err := r.db.Query(c, "select id, name, email, password_hash from users order by id")
	if err != nil {
		return nil, err
	}

	users := make([]models.User, 0)
	for rows.Next() {
		var user models.User
		err := rows.Scan(&user.Id, &user.Name, &user.Email, &user.PasswordHash)
		if err != nil {
			return nil, err
		}

		users = append(users, user)
	}
	if rows.Err() != nil {
		return nil, err
	}

	return users, err
}


func (r *UsersRepository) FindById(c context.Context, id int) (models.User, error){
	var user models.User
	row := r.db.QueryRow(c, "select id, name, email from users where id = $1", id)
	err := row.Scan(&user.Id, &user.Name, &user.Email)
	if err != nil {
		return models.User{}, err
	}

	return user, nil
}

func (r *UsersRepository) Create(c context.Context, user models.User) (int, error) {
	var id int
	err := r.db.QueryRow(c, "insert into users(name, email, password_hash) values($1, $2, $3) returning id", 
							user.Name, user.Email, user.PasswordHash).Scan(&id)

	return id, err

}

func (r *UsersRepository) Update(c context.Context, id int, user models.User) error{
	_, err := r.db.Exec(c, "update users set name=$1, email=$2 where id=$3", user.Name, user.Email, id)


	if err != nil {
		return err
	}

	return nil
}

func (r *UsersRepository) ChangePasswordHash(c context.Context, id int, password string) error{
	_, err := r.db.Exec(c, "update users set password_hash=$1 where id=$2", password, id)

	if err != nil {
		return err
	}

	return nil
}

func (r *UsersRepository) Delete(c context.Context, id int) error{
	_, err := r.db.Exec(c, "delete from users where id=$1", id)
	if err != nil {
		return err
	}
	return nil
}



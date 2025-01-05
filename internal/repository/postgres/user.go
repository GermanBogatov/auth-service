package postgres

import (
	"context"
	"fmt"
	"github.com/GermanBogatov/auth-service/internal/common/apperror"
	"github.com/GermanBogatov/auth-service/internal/common/metrics"
	"github.com/GermanBogatov/auth-service/internal/config"
	"github.com/GermanBogatov/auth-service/internal/entity"
	"github.com/GermanBogatov/auth-service/pkg/postgresql"
	"github.com/GermanBogatov/auth-service/pkg/tracer"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/pkg/errors"
	"strings"
	"time"
)

var _ IUser = &User{}

type IUser interface {
	CreateUser(ctx context.Context, user entity.User) error
	GetUserByID(ctx context.Context, id string) (entity.User, error)
	GetUserByEmailAndPassword(ctx context.Context, email, password string) (entity.User, error)
	DeleteUserByID(ctx context.Context, id string) error
	UpdateUserByID(ctx context.Context, userUpdate entity.UserUpdate) (entity.User, error)
	GetUsers(ctx context.Context, filter entity.Filter) ([]entity.User, error)
	UpdatePrivateUserByID(ctx context.Context, userUpdate entity.UserUpdatePrivate) (entity.User, error)
}

type User struct {
	client postgresql.Client
}

func NewUser(client postgresql.Client) IUser {
	return &User{
		client: client,
	}
}

// CreateUser - создание пользователя
func (u *User) CreateUser(ctx context.Context, user entity.User) error {
	_, span := tracer.StartTrace(ctx, config.SpanPostgresCreateUser)
	defer span.End()
	defer metrics.ObserveRequestDurationPerMethodDB(metrics.Postgres, metrics.CreateUserDb)()

	q := `
	INSERT INTO users 
    	(id,name,surname,email,password,role,created_date) 
    VALUES 
		($1,$2,$3,$4,$5,$6,$7);
		`

	_, err := u.client.Exec(ctx, q, user.ID, user.Name, user.Surname, user.Email, user.Password, user.Role, user.CreatedDate)
	if err != nil {
		metrics.IncRequestTotalDB(metrics.CreateUserDb, metrics.FailStatus)
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == pgerrcode.UniqueViolation {
				return apperror.ErrUserIsExistWithEmail
			}
			return err
		}
	}

	metrics.IncRequestTotalDB(metrics.CreateUserDb, metrics.OkStatus)
	return nil
}

// GetUserByEmailAndPassword - получение пользователя по емайл и паролю
func (u *User) GetUserByEmailAndPassword(ctx context.Context, email, password string) (entity.User, error) {
	_, span := tracer.StartTrace(ctx, config.SpanPostgresGetUserByEmailAndPassword)
	defer span.End()
	defer metrics.ObserveRequestDurationPerMethodDB(metrics.Postgres, metrics.GetUserByEmailAndPasswordDb)()

	q := `
		SELECT id,name,surname,email,password,role,created_date,updated_date 
		FROM users
		WHERE email=$1 AND password=$2;	
		`

	var user entity.User
	err := u.client.QueryRow(ctx, q, email, password).Scan(&user.ID, &user.Name, &user.Surname, &user.Email, &user.Password, &user.Role, &user.CreatedDate, &user.UpdatedDate)
	if err != nil {
		metrics.IncRequestTotalDB(metrics.GetUserByEmailAndPasswordDb, metrics.FailStatus)
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.User{}, apperror.ErrUserNotFound
		}
		return entity.User{}, err
	}

	metrics.IncRequestTotalDB(metrics.GetUserByEmailAndPasswordDb, metrics.OkStatus)
	return user, nil
}

// GetUserByID - получение пользователя по идентификатору
func (u *User) GetUserByID(ctx context.Context, id string) (entity.User, error) {
	_, span := tracer.StartTrace(ctx, config.SpanPostgresGetUserByID)
	defer span.End()
	defer metrics.ObserveRequestDurationPerMethodDB(metrics.Postgres, metrics.GetUserByIDDb)()

	q := `
		SELECT id,name,surname,email,password,role,created_date,updated_date 
		FROM users
		WHERE id=$1;	
		`

	var user entity.User
	err := u.client.QueryRow(ctx, q, id).Scan(&user.ID, &user.Name, &user.Surname, &user.Email, &user.Password, &user.Role, &user.CreatedDate, &user.UpdatedDate)
	if err != nil {
		metrics.IncRequestTotalDB(metrics.GetUserByIDDb, metrics.FailStatus)
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.User{}, apperror.ErrUserNotFound
		}
		return entity.User{}, err
	}

	metrics.IncRequestTotalDB(metrics.GetUserByIDDb, metrics.OkStatus)
	return user, nil
}

func (u *User) DeleteUserByID(ctx context.Context, id string) error {
	_, span := tracer.StartTrace(ctx, config.SpanPostgresDeleteUserByID)
	defer span.End()
	defer metrics.ObserveRequestDurationPerMethodDB(metrics.Postgres, metrics.DeleteUserByIDDb)()

	q := `
	DELETE FROM users 
    WHERE id=$1;`

	_, err := u.client.Exec(ctx, q, id)
	if err != nil {
		metrics.IncRequestTotalDB(metrics.DeleteUserByIDDb, metrics.FailStatus)
		return err
	}

	metrics.IncRequestTotalDB(metrics.DeleteUserByIDDb, metrics.OkStatus)
	return nil
}

// UpdateUserByID - редактирование пользователя
func (u *User) UpdateUserByID(ctx context.Context, userUpdate entity.UserUpdate) (entity.User, error) {
	_, span := tracer.StartTrace(ctx, config.SpanPostgresUpdateUserByID)
	defer span.End()
	defer metrics.ObserveRequestDurationPerMethodDB(metrics.Postgres, metrics.UpdateUserByIDDb)()

	query, args := prepareQueryUpdate(userUpdate)
	var user entity.User
	err := u.client.QueryRow(ctx, query, args...).Scan(&user.ID, &user.Name, &user.Surname, &user.Email, &user.Password, &user.Role, &user.CreatedDate, &user.UpdatedDate)
	if err != nil {
		metrics.IncRequestTotalDB(metrics.UpdateUserByIDDb, metrics.FailStatus)
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == pgerrcode.UniqueViolation {
				return entity.User{}, apperror.ErrUserIsExistWithEmail
			}
			return entity.User{}, err
		}
		return entity.User{}, err
	}

	metrics.IncRequestTotalDB(metrics.UpdateUserByIDDb, metrics.OkStatus)
	return user, nil
}

// prepareQueryUpdate - подготовка запроса для обновления пользователя
func prepareQueryUpdate(user entity.UserUpdate) (string, []interface{}) {
	setValues := make([]string, 0)
	args := make([]interface{}, 0)
	argId := 1
	if user.Name != nil {
		setValues = append(setValues, fmt.Sprintf("name=$%d", argId))
		args = append(args, *user.Name)
		argId++
	}

	if user.Surname != nil {
		setValues = append(setValues, fmt.Sprintf("surname=$%d", argId))
		args = append(args, *user.Surname)
		argId++
	}

	if user.Email != nil {
		setValues = append(setValues, fmt.Sprintf("email=$%d", argId))
		args = append(args, *user.Email)
		argId++
	}

	if user.Password != nil {
		setValues = append(setValues, fmt.Sprintf("password=$%d", argId))
		args = append(args, *user.Password)
		argId++
	}

	setValues = append(setValues, fmt.Sprintf("updated_date=$%d", argId))
	args = append(args, time.Now().UTC())
	argId++

	setQuery := strings.Join(setValues, ", ")
	args = append(args, user.ID)

	query := fmt.Sprintf("UPDATE %s SET %s WHERE id=$%v RETURNING id,name,surname,email,password,role,created_date,updated_date;", "users", setQuery, argId)
	return query, args
}

// GetUsers - получение пользователей
func (u *User) GetUsers(ctx context.Context, filter entity.Filter) ([]entity.User, error) {
	_, span := tracer.StartTrace(ctx, config.SpanPostgresGetUsers)
	defer span.End()
	defer metrics.ObserveRequestDurationPerMethodDB(metrics.Postgres, metrics.GetUsersDb)()
	var q string
	if filter.Role == nil {
		q = fmt.Sprintf(`
			SELECT id,name,surname,email,password,role,created_date,updated_date
			FROM users
			ORDER BY %s %s
			OFFSET %v LIMIT %v;`, filter.Order, filter.Sort, filter.Offset, filter.Limit)
	} else {
		q = fmt.Sprintf(`
			SELECT id,name,surname,email,password,role,created_date,updated_date
			FROM users
			WHERE role = '%s'
			ORDER BY %s %s
			OFFSET %v LIMIT %v;`, *filter.Role, filter.Order, filter.Sort, filter.Offset, filter.Limit)
	}

	rows, err := u.client.Query(ctx, q)
	if err != nil {
		metrics.IncRequestTotalDB(metrics.GetUsersDb, metrics.FailStatus)
		return nil, err
	}

	defer rows.Close()
	users := make([]entity.User, 0, filter.Limit)
	for rows.Next() {
		var user entity.User
		errScan := rows.Scan(&user.ID, &user.Name, &user.Surname, &user.Email, &user.Password, &user.Role, &user.CreatedDate, &user.UpdatedDate)
		if errScan != nil {
			metrics.IncRequestTotalDB(metrics.GetUsersDb, metrics.FailStatus)
			return nil, err
		}
		users = append(users, user)
	}

	metrics.IncRequestTotalDB(metrics.GetUsersDb, metrics.OkStatus)
	return users, nil
}

// UpdatePrivateUserByID - приватное редактирование пользователя
func (u *User) UpdatePrivateUserByID(ctx context.Context, userUpdate entity.UserUpdatePrivate) (entity.User, error) {
	_, span := tracer.StartTrace(ctx, config.SpanPostgresUpdatePrivateUserByID)
	defer span.End()
	defer metrics.ObserveRequestDurationPerMethodDB(metrics.Postgres, metrics.UpdatePrivateUserByIDDb)()

	query, args := prepareQueryUpdatePrivate(userUpdate)
	var user entity.User
	err := u.client.QueryRow(ctx, query, args...).Scan(&user.ID, &user.Name, &user.Surname, &user.Email, &user.Password, &user.Role, &user.CreatedDate, &user.UpdatedDate)
	if err != nil {
		metrics.IncRequestTotalDB(metrics.UpdateUserByIDDb, metrics.FailStatus)
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == pgerrcode.UniqueViolation {
				return entity.User{}, apperror.ErrUserIsExistWithEmail
			}
			return entity.User{}, err
		}
		return entity.User{}, err
	}

	metrics.IncRequestTotalDB(metrics.UpdateUserByIDDb, metrics.OkStatus)
	return user, nil
}

// prepareQueryUpdate - подготовка запроса для обновления пользователя
func prepareQueryUpdatePrivate(user entity.UserUpdatePrivate) (string, []interface{}) {
	setValues := make([]string, 0)
	args := make([]interface{}, 0)
	argId := 1
	if user.Name != nil {
		setValues = append(setValues, fmt.Sprintf("name=$%d", argId))
		args = append(args, *user.Name)
		argId++
	}

	if user.Surname != nil {
		setValues = append(setValues, fmt.Sprintf("surname=$%d", argId))
		args = append(args, *user.Surname)
		argId++
	}

	if user.Email != nil {
		setValues = append(setValues, fmt.Sprintf("email=$%d", argId))
		args = append(args, *user.Email)
		argId++
	}

	if user.Role != nil {
		setValues = append(setValues, fmt.Sprintf("role=$%d", argId))
		args = append(args, *user.Role)
		argId++
	}

	setValues = append(setValues, fmt.Sprintf("updated_date=$%d", argId))
	args = append(args, time.Now().UTC())
	argId++

	setQuery := strings.Join(setValues, ", ")
	args = append(args, user.ID)

	query := fmt.Sprintf("UPDATE %s SET %s WHERE id=$%v RETURNING id,name,surname,email,password,role,created_date,updated_date;", "users", setQuery, argId)
	return query, args
}

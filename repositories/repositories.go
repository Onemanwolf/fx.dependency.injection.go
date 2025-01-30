package repositories

import (
    "context"
    "database/sql"
    "fmt"
    "github.com/microsoft/go-mssqldb/azuread"
    //_ "github.com/denisenkom/go-mssqldb"
    "fx.dependency.injection/models"
)

// UserRepository interface
type UserRepository interface {
    CreateUser(ctx context.Context, user *models.User) error
    GetUserByID(ctx context.Context, id int) (*models.User, error)
    UpdateUser(ctx context.Context, user *models.User) error
    DeleteUser(ctx context.Context, id int) error
}

// AzureSQLUserRepository is an implementation of UserRepository for Azure SQL
type AzureSQLUserRepository struct {
    db *sql.DB
}

var _ UserRepository = (*AzureSQLUserRepository)(nil)

var server = ""
var port = 1433
var database = ""


    // Build connection string


// NewAzureSQLUserRepository creates a new AzureSQLUserRepository
func NewAzureSQLUserRepository() (*AzureSQLUserRepository, error) {
    connString := fmt.Sprintf("server=%s;port=%d;database=%s;fedauth=ActiveDirectoryDefault;", server, port, database)
   // connString := "Server=tcp:domaindatabase.database.windows.net,1433;Initial Catalog=Product;Encrypt=True;TrustServerCertificate=False;Connection Timeout=30;Authentication=\"Active Directory Default\";"
    db, err := sql.Open(azuread.DriverName, connString)
    if err != nil {
        return nil, err
    }

    return &AzureSQLUserRepository{db: db}, nil
}

func (repo *AzureSQLUserRepository) CreateUser(ctx context.Context, user *models.User) error {
    query := "INSERT INTO Users (Name, Email) VALUES (@Name, @Email)"
    _, err := repo.db.ExecContext(ctx, query, sql.Named("Name", user.Name), sql.Named("Email", user.Email))
    return err
}

func (repo *AzureSQLUserRepository) GetUserByID(ctx context.Context, id int) (*models.User, error) {
    query := "SELECT ID, Name, Email FROM Users WHERE ID = @ID"
    row := repo.db.QueryRowContext(ctx, query, sql.Named("ID", id))

    var user models.User
    if err := row.Scan(&user.ID, &user.Name, &user.Email); err != nil {
        return nil, err
    }

    return &user, nil
}

func (repo *AzureSQLUserRepository) UpdateUser(ctx context.Context, user *models.User) error {
    query := "UPDATE Users SET Name = @Name, Email = @Email WHERE ID = @ID"
    _, err := repo.db.ExecContext(ctx, query, sql.Named("Name", user.Name), sql.Named("Email", user.Email), sql.Named("ID", user.ID))
    return err
}

func (repo *AzureSQLUserRepository) DeleteUser(ctx context.Context, id int) error {
    query := "DELETE FROM Users WHERE ID = @ID"
    _, err := repo.db.ExecContext(ctx, query, sql.Named("ID", id))
    return err
}
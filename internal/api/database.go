package api

import (
	"database/sql"
	"log"
	"strings"

	"github.com/google/uuid"
	"github.com/olafstar/salejobs-api/internal/types"
)

type Store struct {
	db *sql.DB
}

func InitDatabase() Store {
	db, err := sql.Open("mysql", "root:root@(127.0.0.1:3306)/test-database?parseTime=true")
	if err != nil {
			log.Fatal(err)
	}
	if err := db.Ping(); err != nil {
			log.Fatal(err)
	}

	// Existing table creation omitted for brevity...

	// Create the advertisements table
	{
			query := `
			CREATE TABLE IF NOT EXISTS advertisements (
					id VARCHAR(36) NOT NULL,
					company_name VARCHAR(255) NOT NULL,
					company_size INT,
					company_website VARCHAR(255),
					company_logo TEXT,
					title VARCHAR(255) NOT NULL,
					experience TEXT NOT NULL,
					skill TEXT NOT NULL,
					description TEXT NOT NULL,
					location_country VARCHAR(100),
					location_city VARCHAR(100),
					location_address TEXT,
					operating_mode VARCHAR(50),
					type_of_work VARCHAR(50),
					apply_email VARCHAR(255),
					apply_url VARCHAR(255),
					consent BOOLEAN,
					contact_name VARCHAR(255),
					contact_email VARCHAR(255),
					contact_phone VARCHAR(50),
					created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
					PRIMARY KEY (id)
			);`
			if _, err := db.Exec(query); err != nil {
					log.Fatal(err)
			}
	}

	{
			query := `
			CREATE TABLE IF NOT EXISTS salary_ranges (
					id INT AUTO_INCREMENT PRIMARY KEY,
					advertisement_id VARCHAR(36) NOT NULL,
					employment_type VARCHAR(100),
					min_salary DECIMAL(10,2),
					max_salary DECIMAL(10,2),
					currency VARCHAR(6),
					FOREIGN KEY (advertisement_id) REFERENCES advertisements(id)
			);`
			if _, err := db.Exec(query); err != nil {
					log.Fatal(err)
			}
	}

	return Store{
			db: db,
	}
}


func ContainsSQLInjection(input string) bool {
	keywords := []string{
		"SELECT", "INSERT", "DELETE", "UPDATE", "DROP", "EXECUTE", "UNION", "FETCH",
		"CHAR", "NCHAR", "VARCHAR", "NVARCHAR", "ALTER",
		"BEGIN", "COMMIT", "ROLLBACK", "CREATE", "DESTROY", "GRANT", "REVOKE", "TRUNCATE",
	}

	for _, keyword := range keywords {
		if strings.Contains(strings.ToUpper(input), keyword) {
			log.Printf("SQL Injection detected due to SQL keyword '%s' in input: %s", keyword, input)
			return true
		}
	}

	return false
}

func (s *Store) CountAdvertisements() (int64, error) {
	var count int64
	query := `SELECT COUNT(*) FROM advertisements`
	err := s.db.QueryRow(query).Scan(&count)
	if err != nil {
			log.Printf("Error querying total advertisements count: %s", err)
			return 0, err
	}
	return count, nil
}

func (s *Store) QueryAdvertisement(id string) (*types.CreateAdvertisementResponse, error) {
	query := `
SELECT id, company_name, company_size, company_website, company_logo, title, experience, skill, description,
			location_country, location_city, location_address, operating_mode, type_of_work, apply_email, apply_url,
			consent, contact_name, contact_email, contact_phone, created_at
FROM advertisements
WHERE id = ?
`
	stmt, err := s.db.Prepare(query)
	if err != nil {
			log.Printf("Error preparing query: %s", err)
			return nil, err
	}
	defer stmt.Close()

	row := stmt.QueryRow(id)
	var ad types.CreateAdvertisementResponse
	err = row.Scan(
			&ad.Id, &ad.Company.Name, &ad.Company.Size, &ad.Company.Website, &ad.Company.Logo,
			&ad.Title, &ad.Experience, &ad.Skill, &ad.Description,
			&ad.Location.Country, &ad.Location.City, &ad.Location.Address,
			&ad.OperatingMode, &ad.TypeOfWork, &ad.ApplyType.Email, &ad.ApplyType.Url,
			&ad.Consent, &ad.Contact.Name, &ad.Contact.Email, &ad.Contact.Phone, &ad.CreatedAt,
	)
	if err != nil {
			log.Printf("Error scanning row: %s", err)
			return nil, err
	}

	// Fetch salary ranges for the specific advertisement
	salaryQuery := `
SELECT employment_type, CAST(min_salary AS SIGNED), CAST(max_salary AS SIGNED), currency
FROM salary_ranges
WHERE advertisement_id = ?
`
	salaryStmt, err := s.db.Prepare(salaryQuery)
	if err != nil {
			log.Printf("Error preparing salary query: %s", err)
			return nil, err
	}
	defer salaryStmt.Close()

	salaryRows, err := salaryStmt.Query(ad.Id)
	if err != nil {
			log.Printf("Error executing salary query: %s", err)
			return nil, err
	}

	for salaryRows.Next() {
			var salary types.SalaryRange
			if err := salaryRows.Scan(&salary.EmploymentType, &salary.MinSalary, &salary.MaxSalary, &salary.Currency); err != nil {
					log.Printf("Error scanning salary row: %s", err)
					return nil, err
			}
			ad.Salary = append(ad.Salary, salary)
	}
	salaryRows.Close()

	if err := row.Err(); err != nil {
			log.Printf("Error during row iteration: %s", err)
			return nil, err
	}

	return &ad, nil
}

func (s *Store) QueryAdvertisementsCards(page, limit int64) ([]types.AdvertisementsCard, error) {
	offset := (page - 1) * limit
	if offset < 0 {
			offset = 0
	}

	query := `
SELECT id, company_name, company_size, company_website, company_logo, title,
		 location_country, location_city, location_address, created_at
FROM advertisements
ORDER BY created_at DESC
LIMIT ? OFFSET ?
`
	stmt, err := s.db.Prepare(query)
	if err != nil {
			log.Printf("Error preparing query: %s", err)
			return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(limit, offset)
	if err != nil {
			log.Printf("Error executing query: %s", err)
			return nil, err
	}
	defer rows.Close()

	var ads []types.AdvertisementsCard
	for rows.Next() {
			var ad types.AdvertisementsCard
			err := rows.Scan(
					&ad.Id, &ad.Company.Name, &ad.Company.Size, &ad.Company.Website, &ad.Company.Logo,
					&ad.Title, &ad.Location.Country, &ad.Location.City, &ad.Location.Address, &ad.CreatedAt,
			)
			if err != nil {
					log.Printf("Error scanning row: %s", err)
					return nil, err
			}

			// Fetch salary ranges for each advertisement
			salaryQuery := `
SELECT employment_type, CAST(min_salary AS SIGNED), CAST(max_salary AS SIGNED), currency
FROM salary_ranges
WHERE advertisement_id = ?
`
			salaryStmt, err := s.db.Prepare(salaryQuery)
			if err != nil {
					log.Printf("Error preparing salary query: %s", err)
					return nil, err
			}

			salaryRows, err := salaryStmt.Query(ad.Id)
			if err != nil {
					log.Printf("Error executing salary query: %s", err)
					return nil, err
			}

			for salaryRows.Next() {
					var salary types.SalaryRange
					if err := salaryRows.Scan(&salary.EmploymentType, &salary.MinSalary, &salary.MaxSalary, &salary.Currency); err != nil {
							log.Printf("Error scanning salary row: %s", err)
							return nil, err
					}
					ad.Salary = append(ad.Salary, salary)
			}
			salaryRows.Close()
			salaryStmt.Close()

			ads = append(ads, ad)
	}

	if err := rows.Err(); err != nil {
			log.Printf("Error during rows iteration: %s", err)
			return nil, err
	}

	return ads, nil
}

func (s *Store) GetJWTUser(username string) (string, error) {
	var passwordHash string
	err := s.db.QueryRow("SELECT password FROM users WHERE username = ?", username).Scan(&passwordHash)
	if err != nil {
		return "", err
	}
	return passwordHash, nil
}

func (s *Store) CreateAdvertisement(ad types.CreateAdvertisementBody) error {
	newUUID, err := uuid.NewUUID()
	if err != nil {
			return err
	}

	stmt, err := s.db.Prepare(`
			INSERT INTO advertisements (
					id, company_name, company_size, company_website, company_logo, title, 
					experience, skill, description, location_country, location_city, 
					location_address, operating_mode, type_of_work, apply_email, 
					apply_url, consent, contact_name, contact_email, contact_phone
			) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?);
	`)
	if err != nil {
			return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(
			newUUID, ad.Company.Name, ad.Company.Size, ad.Company.Website, ad.Company.Logo,
			ad.Title, ad.Experience, ad.Skill, ad.Description,
			ad.Location.Country, ad.Location.City, ad.Location.Address,
			ad.OperatingMode, ad.TypeOfWork, ad.ApplyType.Email, ad.ApplyType.Url,
			ad.Consent, ad.Contact.Name, ad.Contact.Email, ad.Contact.Phone,
	)
	if err != nil {
			return err
	}

	salaryStmt, err := s.db.Prepare(`
			INSERT INTO salary_ranges (advertisement_id, employment_type, min_salary, max_salary, currency) VALUES (?, ?, ?, ?, ?);
	`)
	if err != nil {
			return err
	}
	defer salaryStmt.Close()

	for _, salaryRange := range ad.Salary {
			_, err = salaryStmt.Exec(newUUID, salaryRange.EmploymentType,salaryRange.MinSalary, salaryRange.MaxSalary, salaryRange.Currency)
			if err != nil {
					return err
			}
	}

	return nil
}

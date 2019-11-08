package repository

import (
	. "2019_2_IBAT/internal/pkg/interfaces"
	"fmt"
	"log"
	"strings"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

func (m *DBUserStorage) CreateResume(resumeReg Resume) bool {
	_, err := m.DbConn.Exec("INSERT INTO resumes(id, own_id, first_name, second_name, email, "+
		"region, phone_number, birth_date, sex, citizenship, experience, profession, "+
		"position, wage, education, about)"+
		"VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16);",
		resumeReg.ID, resumeReg.OwnerID, resumeReg.FirstName, resumeReg.SecondName, resumeReg.Email, resumeReg.Region,
		resumeReg.PhoneNumber, resumeReg.BirthDate, resumeReg.Sex, resumeReg.Citizenship, resumeReg.Experience,
		resumeReg.Profession, resumeReg.Position, resumeReg.Wage, resumeReg.Education, resumeReg.About,
	)

	if err != nil {
		fmt.Println("CreateResume: error while creating")
		return false
	}

	return true
}

func (m *DBUserStorage) GetResume(id uuid.UUID) (Resume, error) {

	row := m.DbConn.QueryRowx("SELECT id, own_id, first_name, second_name, email, "+
		"region, phone_number, birth_date, sex, citizenship, experience, profession, "+
		"position, wage, education, about, work_schedule, type_of_employment, email FROM resumes WHERE id = $1;", id,
	)

	var resume Resume
	err := row.StructScan(&resume)
	if err != nil {
		log.Println("GetResume: error while querying")
		return Resume{}, errors.New(InvalidIdMsg)
	}
	log.Println("Storage: GetResume\n Resume:")
	log.Println(resume)

	return resume, nil
}

func (m *DBUserStorage) DeleteResume(id uuid.UUID) error {
	_, err := m.DbConn.Exec("DELETE FROM resumes WHERE id = $1;", id)

	if err != nil {
		fmt.Println("DeleteResume: error while deleting")
		return errors.New(InternalErrorMsg)
	}

	return nil
}

func (m *DBUserStorage) PutResume(resume Resume, userId uuid.UUID, resumeId uuid.UUID) bool {

	_, err := m.DbConn.Exec("UPDATE resumes SET "+
		"first_name = $1, second_name = $2, email = $3, "+
		"region = $4, phone_number = $5, birth_date = $6, sex = $7, citizenship = $8, "+
		"experience = $9, profession = $10, position = $11, wage = $12, education = $13, about = $14 "+
		"WHERE id = $15 AND own_id = $16;",
		resume.FirstName, resume.SecondName, resume.Email, resume.Region, resume.PhoneNumber,
		resume.BirthDate, resume.Sex, resume.Citizenship, resume.Experience, resume.Profession,
		resume.Position, resume.Wage, resume.Education, resume.About, resumeId, userId,
	)

	if err != nil {
		fmt.Println("PutResume: error while changing")
		return false
	}

	return true
}

func (m *DBUserStorage) GetResumes(params map[string]interface{}) ([]Resume, error) {

	resumes := []Resume{}

	log.Printf("Params: %s\n\n", params)
	query := paramsToResumesQuery(params)

	var nmst *sqlx.NamedStmt
	var err error

	if query != "" {
		nmst, err = m.DbConn.PrepareNamed("SELECT id, own_id, first_name, second_name, email, " +
			"region, phone_number, birth_date, sex, citizenship, experience, profession, " +
			"position, wage, education, about, work_schedule, type_of_employment FROM resumes WHERE " + query)

		if err != nil {
			log.Println("GetResumes: error while preparing statement")
			return resumes, errors.New(InternalErrorMsg)
		}
	} else {
		log.Println("GetResumes: query is empty")
	}

	var rows *sqlx.Rows
	if query != "" {
		rows, err = nmst.Queryx(params)
	} else {
		rows, err = m.DbConn.Queryx("SELECT id, own_id, first_name, second_name, email, " +
			"region, phone_number, birth_date, sex, citizenship, experience, profession, " +
			"position, wage, education, about, work_schedule, type_of_employment FROM resumes;",
		)
	}

	if err != nil {
		log.Println("GetVacancies: error while query")
		return resumes, errors.New(InternalErrorMsg)
	}

	for rows.Next() {
		var resume Resume

		_ = rows.StructScan(&resume)

		resumes = append(resumes, resume)
	}

	return resumes, nil
}

func paramsToResumesQuery(params map[string]interface{}) string {
	var query []string

	if params["region"] != nil {
		query = append(query, "region= :region")
	}

	if params["wage_from"] != nil {
		query = append(query, "wage >= :wage_from")
	}

	if params["wage_to"] != nil {
		query = append(query, "wage <= :wage_to")
	}

	if params["experience"] != nil {
		query = append(query, "experience = :experience")
	}

	if params["type_of_employment"] != nil {
		query = append(query, "type_of_employment=:type_of_employment")
	}

	if params["work_schedule"] != nil {
		query = append(query, "work_schedule = :work_schedule")
	}

	str := strings.Join(query, " AND ")

	log.Printf("Query: %s", str)
	return str
}
package models

import (
	"bufio"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gocolly/colly"
	"github.com/likexian/whois-go"
	"github.com/seramirezdev/truora-test/api/entities"
	"io/ioutil"
	"net/http"
	"strings"
)

type DomainModel struct {
	*sql.DB
}

// Devuelve la lista de dominios consultados y guardados en la base de datos
func (domain DomainModel) GetDomains() ([]entities.Domain, error) {
	query := "SELECT domain, servers, servers_changed, ssl_grade, previous_ssl_grade, logo, title, is_down FROM domains ORDER BY id DESC"

	rows, err := domain.DB.Query(query)

	if err != nil {
		fmt.Println("Error en la consulta SQL", err)
		return nil, err
	}

	defer rows.Close()

	domains := []entities.Domain{}

	for rows.Next() {
		domain := entities.Domain{}

		var servers json.RawMessage

		if err := rows.Scan(&domain.Name, &servers, &domain.Servers_changed, &domain.Ssl_grade, &domain.Previous_ssl_grade,
			&domain.Logo, &domain.Title, &domain.Is_down); err != nil {
			fmt.Println("Error leyendo cursor", err)
			return nil, err
		}

		if err := json.Unmarshal(servers, &domain.Servers); err != nil {
			fmt.Println("Error parseando array servers")
			return nil, err
		}

		domains = append(domains, domain)
	}

	return domains, err
}

// Ejecuta las diferentes tarea para obtener la información del dominio consultado
func (domain DomainModel) ConsultDomain(domainSearch string) (entities.Domain, error) {

	// Canales para obtener la información
	dataSSLLab := make(chan entities.DataDomain)
	dataWhoIs := make(chan []string)
	dataHTML := make(chan []string)
	dataDB := make(chan entities.Domain)

	// Se ejecutan la Goritunas para acelerar el proceso de obtención de datos
	go getDataWhoIs(domainSearch, dataWhoIs)
	go getDataHTML(domainSearch, dataHTML)
	go getDataSSLLab(domainSearch, dataSSLLab)
	go getDomainDB(domainSearch, domain.DB, dataDB)

	domainDB := entities.Domain{}
	html := []string{}
	dataOrg := []string{}
	dataDomain := entities.DataDomain{}

	ends := []bool{false, false, false, false}
	// Ciclo para esperar a que se terminen todas las Gorutinas
	for {
		select {
		case domainDB = <-dataDB:
			ends[0] = true
		case html = <-dataHTML:
			ends[1] = true
		case dataOrg = <-dataWhoIs:
			ends[2] = true
		case dataDomain = <-dataSSLLab:
			ends[3] = true
		}

		if ends[0] && ends[1] && ends[2] && ends[3] {
			break
		}
	}

	// Objeto que se usará para devolver la información obtenida
	consultDomain := entities.Domain{}

	// Si es la primera vez que se consulta el dominio entonces se crea en la base de datos
	if domainDB.Title == "" {

		if dataDomain.Status == "ERROR" {
			return consultDomain, errors.New("Error, Dominio inválido")
		}

		consultDomain = entities.Domain{
			Name:               domainSearch,
			Servers:            []entities.Server{},
			Servers_changed:    false,
			Ssl_grade:          "",
			Previous_ssl_grade: "",
			Logo:               html[1],
			Title:              html[0],
			Is_down:            false,
		}

		if err := insertNewDomainDB(consultDomain, domain.DB); err != nil {
			fmt.Println("Error insertando dominio\n", err)
		}

	} else { // Si ya se a guardado en la base de datos entonces se actualiza con la nueva información

		owner := dataOrg[0]
		country := dataOrg[1]

		servers := getListServers(dataDomain, country, owner)
		sslGrade := getSSLGradeDomain(servers)
		serversChanged := false

		if (sslGrade != domainDB.Previous_ssl_grade) || (len(servers) != len(domainDB.Servers)) {
			serversChanged = true
		}

		consultDomain = entities.Domain{
			Servers:            servers,
			Servers_changed:    serversChanged,
			Ssl_grade:          sslGrade,
			Previous_ssl_grade: domainDB.Ssl_grade,
			Title:              html[0],
			Logo:               html[1],
			Is_down:            false,
		}

		if err := updateDomainDB(consultDomain, domainDB.ID, domain.DB); err != nil {
			fmt.Println("Error actualizando dominio\n", err)
		}
	}

	return consultDomain, nil
}

// Consulta en la base de datos si ya se a consultado un dominio previamente
func getDomainDB(domainSearch string, db *sql.DB, dataDB chan entities.Domain) {
	query := "SELECT * FROM domains WHERE domain = $1"

	domain := entities.Domain{}
	var servers json.RawMessage

	if row := db.QueryRow(query, domainSearch); row != nil {

		if err := row.Scan(&domain.ID, &domain.Name, &servers, &domain.Servers_changed, &domain.Ssl_grade,
			&domain.Previous_ssl_grade, &domain.Logo, &domain.Title, &domain.Is_down); err == nil {

			if err := json.Unmarshal(servers, &domain.Servers); err != nil {
				fmt.Println("Error parseando array servers")
			}
		}

	} else {
		fmt.Println("Error en consulta SQL\n")
	}

	dataDB <- domain
}

// Consume el api de SSLLab para traer la información de los servidores
func getDataSSLLab(domainSearch string, dataSSLLab chan entities.DataDomain) {

	response, err := http.Get("https://api.ssllabs.com/api/v3/analyze?host=" + domainSearch)

	if err != nil {
		fmt.Println("No es encuentra la ruta", err)
		return
	}

	responseData, _ := ioutil.ReadAll(response.Body)

	if err := response.Body.Close(); err != nil {
		fmt.Println("Error leyendo respuesta de SSLLab", err)
		return
	}

	var responseObject entities.DataDomain

	if err := json.Unmarshal(responseData, &responseObject); err != nil {
		fmt.Println("Error parseando json de SSLLab", err)
		return
	}

	dataSSLLab <- responseObject
}

// Se usa la libreria Whois para consultar la información del dominio
func getDataWhoIs(domainSearch string, dataWhoIs chan []string) {
	data := []string{"", ""}

	result, err := whois.Whois(domainSearch)
	if err != nil {
		fmt.Println("Error obteniendo datos de WhoIs", err)
		return
	}

	registrar := "Registrant Organization:"
	country := "Registrant Country:"

	scanner := bufio.NewScanner(strings.NewReader(result))

	for scanner.Scan() {
		text := strings.TrimSpace(scanner.Text())
		if strings.Contains(text, registrar) && data[0] == "" {
			data[0] = strings.TrimSpace(strings.Split(text, ":")[1])
		} else if strings.Contains(text, country) && data[1] == "" {
			data[1] = strings.TrimSpace(strings.Split(text, ":")[1])
		}

		if data[0] != "" && data[1] != "" {
			break
		}
	}

	dataWhoIs <- data
}

// Trae el title y el logo del dominio
func getDataHTML(domainSearch string, dataHTML chan []string) {
	data := []string{"", ""}

	c := colly.NewCollector()

	c.OnHTML("head>title", func(e *colly.HTMLElement) {
		data[0] = e.Text
	})

	c.OnHTML("head>link", func(e *colly.HTMLElement) {
		if e.Attr("rel") == "shortcut icon" || e.Attr("rel") == "apple-touch-icon" || e.Attr("rel") == "icon" && data[1] == "" {
			data[1] = e.Attr("href")
		}
	})

	c.OnHTML("head>meta", func(e *colly.HTMLElement) {
		if e.Attr("rel") == "og:image" && data[1] == "" {
			data[1] = e.Attr("href")
		} else if e.Attr("itemprop") == "image" || e.Attr("property") == "og:image" && data[1] == "" {
			data[1] = e.Attr("content")
		}
	})

	if err := c.Visit("http://" + domainSearch); err != nil {
		fmt.Println("Ruta invalida: ", err)
		dataHTML <- []string{"", ""}
		return
	}

	dataHTML <- data
}

// Inserta un nuevo dominio en la base de datos
func insertNewDomainDB(newDomain entities.Domain, db *sql.DB) error {
	insert := "INSERT INTO domains (domain, title, logo, is_down, servers, servers_changed, ssl_grade, previous_ssl_grade) " +
		"VALUES($1, $2, $3, $4, $5, $6, $7, $8)"

	servers, _ := json.Marshal(newDomain.Servers)

	if _, err := db.Exec(insert, newDomain.Name, newDomain.Title, newDomain.Logo, newDomain.Is_down, string(servers),
		newDomain.Servers_changed, newDomain.Ssl_grade, newDomain.Previous_ssl_grade); err != nil {
		return err
	}

	return nil
}

// Actualiza un dominio en la base de datos
func updateDomainDB(updateDomain entities.Domain, id string, db *sql.DB) error {
	insert := "UPDATE domains SET title = $1, logo = $2, is_down = $3, servers = $4, servers_changed = $5, " +
		"ssl_grade = $6, previous_ssl_grade = $7 where id = $8"

	servers, _ := json.Marshal(updateDomain.Servers)

	if _, err := db.Exec(insert, updateDomain.Title, updateDomain.Logo, updateDomain.Is_down, string(servers),
		updateDomain.Servers_changed, updateDomain.Ssl_grade, updateDomain.Previous_ssl_grade, id); err != nil {
		return err
	}

	return nil
}

// Se saca la información necesaria de los servidores obtenidos desde SSLLab
func getListServers(dataDomain entities.DataDomain, country, owner string) []entities.Server {

	servers := []entities.Server{}

	for _, value := range dataDomain.EndPoints {

		/*address := ""
		if value.ServerName != "" {
			address = " - " + value.ServerName
		}*/
		current := entities.Server{
			Address:   value.IpAddress,
			Ssl_grade: value.Grade,
			Country:   country,
			Owner:     owner,
		}

		servers = append(servers, current)
	}

	return servers
}

// Obtener cual es el menor grado de SSL de un dominio
func getSSLGradeDomain(servers []entities.Server) string {

	grades := map[int]string{1: "M", 2: "T", 3: "A-F", 4: "A-", 5: "A", 6: "A+"}
	sslGrade := 7

	// Recorrer todos los servidores del dominio
	for _, server := range servers {
		for index, val := range grades {
			if server.Ssl_grade == val {
				if index < sslGrade {
					sslGrade = index
				}
				break
			}
		}
	}

	return grades[sslGrade]
}

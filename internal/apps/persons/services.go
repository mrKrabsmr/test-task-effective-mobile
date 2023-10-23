package persons

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/mrKrabsmr/commerce-edu-api/internal/apps"
	"github.com/sirupsen/logrus"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"sync"
)

type service struct {
	dao    *dao
	logger *logrus.Logger
}

func newService() *service {
	return &service{
		dao:    newDAO(),
		logger: core.GetLogger(),
	}
}

func (s *service) getFormatList(data url.Values) ([]*Person, *core.Paginate, error) {
	var paginate *core.Paginate
	var p *Person
	var err error

	pageStr := data.Get("page")
	limitStr := data.Get("limit")

	if pageStr != "" || limitStr != "" {
		var page, limit int

		if pageStr != "" {
			page, err = strconv.Atoi(pageStr)
			if err != nil {
				return nil, nil, err
			}
		} else {
			page = 1
		}

		if limitStr != "" {
			limit, err = strconv.Atoi(limitStr)
			if err != nil {
				return nil, nil, err
			}
		} else {
			limit = 10
		}

		paginate = core.NewPaginate(page, limit)
	}

	search := data.Get("search")

	if data.Has("name") || data.Has("surname") ||
		data.Has("patronymic") || data.Has("age") ||
		data.Has("gender") || data.Has("nation") {

		age := -1
		if data.Has("age") {
			age, err = strconv.Atoi(data.Get("age"))
			if err != nil {
				return nil, nil, err
			}
		}

		p = &Person{
			Name:       data.Get("name"),
			Surname:    data.Get("surname"),
			Patronymic: data.Get("patronymic"),
			Age:        age,
			Gender:     data.Get("gender"),
			Nation:     data.Get("nation"),
		}
	}

	s.logger.Debug("SELECT with params", p, search, paginate)

	persons, err := s.dao.getSearchAndFilter(p, search, paginate)
	if err != nil {
		return nil, nil, err
	}

	return persons, paginate, nil
}

func (s *service) createObject(data *PersonRequestDTO) error {
	result := s.getResult(data.Name)
	if result.Err != nil {
		return result.Err
	}

	s.logger.Debug(result)

	p := &Person{
		ID:         uuid.New(),
		Name:       data.Name,
		Surname:    data.Surname,
		Patronymic: data.Patronymic,
		Age:        result.Age,
		Gender:     result.Gender,
		Nation:     result.Nation,
	}

	s.logger.Debug("CREATE OBJECT", p)

	if err := s.dao.create(p); err != nil {
		return err
	}

	return nil
}

func (s *service) updateObject(id uuid.UUID, data *Person) (bool, error) {
	person, err := s.dao.getOne(id)
	if err != nil {
		return false, err
	}

	if person == nil {
		return false, nil
	}

	s.logger.Debug("UPDATE OBJECT", person, data)

	if data.Name != "" {
		person.Name = data.Name
	}

	if data.Surname != "" {
		person.Surname = data.Surname
	}

	if data.Patronymic != "" {
		person.Patronymic = data.Patronymic
	}

	if data.Age != 0 {
		person.Age = data.Age
	}

	if data.Gender != "" {
		person.Gender = data.Gender
	}

	if data.Nation != "" {
		person.Nation = data.Nation
	}

	if err = s.dao.update(person); err != nil {
		return false, err
	}

	return true, nil
}

func (s *service) deleteObject(id uuid.UUID) (bool, error) {
	person, err := s.dao.getOne(id)
	if err != nil {
		return false, err
	}

	if person == nil {
		return false, nil
	}

	if err = s.dao.delete(id); err != nil {
		return false, err
	}

	return true, nil
}

func (s *service) getResult(name string) *ResponseResult {
	res := &ResponseResult{}
	mu, wg := sync.Mutex{}, sync.WaitGroup{}
	wg.Add(3)

	go func() {
		defer wg.Done()
		mu.Lock()
		defer mu.Unlock()
		age, err := s.getAge(name)
		if err != nil {
			if res.Err == nil {
				res.Err = err
			}
			return
		}

		res.Age = age
	}()

	go func() {
		defer wg.Done()
		mu.Lock()
		defer mu.Unlock()
		gender, err := s.getGender(name)
		if err != nil {
			if res.Err == nil {
				res.Err = err
			}
			return
		}

		res.Gender = gender
	}()

	go func() {
		defer wg.Done()
		mu.Lock()
		defer mu.Unlock()
		nation, err := s.getNation(name)
		if err != nil {
			if res.Err == nil {
				res.Err = err
			}
			return
		}

		res.Nation = nation
	}()

	wg.Wait()

	return res
}

func (s *service) getAge(name string) (int, error) {
	u := fmt.Sprintf("https://api.agify.io/?name=%s", name)
	data, err := s.getData(u)
	if err != nil {
		return -1, err
	}

	age, ok := data["age"]
	if !ok {
		return -1, fmt.Errorf("no age available")
	}

	a := age.(float64)

	return int(a), nil
}

func (s *service) getGender(name string) (string, error) {
	u := fmt.Sprintf("https://api.genderize.io/?name=%s", name)
	data, err := s.getData(u)
	if err != nil {
		return "", err
	}

	gender, ok := data["gender"]
	if !ok {
		return "", fmt.Errorf("no gender available")
	}

	return gender.(string), nil
}

func (s *service) getNation(name string) (string, error) {
	u := fmt.Sprintf("https://api.nationalize.io/?name=%s", name)
	data, err := s.getData(u)
	if err != nil {
		return "", err
	}

	country, ok := data["country"]
	if !ok {
		return "", fmt.Errorf("no country available")
	}

	d := country.([]interface{})
	nation := ""
	for _, f := range d {
		g := f.(map[string]interface{})
		if c, ok := g["country_id"]; !ok {
			return "", fmt.Errorf("no nation available")
		} else {
			nation = c.(string)
			break
		}
	}

	return nation, err
}

func (s *service) getData(u string) (map[string]interface{}, error) {
	var data map[string]interface{}

	resp, err := http.Get(u)
	if err != nil {
		return nil, err
	}

	d, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if err = json.Unmarshal(d, &data); err != nil {
		return nil, err
	}

	s.logger.Debug(data)

	return data, nil
}

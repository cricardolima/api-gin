package main

import (
	"api-gin/controllers"
	"api-gin/database"
	"api-gin/models"
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
)

var ID int

func SetupTests() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	routes := gin.Default()
	return routes
}

func MockStudent() {
	aluno := models.Aluno{Nome: "Mocked Student", CPF: "12345678910", RG: "123456789"}
	database.DB.Create(&aluno)
	ID = int(aluno.ID)
}

func DeleteMockStudent() {
	var aluno models.Aluno
	database.DB.Delete(&aluno, ID)
}

func TestGreetingRoute(t *testing.T) {
	r := SetupTests()
	r.GET("/:nome", controllers.Saudacao)
	req, _ := http.NewRequest("GET", "/Fulano4", nil)
	response := httptest.NewRecorder()
	r.ServeHTTP(response, req)
	assert.Equal(t, http.StatusOK, response.Code, "they should be equal")
	mockResponse := `{"API diz:":"E ai Fulano4, tudo beleza?"}`
	responseBody, _ := io.ReadAll(response.Body)
	assert.Equal(t, mockResponse, string(responseBody))
}

func TestGetAll(t *testing.T) {
	database.ConectaComBancoDeDados()
	MockStudent()
	defer DeleteMockStudent()
	r := SetupTests()
	r.GET("/alunos", controllers.ExibeTodosAlunos)
	req, _ := http.NewRequest("GET", "/alunos", nil)
	response := httptest.NewRecorder()
	r.ServeHTTP(response, req)
	assert.Equal(t, http.StatusOK, response.Code)
}

func TestGeyByCPF(t *testing.T) {
	database.ConectaComBancoDeDados()
	MockStudent()
	defer DeleteMockStudent()
	r := SetupTests()
	r.GET("/alunos/cpf/:cpf", controllers.BuscaAlunoPorCPF)
	req, _ := http.NewRequest("GET", "/alunos/cpf/12345678910", nil)
	response := httptest.NewRecorder()
	r.ServeHTTP(response, req)
	assert.Equal(t, http.StatusOK, response.Code)
}

func TestGetById(t *testing.T) {
	database.ConectaComBancoDeDados()
	MockStudent()
	defer DeleteMockStudent()
	r := SetupTests()
	r.GET("/alunos/:id", controllers.BuscaAlunoPorID)
	path := "/alunos/" + strconv.Itoa(ID)
	req, _ := http.NewRequest("GET", path, nil)
	response := httptest.NewRecorder()
	r.ServeHTTP(response, req)
	var mockedStudent models.Aluno
	json.Unmarshal(response.Body.Bytes(), &mockedStudent)
	assert.Equal(t, "Mocked Student", mockedStudent.Nome)
	assert.Equal(t, "12345678910", mockedStudent.CPF)
	assert.Equal(t, "123456789", mockedStudent.RG)
	assert.Equal(t, http.StatusOK, response.Code)
}

func TestDeleteStudent(t *testing.T) {
	database.ConectaComBancoDeDados()
	MockStudent()
	r := SetupTests()
	r.DELETE("/alunos/:id", controllers.DeletaAluno)
	path := "/alunos/" + strconv.Itoa(ID)
	req, _ := http.NewRequest("DELETE", path, nil)
	response := httptest.NewRecorder()
	r.ServeHTTP(response, req)
	assert.Equal(t, http.StatusOK, response.Code)
}

func TestPatchStudents(t *testing.T) {
	database.ConectaComBancoDeDados()
	MockStudent()
	defer DeleteMockStudent()
	r := SetupTests()
	r.PATCH("/alunos/:id", controllers.EditaAluno)
	student := models.Aluno{Nome: "Nome do Aluno Teste", CPF: "47815428912", RG: "123456788"}
	jsonValue, _ := json.Marshal(student)
	path := "/alunos/" + strconv.Itoa(ID)
	req, _ := http.NewRequest("PATCH", path, bytes.NewBuffer(jsonValue))
	response := httptest.NewRecorder()
	r.ServeHTTP(response, req)
	var mockedStudent models.Aluno
	json.Unmarshal(response.Body.Bytes(), &mockedStudent)
	assert.Equal(t, "47815428912", mockedStudent.CPF)
	assert.Equal(t, "123456788", mockedStudent.RG)
	assert.Equal(t, "Nome do Aluno Teste", mockedStudent.Nome)
}

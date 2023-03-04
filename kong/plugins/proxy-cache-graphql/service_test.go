package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type ServiceSuite struct {
	suite.Suite
	svc *Service
}

func (suite *ServiceSuite) SetupSuite() {
	suite.svc = NewService()
}

func (suite *ServiceSuite) TestGetAndNormalizeGraphQLAst() {
	requestBody := "{\"query\":\"query Query{country(code: \\\"TL\\\"){native,emoji,name,capital}}\"}"
	expectedAst, err := suite.svc.GetAndNormalizeGraphQLAst([]byte(requestBody))
	if err != nil {
		suite.T().Error(err.Error())
	}

	requestBodyList := []string{
		"{\"query\":\"query Query{country(code: \\\"TL\\\"){name,emoji,capital,native}}\"}",
		"{\"query\":\"query Query{country(code: \\\"TL\\\"){name,capital,native,emoji}}\"}",
		"{\"query\":\"query Query{country(code: \\\"TL\\\"){emoji,name,capital,native}}\"}",
		"{\"query\":\"query Query{country(code: \\\"TL\\\"){emoji,native,name,capital}}\"}",
	}

	for _, requestBodyTC := range requestBodyList {
		actualAst, err := suite.svc.GetAndNormalizeGraphQLAst([]byte(requestBodyTC))
		if err != nil {
			suite.T().Error(err.Error())
		}
		assert.Equal(suite.T(), string(actualAst), string(expectedAst))
	}
}

func TestServiceSuite(t *testing.T) {
	suite.Run(t, new(ServiceSuite))
}

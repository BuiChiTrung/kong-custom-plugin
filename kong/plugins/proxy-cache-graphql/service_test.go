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
	type TestCase struct {
		requestBody            string
		similarRequestBodyList []string
	}

	testcases := []TestCase{
		// 0. Change the order of field & add comment
		{
			requestBody: "query Query{country(code: \"VN\"){native,emoji,name,capital}}",
			similarRequestBodyList: []string{
				"query Query{country(code: \"VN\"){name,emoji,capital,native}}",
				"query Query{country(code: \"VN\"){name,capital,native,emoji}}",
				"query Query{country(code: \"VN\"){emoji,name,capital,native}}",
				"query Query{country(code: \"VN\"){emoji,native,name,capital}}",
				"query Query{country(code: \"VN\") {\n    native,    # country name\n    capital,\n    emoji,    \n    name,\n    #languages {code,name}\n}}",
			},
		},
		// 1. Change the order of argument & field
		{
			requestBody: "query Repository {\n  repository(followRenames: false,name: \"kong-custom-plugin\",owner: \"BuiChiTrung\",) {\n    owner {\n      id\n      avatarUrl\n      login  \n    }  \n    isPrivate\n    createdAt\n    autoMergeAllowed\n    allowUpdateBranch\n    id\n  }\n}\n",
			similarRequestBodyList: []string{
				"query Repository {\n  repository(followRenames: false,owner: \"BuiChiTrung\",name: \"kong-custom-plugin\") {\n    owner {\n      id\n      avatarUrl\n      login  \n    }  \n    isPrivate\n    createdAt\n    autoMergeAllowed\n    allowUpdateBranch\n    id\n  }\n}\n",
				"query Repository {\n  repository(name: \"kong-custom-plugin\",followRenames: false,owner: \"BuiChiTrung\") {\n    owner {\n      avatarUrl\n      id\n      login  \n    }  \n    autoMergeAllowed\n    allowUpdateBranch\n    id\n    isPrivate\n    createdAt\n  }\n}\n",
			},
		},
		// 2. Omit, change operation name, type
		{
			requestBody: "query Query{country(code: \"VN\"){native,emoji,name,capital}}",
			similarRequestBodyList: []string{
				"query AnotherQuery{country(code: \"VN\"){name,emoji,capital,native}}",
				"query {country(code: \"VN\"){name,capital,native,emoji}}",
				"{country(code: \"VN\"){emoji,name,capital,native}}",
			},
		},

		// 3. Change the order of variable & field
		{
			requestBody: "query Repository($name: String!, $owner: String!, $followRenames: Boolean) {\n  repository(name: $name, owner: $owner, followRenames: $followRenames) {\n    allowUpdateBranch\n    autoMergeAllowed\n    createdAt\n    id\n    isPrivate\n    owner {\n      avatarUrl\n      id\n      login\n      url\n    }\n  }\n}",
			similarRequestBodyList: []string{
				"query Repository($name: String!, $followRenames: Boolean, $owner: String!) {\n  repository(owner: $owner, followRenames: $followRenames, name: $name) {\n    allowUpdateBranch\n    autoMergeAllowed\n    createdAt\n    id\n    isPrivate\n    owner {\n      avatarUrl\n      id\n      login\n      url\n    }\n  }\n}",
				"query Repository($followRenames: Boolean, $name: String!, $owner: String!) {\n  repository(owner: $owner, name: $name, followRenames: $followRenames) {\n    allowUpdateBranch\n    autoMergeAllowed\n    id\n    createdAt\n    owner {\n      avatarUrl\n      id\n      login\n      url\n    }\n    isPrivate\n  }\n}",
			},
		},
		// 3. Mutation

		// 4. Fragment
	}

	skippedTestCases := []int{2}

	for i := 0; i < len(testcases); i++ {

		isSkip := false
		for _, skippedTestCase := range skippedTestCases {
			if i == skippedTestCase {
				isSkip = true
			}
		}
		if isSkip {
			continue
		}

		expectedAst, err := suite.svc.GetAndNormalizeGraphQLAst(testcases[i].requestBody)
		if err != nil {
			suite.T().Error(err.Error())
		}

		for _, similarRequestBody := range testcases[i].similarRequestBodyList {
			actualAst, err := suite.svc.GetAndNormalizeGraphQLAst(similarRequestBody)
			if err != nil {
				suite.T().Error(err.Error())
			}
			assert.Equal(suite.T(), getObjJSONString(expectedAst), getObjJSONString(actualAst))
		}

	}
}

func (suite *ServiceSuite) TestNormalizeGraphQLVariable() {
	type TestCase struct {
		variableMp            map[string]interface{}
		similarVariableMpList []map[string]interface{}
	}

	testcases := []TestCase{
		{
			variableMp: map[string]interface{}{"name": "kong-custom-plugin", "owner": "BuiChiTrung", "followRenames": true},
			similarVariableMpList: []map[string]interface{}{
				{"name": "kong-custom-plugin", "followRenames": true, "owner": "BuiChiTrung"},
				{"owner": "BuiChiTrung", "name": "kong-custom-plugin", "followRenames": true},
				{"owner": "BuiChiTrung", "followRenames": true, "name": "kong-custom-plugin"},
			},
		},
	}

	for i := 0; i < len(testcases); i++ {
		expectedVariableStr := suite.svc.NormalizeGraphQLVariable(testcases[i].variableMp)

		for _, similarVariableMp := range testcases[i].similarVariableMpList {
			actualVariableStr := suite.svc.NormalizeGraphQLVariable(similarVariableMp)
			assert.Equal(suite.T(), expectedVariableStr, actualVariableStr)
		}
	}
}

func TestServiceSuite(t *testing.T) {
	suite.Run(t, new(ServiceSuite))
}

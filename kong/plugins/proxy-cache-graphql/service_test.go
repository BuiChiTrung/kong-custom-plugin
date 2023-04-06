package main

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"reflect"
	"testing"
)

type ServiceSuite struct {
	suite.Suite
	svc *Service
}

func TestServiceSuite(t *testing.T) {
	suite.Run(t, new(ServiceSuite))
}

func (suite *ServiceSuite) SetupSuite() {
	suite.svc = &Service{}
}

func (suite *ServiceSuite) TestGenerateCacheKey() {
	type TestCase struct {
		testDesc               string
		testID                 int
		requestBody            string
		similarRequestBodyList []string
		shouldCached           bool
	}

	testCases := []TestCase{
		{
			testID:      0,
			testDesc:    "Github",
			requestBody: "{\"query\":\"query Repository($name: String!, $owner: String!, $followRenames: Boolean) {\\n  repository(name: $name, owner: $owner, followRenames: $followRenames) {\\n    allowUpdateBranch\\n    autoMergeAllowed\\n    id\\n    createdAt\\n    owner {\\n      avatarUrl\\n      id\\n      login\\n      url\\n    }\\n    isPrivate\\n  }\\n}\\n\",\"variables\":{\"name\":\"kong-custom-plugin\",\"owner\":\"BuiChiTrung\",\"followRenames\":true}}",
			similarRequestBodyList: []string{
				"{\"query\":\"query Repository($name: String!, $followRenames: Boolean, $owner: String!) {\\n  repository(owner: $owner, followRenames: $followRenames, name: $name) {\\n    allowUpdateBranch\\n    autoMergeAllowed\\n    createdAt\\n    id\\n    isPrivate\\n    owner {\\n      avatarUrl\\n      id\\n      login\\n      url\\n    }\\n  }\\n}\\n\",\"variables\":{\"name\":\"kong-custom-plugin\",\"owner\":\"BuiChiTrung\",\"followRenames\":true}}",
				"{\"query\":\"query Repository($followRenames: Boolean, $owner: String!, $name: String!) {\\n  repository(name: $name, owner: $owner, followRenames: $followRenames) {\\n    owner {\\n      avatarUrl\\n      id\\n      login\\n      url\\n    }\\n    allowUpdateBranch\\n    autoMergeAllowed\\n    createdAt\\n    id\\n    isPrivate\\n  }\\n}\\n\",\"variables\":{\"name\":\"kong-custom-plugin\",\"followRenames\":true,\"owner\":\"BuiChiTrung\"}}",
			},
			shouldCached: true,
		},
		{
			testID:      1,
			testDesc:    "Country",
			requestBody: "{\"query\":\"query {\\n    country(code: \\\"VN\\\") {\\n        native,\\n        capital,\\n        emoji,    \\n        name,\\n    }\\n}\",\"variables\":{}}",
			similarRequestBodyList: []string{
				"{\"query\":\"query {\\n    country(code: \\\"VN\\\") {\\n        emoji,    \\n        name,\\n        native,\\n        capital, #aaaa\\n    }\\n}\",\"variables\":{}}",
				"{\"query\":\"query Countrya {\\n    country(code: \\\"VN\\\") {\\n        name,native,emoji\\n        capital, #aaaa\\n    }\\n}\",\"variables\":{}}",
			},
			shouldCached: true,
		},
		{
			testID:      2,
			testDesc:    "Multiple query",
			requestBody: "{\"query\":\"query Repository($login: String!) {\\n  user(login: $login) {\\n    bio\\n    avatarUrl\\n  }\\n  repository(name: \\\"kong-custom-plugin\\\",followRenames: false,owner: \\\"BuiChiTrung\\\") {\\n    owner {\\n      avatarUrl\\n      id\\n      login  \\n    }  \\n    autoMergeAllowed\\n    allowUpdateBranch\\n    id\\n    isPrivate\\n    createdAt\\n  }\\n}\\n\",\"variables\":{\"name\":\"kong-custom-plugin\",\"followRenames\":true,\"owner\":\"BuiChiTrung\",\"login\":\"octocat\"}}",
			similarRequestBodyList: []string{
				"{\"query\":\"query Repository($login: String!) {\\n  repository(name: \\\"kong-custom-plugin\\\",followRenames: false,owner: \\\"BuiChiTrung\\\") {\\n    autoMergeAllowed\\n    allowUpdateBranch\\n    id\\n    isPrivate\\n    createdAt  \\n    owner {\\n      avatarUrl\\n      id\\n      login  \\n    }  \\n  }\\n  user(login: $login) {\\n    bio\\n    avatarUrl\\n  }\\n}\\n\",\"variables\":{\"name\":\"kong-custom-plugin\",\"followRenames\":true,\"login\":\"octocat\",\"owner\":\"BuiChiTrung\"}}",
				"{\"query\":\"query Repository($login: String!) {\\n  repository(followRenames: false,owner: \\\"BuiChiTrung\\\",name: \\\"kong-custom-plugin\\\") {\\n    autoMergeAllowed\\n    allowUpdateBranch\\n    id\\n    isPrivate\\n    createdAt  \\n    owner {\\n      id\\n      login  \\n      avatarUrl\\n    }  \\n  }\\n  user(login: $login) {\\n    avatarUrl\\n    bio\\n  }\\n}\\n\",\"variables\":{\"name\":\"kong-custom-plugin\",\"followRenames\":true,\"login\":\"octocat\",\"owner\":\"BuiChiTrung\"}}",
			},
			shouldCached: true,
		},
		{
			testID:       3,
			testDesc:     "Shouldn't cache mutation operation",
			requestBody:  "{\"query\":\"mutation AddReactionToIssue {\\n  addReaction(input:{subjectId:\\\"I_kwDOABPHjc5fqsV6\\\",content:HOORAY}) {\\n    reaction {\\n      content\\n    }\\n    subject {\\n      id\\n    }\\n  }\\n}\",\"variables\":{}}",
			shouldCached: false,
		},
	}

	skippedTestCases := []int{}

	for _, testCase := range testCases {
		isSkip := false
		for _, skippedTestCase := range skippedTestCases {
			if testCase.testID == skippedTestCase {
				isSkip = true
			}
		}
		if isSkip {
			continue
		}

		expectedCacheKey, shouldCached, err := suite.svc.GenerateCacheKey(testCase.requestBody, "", "")
		if err != nil {
			suite.T().Error(err.Error())
		}

		assert.Equal(suite.T(), testCase.shouldCached, shouldCached)

		for _, similarRequestBody := range testCase.similarRequestBodyList {
			actualCacheKey, _, err := suite.svc.GenerateCacheKey(similarRequestBody, "", "")
			if err != nil {
				suite.T().Error(err.Error())
			}
			suite.T().Log(expectedCacheKey, actualCacheKey)
			assert.Equal(suite.T(), expectedCacheKey, actualCacheKey)
		}
	}
}

func (suite *ServiceSuite) TestGetAndNormalizeGraphQLAst() {
	type TestCase struct {
		testDesc                string
		testID                  int
		graphQLQuery            string
		similarGraphQLQueryList []string
	}

	testCases := []TestCase{
		{
			testID:       0,
			testDesc:     "Change the order of field & add comment",
			graphQLQuery: "query Query{country(code: \"VN\"){native,emoji,name,capital}}",
			similarGraphQLQueryList: []string{
				"query Query{country(code: \"VN\"){name,emoji,capital,native}}",
				"query Query{country(code: \"VN\"){name,capital,native,emoji}}",
				"query Query{country(code: \"VN\"){emoji,name,capital,native}}",
				"query Query{country(code: \"VN\"){emoji,native,name,capital}}",
				"query Query{country(code: \"VN\") {\n    native,    # country name\n    capital,\n    emoji,    \n    name,\n    #languages {code,name}\n}}",
			},
		},
		{
			testID:       1,
			testDesc:     "Change the order of argument & field",
			graphQLQuery: "query Repository {\n  repository(followRenames: false,name: \"kong-custom-plugin\",owner: \"BuiChiTrung\",) {\n    owner {\n      id\n      avatarUrl\n      login  \n    }  \n    isPrivate\n    createdAt\n    autoMergeAllowed\n    allowUpdateBranch\n    id\n  }\n}\n",
			similarGraphQLQueryList: []string{
				"query Repository {\n  repository(followRenames: false,owner: \"BuiChiTrung\",name: \"kong-custom-plugin\") {\n    owner {\n      id\n      avatarUrl\n      login  \n    }  \n    isPrivate\n    createdAt\n    autoMergeAllowed\n    allowUpdateBranch\n    id\n  }\n}\n",
				"query Repository {\n  repository(name: \"kong-custom-plugin\",followRenames: false,owner: \"BuiChiTrung\") {\n    owner {\n      avatarUrl\n      id\n      login  \n    }  \n    autoMergeAllowed\n    allowUpdateBranch\n    id\n    isPrivate\n    createdAt\n  }\n}\n",
			},
		},
		{
			testID:       2,
			testDesc:     "Omit, change operation name, type",
			graphQLQuery: "query Query{country(code: \"VN\"){native,emoji,name,capital}}",
			similarGraphQLQueryList: []string{
				"query AnotherQuery{country(code: \"VN\"){name,emoji,capital,native}}",
				"query {country(code: \"VN\"){name,capital,native,emoji}}",
				"{country(code: \"VN\"){emoji,name,capital,native}}",
			},
		},
		{
			testID:       3,
			testDesc:     "Change the order of variable & field",
			graphQLQuery: "query Repository($name: String!, $owner: String!, $followRenames: Boolean) {\n  repository(name: $name, owner: $owner, followRenames: $followRenames) {\n    allowUpdateBranch\n    autoMergeAllowed\n    createdAt\n    id\n    isPrivate\n    owner {\n      avatarUrl\n      id\n      login\n      url\n    }\n  }\n}",
			similarGraphQLQueryList: []string{
				"query Repository($name: String!, $followRenames: Boolean, $owner: String!) {\n  repository(owner: $owner, followRenames: $followRenames, name: $name) {\n    allowUpdateBranch\n    autoMergeAllowed\n    createdAt\n    id\n    isPrivate\n    owner {\n      avatarUrl\n      id\n      login\n      url\n    }\n  }\n}",
				"query Repository($followRenames: Boolean, $name: String!, $owner: String!) {\n  repository(owner: $owner, name: $name, followRenames: $followRenames) {\n    allowUpdateBranch\n    autoMergeAllowed\n    id\n    createdAt\n    owner {\n      avatarUrl\n      id\n      login\n      url\n    }\n    isPrivate\n  }\n}",
			},
		},
		// 3. Mutation

		// 4. Fragment
	}

	skippedTestCases := []int{}

	for _, testCase := range testCases {
		isSkip := false
		for _, skippedTestCase := range skippedTestCases {
			if testCase.testID == skippedTestCase {
				isSkip = true
			}
		}
		if isSkip {
			continue
		}

		expectedAst, err := suite.svc.GetGraphQLAst(testCase.graphQLQuery)
		if err != nil {
			suite.T().Error(err.Error())
		}
		suite.svc.NormalizeOperationName(expectedAst)
		suite.svc.NormalizeGraphQLAST(reflect.ValueOf(expectedAst).Elem())

		for _, similarGraphQLQuery := range testCase.similarGraphQLQueryList {
			actualAst, err := suite.svc.GetGraphQLAst(similarGraphQLQuery)
			if err != nil {
				suite.T().Error(err.Error())
			}
			suite.svc.NormalizeOperationName(actualAst)
			suite.svc.NormalizeGraphQLAST(reflect.ValueOf(actualAst).Elem())
			//fmt.Println(getObjJSONString(expectedAst))
			//fmt.Println(getObjJSONString(actualAst))
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

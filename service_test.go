package spreadsheet

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

const spreadsheetID = "1mYiA2T4_QTFUkAXk0BE3u7snN2o5FgSRqxmRrn_Dzh4"

type TestSuite struct {
	suite.Suite
	service *Service
}

func (suite *TestSuite) SetupSuite() {
	var err error
	suite.service, err = NewService()
	suite.Require().NoError(err)
}

func (suite *TestSuite) TestFetchSpreadsheet() {
	spreadsheet, err := suite.service.FetchSpreadsheet(spreadsheetID)
	suite.Require().NoError(err)
	suite.Equal(spreadsheetID, spreadsheet.ID)
	suite.Require().Equal(2, len(spreadsheet.Sheets))

	sheet := spreadsheet.Sheets[0]
	suite.Equal(uint(0), sheet.Properties.ID)
	suite.Equal("TestSheet", sheet.Properties.Title)
	suite.Equal(uint(0), sheet.Properties.Index)
	suite.Equal("GRID", sheet.Properties.SheetType)
	suite.True(len(sheet.Rows) >= 3)
	suite.True(len(sheet.Columns) >= 3)
	suite.Equal(uint(2), sheet.Rows[1][2].Column)
}

func (suite *TestSuite) TestSyncSheet() {
	spreadsheet, err := suite.service.FetchSpreadsheet(spreadsheetID)
	suite.Require().NoError(err)
	sheet, err := spreadsheet.SheetByTitle("TestSheet")
	suite.Require().NoError(err)
	sheet.Update(1, 6, "=SUM(D1:D2)")
	err = suite.service.SyncSheet(sheet)
	suite.NoError(err)
}

func (suite *TestSuite) TestDeleteRows() {
	spreadsheet, err := suite.service.FetchSpreadsheet(spreadsheetID)
	suite.Require().NoError(err)
	sheet, err := spreadsheet.SheetByTitle("TestSheet2")
	suite.Require().NoError(err)
	rowCount := sheet.Properties.GridProperties.RowCount

	err = sheet.DeleteRows(0, 1)
	suite.NoError(err)
	suite.Equal(rowCount-1, sheet.Properties.GridProperties.RowCount)
}

func TestRun(t *testing.T) {
	suite.Run(t, new(TestSuite))
}
